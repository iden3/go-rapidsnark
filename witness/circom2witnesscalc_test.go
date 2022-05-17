package witness

import (
	"io/fs"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultFileMode = 0644

func TestCircom2CalculateWitness(t *testing.T) {
	wasmBytes, err := ioutil.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := ioutil.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witness, err := calc.CalculateWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witness)
}

func TestCircom2CalculateBinWitness(t *testing.T) {
	wasmBytes, err := ioutil.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := ioutil.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	witnessBytes, err := calc.CalculateBinWitness(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, witnessBytes)
}

func TestCircom2CalculateWTNSBin(t *testing.T) {
	wasmBytes, err := ioutil.ReadFile("test_files/circom2/circuit.wasm")
	require.NoError(t, err)

	inputBytes, err := ioutil.ReadFile("test_files/circom2/input.json")
	require.NoError(t, err)

	calc, err := NewCircom2WitnessCalculator(wasmBytes, true)
	require.NoError(t, err)
	require.NotEmpty(t, calc)

	inputs, err := ParseInputs(inputBytes)
	require.NoError(t, err)

	wtnsBytes, err := calc.CalculateWTNSBin(inputs, true)
	require.NoError(t, err)
	require.NotEmpty(t, wtnsBytes)

	_ = ioutil.WriteFile("test_files/circom2/witness.wtns", wtnsBytes, fs.FileMode(defaultFileMode))
}
