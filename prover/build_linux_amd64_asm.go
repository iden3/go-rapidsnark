//go:build linux && amd64 && !dynamic && rapidsnark_asm && !rapidsnark_noasm

package prover

// #cgo CFLAGS: -DUSE_VENDORED_RAPIDSNARK
// #cgo LDFLAGS: ${SRCDIR}/rapidsnark_vendor/librapidsnark-linux-amd64.a ${SRCDIR}/rapidsnark_vendor/libgmp-linux-amd64.a -lstdc++ -fopenmp
import "C"
