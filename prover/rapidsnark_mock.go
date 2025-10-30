//go:build prover_disabled

package prover

import (
	"errors"

	"github.com/iden3/go-rapidsnark/types"
)

// Groth16Prover generates proof and returns proof and pubsignals as types.ZKProof
func Groth16Prover(zkey []byte, witness []byte) (proof *types.ZKProof, err error) {
	return nil, errors.New("prover disabled: 'prover_disabled' build flag was passed")
}

// Groth16ProverRaw generates proof and returns proof and pubsignals as json string
func Groth16ProverRaw(zkey []byte, witness []byte) (proof string, publicInputs string, err error) {
	return "", "", errors.New("prover disabled: 'prover_disabled' build flag was passed")
}
