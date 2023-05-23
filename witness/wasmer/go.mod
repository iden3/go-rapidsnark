module github.com/iden3/go-rapidsnark/witness/v2/wasmer

go 1.18

require (
	github.com/iden3/go-iden3-crypto v0.0.15
	github.com/iden3/go-rapidsnark/witness/v2 v2.0.0-20230523125954-fcfab2575c4d
	github.com/iden3/wasmer-go v0.0.1
)

require golang.org/x/sys v0.6.0 // indirect

replace github.com/iden3/go-rapidsnark/witness/v2 => ../
