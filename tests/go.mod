module github.com/iden3/go-rapidsnark/tests

go 1.18

require (
	github.com/iden3/go-rapidsnark/prover v0.0.10
	github.com/iden3/go-rapidsnark/verifier v0.0.5
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/iden3/go-iden3-crypto v0.0.15 // indirect
	github.com/iden3/go-rapidsnark/types v0.0.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/iden3/go-rapidsnark/prover => ../prover
	github.com/iden3/go-rapidsnark/types => ../types
	github.com/iden3/go-rapidsnark/verifier => ../verifier
)
