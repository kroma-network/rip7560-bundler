package client

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/fees"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/filter"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"github.com/stackup-wallet/stackup-bundler/pkg/state"
)

// GetRip7560TxReceiptFunc is a general interface for fetching a Rip-7560 transaction Receipt given a txHash,
// EntryPoint address, and block range.
type GetRip7560TxReceiptFunc = func(txHash string) (*types.Receipt, error)

func getRip7560TxReceiptNotx() GetRip7560TxReceiptFunc {
	return func(txHash string) (*types.Receipt, error) {
		return nil, nil
	}
}

// GetRip7560TransactionReceiptWithEthClient returns an implementation of GetRip7560UserTxReceiptFunc that relies on an eth
// client to fetch a Rip-7560 transactionReceipt.
func GetRip7560TransactionReceiptWithEthClient(eth *ethclient.Client) GetRip7560TxReceiptFunc {
	return func(txHash string) (*types.Receipt, error) {
		return filter.GetRip7560TransactionReceipt(eth, txHash)
	}
}

// GetGasPricesFunc is a general interface for fetching values for maxFeePerGas and maxPriorityFeePerGas.
type GetGasPricesFunc = func() (*fees.GasPrices, error)

func getGasPricesNotx() GetGasPricesFunc {
	return func() (*fees.GasPrices, error) {
		return &fees.GasPrices{
			MaxFeePerGas:         big.NewInt(0),
			MaxPriorityFeePerGas: big.NewInt(0),
		}, nil
	}
}

// GetGasPricesWithEthClient returns an implementation of GetGasPricesFunc that relies on an eth client to
// fetch values for maxFeePerGas and maxPriorityFeePerGas.
func GetGasPricesWithEthClient(eth *ethclient.Client) GetGasPricesFunc {
	return func() (*fees.GasPrices, error) {
		return fees.NewGasPrices(eth)
	}
}

// GetGasEstimateFunc is a general interface for fetching an estimate for verificationGasLimit and
// callGasLimit given a Rip-7560 transaction.
type GetGasEstimateFunc = func(
	aaTxArgs *transaction.TransactionArgs,
	sos state.OverrideSet,
) (verificationGas uint64, callGas uint64, err error)

func getGasEstimateNoop() GetGasEstimateFunc {
	return func(
		aaTxArgs *transaction.TransactionArgs,
		sos state.OverrideSet,
	) (verificationGas uint64, callGas uint64, err error) {
		return 0, 0, nil
	}
}

// GetGasEstimateWithEthClient returns an implementation of GetGasEstimateFunc that relies on an eth client to
// fetch an estimate for verificationGasLimit and callGasLimit.
func GetGasEstimateWithEthClient(
	rpc *rpc.Client,
	chain *big.Int,
	maxGasLimit *big.Int,
) GetGasEstimateFunc {
	return func(
		aaTxArgs *transaction.TransactionArgs,
		sos state.OverrideSet,
	) (verificationGas uint64, callGas uint64, err error) {
		// same as ethapi/rip7560api/Rip7560UsedGas
		type Rip7560UsedGas struct {
			ValidationGas hexutil.Uint64 `json:"validationGas"`
			ExecutionGas  hexutil.Uint64 `json:"executionGas"`
		}

		var res Rip7560UsedGas
		if err := rpc.CallContext(context.Background(), &res, "eth_estimateRip7560TransactionGas", aaTxArgs, "latest", sos); err != nil {
			return 0, 0, err
		}
		return uint64(res.ValidationGas), uint64(res.ExecutionGas), nil
	}
}

func MapToTransactionArgs(input map[string]interface{}) (transaction.TransactionArgs, error) {
	var txArgs transaction.TransactionArgs
	data, err := json.Marshal(input)
	if err != nil {
		return txArgs, err
	}
	err = json.Unmarshal(data, &txArgs)
	if err != nil {
		return txArgs, err
	}
	return txArgs, nil
}
