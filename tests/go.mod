module github.com/iden3/go-rapidsnark/tests

go 1.18

require (
	github.com/iden3/go-rapidsnark/prover v0.0.4
	github.com/iden3/go-rapidsnark/verifier v0.0.1
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ethereum/go-ethereum v1.10.17 // indirect
	github.com/iden3/go-iden3-crypto v0.0.13 // indirect
	github.com/iden3/go-rapidsnark/types v0.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/iden3/go-rapidsnark/prover => ../prover
	github.com/iden3/go-rapidsnark/verifier => ../verifier
)
