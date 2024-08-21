// Package client provides the mediator for processing incoming AA Transactions to the bundler.
package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/notx"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"github.com/stackup-wallet/stackup-bundler/pkg/state"
	"math/big"
)

// Client controls the end to end process of adding incoming AA transactions to the mempool. It also
// implements the required RPC methods as specified in RIP-7560.
type Client struct {
	mempool             *mempool.Mempool
	chainID             *big.Int
	rip7560TxHandler    modules.Rip7560TxHandlerFunc
	logger              logr.Logger
	getRip7560TxReceipt GetRip7560TxReceiptFunc
	getGasPrices        GetGasPricesFunc
	getGasEstimate      GetGasEstimateFunc
	txLookupLimit       uint64
}

// New initializes a new RIP-7560 client which can be extended with modules for validating Transactions
// that are allowed to be added to the mempool.
func New(
	mempool *mempool.Mempool,
	chainID *big.Int,
	txLookupLimit uint64,
) *Client {
	return &Client{
		mempool:             mempool,
		chainID:             chainID,
		rip7560TxHandler:    notx.Rip7560TxHandler,
		logger:              logger.NewZeroLogr().WithName("client"),
		getRip7560TxReceipt: getRip7560TxReceiptNotx(),
		getGasPrices:        getGasPricesNotx(),
		getGasEstimate:      getGasEstimateNoop(),
		txLookupLimit:       txLookupLimit,
	}
}

// UseLogger defines the logger object used by the Client instance based on the go-logr/logr interface.
func (i *Client) UseLogger(logger logr.Logger) {
	i.logger = logger.WithName("client")
}

// UseModules defines the UserOpHandlers to process a userOp after it has gone through the standard checks.
func (i *Client) UseModules(handlers ...modules.Rip7560TxHandlerFunc) {
	i.rip7560TxHandler = modules.ComposeUserOpHandlerFunc(handlers...)
}

// SetGetRip7560TransactionReceiptFunc defines a general function for fetching a TransactionReceipt given a userOpHash and
// EntryPoint address. This function is called in *Client.GetRip7560TransactionReceipt.
func (i *Client) SetGetRip7560TransactionReceiptFunc(fn GetRip7560TxReceiptFunc) {
	i.getRip7560TxReceipt = fn
}

// SetGetGasPricesFunc defines a general function for fetching values for maxFeePerGas and
// maxPriorityFeePerGas. This function is called in *Client.EstimateRip7560TransactionGas if given fee values are
// 0.
func (i *Client) SetGetGasPricesFunc(fn GetGasPricesFunc) {
	i.getGasPrices = fn
}

// SetGetGasEstimateFunc defines a general function for fetching an estimate for verificationGasLimit and
// callGasLimit given a userOp and EntryPoint address. This function is called in
// *Client.EstimateRip7560TransactionGas.
func (i *Client) SetGetGasEstimateFunc(fn GetGasEstimateFunc) {
	i.getGasEstimate = fn
}

// SendRip7560Transaction implements the method call for eth_sendRip7560Transaction.
// It returns true if Rip7560Transaction was accepted otherwise returns an error.
func (i *Client) SendRip7560Transaction(txArgs *transaction.TransactionArgs) (string, error) {
	// Init logger
	l := i.logger.WithName("eth_sendRip7560Transaction")
	l = l.WithValues("chain_id", i.chainID.String())

	tx := txArgs.ToTransaction()
	l = l.WithValues("txHash", tx.Hash())

	// Run through client module stack.
	ctx, err := modules.NewTxHandlerContext(
		txArgs,
		i.chainID,
		i.mempool,
	)
	if err != nil {
		l.Error(err, "eth_sendRip7560Transaction error")
		return "", err
	}
	if err := i.rip7560TxHandler(ctx); err != nil {
		l.Error(err, "eth_sendRip7560Transaction error")
		return "", err
	}

	// Add Rip-7560 transaction to mempool.
	if err := i.mempool.AddTx(ctx.Tx); err != nil {
		l.Error(err, "eth_sendRip7560Transaction error")
		return "", err
	}

	l.Info("eth_sendRip7560Transaction ok")
	return tx.Hash().String(), nil
}

// EstimateRip7560TransactionGas returns estimates for PreVerificationGas, VerificationGasLimit, and CallGasLimit
// given a UserOperation, EntryPoint address, and state OverrideSet. The signature field and current gas
// values will not be validated although there should be dummy values in place for the most reliable results
// (e.g. a signature with the correct length).
func (i *Client) EstimateRip7560TransactionGas(
	txArgs *transaction.TransactionArgs,
	os map[string]any,
) (*gas.GasEstimates, error) {
	// Init logger
	l := i.logger.WithName("eth_estimateRip7560TransactionGas")

	hash := txArgs.ToTransaction().Hash()
	l = l.WithValues("rip7560Tx_hash", hash)

	// Parse state override set.
	sos, err := state.ParseOverrideData(os)
	if err != nil {
		l.Error(err, "eth_estimateRip7560TransactionGas error")
		return nil, err
	}

	// Override op with suggested gas prices if maxFeePerGas is 0. This allows for more reliable gas
	// estimations upstream. The default balance override also ensures simulations won't revert on
	// insufficient funds.
	if txArgs.MaxFeePerGas.ToInt().Cmp(common.Big0) != 1 {
		gp, err := i.getGasPrices()
		if err != nil {
			l.Error(err, "eth_estimateRip7560TransactionGas error")
			return nil, err
		}
		gpMaxFeePerGas := hexutil.Big(*gp.MaxFeePerGas)
		gpMaxPriorityFeePerGas := hexutil.Big(*gp.MaxPriorityFeePerGas)
		txArgs.MaxFeePerGas = &gpMaxFeePerGas
		txArgs.MaxPriorityFeePerGas = &gpMaxPriorityFeePerGas
	}

	// Estimate gas limits
	vg, cg, err := i.getGasEstimate(txArgs, sos)
	if err != nil {
		l.Error(err, "eth_estimateRip7560TransactionGas error")
		return nil, err
	}

	l.Info("eth_estimateRip7560TransactionGas ok")
	return &gas.GasEstimates{
		VerificationGasLimit: big.NewInt(int64(vg)),
		CallGasLimit:         big.NewInt(int64(cg)),
	}, nil
}

// GetRip7560TransactionReceipt returns RIP-7560 transaction receipt based on a tx hash
func (i *Client) GetRip7560TransactionReceipt(
	txArgs *transaction.TransactionArgs,
) (*types.Receipt, error) {
	// Init logger
	l := i.logger.WithName("eth_getRip7560TransactionReceipt").WithValues("rip7560transaction")

	receipt, err := i.getRip7560TxReceipt(txArgs.ToTransaction().Hash().String(), i.txLookupLimit)
	if err != nil {
		l.Error(err, "getRip7560TransactionReceipt error")
	}

	l.Info("eth_getRip7560TransactionReceipt ok")
	return receipt, nil
}

// ChainID implements the method call for eth_chainId. It returns the current chainID used by the client.
// This method is used to validate that the client's chainID is in sync with the caller.
func (i *Client) ChainID() (string, error) {
	return hexutil.EncodeBig(i.chainID), nil
}
