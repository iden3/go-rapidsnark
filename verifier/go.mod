module github.com/iden3/go-rapidsnark/verifier

go 1.18

replace github.com/iden3/go-rapidsnark/types => ../types

require (
	github.com/iden3/go-iden3-crypto v0.0.13
	github.com/iden3/go-rapidsnark/types v0.0.0-00010101000000-000000000000
)

require (
	github.com/ethereum/go-ethereum v1.10.17 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)
