package witness

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/big"
	"reflect"
)

// swap the order of the bytes in a slice.  This allows flipping the endianness.
func swap(b []byte) []byte {
	bs := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		bs[len(b)-1-i] = b[i]
	}
	return bs
}

// parseInput is a recurisve helper function for ParseInputs
func parseInput(v interface{}) (interface{}, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		n, ok := new(big.Int).SetString(v.(string), 0)
		if !ok {
			return nil, fmt.Errorf("Error parsing input %v", v)
		}
		return n, nil
	case reflect.Float64:
		return new(big.Int).SetInt64(int64(v.(float64))), nil
	case reflect.Slice:
		res := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			var err error
			res[i], err = parseInput(rv.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("Error parsing input %v: %w", v, err)
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("Unexpected type for input %v: %T", v, v)
	}
}

// ParseInputs parses WitnessCalc inputs from JSON that consist of a map of
// types which contain a recursive combination of: numbers, base-10 encoded
// numbers in string format, arrays.
func ParseInputs(inputsJSON []byte) (map[string]interface{}, error) {
	inputsRAW := make(map[string]interface{})
	if err := json.Unmarshal(inputsJSON, &inputsRAW); err != nil {
		return nil, err
	}
	inputs := make(map[string]interface{})
	for inputName, inputValue := range inputsRAW {
		v, err := parseInput(inputValue)
		if err != nil {
			return nil, err
		}
		inputs[inputName] = v
	}
	return inputs, nil
}

// _flatSlice is a recursive helper function for flatSlice.
func _flatSlice(acc *[]*big.Int, v interface{}) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			_flatSlice(acc, rv.Index(i).Interface())
		}
	default:
		*acc = append(*acc, v.(*big.Int))
	}
}

// flatSlice takes a structure that contains a recursive combination of slices
// and *big.Int and flattens it into a single slice.
func flatSlice(v interface{}) []*big.Int {
	res := make([]*big.Int, 0)
	_flatSlice(&res, v)
	return res
}

// fnvHash returns the 64 bit FNV-1a hash split into two 32 bit values: (MSB, LSB)
func fnvHash(s string) (int32, int32) {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	h := hash.Sum64()
	return int32(h >> 32), int32(h & 0xffffffff)
}

// fnvHash returns the 64 bit FNV-1a hash split into two 32 bit values: (MSB, LSB)
func fnvHashUint(s string) (uint32, uint32) {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	h := hash.Sum64()
	return uint32(h >> 32), uint32(h & 0xffffffff)
}
