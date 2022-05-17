//go:build !dynamic

package prover

// #cgo CFLAGS: -DUSE_VENDORED_RAPIDSNARK
// #cgo LDFLAGS: ${SRCDIR}/rapidsnark_vendor/librapidsnark-darwin-arm64.a ${SRCDIR}/rapidsnark_vendor/libgmp-darwin-arm64.a -lstdc++
import "C"
