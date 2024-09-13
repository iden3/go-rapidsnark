module github.com/iden3/go-rapidsnark/witness/test-wasm-impls

go 1.21

toolchain go1.23.1

require (
	github.com/iden3/go-rapidsnark/witness/v2 v2.0.0
	github.com/iden3/go-rapidsnark/witness/wasmer v0.0.0
	github.com/iden3/go-rapidsnark/witness/wazero v0.0.0
	github.com/stretchr/testify v1.8.2

)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/iden3/go-iden3-crypto v0.0.15 // indirect
	github.com/iden3/wasmer-go v0.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/iden3/go-rapidsnark/witness/v2 => ../
	github.com/iden3/go-rapidsnark/witness/wasmer => ../wasmer
	github.com/iden3/go-rapidsnark/witness/wazero => ../wazero
)
