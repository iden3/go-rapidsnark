package tests

import (
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/stretchr/testify/assert"
)

func Test_Groth16Prover(t *testing.T) {
	var provingKey, verificationKey, witness []byte
	var err error

	provingKey, err = os.ReadFile("./testdata/circuit_final.zkey")
	assert.Nil(t, err)

	witness, err = os.ReadFile("./testdata/witness.wtns")
	assert.Nil(t, err)

	verificationKey, err = os.ReadFile("./testdata/verification_key.json")
	assert.Nil(t, err)

	assert.NoError(t, err)

	proof, err := prover.Groth16Prover(provingKey, witness)
	assert.NoError(t, err)

	err = verifier.VerifyGroth16(*proof, verificationKey)
	assert.NoError(t, err)
}

func Benchmark(b *testing.B) {
	var provingKey, verificationKey, witness []byte
	var err error

	provingKey, err = os.ReadFile("./testdata/circuit_final.zkey")
	assert.Nil(b, err)

	witness, err = os.ReadFile("./testdata/witness.wtns")
	assert.Nil(b, err)

	verificationKey, err = os.ReadFile("./testdata/verification_key.json")
	assert.Nil(b, err)

	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		proof, err := prover.Groth16Prover(provingKey, witness)
		assert.NoError(b, err)

		err = verifier.VerifyGroth16(*proof, verificationKey)
		assert.NoError(b, err)
	}
}
