package prover

/*
#include <stdlib.h>
#include "select_rapidsnark.h"
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/iden3/go-rapidsnark/types"
)

const bufferSize = 16384
const MaxBufferSize = 10485760

// Groth16Prover generates proof and returns proof and pubsignals as types.ZKProof
func Groth16Prover(zkey []byte,
	witness []byte) (proof *types.ZKProof, err error) {
	proofStr, pubSignalsStr, err := Groth16ProverRaw(zkey, witness)
	if err != nil {
		return nil, err
	}
	var proofData types.ProofData
	var pubSignals []string

	err = json.Unmarshal([]byte(proofStr), &proofData)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(pubSignalsStr), &pubSignals)
	if err != nil {
		return nil, err
	}

	return &types.ZKProof{Proof: &proofData, PubSignals: pubSignals}, nil
}

// Groth16ProverRaw generates proof and returns proof and pubsignals as json string
func Groth16ProverRaw(zkey []byte,
	witness []byte) (proof string, publicInputs string, err error) {
	if len(zkey) == 0 {
		return "", "", errors.New("zkey is empty")
	}
	if len(witness) == 0 {
		return "", "", errors.New("witness is empty")
	}

	proofBufSize := bufferSize
	publicBufSize := bufferSize
	var proofBuffer = make([]byte, proofBufSize)
	var publicBuffer = make([]byte, publicBufSize)

	const errorBufSize = 4096
	var errorMessage [errorBufSize]byte

	for {
		r := C.groth16_prover(
			unsafe.Pointer(&zkey[0]), C.ulong(len(zkey)),
			unsafe.Pointer(&witness[0]), C.ulong(len(witness)),
			(*C.char)(unsafe.Pointer(&proofBuffer[0])), (*C.ulong)(unsafe.Pointer(&proofBufSize)),
			(*C.char)(unsafe.Pointer(&publicBuffer[0])), (*C.ulong)(unsafe.Pointer(&publicBufSize)),
			(*C.char)(unsafe.Pointer(&errorMessage[0])), errorBufSize)

		if r != 0 {
			idx := bytes.IndexByte(errorMessage[:], 0)
			if idx == -1 {
				idx = len(errorMessage)
			}
			return "", "", fmt.Errorf(
				"error generating proof. Code: %v. Message: %v",
				r, string(errorMessage[:idx]))
		}

		// if true, enlarge buffers and repeat.
		repeat := false

		idx := bytes.IndexByte(proofBuffer[:], 0)
		if idx == -1 {
			if proofBufSize >= MaxBufferSize {
				return "", "", errors.New("proof is too large")
			}
			proofBufSize *= 2
			if proofBufSize >= MaxBufferSize {
				proofBufSize = MaxBufferSize
			}
			proofBuffer = make([]byte, proofBufSize)
			repeat = true
		} else {
			proof = string(proofBuffer[:idx])
		}

		idx = bytes.IndexByte(publicBuffer[:], 0)
		if idx == -1 {
			if publicBufSize >= MaxBufferSize {
				return "", "", errors.New("public inputs is too large")
			}
			publicBufSize *= 2
			if publicBufSize >= MaxBufferSize {
				publicBufSize = MaxBufferSize
			}
			publicBuffer = make([]byte, publicBufSize)
			repeat = true
		} else {
			publicInputs = string(publicBuffer[:idx])
		}

		if repeat {
			continue
		}

		return
	}
}
