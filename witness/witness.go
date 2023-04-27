package witness

import "math/big"

type Option func(cfg *calcConfig)

func WithWasmEngine(calculator func([]byte) (WitnessCalculator, error)) Option {
	return func(cfg *calcConfig) {
		cfg.wasmEngine = calculator
	}
}

type WitnessCalculator interface {
	CalculateWitness(inputs map[string]interface{},
		sanityCheck bool) ([]*big.Int, error)
	CalculateBinWitness(inputs map[string]interface{},
		sanityCheck bool) ([]byte, error)
	CalculateWTNSBin(inputs map[string]interface{},
		sanityCheck bool) ([]byte, error)
}

type calcConfig struct {
	wasmEngine func([]byte) (WitnessCalculator, error)
}

func NewCalc(wasm []byte, ops ...Option) (WitnessCalculator, error) {
	var config calcConfig
	for _, op := range ops {
		op(&config)
	}
	if config.wasmEngine == nil {
		config.wasmEngine = NewCircom2WZWitnessCalculator
	}
	return config.wasmEngine(wasm)
}
