package methods

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	CreateAccountMethod = abi.NewMethod(
		"createAccount",
		"createAccount",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "owner", Type: address},
			{Name: "salt", Type: uint256},
		},
		abi.Arguments{
			{Type: address},
		},
	)
	CreateAccountSelector = hexutil.Encode(CreateAccountMethod.ID)
)
