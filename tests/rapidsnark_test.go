package tests

import (
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/stretchr/testify/require"
)

func Test_Groth16Prover(t *testing.T) {
	var provingKey, verificationKey, witness []byte
	var err error

	provingKey, err = os.ReadFile("./testdata/circuit_final.zkey")
	require.NoError(t, err)

	witness, err = os.ReadFile("./testdata/witness.wtns")
	require.NoError(t, err)

	verificationKey, err = os.ReadFile("./testdata/verification_key.json")
	require.NoError(t, err)

	proof, err := prover.Groth16Prover(provingKey, witness)
	require.NoError(t, err)

	err = verifier.VerifyGroth16(*proof, verificationKey)
	require.NoError(t, err)
}

func Benchmark(b *testing.B) {
	var provingKey, verificationKey, witness []byte
	var err error

	provingKey, err = os.ReadFile("./testdata/circuit_final.zkey")
	require.NoError(b, err)

	witness, err = os.ReadFile("./testdata/witness.wtns")
	require.NoError(b, err)

	verificationKey, err = os.ReadFile("./testdata/verification_key.json")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		proof, err := prover.Groth16Prover(provingKey, witness)
		require.NoError(b, err)
		require.NotEmpty(b, proof)
		require.NotEmpty(b, verificationKey)

		err = verifier.VerifyGroth16(*proof, verificationKey)
		require.NoError(b, err)
	}
}
