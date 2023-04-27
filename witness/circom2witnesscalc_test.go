package witness

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCircom2CalculateWitness(t *testing.T) {
	wasmBytes, err := os.ReadFile("testdata/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("testdata/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
}

func TestCircom2CalculateBinWitness(t *testing.T) {
	wasmBytes, err := os.ReadFile("testdata/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("testdata/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
}

func TestCircom2CalculateWTNSBin(t *testing.T) {
	wasmBytes, err := os.ReadFile("testdata/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("testdata/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)

	//_ = ioutil.WriteFile("testdata/circom2/witness.wtns", wtnsBytes, fs.FileMode(defaultFileMode))
}

// TestCircom2CalculateWitness210 tests the calculation of the witness for the circom 2.1.0
func TestCircom2CalculateWitness210(t *testing.T) {
	wasmBytes, err := os.ReadFile("testdata/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("testdata/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
}

// TestCircom2CalculateBinWitness210 tests the calculation of the witness for the circom 2.1.0
func TestCircom2CalculateBinWitness210(t *testing.T) {
	wasmBytes, err := os.ReadFile("testdata/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("testdata/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
}

// TestCircom2CalculateWTNSBin210 tests the calculation of the witness for the circom 2.1.0
func TestCircom2CalculateWTNSBin210(t *testing.T) {
	wasmBytes, err := os.ReadFile("testdata/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("testdata/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)

	//_ = ioutil.WriteFile("testdata/circom2_1_0/witness.wtns", wtnsBytes, fs.FileMode(defaultFileMode))
}
