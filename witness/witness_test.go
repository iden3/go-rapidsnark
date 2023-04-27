package witness

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEngines(t *testing.T) {
	engineTestCases := []struct {
		title  string
		engine func([]byte) (WitnessCalculator, error)
	}{
		{
			title:  "Wazero",
			engine: NewCircom2WZWitnessCalculator,
		},
		{
			title:  "Wasmer",
			engine: NewCircom2WitnessCalculator,
		},
		{
			title:  "empty",
			engine: nil,
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

					var ops []Option
					if engTC.engine != nil {
						ops = append(ops, WithWasmEngine(engTC.engine))
					}
					calc, err := NewCalc(wasmBytes, ops...)
					require.NoError(t, err)

					inputs, err := ParseInputs(inputBytes)
					require.NoError(t, err)

					t.Run("CalculateWitness", func(t *testing.T) {
						witness, err := calc.CalculateWitness(inputs, true)
						require.NoError(t, err)
						require.NotEmpty(t, witness)
						require.Equal(t, circomTC.wantWtnsHex,
							hashInts(witness))
					})

					t.Run("CalculateBinWitness", func(t *testing.T) {
						witness, err := calc.CalculateBinWitness(inputs, true)
						require.NoError(t, err)
						require.NotEmpty(t, witness)
						require.Equal(t, circomTC.wantBinWtnsHex,
							hashBytes(witness))
					})

					t.Run("CalculateWTNSBin", func(t *testing.T) {
						witness, err := calc.CalculateWTNSBin(inputs, true)
						require.NoError(t, err)
						require.NotEmpty(t, witness)
						require.Equal(t, circomTC.wantWTNSBinHex,
							hashBytes(witness))
					})
				})
			}
		})
	}
}
