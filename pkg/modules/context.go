package modules

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
)

type PendingRemovalItem struct {
	Tx     *transaction.TransactionArgs
	Reason string
}

// BatchHandlerCtx is the object passed to BatchHandler functions during the Bundler's Run process. It
// also contains a Data field for adding arbitrary key-value pairs to the context. These values will be
// logged by the Bundler at the end of each run.
type BatchHandlerCtx struct {
	Batch          []*transaction.TransactionArgs
	PendingRemoval []*PendingRemovalItem
	ChainID        *big.Int
	BaseFee        *big.Int
	Tip            *big.Int
	GasPrice       *big.Int
	Data           map[string]any
}

// NewBatchHandlerContext creates a new BatchHandlerCtx using a copy of the given batch.
func NewBatchHandlerContext(
	batch []*transaction.TransactionArgs,
	chainID *big.Int,
	baseFee *big.Int,
	tip *big.Int,
	gasPrice *big.Int,
) *BatchHandlerCtx {
	var batchAppended []*transaction.TransactionArgs
	batchAppended = append(batchAppended, batch...)

	return &BatchHandlerCtx{
		Batch:          batchAppended,
		PendingRemoval: []*PendingRemovalItem{},
		ChainID:        chainID,
		BaseFee:        baseFee,
		Tip:            tip,
		GasPrice:       gasPrice,
		Data:           make(map[string]any),
	}
}

// MarkTxIndexForRemoval will remove the tx by index from the batch and add it to the pending removal array.
// This should be used for txs that are not to be included on-chain and dropped from the mempool.
func (c *BatchHandlerCtx) MarkTxIndexForRemoval(index int, reason string) {
	var batch []*transaction.TransactionArgs
	var tx *transaction.TransactionArgs
	for i, curr := range c.Batch {
		if i == index {
			tx = curr
		} else {
			batch = append(batch, curr)
		}
	}
	if tx == nil {
		return
	}

	c.Batch = batch
	c.PendingRemoval = append(c.PendingRemoval, &PendingRemovalItem{
		Tx:     tx,
		Reason: reason,
	})
}

// TxHandlerCtx is the object passed to Rip7560TxHandler functions during the Client's SendRip7560Transaction
// process.
type TxHandlerCtx struct {
	Tx                  *transaction.TransactionArgs
	ChainID             *big.Int
	pendingSenderTxs    []*transaction.TransactionArgs
	pendingDeployerTxs  []*transaction.TransactionArgs
	pendingPaymasterTxs []*transaction.TransactionArgs
}

// NewTxHandlerContext creates a new TxHandlerCtx using a given tx.
func NewTxHandlerContext(
	txArgs *transaction.TransactionArgs,
	chainID *big.Int,
	mem *mempool.Mempool,
) (*TxHandlerCtx, error) {
	// Fetch any pending Transactions in the mempool by entity
	pso, err := mem.GetTxs(txArgs.GetSender())
	if err != nil {
		return nil, err
	}
	pdo, err := mem.GetTxs(txArgs.GetDeployer())
	if err != nil {
		return nil, err
	}
	ppo, err := mem.GetTxs(txArgs.GetPaymaster())
	if err != nil {
		return nil, err
	}

	return &TxHandlerCtx{
		Tx:                  txArgs,
		ChainID:             chainID,
		pendingSenderTxs:    pso,
		pendingDeployerTxs:  pdo,
		pendingPaymasterTxs: ppo,
	}, nil
}

func (c *TxHandlerCtx) GetRip7560Transaction() *types.Rip7560AccountAbstractionTx {
	return c.Tx.ToTransaction().Rip7560TransactionData()
}

func (c *TxHandlerCtx) GetSender() common.Address {
	if c.Tx.ToTransaction().Rip7560TransactionData().Sender == nil {
		return common.Address{}
	}
	return *c.Tx.ToTransaction().Rip7560TransactionData().Sender
}

func (c *TxHandlerCtx) GetPaymaster() common.Address {
	if c.Tx.ToTransaction().Rip7560TransactionData().Paymaster == nil {
		return common.Address{}
	}
	return *c.Tx.ToTransaction().Rip7560TransactionData().Paymaster
}

func (c *TxHandlerCtx) GetDeployer() common.Address {
	if c.Tx.ToTransaction().Rip7560TransactionData().Deployer == nil {
		return common.Address{}
	}
	return *c.Tx.ToTransaction().Rip7560TransactionData().Deployer
}

func (c *TxHandlerCtx) GetPaymasterData() []byte {
	return c.Tx.ToTransaction().Rip7560TransactionData().PaymasterData
}

func (c *TxHandlerCtx) GetDeployerData() []byte {
	return c.Tx.ToTransaction().Rip7560TransactionData().DeployerData
}

// GetPendingSenderTxs returns all pending Rip-7560 transactions in the mempool by the same sender.
func (c *TxHandlerCtx) GetPendingSenderTxs() []*transaction.TransactionArgs {
	return c.pendingSenderTxs
}

// GetPendingFactoryTxs returns all pending Rip-7560 transactions in the mempool by the same factory.
func (c *TxHandlerCtx) GetPendingFactoryTxs() []*transaction.TransactionArgs {
	return c.pendingDeployerTxs
}

// GetPendingPaymasterTxs returns all pending Rip-7560 transactions in the mempool by the same paymaster.
func (c *TxHandlerCtx) GetPendingPaymasterTxs() []*transaction.TransactionArgs {
	return c.pendingPaymasterTxs
}
