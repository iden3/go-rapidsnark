# go-rapidsnark

A go wrapper for the RapidSNARK C++ library.

# Get module

## Build using pre-built vendored libraries.

For few architectures, you can use pre-built vendored libraries. And just use
this module without any externally built dependencies.

Supported architectures are:
* MacOS x86_64
* MacOS Apple Silicon M1

## Build using custom libraries.

You need `gmp` and `rapidsnark` libraries available on build host.

Supposed all needed files are in following directories:
* `prover.h` is located in `${HOME}/src/rapidsnark/src`
* `librapidsnark.a` is located in `${HOME}/src/rapidsnark/build_prover/src`
* `libgmp.a` is located in `${HOME}/src/rapidsnark/depends/gmp/package/lib`

```shell
export CGO_CFLAGS="-I${HOME}/src/rapidsnark/src" 
export CGO_LDFLAGS="-L${HOME}/src/rapidsnark/build_prover/src -L${HOME}/src/rapidsnark/depends/gmp/package/lib"
go build -tags dynamic
```

Tag `dynamic` is required to exclude usage of vendored libraries.
