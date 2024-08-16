package methods

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	ValidatePaymasterTransactionMethod = abi.NewMethod(
		"validatePaymasterTransaction",
		"validatePaymasterTransaction",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "version", Type: uint256},
			{Name: "txHash", Type: bytes32},
			{Name: "transaction", Type: bytes},
		},
		abi.Arguments{
			{Type: bytes},
		},
	)
	ValidatePaymasterTransactionSelector = hexutil.Encode(ValidatePaymasterTransactionMethod.ID)
)

type validatePaymasterTransactionOutput struct {
	Context []byte
}

func DecodevalidatePaymasterTransactionOutputOutput(ret any) (*validatePaymasterTransactionOutput, error) {
	hex, ok := ret.(string)
	if !ok {
		return nil, errors.New("validatePaymasterTransactionOutput: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, fmt.Errorf("validatePaymasterTransactionOutput: %s", err)
	}

	args, err := ValidatePaymasterTransactionMethod.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("validatePaymasterTransactionOutput: %s", err)
	}
	if len(args) != 2 {
		return nil, fmt.Errorf(
			"validatePaymasterTransactionOutput: invalid args length: expected 2, got %d",
			len(args),
		)
	}

	ctx, ok := args[0].([]byte)
	if !ok {
		return nil, errors.New("validatePaymasterTransactionOutput: cannot assert type: hex is not of type string")
	}

	return &validatePaymasterTransactionOutput{
		Context: ctx,
	}, nil
}
