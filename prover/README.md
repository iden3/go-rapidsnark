# go-rapidsnark prover

A go wrapper for the RapidSNARK C++ library.

To use this module you either need to have RapidSNARK library available on
your build host or have one of the supported architectures to use vendored
pre-built libraries for convenience.

## Dependencies
* C/C++ compiler and standard libraries available
* glibc version >= 2.31
* OpenMP

On Ubuntu it would be enough to install `build-essential` and `libomp-dev` packages.

If you build your project using `golang` Docker container, all tools are already installed.

To run compiled project on Alpine linux you would need to install there `libstdc++`, `gcompat` and `libgomp` packages.

## Build using pre-built vendored RapidSNARK library.

For the following architectures, pre-built vendored libraries are included:
* MacOS x86_64
* MacOS ARM64 Apple Silicon
* Linux x86_64
* Linux ARM64 v8

## Performance optimization on x86_64 hardware

Rapidsnark has optimization for recent x86_64 processors that gives ~2x speed boost, but older hardware may lack support for ADX and BMI2 instruction sets used.

MacOS build has it enabled. But for linux we disabled the optimization (at least for now, because GitHub Actions may use old hardware).

To **enable** optimization on linux use `rapidsnark_asm` build tag:

```shell
go build -tags rapidsnark_asm
go test -tags rapidsnark_asm
```

In the future we may change default behaviour, so to force disable optimizations use `rapidsnark_noasm` build tag.

## Performance optimization on arm64 hardware
We used NEON instruction set, and it is always enabled, so no build tags are needed.

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
