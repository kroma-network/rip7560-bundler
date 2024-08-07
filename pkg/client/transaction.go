package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type UserOperationArgs struct {
	From                 string `json:"from"`
	SubType              string `json:"subType"`
	To                   string `json:"to"`
	GasPrice             string `json:"gasPrice"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	Data                 string `json:"data"`
	Sender               string `json:"sender"`
	PaymasterData        string `json:"paymasterData"`
	DeployerData         string `json:"deployerData"`
	BuilderFee           string `json:"builderFee"`
	Gas                  string `json:"gas"`
	ValidationGas        string `json:"validationGas"`
	PaymasterGas         string `json:"paymasterGas"`
	PostOpGas            string `json:"postOpGas"`
	Signature            string `json:"signature"`
	BigNonce             string `json:"bigNonce"`
}

func CreateUserOperationArgs(userOp *userop.UserOperation) UserOperationArgs {
	// [RIP-7560] hard-coded fixed config
	from := "0x0000000000000000000000000000000000007560"
	to := common.Address{}
	subType := "0x1"
	return UserOperationArgs{
		From:                 from,
		SubType:              subType,
		To:                   hexutil.Encode(to.Bytes()),
		GasPrice:             hexutil.EncodeBig(userOp.MaxFeePerGas),
		MaxFeePerGas:         hexutil.EncodeBig(userOp.MaxFeePerGas),
		MaxPriorityFeePerGas: hexutil.EncodeBig(userOp.MaxPriorityFeePerGas),
		Data:                 hexutil.Encode(userOp.CallData),
		Sender:               hexutil.Encode(userOp.Sender.Bytes()),
		PaymasterData:        hexutil.Encode(userOp.PaymasterAndData),
		DeployerData:         hexutil.Encode(userOp.InitCode),
		BuilderFee:           hexutil.EncodeBig(userOp.PreVerificationGas),
		Gas:                  hexutil.EncodeBig(userOp.CallGasLimit),
		ValidationGas:        hexutil.EncodeBig(userOp.VerificationGasLimit),
		PaymasterGas:         hexutil.EncodeBig(userOp.PaymasterVerificationGasLimit),
		PostOpGas:            hexutil.EncodeBig(userOp.PaymasterPostOpGasLimit),
		Signature:            hexutil.Encode(userOp.Signature),
		BigNonce:             hexutil.EncodeBig(userOp.Nonce),
	}
}
