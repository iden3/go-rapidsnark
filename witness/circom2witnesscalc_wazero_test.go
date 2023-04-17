package witness

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWZCircom2CalculateWitness(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
	require.Equal(t, "c1780821352c069392e9d0fab4330531", hashInts(witness))
}

func TestWZCircom2CalculateBinWitness(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
	require.Equal(t, "d2c0486d7fd6f0715d04d535765f028b",
		hashBytes(witnessBytes))
}

func TestWZCircom2CalculateWTNSBin(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)
	require.Equal(t, "1709fbda942dabed641044f39b466e94",
		hashBytes(wtnsBytes))

}

// TestWZCircom2CalculateWitness210 tests the calculation of the witness for the circom 2.1.0
func TestWZCircom2CalculateWitness210(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
	require.Equal(t, "c0a2b43f5a333310c2bb8d357db46d3b", hashInts(witness))
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

// TestWZCircom2CalculateBinWitness210 tests the calculation of the witness
// for the circom 2.1.0
func TestWZCircom2CalculateBinWitness210(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
	require.Equal(t, "2b38b66035d8e923eacc028ea0f1dad2",
		hashBytes(witnessBytes))
}

// TestWZCircom2CalculateWTNSBin210 tests the calculation of the witness
// for the circom 2.1.0
func TestWZCircom2CalculateWTNSBin210(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)
	require.Equal(t, "75c5682a7195c20868b59d6580852fce",
		hashBytes(wtnsBytes))
}

// TestWZCircom2CalculateWTNSBin210 tests the calculation of the witness
// for the circom 2.1.0
func TestWZCircom2CalculateWTNSBin210_Error(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WZWitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)
	wrongSmtRoot, ok := big.NewInt(0).SetString(
		"23891407091237035626910338386637210028103224489833886255774452947213913989795",
		10)
	require.True(t, ok)
	inputs["globalSmtRoot"] = wrongSmtRoot

	_, err = calc.CalculateWTNSBin(inputs, true)
	require.EqualError(t, err, `error code: 4: Assert Failed.
Error in template ForceEqualIfEnabled_234 line: 56
Error in template SMTVerifier_235 line: 134
Error in template AuthV2_347 line: 93`)
}
