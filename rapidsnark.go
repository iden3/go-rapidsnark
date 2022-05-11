package rapidsnark

/*
#include <stdlib.h>
#include "select_rapidsnark.h"
*/
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

const bufferSize = 16384
const MaxBufferSize = 10485760

func Groth16Prover(zkey []byte,
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
			(*C.char)(unsafe.Pointer(&proofBuffer[0])), C.ulong(proofBufSize),
			(*C.char)(unsafe.Pointer(&publicBuffer[0])), C.ulong(publicBufSize),
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
