//go:build !prover_disabled

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

	const errorBufSize = 4096
	errorMessage := make([]byte, errorBufSize)

	zkeyPointer := C.CBytes(zkey)
	wtnsPointer := C.CBytes(witness)

	errorMessagePointer := C.CString(string(errorMessage))

	defer func() {
		C.free(zkeyPointer)
		C.free(wtnsPointer)
		C.free(unsafe.Pointer(errorMessagePointer))
	}()

	for {
		var proofBuffer = make([]byte, proofBufSize)
		var publicBuffer = make([]byte, publicBufSize)

		//proofBufferPointer := unsafe.Pointer(&proofBuffer[0])
		proofBufferPointer := C.CString(string(proofBuffer))
		proofBufSizePointer := unsafe.Pointer(&proofBufSize)
		publicBufferPointer := C.CString(string(publicBuffer))
		publicBufSizePointer := unsafe.Pointer(&publicBufSize)

		defer func() {
			C.free(unsafe.Pointer(proofBufferPointer))
			C.free(unsafe.Pointer(publicBufferPointer))
		}()

		r := C.groth16_prover(
			zkeyPointer, C.ulong(len(zkey)),
			wtnsPointer, C.ulong(len(witness)),
			proofBufferPointer, (*C.ulong)(proofBufSizePointer),
			publicBufferPointer, (*C.ulong)(publicBufSizePointer),
			errorMessagePointer, errorBufSize)

		proofBuffer = []byte(C.GoString(proofBufferPointer))
		publicBuffer = []byte(C.GoString(publicBufferPointer))
		errorMessage = []byte(C.GoString(errorMessagePointer))

		if r != 0 && r != 2 {
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

		if len(proofBuffer) == 0 {
			if proofBufSize >= MaxBufferSize {
				return "", "", errors.New("proof is too large")
			}
			proofBufSize *= 2
			if proofBufSize >= MaxBufferSize {
				proofBufSize = MaxBufferSize
			}
			repeat = true
		} else {
			proof = string(proofBuffer)
		}

		if len(publicBuffer) == 0 {
			if publicBufSize >= MaxBufferSize {
				return "", "", errors.New("public inputs is too large")
			}
			publicBufSize *= 2
			if publicBufSize >= MaxBufferSize {
				publicBufSize = MaxBufferSize
			}
			repeat = true
		} else {
			publicInputs = string(publicBuffer)
		}

		if repeat {
			continue
		}

		return
	}
}
