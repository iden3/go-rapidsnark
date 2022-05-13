# go-rapidsnark

A go wrapper for the RapidSNARK C++ library.

To use this module you either need to have RapidSNARK library available on
your build host or have one of the supported architectures to use vendored
pre-built library for convenience use.

## Build using pre-built vendored RapidSNARK library.

For few architectures, pre-built vendored libraries are included. And just use
this module without any externally built dependencies except standard C/C++
libraries.

Supported architectures are:
* MacOS x86_64
* MacOS Apple Silicon M1
* Linux x86_64
* Linux ARM64 v8

Minimum glibc version that should be available on build host is 2.31.

Also, you need a C/C++ compiler and standard libraries available. On Ubuntu it
would be enough to install `build-essential` package. If you build your project
using `golang` Docker container, all tools are already installed.

## Build using custom RapidSNARK library.

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

## Examples

Library usage example is available in [`/cmd/proof/`](cmd/proof) directory.
