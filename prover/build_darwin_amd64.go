//go:build !dynamic

package prover

// #cgo CFLAGS: -DUSE_VENDORED_RAPIDSNARK
// #cgo LDFLAGS: ${SRCDIR}/rapidsnark_vendor/librapidsnark-darwin-amd64.a ${SRCDIR}/rapidsnark_vendor/libgmp-darwin-amd64.a -lstdc++
import "C"
