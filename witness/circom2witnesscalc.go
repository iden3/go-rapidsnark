package witness

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/wasmerio/wasmer-go/wasmer"
)

// Circom2WitnessCalculator is the object that allows performing witness calculation
// from signal inputs using the WitnessCalc WASM module.
type Circom2WitnessCalculator struct {
	instance            *wasmer.Instance
	store               *wasmer.Store
	sanityCheck         bool
	n32                 int32
	version             int32
	witnessSize         int32
	prime               *big.Int
	init                wasmer.NativeFunction
	getFieldNumLen32    wasmer.NativeFunction
	getInputSignalSize  wasmer.NativeFunction
	getInputSize        wasmer.NativeFunction
	getRawPrime         wasmer.NativeFunction
	getVersion          wasmer.NativeFunction
	getWitness          wasmer.NativeFunction
	readSharedRWMemory  wasmer.NativeFunction
	setInputSignal      wasmer.NativeFunction
	writeSharedRWMemory wasmer.NativeFunction
	getMessageChar      wasmer.NativeFunction
	exception           error
	errStr              bytes.Buffer
	msgStr              bytes.Buffer
}

// NewCircom2WitnessCalculator creates a new WitnessCalculator from the WitnessCalc
// loaded WASM module in the runtime.
func NewCircom2WitnessCalculator(wasmBytes []byte, sanityCheck bool) (*Circom2WitnessCalculator, error) {
	wc := Circom2WitnessCalculator{}
	wc.sanityCheck = sanityCheck

	engine := wasmer.NewEngine()
	wc.store = wasmer.NewStore(engine)

	// Compiles the module
	module, _ := wasmer.NewModule(wc.store, wasmBytes)

	limits, err := wasmer.NewLimits(2000, 100000)
	if err != nil {
		return nil, err
	}

	memType := wasmer.NewMemoryType(limits)

	memory := wasmer.NewMemory(wc.store, memType)

	// Instantiates the module
	importObject := wasmer.NewImportObject()

	importObject.Register("env", map[string]wasmer.IntoExtern{
		"memory": memory,
	})

	importObject.Register("runtime", map[string]wasmer.IntoExtern{
		"exceptionHandler":   wc.getExceptionHandler(),
		"showSharedRWMemory": wc.getShowSharedRWMemoryHandler(),
		"log":                wc.getLogHandler(),
		"printErrorMessage":  wc.printErrorMessageHandler(),
		"writeBufferMessage": wc.writeBufferMessageHandler(),
	})

	wc.instance, err = wasmer.NewInstance(module, importObject)
	if err != nil {
		return nil, err
	}

	// Gets the `init` exported function from the WebAssembly instance.
	init, err := wc.instance.Exports.GetFunction("init")
	if err != nil {
		return nil, err
	}

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	_, err = init(1)
	if err != nil {
		return nil, err
	}

	getFieldNumLen32, err := wc.instance.Exports.GetFunction("getFieldNumLen32")
	if err != nil {
		return nil, err
	}
	n32raw, err := getFieldNumLen32()
	if err != nil {
		return nil, err
	}
	wc.n32 = n32raw.(int32)

	// this function is missing in wasm files generated with circom version prior to v2.0.4
	getInputSignalSize, _ := wc.instance.Exports.GetFunction("getInputSignalSize")

	getInputSize, err := wc.instance.Exports.GetFunction("getInputSize")
	if err != nil {
		return nil, err
	}

	getRawPrime, err := wc.instance.Exports.GetFunction("getRawPrime")
	if err != nil {
		return nil, err
	}

	getVersion, err := wc.instance.Exports.GetFunction("getVersion")
	if err != nil {
		return nil, err
	}

	version, err := getVersion()
	if err != nil {
		return nil, err
	}

	getWitness, err := wc.instance.Exports.GetFunction("getWitness")
	if err != nil {
		return nil, err
	}

	getWitnessSize, err := wc.instance.Exports.GetFunction("getWitnessSize")
	if err != nil {
		return nil, err
	}

	witnessSize, err := getWitnessSize()
	if err != nil {
		return nil, err
	}

	setInputSignal, err := wc.instance.Exports.GetFunction("setInputSignal")
	if err != nil {
		return nil, err
	}

	readSharedRWMemory, err := wc.instance.Exports.GetFunction("readSharedRWMemory")
	if err != nil {
		return nil, err
	}

	writeSharedRWMemory, err := wc.instance.Exports.GetFunction("writeSharedRWMemory")
	if err != nil {
		return nil, err
	}

	getMessageChar, err := wc.instance.Exports.GetFunction("getMessageChar")
	if err != nil {
		return nil, err
	}

	//get prime number
	_, err = getRawPrime()
	if err != nil {
		return nil, err
	}
	primeArr := make([]uint32, wc.n32)
	for j := 0; j < int(wc.n32); j++ {
		val, err := readSharedRWMemory(int32(j))
		if err != nil {
			return nil, err
		}
		primeArr[int(wc.n32)-1-j] = uint32(val.(int32))
	}
	prime := fromArray32(primeArr)

	wc.version = version.(int32)
	wc.witnessSize = witnessSize.(int32)
	wc.prime = prime
	wc.init = init
	wc.getFieldNumLen32 = getFieldNumLen32
	wc.getInputSignalSize = getInputSignalSize
	wc.getInputSize = getInputSize
	wc.getRawPrime = getRawPrime
	wc.getWitness = getWitness
	wc.getVersion = getVersion
	wc.setInputSignal = setInputSignal
	wc.readSharedRWMemory = readSharedRWMemory
	wc.writeSharedRWMemory = writeSharedRWMemory
	wc.getMessageChar = getMessageChar

	return &wc, nil
}

