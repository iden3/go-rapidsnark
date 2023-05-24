package witness

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/witness/v2"
	"github.com/iden3/go-rapidsnark/witness/wasmer"
	"github.com/iden3/go-rapidsnark/witness/wazero"
	"github.com/stretchr/testify/require"
)

func TestEngines(t *testing.T) {
	engineTestCases := []struct {
		title   string
		engine  func(code []byte) (witness.CalculatorImpl, error)
		wantErr string
	}{
		{
			title:  "Wazero",
			engine: wazero.NewCircom2WZWitnessCalculator,
		},
		{
			title:  "Wasmer",
			engine: wasmer.NewCircom2WitnessCalculator,
		},
		{
			title:   "empty",
			wantErr: "witness calculator wasm engine not set",
		},
	}

	circomTestCases := []struct {
		wasmFile       string
		inputs         string
		wantWtnsHex    string
		wantBinWtnsHex string
		wantWTNSBinHex string
	}{
		{
			wasmFile:       "testdata/circom2/circuit.wasm",
			inputs:         "testdata/circom2/input.json",
			wantWtnsHex:    "c1780821352c069392e9d0fab4330531",
			wantBinWtnsHex: "d2c0486d7fd6f0715d04d535765f028b",
			wantWTNSBinHex: "1709fbda942dabed641044f39b466e94",
		},
		{
			wasmFile:       "testdata/circom2_1_0/circuit.wasm",
			inputs:         "testdata/circom2_1_0/input.json",
			wantWtnsHex:    "c0a2b43f5a333310c2bb8d357db46d3b",
			wantBinWtnsHex: "2b38b66035d8e923eacc028ea0f1dad2",
			wantWTNSBinHex: "75c5682a7195c20868b59d6580852fce",
		},
	}

	for i := range engineTestCases {
		engTC := engineTestCases[i]
		t.Run(engTC.title, func(t *testing.T) {
			for _, circomTC := range circomTestCases {
				t.Run(circomTC.wasmFile, func(t *testing.T) {
					wasmBytes, err := os.ReadFile(circomTC.wasmFile)
					require.NoError(t, err)
					inputBytes, err := os.ReadFile(circomTC.inputs)
					require.NoError(t, err)

					var ops []witness.Option
					if engTC.engine != nil {
						ops = append(ops, witness.WithWasmEngine(engTC.engine))
					}
					calc, err := witness.NewCalculator(wasmBytes, ops...)
					if engTC.wantErr != "" {
						require.EqualError(t, err, engTC.wantErr)
						return
					}

					require.NoError(t, err)

					inputs, err := witness.ParseInputs(inputBytes)
					require.NoError(t, err)

					t.Run("CalculateWitness", func(t *testing.T) {
						wtns, err2 := calc.CalculateWitness(inputs, true)
						require.NoError(t, err2)
						require.NotEmpty(t, wtns)
						require.Equal(t, circomTC.wantWtnsHex, hashInts(wtns))
					})

					t.Run("CalculateBinWitness", func(t *testing.T) {
						wtns, err2 := calc.CalculateBinWitness(inputs, true)
						require.NoError(t, err2)
						require.NotEmpty(t, wtns)
						require.Equal(t, circomTC.wantBinWtnsHex,
							hashBytes(wtns))
					})

					t.Run("CalculateWTNSBin", func(t *testing.T) {
						wtns, err2 := calc.CalculateWTNSBin(inputs, true)
						require.NoError(t, err2)
						require.NotEmpty(t, wtns)
						require.Equal(t, circomTC.wantWTNSBinHex,
							hashBytes(wtns))
					})
				})
			}
		})
	}
}

func hashInts(in []*big.Int) string {
	h := md5.New()
	for _, i := range in {
		h.Write(i.Bytes())
	}
	return hex.EncodeToString(h.Sum(nil))
}

func hashBytes(in []byte) string {
	h := md5.New()
	n, err := h.Write(in)
	if err != nil {
		panic(err)
	}
	if n != len(in) {
		panic("incorrect size")
	}
	return hex.EncodeToString(h.Sum(nil))
}
