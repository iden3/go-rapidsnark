package verifier

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/constants"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier/bn256"
)

// VerifyGroth16 performs a verification of zkp  based on verification key and public inputs
func VerifyGroth16(zkProof types.ZKProof, verificationKey []byte) error {

	// 1. cast external proof data to internal model.
	p, err := parseProofData(*zkProof.Proof)
	if err != nil {
		return err
	}

	// 2. cast external verification key data to internal model.
	var vkStr vkJSON
	err = json.Unmarshal(verificationKey, &vkStr)
	if err != nil {
		return err
	}
	vkKey, err := parseVK(vkStr)
	if err != nil {
		return err
	}

	// 2. cast external public inputs data to internal model.
	pubSignals, err := stringsToArrayBigInt(zkProof.PubSignals)
	if err != nil {
		return err
	}

	return verify(vkKey, p, pubSignals)
}

// verify performs the verification the Groth16 zkSNARK proofs
func verify(vk *vk, proof proofPairingData, inputs []*big.Int) error {
	if len(inputs)+1 != len(vk.IC) {
		return fmt.Errorf("len(inputs)+1 != len(vk.IC)")
	}
	vkX := new(bn256.G1).ScalarBaseMult(big.NewInt(0))
	for i := 0; i < len(inputs); i++ {
		// check input inside field
		if inputs[i].Cmp(constants.Q) != -1 {
			return fmt.Errorf("input value is not in the fields")
		}
		vkX = new(bn256.G1).Add(vkX, new(bn256.G1).ScalarMult(vk.IC[i+1], inputs[i]))
	}
	vkX = new(bn256.G1).Add(vkX, vk.IC[0])

	g1 := []*bn256.G1{proof.A, new(bn256.G1).Neg(vk.Alpha), vkX.Neg(vkX), new(bn256.G1).Neg(proof.C)}
	g2 := []*bn256.G2{proof.B, vk.Beta, vk.Gamma, vk.Delta}

	res := bn256.PairingCheck(g1, g2)
	if !res {
		return fmt.Errorf("invalid proofs")
	}
	return nil
}
