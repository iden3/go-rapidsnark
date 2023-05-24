package witness

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"github.com/iden3/go-iden3-crypto/utils"
)

type Option func(cfg *calcConfig)

func WithWasmEngine(calculator func([]byte) (CalculatorImpl, error)) Option {
	return func(cfg *calcConfig) {
		cfg.wasmEngine = calculator
	}
}

type CalculatorImpl interface {
	Calculate(inputs map[string]interface{},
		sanityCheck bool) (wtns Witness, err error)
}

type Calculator interface {
	CalculateWitness(inputs map[string]interface{},
		sanityCheck bool) ([]*big.Int, error)
	CalculateBinWitness(inputs map[string]interface{},
		sanityCheck bool) ([]byte, error)
	CalculateWTNSBin(inputs map[string]interface{},
		sanityCheck bool) ([]byte, error)
}

type calcConfig struct {
	wasmEngine func([]byte) (CalculatorImpl, error)
}

type calc struct {
	wc CalculatorImpl
}

func (c *calc) CalculateWitness(inputs map[string]interface{},
	sanityCheck bool) ([]*big.Int, error) {

	wtns, err := c.wc.Calculate(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}
	return wtns.Witness, nil
}

func (c *calc) CalculateBinWitness(inputs map[string]interface{},
	sanityCheck bool) ([]byte, error) {

	wtns, err := c.wc.Calculate(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.Grow(wtns.N32 * 4 * len(wtns.Witness))
	for _, i := range wtns.Witness {
		bs := utils.SwapEndianness(i.Bytes())
		b.Write(bs)
		if len(bs) < wtns.N32*4 {
			for j := 0; j < (wtns.N32*4)-len(bs); j++ {
				b.WriteByte(0)
			}
		}
	}

	return b.Bytes(), nil
}

func (c *calc) CalculateWTNSBin(inputs map[string]interface{},
	sanityCheck bool) ([]byte, error) {

	wtns, err := c.wc.Calculate(inputs, sanityCheck)
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)

	n8 := wtns.N32 * 4
	idSection2length := n8 * len(wtns.Witness)

	totalLn := 4 + 4 + 4 + 4 + 8 + 4 + n8 + 4 + 4 + 8 + idSection2length
	buff.Grow(totalLn)

	// wtns
	_, _ = buff.Write([]byte("wtns"))

	//version 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	//number of sections: 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	//id section 1
	_ = binary.Write(buff, binary.LittleEndian, uint32(1))

	//id section 1 length in 64bytes
	idSection1length := 8 + n8
	_ = binary.Write(buff, binary.LittleEndian, uint64(idSection1length))

	//this.n32
	_ = binary.Write(buff, binary.LittleEndian, uint32(n8))

	err = writeInt(buff, wtns.Prime, n8)
	if err != nil {
		return nil, err
	}

	// witness size
	_ = binary.Write(buff, binary.LittleEndian, uint32(len(wtns.Witness)))

	//id section 2
	_ = binary.Write(buff, binary.LittleEndian, uint32(2))

	// section 2 length
	_ = binary.Write(buff, binary.LittleEndian, uint64(idSection2length))

	for _, i := range wtns.Witness {
		err = writeInt(buff, i, n8)
		if err != nil {
			return nil, err
		}
	}

	return buff.Bytes(), nil
}

func writeInt(out io.Writer, i *big.Int, bytesLn int) error {
	bs := utils.SwapEndianness(i.Bytes())
	_, err := out.Write(bs)
	if err != nil {
		return err
	}
	if len(bs) < bytesLn {
		_, err = out.Write(make([]byte, bytesLn-len(bs)))
	}

	return err
}

func NewCalculator(wasm []byte, ops ...Option) (Calculator, error) {
	var config calcConfig
	for _, op := range ops {
		op(&config)
	}
	if config.wasmEngine == nil {
		return nil, errors.New("witness calculator wasm engine not set")
	}
	wc, err := config.wasmEngine(wasm)
	if err != nil {
		return nil, err
	}
	return &calc{wc: wc}, nil
}

type Witness struct {
	// number of int32 values required to represent the *big.Int
	N32     int
	Prime   *big.Int
	Witness []*big.Int
}
