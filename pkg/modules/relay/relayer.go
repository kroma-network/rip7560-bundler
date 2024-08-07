// Package relay implements a module for private bundlers to send batches to the EntryPoint through regular
// EOA transactions.
package relay

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/transaction"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Relayer provides a module that can relay batches with a regular EOA. Relaying batches to the EntryPoint
// through a regular transaction comes with several important notes:
//
//   - The bundler will NOT be operating as a block builder.
//   - This opens the bundler up to frontrunning.
//
// This module only works in the case of a private mempool and will not work in the P2P case where ops are
// propagated through the network and it is impossible to prevent collisions from multiple bundlers trying to
// relay the same ops.
type Relayer struct {
	eoa         *signer.EOA
	eth         *ethclient.Client
	chainID     *big.Int
	rpc         *rpc.Client
	logger      logr.Logger
	waitTimeout time.Duration
}

// New initializes a new EOA relayer for sending batches to the EntryPoint.
func New(
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	rpc *rpc.Client,
	l logr.Logger,
) *Relayer {
	return &Relayer{
		eoa:         eoa,
		eth:         eth,
		chainID:     chainID,
		rpc:         rpc,
		logger:      l.WithName("relayer"),
		waitTimeout: DefaultWaitTimeout,
	}
}

// SetWaitTimeout sets the total time to wait for a transaction to be included. When a timeout is reached, the
// BatchHandler will throw an error if the transaction has not been included or has been included but with a
// failed status.
//
// The default value is 30 seconds. Setting the value to 0 will skip waiting for a transaction to be included.
func (r *Relayer) SetWaitTimeout(timeout time.Duration) {
	r.waitTimeout = timeout
}

func (r *Relayer) SendUserOperationRip7560() modules.BatchHandlerFunc {
	// [RIP-7560] hard-coded bundler config
	creationBlock := new(big.Int).SetUint64(0)
	expectedRevenue := new(big.Int).SetUint64(0)
	bundlerId := "1"
	return func(ctx *modules.BatchHandlerCtx) error {
		transactionArgs := r.BuildTransactionArgs(ctx.Batch)
		return r.sendTransactionBundle(transactionArgs, creationBlock, expectedRevenue, bundlerId)
	}
}

func (r *Relayer) BuildTransactionArgs(batch []*userop.UserOperation) []client.UserOperationArgs {
	var transactionArgs []client.UserOperationArgs
	for _, userOp := range batch {
		txArgs := client.CreateUserOperationArgs(userOp)
		transactionArgs = append(transactionArgs, txArgs)
	}
	return transactionArgs
}

func (r *Relayer) sendTransactionBundle(transactionArgs []client.UserOperationArgs, creationBlock, expectedRevenue *big.Int, bundlerId string) error {
	var out any
	if err := r.rpc.Call(&out, "eth_sendRip7560TransactionsBundle", &transactionArgs, creationBlock, expectedRevenue, bundlerId); err != nil {
		return err
	}
	r.logger.Info("eth_sendRip7560TransactionsBundle", out)
	return nil
}

// SendUserOperation returns a BatchHandler that is used by the Bundler to send batches in a regular EOA
// transaction.
func (r *Relayer) SendUserOperation() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		opts := transaction.Opts{
			EOA:         r.eoa,
			Eth:         r.eth,
			ChainID:     ctx.ChainID,
			EntryPoint:  ctx.EntryPoint,
			Batch:       ctx.Batch,
			BaseFee:     ctx.BaseFee,
			Tip:         ctx.Tip,
			GasPrice:    ctx.GasPrice,
			GasLimit:    0,
			WaitTimeout: r.waitTimeout,
		}
		// Estimate gas for handleOps() and drop all userOps that cause unexpected reverts.
		for len(ctx.Batch) > 0 {
			est, revert, err := transaction.EstimateHandleOpsGas(&opts)

			if err != nil {
				return err
			} else if revert != nil {
				ctx.MarkOpIndexForRemoval(revert.OpIndex, revert.Reason)
			} else {
				opts.GasLimit = est
				break
			}
		}

		// Call handleOps() with gas estimate. Any userOps that cause a revert at this stage will be
		// caught and dropped in the next iteration.
		if len(ctx.Batch) > 0 {
			if txn, err := transaction.HandleOps(&opts); err != nil {
				return err
			} else {
				ctx.Data["txn_hash"] = txn.Hash().String()
			}
		}

		return nil
	}
}
