package transaction

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var (
	validate = validator.New()
	onlyOnce = sync.Once{}

	ErrBadTransactionArgsData = errors.New("cannot decode TransactionArgs")
)

func exactFieldMatch(mapKey, fieldName string) bool {
	return mapKey == fieldName
}

func decodeTxTypes(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	// String to common.Address conversion
	if f.Kind() == reflect.String && t.Kind() == reflect.Array {
		return common.HexToAddress(data.(string)), nil
	}

	// String to big.Int conversion
	if f.Kind() == reflect.String && t.Kind() == reflect.Struct {
		n := new(big.Int)
		n, ok := n.SetString(data.(string), 0)
		if !ok {
			return nil, errors.New("bigInt conversion failed")
		}
		return n, nil
	}

	// Float64 to big.Int conversion
	if f.Kind() == reflect.Float64 && t.Kind() == reflect.Struct {
		n, ok := data.(float64)
		if !ok {
			return nil, errors.New("bigInt conversion failed")
		}
		return big.NewInt(int64(n)), nil
	}

	// String to []byte conversion
	if f.Kind() == reflect.String && t.Kind() == reflect.Slice {
		byteStr := data.(string)
		if len(byteStr) < 2 || byteStr[:2] != "0x" {
			return nil, errors.New("not byte string")
		}

		b, err := hex.DecodeString(byteStr[2:])
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	if f.Kind() == reflect.String && t.Kind() == reflect.Ptr {
		strData, ok := data.(string)
		if !ok {
			return data, nil
		}
		switch t.Elem() {
		// String to *common.Address conversion
		case reflect.TypeOf(common.Address{}):
			addressVal := common.HexToAddress(strData)
			return addressVal, nil
		// String to *hexutil.Uint64 conversion
		case reflect.TypeOf(hexutil.Uint64(0)):
			uintVal := new(hexutil.Uint64)
			err := uintVal.UnmarshalText([]byte(strData))
			if err != nil {
				return nil, err
			}
			return uintVal, nil

		// String to *hexutil.Big conversion
		case reflect.TypeOf(hexutil.Big{}):
			bigVal := new(hexutil.Big)
			err := bigVal.UnmarshalText([]byte(strData))
			if err != nil {
				return nil, err
			}
			return bigVal, nil

		// // String to *hexutil.Bytes conversion
		case reflect.TypeOf(hexutil.Bytes{}):
			bytesVal := new(hexutil.Bytes)
			err := bytesVal.UnmarshalText([]byte(strData))
			if err != nil {
				return nil, err
			}
			return bytesVal, nil
		}
	}

	//// String to hexutil.Uint64 conversion
	//if f == reflect.String && t == reflect.TypeOf(new(hexutil.Uint64)).Kind() {
	//	str, ok := data.(string)
	//	if !ok {
	//		return data, nil
	//	}
	//	uintVal := new(hexutil.Uint64)
	//	err := uintVal.UnmarshalText([]byte(str))
	//	return &uintVal, err
	//}
	//
	//// String to hexutil.Big conversion
	//if f == reflect.String && t == reflect.TypeOf(new(hexutil.Big)).Kind() {
	//	str, ok := data.(string)
	//	if !ok {
	//		return data, nil
	//	}
	//	bigVal := new(hexutil.Big)
	//	err := bigVal.UnmarshalText([]byte(str))
	//	return &bigVal, err
	//}
	//
	//// String to hexutil.Bytes conversion
	//if f == reflect.String && t == reflect.TypeOf(new(hexutil.Bytes)).Kind() {
	//	if t != reflect.TypeOf(hexutil.Bytes{}).Kind() {
	//		return data, nil
	//	}
	//	str, ok := data.(string)
	//	if !ok {
	//		return data, nil
	//	}
	//	bytesVal := new(hexutil.Bytes)
	//	err := bytesVal.UnmarshalText([]byte(str))
	//	return &bytesVal, err
	//}

	return data, nil
}

func validateAddressType(field reflect.Value) interface{} {
	value, ok := field.Interface().(common.Address)
	if !ok || value == common.HexToAddress("0x") {
		return nil
	}

	return field
}

func validateBigIntType(field reflect.Value) interface{} {
	value, ok := field.Interface().(big.Int)
	if !ok || value.Cmp(big.NewInt(0)) == -1 {
		return nil
	}

	return field
}

// New decodes a map into a TransactionArgs object and validates all the fields are correctly typed.
func New(data map[string]any) (*TransactionArgs, error) {
	var txArgs TransactionArgs

	// Convert map to struct
	config := &mapstructure.DecoderConfig{
		DecodeHook: decodeTxTypes,
		Result:     &txArgs,
		ErrorUnset: false,
		TagName:    "json",
		//MatchName:  exactFieldMatch,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(data); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBadTransactionArgsData, err)
	}

	// Validate struct
	onlyOnce.Do(func() {
		validate.RegisterCustomTypeFunc(validateAddressType, common.Address{})
		validate.RegisterCustomTypeFunc(validateBigIntType, big.Int{})
	})
	err = validate.Struct(txArgs)
	if err != nil {
		return nil, err
	}

	return &txArgs, nil
}