// CalculateWitness calculates the witness given the inputs.
func (wc *Circom2WitnessCalculator) CalculateWitness(inputs map[string]interface{}, sanityCheck bool) ([]*big.Int, error) {

	w := make([]*big.Int, wc.witnessSize)

	err := wc.doCalculateWitness(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(wc.witnessSize); i++ {
		_, err := wc.getWitness(i)
		if err != nil {
			return nil, err
		}
		arr := make([]uint32, wc.n32)
		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory(int32(j))
			if err != nil {
				return nil, err
			}
			arr[int(wc.n32)-1-j] = uint32(val.(int32))
		}
		w[i] = fromArray32(arr)
	}

	return w, nil
}

// CalculateBinWitness calculates the witness in binary given the inputs.
func (wc *Circom2WitnessCalculator) CalculateBinWitness(inputs map[string]interface{}, sanityCheck bool) ([]byte, error) {
	buff := new(bytes.Buffer)

	err := wc.doCalculateWitness(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(wc.witnessSize); i++ {
		_, err := wc.getWitness(i)
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory(j)
			if err != nil {
				return nil, err
			}
			_ = binary.Write(buff, binary.LittleEndian, uint32(val.(int32)))
		}
	}

	return buff.Bytes(), nil
}

// CalculateWTNSBin calculates the witness in binary given the inputs.
func (wc *Circom2WitnessCalculator) CalculateWTNSBin(inputs map[string]interface{}, sanityCheck bool) ([]byte, error) {
	buff := new(bytes.Buffer)

	err := wc.doCalculateWitness(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	buff.Grow(int(wc.witnessSize*wc.n32 + wc.n32 + 11))

	// wtns
	_ = buff.WriteByte('w')
	_ = buff.WriteByte('t')
	_ = buff.WriteByte('n')
	_ = buff.WriteByte('s')

	//version 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	//number of sections: 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	//id section 1
	_ = binary.Write(buff, binary.LittleEndian, uint32(1))

	n8 := wc.n32 * 4
	//id section 1 length in 64bytes
	idSection1length := 8 + n8
	_ = binary.Write(buff, binary.LittleEndian, uint64(idSection1length))

	//this.n32
	_ = binary.Write(buff, binary.LittleEndian, uint32(n8))

	//prime number
	_, err = wc.getRawPrime()
	if err != nil {
		return nil, err
	}

	for j := 0; j < int(wc.n32); j++ {
		val, err := wc.readSharedRWMemory(int32(j))
		if err != nil {
			return nil, err
		}
		_ = binary.Write(buff, binary.LittleEndian, uint32(val.(int32)))
	}

	// witness size
	_ = binary.Write(buff, binary.LittleEndian, uint32(wc.witnessSize))

	//id section 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	// section 2 length
	idSection2length := n8 * wc.witnessSize
	_ = binary.Write(buff, binary.LittleEndian, uint64(idSection2length))

	for i := 0; i < int(wc.witnessSize); i++ {
		_, err := wc.getWitness(i)
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory(j)
			if err != nil {
				return nil, err
			}
			_ = binary.Write(buff, binary.LittleEndian, uint32(val.(int32)))
		}
	}

	return buff.Bytes(), nil
}

