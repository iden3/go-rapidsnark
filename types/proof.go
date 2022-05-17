package types

// ProofData is structure that represents SnarkJS library result of proof generation
type ProofData struct {
	A        []string   `json:"pi_a"`
	B        [][]string `json:"pi_b"`
	C        []string   `json:"pi_c"`
	Protocol string     `json:"protocol"`
}

// ZKProof is proof data with public signals
type ZKProof struct {
	Proof      *ProofData `json:"proof"`
	PubSignals []string   `json:"pub_signals"`
}
