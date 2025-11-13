//go:build dynamic && !prover_disabled

package prover

// #cgo LDFLAGS: -lrapidsnark -lgmp -lstdc++
import "C"
