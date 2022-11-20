# go-rapidsnark witness calculator

Calculates witness, that can be passed to a prover ([snarkjs](https://github.com/iden3/snarkjs), [go-rapidsnark/prover](/prover) or [rapidsnark](https://github.com/iden3/rapidsnark)) to generate a zero-knowledge proof.

## Installation

```
go get https://github.com/iden3/go-rapidsnark/witness
```

## Dependencies

This package depends on wasmer shared library, which needs to be copied from [wasmer-go](https://github.com/wasmerio/wasmer-go/tree/master/wasmer/packaged/lib) module source code.
E.g. 

