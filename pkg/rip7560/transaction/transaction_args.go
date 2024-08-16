package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// TransactionArgs represents the arguments to construct a new transaction
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to" rlp:"nil"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty" rlp:"nil"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty" rlp:"nil"`

	// Introduced by RIP-7560 Transaction
	Sender        *common.Address `json:"sender"`
	Signature     *hexutil.Bytes  `json:"signature"`
	Paymaster     *common.Address `json:"paymaster,omitempty"     rlp:"nil"`
	PaymasterData *hexutil.Bytes  `json:"paymasterData,omitempty" rlp:"nil"`
	Deployer      *common.Address `json:"deployer,omitempty"      rlp:"nil"`
	DeployerData  *hexutil.Bytes  `json:"deployerData,omitempty"  rlp:"nil"`
	BuilderFee    *hexutil.Big    `json:"builderFee"`
	ValidationGas *hexutil.Uint64 `json:"verificationGasLimit"`
	PaymasterGas  *hexutil.Uint64 `json:"paymasterVerificationGasLimit"`
	PostOpGas     *hexutil.Uint64 `json:"paymasterPostOpGasLimit"`
	BigNonce      *hexutil.Big    `json:"bigNonce"`
}

// from retrieves the transaction sender address.
func (args *TransactionArgs) from() common.Address {
	if args.From == nil {
		return common.Address{}
	}
	return *args.From
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *TransactionArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

func (args *TransactionArgs) sender() *common.Address {
	if args.Sender != nil {
		return args.Sender
	}
	return nil
}

func (args *TransactionArgs) signature() []byte {
	if args.Signature != nil {
		return *args.Signature
	}
	return nil
}

func (args *TransactionArgs) paymaster() *common.Address {
	if args.Paymaster != nil {
		return args.Paymaster
	}
	return nil
}

func (args *TransactionArgs) paymasterData() []byte {
	if args.PaymasterData != nil {
		return *args.PaymasterData
	}
	return nil
}

func (args *TransactionArgs) paymasterGas() uint64 {
	if args.PaymasterGas != nil {
		return uint64(*args.PaymasterGas)
	}
	return 0
}

func (args *TransactionArgs) validationGas() uint64 {
	if args.ValidationGas != nil {
		return uint64(*args.ValidationGas)
	}
	return 0
}

func (args *TransactionArgs) postOpGas() uint64 {
	if args.PostOpGas != nil {
		return uint64(*args.PostOpGas)
	}
	return 0
}

func (args *TransactionArgs) deployer() *common.Address {
	if args.Deployer != nil {
		return args.Deployer
	}
	return nil
}

func (args *TransactionArgs) deployerData() []byte {
	if args.DeployerData != nil {
		return *args.DeployerData
	}
	return nil
}

func (args *TransactionArgs) GetSender() common.Address {
	if args.Sender != nil {
		return *args.Sender
	}
	return common.Address{}
}

func (args *TransactionArgs) GetSignature() []byte {
	if args.Signature != nil {
		return *args.Signature
	}
	return nil
}

func (args *TransactionArgs) GetPaymaster() common.Address {
	if args.Paymaster != nil {
		return *args.Paymaster
	}
	return common.Address{}
}

func (args *TransactionArgs) GetPaymasterData() []byte {
	if args.PaymasterData != nil {
		return *args.PaymasterData
	}
	return nil
}

func (args *TransactionArgs) GetPaymasterGas() uint64 {
	if args.PaymasterGas != nil {
		return uint64(*args.PaymasterGas)
	}
	return 0
}

func (args *TransactionArgs) GetValidationGas() uint64 {
	if args.ValidationGas != nil {
		return uint64(*args.ValidationGas)
	}
	return 0
}

func (args *TransactionArgs) GetPostOpGas() uint64 {
	if args.PostOpGas != nil {
		return uint64(*args.PostOpGas)
	}
	return 0
}

func (args *TransactionArgs) GetDeployer() common.Address {
	if args.Deployer != nil {
		return *args.Deployer
	}
	return common.Address{}
}

func (args *TransactionArgs) GetDeployerData() []byte {
	if args.DeployerData != nil {
		return *args.DeployerData
	}
	return nil
}

// GetDynamicGasPrice returns the effective gas price paid by the RIP-7560 transaction given a basefee.
// If basefee is nil, it will assume a value of 0.
func (args *TransactionArgs) GetDynamicGasPrice(basefee *big.Int) *big.Int {
	bf := basefee
	if bf == nil {
		bf = big.NewInt(0)
	}

	var maxPriorityFeePerGas *big.Int
	var maxFeePerGas *big.Int
	maxPriorityFeePerGas = (*big.Int)(args.MaxPriorityFeePerGas)
	maxFeePerGas = (*big.Int)(args.MaxFeePerGas)
	gp := big.NewInt(0).Add(bf, maxPriorityFeePerGas)
	if gp.Cmp(maxFeePerGas) == 1 {
		return maxFeePerGas
	}
	return gp
}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) toTransaction() *types.Transaction {
	var data types.TxData
	if args.Sender == nil {
		log.Error("RIP-7560 transaction Sender not exists")
		return nil
	}
	al := types.AccessList{}
	if args.AccessList != nil {
		al = *args.AccessList
	}
	rip7560Tx := types.Rip7560AccountAbstractionTx{
		To:         &common.Address{},
		ChainID:    (*big.Int)(args.ChainID),
		Gas:        uint64(*args.Gas),
		GasFeeCap:  (*big.Int)(args.MaxFeePerGas),
		GasTipCap:  (*big.Int)(args.MaxPriorityFeePerGas),
		Value:      (*big.Int)(args.Value),
		Data:       args.data(),
		AccessList: al,
		// RIP-7560 parameters
		Sender:        args.Sender,
		Signature:     args.signature(),
		Paymaster:     args.paymaster(),
		PaymasterData: args.paymasterData(),
		Deployer:      args.deployer(),
		DeployerData:  args.deployerData(),
		BuilderFee:    (*big.Int)(args.BuilderFee),
		ValidationGas: args.validationGas(),
		PaymasterGas:  args.paymasterGas(),
		PostOpGas:     args.postOpGas(),
		// RIP-7712 parameter
		BigNonce: (*big.Int)(args.BigNonce),
	}
	data = &rip7560Tx
	hash := types.NewTx(data).Hash()
	log.Info("RIP-7560 transaction created", "sender", rip7560Tx.Sender.Hex(), "hash", hash)
	return types.NewTx(data)
}

// ToTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) ToTransaction() *types.Transaction {
	return args.toTransaction()
}
