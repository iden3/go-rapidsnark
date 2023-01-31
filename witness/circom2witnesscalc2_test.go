package witness_test

import (
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/witness"
	"github.com/stretchr/testify/require"
)

func TestCircom2CalculateWitnessWazero(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := witness.NewCircom2WitnessCalculator2(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := witness.ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
}

func TestCircom2CalculateBinWitnessWazero(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := witness.NewCircom2WitnessCalculator2(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := witness.ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
}

func TestCircom2CalculateWTNSBinWazero(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := witness.NewCircom2WitnessCalculator2(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := witness.ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)

	//_ = ioutil.WriteFile("test_files/circom2/witness.wtns", wtnsBytes, fs.FileMode(defaultFileMode))
}

// TestCircom2CalculateWitness210 tests the calculation of the witness for the circom 2.1.0
func TestCircom2CalculateWitness210Wazero(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := witness.NewCircom2WitnessCalculator2(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := witness.ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
}

// TestCircom2CalculateBinWitness210 tests the calculation of the witness for the circom 2.1.0
func TestCircom2CalculateBinWitness210Wazero(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := witness.NewCircom2WitnessCalculator2(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := witness.ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
}

// TestCircom2CalculateWTNSBin210 tests the calculation of the witness for the circom 2.1.0
func TestCircom2CalculateWTNSBin210Wazero(t *testing.T) {
	wasmBytes, err := os.ReadFile("test_files/circom2_1_0/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := os.ReadFile("test_files/circom2_1_0/input.json")
	require.NoError(t, err)

	calc, err := witness.NewCircom2WitnessCalculator2(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := witness.ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)

	//_ = ioutil.WriteFile("test_files/circom2_1_0/witness.wtns", wtnsBytes, fs.FileMode(defaultFileMode))
}
