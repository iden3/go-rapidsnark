module github.com/iden3/go-rapidsnark/witness/wazero

go 1.18

require (
	github.com/iden3/go-iden3-crypto v0.0.15
	github.com/iden3/go-rapidsnark/witness v0.0.7-0.20230523122916-060a2d3d4a85
	github.com/tetratelabs/wazero v1.1.0
)

require golang.org/x/sys v0.6.0 // indirect

replace github.com/iden3/go-rapidsnark/witness => ../
