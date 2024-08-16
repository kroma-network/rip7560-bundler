package gasprice_test

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
)

// TestFilterUnderpricedDynamic verifies that FilterUnderpriced will remove all UserOperations from a batch
// where the effective gas price is less than the expected bundler transaction's gas price.
func TestFilterUnderpricedDynamic(t *testing.T) {
	bf := big.NewInt(4)
	tip := big.NewInt(1)

	tx1 := testutils.MockValidInitRip7560Tx()
	*tx1.MaxFeePerGas = hexutil.Big(*big.NewInt(4))
	*tx1.MaxPriorityFeePerGas = hexutil.Big(*big.NewInt(3))

	tx2 := testutils.MockValidInitRip7560Tx()
	*tx2.Sender = testutils.ValidAddress2
	*tx2.MaxFeePerGas = hexutil.Big(*big.NewInt(5))
	*tx2.MaxPriorityFeePerGas = hexutil.Big(*big.NewInt(2))

	tx3 := testutils.MockValidInitRip7560Tx()
	*tx3.Sender = testutils.ValidAddress3
	*tx3.MaxFeePerGas = hexutil.Big(*big.NewInt(6))
	*tx3.MaxPriorityFeePerGas = hexutil.Big(*big.NewInt(1))

	ctx := modules.NewBatchHandlerContext(
		[]*transaction.TransactionArgs{tx1, tx2, tx3},
		testutils.ChainID,
		bf,
		tip,
		big.NewInt(10),
	)
	if err := gasprice.FilterUnderpriced()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 2 {
		t.Fatalf("got length %d, want 2", len(ctx.Batch))
	} else if !testutils.IsTxsEqual(ctx.Batch[0], tx2) {
		t.Fatal("incorrect order: first op out of place")
	} else if !testutils.IsTxsEqual(ctx.Batch[1], tx3) {
		t.Fatal("incorrect order: second op out of place")
	}
}

// TestFilterUnderpricedGasPrice verifies that FilterUnderpriced will remove all UserOperations from a batch
// where the MaxFeePerGas is less than the context GasPrice.
func TestFilterUnderpricedGasPrice(t *testing.T) {
	tx1 := testutils.MockValidInitRip7560Tx()
	*tx1.MaxFeePerGas = hexutil.Big(*big.NewInt(4))
	*tx1.MaxPriorityFeePerGas = hexutil.Big(*big.NewInt(4))

	tx2 := testutils.MockValidInitRip7560Tx()
	*tx2.Sender = testutils.ValidAddress2
	*tx2.MaxFeePerGas = hexutil.Big(*big.NewInt(5))
	*tx2.MaxPriorityFeePerGas = hexutil.Big(*big.NewInt(5))

	tx3 := testutils.MockValidInitRip7560Tx()
	*tx3.Sender = testutils.ValidAddress3
	*tx3.MaxFeePerGas = hexutil.Big(*big.NewInt(6))
	*tx3.MaxPriorityFeePerGas = hexutil.Big(*big.NewInt(6))

	ctx := modules.NewBatchHandlerContext(
		[]*transaction.TransactionArgs{tx1, tx2, tx3},
		testutils.ChainID,
		nil,
		nil,
		big.NewInt(5),
	)
	if err := gasprice.FilterUnderpriced()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 2 {
		t.Fatalf("got length %d, want 2", len(ctx.Batch))
	} else if !testutils.IsTxsEqual(ctx.Batch[0], tx2) {
		t.Fatal("incorrect order: first op out of place")
	} else if !testutils.IsTxsEqual(ctx.Batch[1], tx3) {
		t.Fatal("incorrect order: second op out of place")
	}
}
