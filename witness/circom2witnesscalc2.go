package witness

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Circom2WitnessCalculator2 is the object that allows performing witness calculation
// from signal inputs using the WitnessCalc WASM module.
type Circom2WitnessCalculator2 struct {
	runtime             wazero.Runtime
	mod                 *api.Module
	sanityCheck         bool
	n32                 int32
	version             int32
	witnessSize         int32
	prime               *big.Int
	init                api.Function
	getFieldNumLen32    api.Function
	getInputSignalSize  api.Function
	getInputSize        api.Function
	getRawPrime         api.Function
	getVersion          api.Function
	getWitness          api.Function
	readSharedRWMemory  api.Function
	setInputSignal      api.Function
	writeSharedRWMemory api.Function
	getMessageChar      api.Function
	exception           error
	errStr              bytes.Buffer
	msgStr              bytes.Buffer
}

// NewCircom2WitnessCalculator2 creates a new WitnessCalculator from the WitnessCalc
// loaded WASM module in the runtime.
func NewCircom2WitnessCalculator2(wasmBytes []byte, sanityCheck bool) (*Circom2WitnessCalculator2, error) {
	wc := Circom2WitnessCalculator2{}
	wc.sanityCheck = sanityCheck

	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	wc.runtime = wazero.NewRuntime(ctx)
	//defer func() {
	//	_ = wc.runtime.Close(ctx) // This closes everything this Runtime created.
	//}()

	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	//wasi_snapshot_preview1.MustInstantiate(ctx, wc.runtime)

	_, err := wc.runtime.NewHostModuleBuilder("runtime").
		NewFunctionBuilder().
		WithFunc(wc.getExceptionHandler()).
		Export("exceptionHandler").
		NewFunctionBuilder().
		WithFunc(wc.getShowSharedRWMemoryHandler()).
		Export("showSharedRWMemory").
		NewFunctionBuilder().
		WithFunc(wc.getPrintErrorMessageHandler()).
		Export("printErrorMessage").
		NewFunctionBuilder().
		WithFunc(wc.getWriteBufferMessageHandler()).
		Export("writeBufferMessage").
		Instantiate(ctx)

	// Instantiate the guest Wasm into the same runtime. It exports the `add`
	// function, implemented in WebAssembly.
	mod, err := wc.runtime.InstantiateModuleFromBinary(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	//importObject.Register("runtime", map[string]wasmer.IntoExtern{
	//	"exceptionHandler":   wc.getExceptionHandler(),
	//	"showSharedRWMemory": wc.getShowSharedRWMemoryHandler(),
	//	"log":                wc.getLogHandler(),
	//	"printErrorMessage":  wc.printErrorMessageHandler(),
	//	"writeBufferMessage": wc.writeBufferMessageHandler(),
	//})

	// Gets the `init` exported function from the WebAssembly instance.
	init := mod.ExportedFunction("init")
	if init == nil {
		return nil, errors.New("no function init")
	}

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	_, err = init.Call(ctx, 1)
	if err != nil {
		return nil, err
	}

	getFieldNumLen32 := mod.ExportedFunction("getFieldNumLen32")
	if getFieldNumLen32 == nil {
		return nil, errors.New("no function getFieldNumLen32")
	}
	n32raw, err := getFieldNumLen32.Call(ctx)
	if err != nil {
		return nil, err
	}
	wc.n32 = int32(n32raw[0])

	// this function is missing in wasm files generated with circom version prior to v2.0.4
	getInputSignalSize := mod.ExportedFunction("getInputSignalSize")

	getInputSize := mod.ExportedFunction("getInputSize")
	if getInputSize == nil {
		return nil, errors.New("no function getInputSize")
	}

	getRawPrime := mod.ExportedFunction("getRawPrime")
	if getRawPrime == nil {
		return nil, errors.New("no function getRawPrime")
	}

	getVersion := mod.ExportedFunction("getVersion")
	if getVersion == nil {
		return nil, errors.New("no function getVersion")
	}

	version, err := getVersion.Call(ctx)
	if err != nil {
		return nil, err
	}

	getWitness := mod.ExportedFunction("getWitness")
	if getWitness == nil {
		return nil, errors.New("no function getWitness")
	}

	getWitnessSize := mod.ExportedFunction("getWitnessSize")
	if getWitnessSize == nil {
		return nil, errors.New("no function getWitnessSize")
	}

	witnessSize, err := getWitnessSize.Call(ctx)
	if err != nil {
		return nil, err
	}

	setInputSignal := mod.ExportedFunction("setInputSignal")
	if setInputSignal == nil {
		return nil, errors.New("no function setInputSignal")
	}

	readSharedRWMemory := mod.ExportedFunction("readSharedRWMemory")
	if readSharedRWMemory == nil {
		return nil, errors.New("no function readSharedRWMemory")
	}

	writeSharedRWMemory := mod.ExportedFunction("writeSharedRWMemory")
	if writeSharedRWMemory == nil {
		return nil, errors.New("no function writeSharedRWMemory")
	}

	getMessageChar := mod.ExportedFunction("getMessageChar")
	if getMessageChar == nil {
		return nil, errors.New("no function getMessageChar")
	}

	//get prime number
	_, err = getRawPrime.Call(ctx)
	if err != nil {
		return nil, err
	}
	primeArr := make([]uint32, wc.n32)
	for j := 0; j < int(wc.n32); j++ {
		val, err := readSharedRWMemory.Call(ctx, uint64(j))
		if err != nil {
			return nil, err
		}
		primeArr[int(wc.n32)-1-j] = uint32(val[0])
	}
	prime := fromArray32(primeArr)

	wc.version = int32(version[0])
	wc.witnessSize = int32(witnessSize[0])
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
func (wc *Circom2WitnessCalculator2) CalculateWitness(inputs map[string]interface{}, sanityCheck bool) ([]*big.Int, error) {

	w := make([]*big.Int, wc.witnessSize)

	err := wc.doCalculateWitness(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	for i := 0; i < int(wc.witnessSize); i++ {
		_, err := wc.getWitness.Call(ctx, uint64(i))
		if err != nil {
			return nil, err
		}
		arr := make([]uint32, wc.n32)
		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory.Call(ctx, uint64(j))
			if err != nil {
				return nil, err
			}
			arr[int(wc.n32)-1-j] = uint32(val[0])
		}
		w[i] = fromArray32(arr)
	}

	return w, nil
}

// CalculateBinWitness calculates the witness in binary given the inputs.
func (wc *Circom2WitnessCalculator2) CalculateBinWitness(inputs map[string]interface{}, sanityCheck bool) ([]byte, error) {
	buff := new(bytes.Buffer)

	err := wc.doCalculateWitness(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	for i := 0; i < int(wc.witnessSize); i++ {
		_, err := wc.getWitness.Call(ctx, uint64(i))
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory.Call(ctx, uint64(j))
			if err != nil {
				return nil, err
			}
			_ = binary.Write(buff, binary.LittleEndian, uint32(val[0]))
		}
	}

	return buff.Bytes(), nil
}

// CalculateWTNSBin calculates the witness in binary given the inputs.
func (wc *Circom2WitnessCalculator2) CalculateWTNSBin(inputs map[string]interface{}, sanityCheck bool) ([]byte, error) {
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

	ctx := context.Background()

	//prime number
	_, err = wc.getRawPrime.Call(ctx)
	if err != nil {
		return nil, err
	}

	for j := 0; j < int(wc.n32); j++ {
		val, err := wc.readSharedRWMemory.Call(ctx, uint64(j))
		if err != nil {
			return nil, err
		}
		_ = binary.Write(buff, binary.LittleEndian, uint32(val[0]))
	}

	// witness size
	_ = binary.Write(buff, binary.LittleEndian, uint32(wc.witnessSize))

	//id section 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	// section 2 length
	idSection2length := n8 * wc.witnessSize
	_ = binary.Write(buff, binary.LittleEndian, uint64(idSection2length))

	for i := 0; i < int(wc.witnessSize); i++ {
		_, err := wc.getWitness.Call(ctx, uint64(i))
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory.Call(ctx, uint64(j))
			if err != nil {
				return nil, err
			}
			_ = binary.Write(buff, binary.LittleEndian, uint32(val[0]))
		}
	}

	return buff.Bytes(), nil
}

// CalculateWitness calculates the witness given the inputs.
func (wc *Circom2WitnessCalculator2) doCalculateWitness(inputs map[string]interface{}, sanityCheck bool) (funcErr error) {
	//input is assumed to be a map from signals to arrays of bigInts
	sanityCheckVal := uint64(0)
	if sanityCheck {
		sanityCheckVal = 1
	}

	wc.exception = nil
	wc.errStr.Reset()
	wc.msgStr.Reset()

	ctx := context.Background()

	_, err := wc.init.Call(ctx, sanityCheckVal)
	if err != nil {
		return err
	}

	// overwrite return error if there was an exception during execution
	defer func() {
		if wc.exception != nil {
			funcErr = wc.exception
		}
		_ = wc.runtime.Close(context.Background()) // This closes everything this Runtime created.
	}()

	inputCounter := 0
	for inputName, inputValue := range inputs {
		hMSB, hLSB := fnvHash(inputName)
		fSlice := flatSlice(inputValue)

		if wc.getInputSignalSize != nil {
			signalSize, err := wc.getInputSignalSize.Call(ctx, uint64(uint32(hMSB)), uint64(uint32(hLSB)))
			if err != nil {
				return err
			}

			if signalSize[0] < 0 {
				return fmt.Errorf("signal %s not found", inputName)
			}
			if len(fSlice) < int(signalSize[0]) {
				return fmt.Errorf("not enough values for input signal %s", inputName)
			}
			if len(fSlice) > int(signalSize[0]) {
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
				_, err := wc.writeSharedRWMemory.Call(ctx, uint64(j), uint64(arrFr[int(wc.n32)-1-j]))
				if err != nil {
					return err
				}
			}
			_, err = wc.setInputSignal.Call(ctx, uint64(hMSB), uint64(hLSB), uint64(i))
			if err != nil {
				return err
			}
			inputCounter++
		}
	}
	inputSize, err := wc.getInputSize.Call(ctx)
	if inputCounter < int(inputSize[0]) {
		return fmt.Errorf("not all inputs have been set: only %d out of %d", inputCounter, inputSize)
	}
	return nil
}

func (wc *Circom2WitnessCalculator2) getExceptionHandler() func(uint32) {
	return func(code uint32) {
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
}

func (wc *Circom2WitnessCalculator2) getShowSharedRWMemoryHandler() func() {
	return func() {
		arr := make([]uint32, wc.n32)
		ctx := context.Background()
		for j := 0; j < int(wc.n32); j++ {
			val, err := wc.readSharedRWMemory.Call(ctx, uint64(j))
			if err != nil {
				//panic(err)
			}
			arr[int(wc.n32)-1-j] = uint32(val[0])
		}

		// If we've buffered other content, put a space in between the items
		if wc.msgStr.Len() > 0 {
			wc.msgStr.WriteString(" ")
		}
		// Then append the value to the message we are creating
		wc.msgStr.WriteString(fromArray32(arr).String())
	}
}

//func (wc *Circom2WitnessCalculator2) getLogHandler() wasmer.IntoExtern {
//	function := wasmer.NewFunction(
//		wc.store,
//		wasmer.NewFunctionType(
//			wasmer.NewValueTypes(),
//			wasmer.NewValueTypes(),
//		),
//		func(args []wasmer.Value) ([]wasmer.Value, error) {
//			return []wasmer.Value{}, nil
//		},
//	)
//	return function
//}

func (wc *Circom2WitnessCalculator2) getMessage() (string, error) {
	message := ""
	ctx := context.Background()
	c, err := wc.getMessageChar.Call(ctx)
	if err != nil {
		return "", err
	}
	for len(c) > 0 && c[0] != 0 {
		message += string(int32(c[0]))
		c, err = wc.getMessageChar.Call(ctx)
		if err != nil {
			return message, err
		}
	}
	return message, nil
}

func (wc *Circom2WitnessCalculator2) getPrintErrorMessageHandler() func() {
	return func() {
		message, _ := wc.getMessage()
		wc.errStr.WriteString(message + "\n")
	}
}

func (wc *Circom2WitnessCalculator2) getWriteBufferMessageHandler() func() {
	return func() {
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
	}
}

func toArray32_(s *big.Int, size int) ([]uint32, error) {
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

func fromArray32_(arr []uint32) *big.Int {
	res := new(big.Int)
	radix := big.NewInt(0x100000000)
	for i := 0; i < len(arr); i++ {
		res.Mul(res, radix)
		res.Add(res, big.NewInt(int64(arr[i])))
	}
	return res
}