// CalculateWitness calculates the witness given the inputs.
func (wc *Circom2WitnessCalculator) doCalculateWitness(inputs map[string]interface{}, sanityCheck bool) (funcErr error) {
	//input is assumed to be a map from signals to arrays of bigInts
	sanityCheckVal := int32(0)
	if sanityCheck {
		sanityCheckVal = 1
	}

	wc.exception = nil
	wc.errStr.Reset()
	wc.msgStr.Reset()

	_, err := wc.init(sanityCheckVal)
	if err != nil {
		return err
	}

	// overwrite return error if there was an exception during execution
	defer func() {
		if wc.exception != nil {
			funcErr = wc.exception
		}
	}()

	inputCounter := 0
	for inputName, inputValue := range inputs {
		hMSB, hLSB := fnvHash(inputName)
		fSlice := flatSlice(inputValue)

		if wc.getInputSignalSize != nil {
			signalSize, err := wc.getInputSignalSize(hMSB, hLSB)
			if err != nil {
				return err
			}

			if signalSize.(int32) < 0 {
				return fmt.Errorf("signal %s not found", inputName)
			}
			if len(fSlice) < int(signalSize.(int32)) {
				return fmt.Errorf("not enough values for input signal %s", inputName)
			}
			if len(fSlice) > int(signalSize.(int32)) {
				return fmt.Errorf("too many values for input signal %s", inputName)
			}
		}

		for i := 0; i < len(fSlice); i++ {
			// doing val = (val + prime) % prime
			val := new(big.Int)
			val = val.Add(fSlice[i], wc.prime)
			val = val.Mod(val, wc.prime)
			arrFr, err := toArray32(val, int(wc.n32))
			if err != nil {
				return err
			}
			for j := 0; j < int(wc.n32); j++ {
				_, err := wc.writeSharedRWMemory(j, int32(arrFr[int(wc.n32)-1-j]))
				if err != nil {
					return err
				}
			}
			_, err = wc.setInputSignal(hMSB, hLSB, i)
			if err != nil {
				return err
			}
			inputCounter++
		}
	}
	inputSize, err := wc.getInputSize()
	if inputCounter < int(inputSize.(int32)) {
		return fmt.Errorf("not all inputs have been set: only %d out of %d", inputCounter, inputSize)
	}
	return nil
}

