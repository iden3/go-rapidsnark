package witness

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
)

// parseInput is a recursive helper function for ParseInputs
func parseInput(v interface{}) (interface{}, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		n, ok := new(big.Int).SetString(v.(string), 0)
		if !ok {
			return nil, fmt.Errorf("error parsing input %v", v)
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
				return nil, fmt.Errorf("error parsing input %v: %w", v, err)
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("unexpected type for input %v: %T", v, v)
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
