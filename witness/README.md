# go-rapidsnark witness calculator

Calculates witness, that can be passed to a prover ([snarkjs](https://github.com/iden3/snarkjs), [go-rapidsnark/prover](/prover) or [rapidsnark](https://github.com/iden3/rapidsnark)) to generate a zero-knowledge proof.

## Installation

```
go get github.com/iden3/go-rapidsnark/witness/v2
```

## Dependencies

This package depends on wasmer shared library, which needs to be copied from [wasmer-go](https://github.com/wasmerio/wasmer-go/tree/master/wasmer/packaged/lib) module source code.
E.g. to run compiled project on Alpine linux you would need to copy `/go/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-amd64/libwasmer.so` from the build host/container.