func (wc *Circom2WitnessCalculator) getExceptionHandler() wasmer.IntoExtern {
	function := wasmer.NewFunction(
		wc.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32), // one i32 argument
			wasmer.NewValueTypes(),           // zero results
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			if len(args) > 0 {
				code := args[0].I32()
				var errStr string
				if code == 1 {
					errStr = "Signal not found"
				} else if code == 2 {
					errStr = "Too many signals set"
				} else if code == 3 {
					errStr = "Signal already set"
				} else if code == 4 {
					errStr = "Assert Failed"
				} else if code == 5 {
					errStr = "Not enough memory"
				} else if code == 6 {
					errStr = "Input signal array access exceeds the size"
				} else {
					errStr = "Unknown error"
				}
				// Append stack trace to error message
				if wc.errStr.Len() > 0 {
					errStr += ".\n" + wc.errStr.String()
				}
				// returning error here crashes wasmer for all following witness calculation calls,
				// so we have to use a field to pass exception to the outside world
				wc.exception = errors.New(errStr)
				//fmt.Println(errStr)
				//return nil, errors.New(errStr)
			}
			return []wasmer.Value{}, nil
		},
	)
	return function
}

func (wc *Circom2WitnessCalculator) getShowSharedRWMemoryHandler() wasmer.IntoExtern {
	function := wasmer.NewFunction(
		wc.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			arr := make([]uint32, wc.n32)
			for j := 0; j < int(wc.n32); j++ {
				val, err := wc.readSharedRWMemory(int32(j))
				if err != nil {
					return nil, err
				}
				arr[int(wc.n32)-1-j] = uint32(val.(int32))
			}

			// If we've buffered other content, put a space in between the items
			if wc.msgStr.Len() > 0 {
				wc.msgStr.WriteString(" ")
			}
			// Then append the value to the message we are creating
			wc.msgStr.WriteString(fromArray32(arr).String())
			return []wasmer.Value{}, nil
		},
	)
	return function
}

func (wc *Circom2WitnessCalculator) getLogHandler() wasmer.IntoExtern {
	function := wasmer.NewFunction(
		wc.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			return []wasmer.Value{}, nil
		},
	)
	return function
}

func (wc *Circom2WitnessCalculator) getMessage() (string, error) {
	message := ""
	c, err := wc.getMessageChar()
	if err != nil {
		return "", err
	}
	for c.(int32) != 0 {
		message += string(c.(int32))
		c, err = wc.getMessageChar()
		if err != nil {
			return message, err
		}
	}
	return message, nil
}

func (wc *Circom2WitnessCalculator) printErrorMessageHandler() wasmer.IntoExtern {
	function := wasmer.NewFunction(
		wc.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(), // zero results
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			message, _ := wc.getMessage()
			wc.errStr.WriteString(message + "\n")
			return []wasmer.Value{}, nil
		},
	)
	return function
}

func (wc *Circom2WitnessCalculator) writeBufferMessageHandler() wasmer.IntoExtern {
	function := wasmer.NewFunction(
		wc.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(),
			wasmer.NewValueTypes(), // zero results
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			msg, _ := wc.getMessage()
			// Any calls to `log()` will always end with a `\n`, so that's when we print and reset
			if msg == "\n" {
				fmt.Println(wc.msgStr.String())
				wc.msgStr.Reset()
			} else {
				// If we've buffered other content, put a space in between the items
				if wc.msgStr.Len() > 0 {
					wc.msgStr.WriteString(" ")
				}
				// Then append the message to the message we are creating
				wc.msgStr.WriteString(msg)
			}
			return []wasmer.Value{}, nil
		},
	)
	return function
}

func toArray32(s *big.Int, size int) ([]uint32, error) {
	res := make([]uint32, size)
	rem := s

	radix := big.NewInt(0x100000000)
	zero := big.NewInt(0)
	i := size - 1
	// while not zero rem
	for rem.Cmp(zero) != 0 {
		res[i] = uint32(new(big.Int).Mod(rem, radix).Uint64())
		rem.Div(rem, radix)
		i--
	}
	return res, nil
}

func fromArray32(arr []uint32) *big.Int {
	res := new(big.Int)
	radix := big.NewInt(0x100000000)
	for i := 0; i < len(arr); i++ {
		res.Mul(res, radix)
		res.Add(res, big.NewInt(int64(arr[i])))
	}
	return res
}
