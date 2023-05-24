package wazero

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/iden3/go-iden3-crypto/constants"
	"github.com/iden3/go-rapidsnark/witness/v2"
	wz "github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Circom2WZWitnessCalculator struct {
	runtime        wz.Runtime
	modRuntime     api.Module
	compiledModule wz.CompiledModule
}

func NewCircom2WZWitnessCalculator(
	wasmBytes []byte) (witness.CalculatorImpl, error) {

	runtime := wz.NewRuntime(context.Background())

	ctx := context.Background()
	modRuntime, err := runtime.NewHostModuleBuilder("runtime").
		NewFunctionBuilder().
		WithGoFunction(
			api.GoFunc(exceptionHandler),
			[]api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		Export("exceptionHandler").
		NewFunctionBuilder().
		WithGoModuleFunction(
			api.GoModuleFunc(printErrorMessage),
			[]api.ValueType{}, []api.ValueType{}).
		Export("printErrorMessage").
		NewFunctionBuilder().
		WithGoModuleFunction(
			api.GoModuleFunc(writeBufferMessage),
			[]api.ValueType{}, []api.ValueType{}).
		Export("writeBufferMessage").
		NewFunctionBuilder().
		WithGoModuleFunction(
			api.GoModuleFunc(showSharedRWMemory),
			[]api.ValueType{}, []api.ValueType{}).
		Export("showSharedRWMemory").
		Instantiate(ctx)
	if err != nil {
		return nil, err
	}

	compiledModule, err := runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	return &Circom2WZWitnessCalculator{
		runtime:        runtime,
		modRuntime:     modRuntime,
		compiledModule: compiledModule,
	}, nil
}

func (w *Circom2WZWitnessCalculator) Close() error {
	ctx := context.Background()

	err := w.compiledModule.Close(ctx)

	err2 := w.modRuntime.Close(ctx)
	if err == nil {
		err = err2
	}

	err2 = w.runtime.Close(ctx)
	if err == nil {
		err = err2
	}

	return err
}

func (w *Circom2WZWitnessCalculator) doCalculateWitness(ctx context.Context,
	wCtx witnessCtx, inputs map[string]any, sanityCheck bool) (err error) {

	if err = wCtx.init(ctx, sanityCheck); err != nil {
		return err
	}

	inputCntr := 0

	var arrFr []int32
	var signalSize int32
	for k := range inputs {
		hMSB, hLSB := fnvHash(k)
		signalSize, err = wCtx.getInputSignalSize(ctx, hMSB, hLSB)
		if err != nil {
			return err
		}
		fArr := make([]*big.Int, 0, signalSize)
		fArr, err = flatSlice2(fArr, inputs[k])
		if err != nil {
			return err
		}
		if len(fArr) != int(signalSize) {
			return errors.New("signal size mismatch")
		}

		for i := range fArr {
			arrFr = encodeInt(fArr[i], int(wCtx.n32))
			for j := range arrFr {
				err = wCtx.writeSharedRWMemory(ctx,
					int32(j), arrFr[int(wCtx.n32)-1-j])
				if err != nil {
					return err
				}
			}
			err = wCtx.setInputSignal(ctx, hMSB, hLSB, int32(i))
			if err != nil {
				return err
			}
			inputCntr++
		}
	}

	if wCtx.inputSize != int32(inputCntr) {
		return errors.New("input size mismatch")
	}

	return nil
}

type witnessCtx struct {
	n32                int32
	inputSize          int32
	witnessSize        int32
	init               func(ctx context.Context, sanityCheck bool) error
	getInputSignalSize func(ctx context.Context, hMSB, hLSB int32) (int32,
		error)
	writeSharedRWMemory func(ctx context.Context, sigIdx, data int32) error
	setInputSignal      func(ctx context.Context, hMSB, hLSB, z int32) error
	readSharedRWMemory  func(ctx context.Context, i int32) (int32, error)
	getWitness          func(ctx context.Context, i int32) error
	getRawPrime         func(ctx context.Context) error
}

func (wCtx *witnessCtx) prime(ctx context.Context) (*big.Int, error) {

	err := wCtx.getRawPrime(ctx)
	if err != nil {
		return nil, err
	}

	return wCtx.readInt(ctx)
}

func (wCtx *witnessCtx) readInt(ctx context.Context) (*big.Int, error) {
	arr := make([]uint32, wCtx.n32)
	for j := 0; j < int(wCtx.n32); j++ {
		val, err := wCtx.readSharedRWMemory(ctx, int32(j))
		if err != nil {
			return nil, err
		}
		arr[int(wCtx.n32)-1-j] = uint32(val)
	}

	return fromArray32(arr), nil
}

func calculateWtnsCtx(ctx context.Context,
	instance api.Module) (witnessCtx, error) {

	var wCtx witnessCtx
	var wResult []uint64
	var err error

	getFieldNumLen32 := instance.ExportedFunction("getFieldNumLen32")
	wResult, err = getFieldNumLen32.Call(ctx)
	if err != nil {
		return wCtx, err
	}
	wCtx.n32 = api.DecodeI32(wResult[0])

	wResult, err = instance.ExportedFunction("getInputSize").Call(ctx)
	if err != nil {
		return wCtx, err
	}
	wCtx.inputSize = api.DecodeI32(wResult[0])

	wResult, err = instance.ExportedFunction("getWitnessSize").Call(ctx)
	if err != nil {
		return wCtx, err
	}
	wCtx.witnessSize = api.DecodeI32(wResult[0])

	_init := instance.ExportedFunction("init")
	wCtx.init = func(ctx context.Context, sanityCheck bool) error {
		sch := int32(0)
		if sanityCheck {
			sch = 1
		}
		_, err = _init.Call(ctx, api.EncodeI32(sch))
		return err
	}

	_getInputSignalSize := instance.ExportedFunction("getInputSignalSize")
	wCtx.getInputSignalSize = func(ctx context.Context,
		hMSB, hLSB int32) (int32, error) {

		res, err2 := _getInputSignalSize.Call(ctx,
			api.EncodeI32(hMSB), api.EncodeI32(hLSB))
		if err2 != nil {
			return 0, err2
		}
		return api.DecodeI32(res[0]), nil
	}

	_writeSharedRWMemory := instance.ExportedFunction("writeSharedRWMemory")
	wCtx.writeSharedRWMemory = func(ctx context.Context,
		sigIdx, data int32) error {

		_, err2 := _writeSharedRWMemory.Call(ctx,
			api.EncodeI32(sigIdx), api.EncodeI32(data))
		return err2
	}

	_setInputSignal := instance.ExportedFunction("setInputSignal")
	wCtx.setInputSignal = func(ctx context.Context, hMSB, hLSB, i int32) error {
		_, err2 := _setInputSignal.Call(ctx,
			api.EncodeI32(hMSB), api.EncodeI32(hLSB), api.EncodeI32(i))
		return err2
	}

	_readSharedRWMemory := instance.ExportedFunction("readSharedRWMemory")
	wCtx.readSharedRWMemory = func(ctx context.Context,
		i int32) (int32, error) {

		res, err2 := _readSharedRWMemory.Call(ctx, api.EncodeI32(i))
		if err2 != nil {
			return 0, err2
		}
		return api.DecodeI32(res[0]), nil
	}

	_getWitness := instance.ExportedFunction("getWitness")
	wCtx.getWitness = func(ctx context.Context, i int32) error {
		_, err2 := _getWitness.Call(ctx, api.EncodeI32(i))
		return err2
	}

	_getRawPrime := instance.ExportedFunction("getRawPrime")
	wCtx.getRawPrime = func(ctx context.Context) error {
		_, err2 := _getRawPrime.Call(ctx)
		return err2
	}

	return wCtx, nil
}

func flatSlice2(arr []*big.Int, v any) ([]*big.Int, error) {
	switch vt := v.(type) {
	case string:
		i, ok := new(big.Int).SetString(vt, 10)
		if !ok {
			return nil, fmt.Errorf("can't parse string as int: %v", vt)
		}
		i.Rem(i, constants.Q)
		return append(arr, i), nil
	case *big.Int:
		i := new(big.Int).Set(vt)
		i.Rem(i, constants.Q)
		return append(arr, i), nil
	case []any:
		for _, e := range vt {
			var err error
			arr, err = flatSlice2(arr, e)
			if err != nil {
				return nil, err
			}
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("invalid type: %T", v)
	}
}

func encodeInt(i *big.Int, ln int) []int32 {
	i = new(big.Int).Set(i)
	arr := make([]int32, 0, ln)
	radix := big.NewInt(int64(math.MaxUint32) + 1)
	for j := 0; j < ln; j++ {
		arr = append(arr, int32(new(big.Int).Rem(i, radix).Int64()))
		i.Div(i, radix)
	}
	reverse(arr)
	return arr
}

func reverse[T any](a []T) {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
}

type ctxCloser interface {
	Close(ctx context.Context) error
}

func closeWithErrOrLog(ctx context.Context, c ctxCloser, err *error) {
	err2 := c.Close(ctx)
	if err2 != nil {
		if *err == nil {
			*err = err2
		} else {
			log.Printf("error closing instance: %v", err2)
		}
	}

}

type witnessCtxState struct {
	errStrs   []string
	msgStrs   []string
	errorCode int32
	errs      []error
}

func (s *witnessCtxState) errMessage() string {
	switch s.errorCode {
	case 0:
		return "OK"
	case 1:
		return "Signal not found."
	case 2:
		return "Too many signals set."
	case 3:
		return "Signal already set."
	case 4:
		return "Assert Failed."
	case 5:
		return "Not enough memory."
	case 6:
		return "Input signal array access exceeds the size."
	default:
		return "Unknown error."
	}
}

func (s *witnessCtxState) err() error {
	if len(s.errs) == 0 && s.errorCode == 0 {
		return nil
	}

	var errLines []string
	errLines = append(errLines,
		fmt.Sprintf("error code: %v: %v", s.errorCode, s.errMessage()))
	errLines = append(errLines, s.errStrs...)
	for i := range s.errs {
		errLines = append(errLines,
			fmt.Sprintf("Err #%v", i),
			s.errs[i].Error())
	}
	return errors.New(strings.Join(errLines, "\n"))
}

type wtnsCtxKey string

func withWtnsCtx(ctx context.Context,
	wCtxState *witnessCtxState) context.Context {

	return context.WithValue(ctx, wtnsCtxKey("wtnsCtx"), wCtxState)
}

func fromWtnsCtx(ctx context.Context) *witnessCtxState {
	v := ctx.Value(wtnsCtxKey("wtnsCtx"))
	if v == nil {
		return nil
	}
	return v.(*witnessCtxState)
}

func getMessage(ctx context.Context, m api.Module) (string, error) {
	var buf bytes.Buffer
	max := 4048
	for {
		data, err := m.ExportedFunction("getMessageChar").Call(ctx)
		if err != nil {
			return "", err
		}
		b := byte(api.DecodeI32(data[0]))
		if b == 0 {
			return buf.String(), nil
		}
		buf.WriteByte(b)

		max--
		if max == 0 {
			return "", errors.New("max iterations reached")
		}
	}
}

func exceptionHandler(ctx context.Context, params []uint64) {
	wtnsCtx := fromWtnsCtx(ctx)
	wtnsCtx.errorCode = api.DecodeI32(params[0])
}

func printErrorMessage(ctx context.Context, m api.Module, _ []uint64) {
	wtnsCtx := fromWtnsCtx(ctx)

	msg, err := getMessage(ctx, m)
	if err != nil {
		wtnsCtx.errs = append(wtnsCtx.errs, err)
		return
	}

	wtnsCtx.errStrs = append(wtnsCtx.errStrs, msg)
}

func writeBufferMessage(ctx context.Context, m api.Module, _ []uint64) {
	wtnsCtx := fromWtnsCtx(ctx)

	msg, err := getMessage(ctx, m)
	if err != nil {
		wtnsCtx.errs = append(wtnsCtx.errs, err)
		return
	}

	if msg == "\n" {
		log.Print(strings.Join(wtnsCtx.msgStrs, " "))
		wtnsCtx.msgStrs = wtnsCtx.msgStrs[:0]
	} else {
		wtnsCtx.msgStrs = append(wtnsCtx.msgStrs, msg)
	}
}

func showSharedRWMemory(ctx context.Context, m api.Module, _ []uint64) {
	printSharedRWMemory(ctx, m)
}

func printSharedRWMemory(ctx context.Context, m api.Module) {
	wtnsCtx := fromWtnsCtx(ctx)

	data, err := m.ExportedFunction("getFieldNumLen32").Call(ctx)
	if err != nil {
		wtnsCtx.errs = append(wtnsCtx.errs, err)
		return
	}

	sharedRwMemorySize := int(api.DecodeI32(data[0]))
	var arr = make([]uint32, sharedRwMemorySize)

	for j := 0; j < sharedRwMemorySize; j++ {
		data, err = m.ExportedFunction("readSharedRWMemory").
			Call(ctx, uint64(j))
		if err != nil {
			wtnsCtx.errs = append(wtnsCtx.errs, err)
			return
		}
		arr[j] = uint32(api.DecodeI32(data[0]))
	}

	wtnsCtx.msgStrs = append(wtnsCtx.msgStrs, fromArray32(arr).Text(10))
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

// fnvHash returns the 64 bit FNV-1a hash split into two 32 bit values: (MSB, LSB)
func fnvHash(s string) (int32, int32) {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	h := hash.Sum64()
	return int32(h >> 32), int32(h & 0xffffffff)
}

// Calculate calculates the witness given the inputs.
func (wc *Circom2WZWitnessCalculator) Calculate(inputs map[string]interface{},
	sanityCheck bool) (wtns witness.Witness, err error) {

	wCtxState := &witnessCtxState{}
	ctx := withWtnsCtx(context.Background(), wCtxState)

	cfg := wz.NewModuleConfig()
	var instance api.Module
	instance, err = wc.runtime.InstantiateModule(ctx, wc.compiledModule, cfg)
	if err != nil {
		return wtns, err
	}
	defer closeWithErrOrLog(ctx, instance, &err)

	var wCtx witnessCtx
	wCtx, err = calculateWtnsCtx(ctx, instance)
	if err != nil {
		return wtns, err
	}

	wtns.N32 = int(wCtx.n32)

	err = wc.doCalculateWitness(ctx, wCtx, inputs, sanityCheck)
	if err != nil {
		return wtns, err
	}

	wtns.Witness = make([]*big.Int, wCtx.witnessSize)

	for i := 0; i < int(wCtx.witnessSize); i++ {
		err = wCtx.getWitness(ctx, int32(i))
		if err != nil {
			return wtns, err
		}
		wtns.Witness[i], err = wCtx.readInt(ctx)
	}

	wtns.Prime, err = wCtx.prime(ctx)
	if err != nil {
		return wtns, err
	}

	return wtns, wCtxState.err()
}
