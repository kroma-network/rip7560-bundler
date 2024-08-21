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

// TestSortByGasPriceBaseDynamic verifies that SortByGasPrice sorts the Rip-7560 transactions in a batch by highest
// effective Gas Price first.
func TestSortByGasPriceBaseDynamic(t *testing.T) {
	bf := big.NewInt(3)
	tip := big.NewInt(0)

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
		big.NewInt(6),
	)
	if err := gasprice.SortByGasPrice()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 3 {
		t.Fatalf("got length %d, want 3", len(ctx.Batch))
	} else if !testutils.IsTxsEqual(ctx.Batch[0], tx2) {
		t.Fatal("incorrect order: first tx out of place")
	} else if !testutils.IsTxsEqual(ctx.Batch[1], tx1) {
		t.Fatal("incorrect order: second tx out of place")
	} else if !testutils.IsTxsEqual(ctx.Batch[2], tx3) {
		t.Fatal("incorrect order: third tx out of place")
	}
}

// TestSortByGasPriceLegacy verifies that SortByGasPrice sorts the Rip-7560 transactions in a batch by highest
// MaxFeePerGas if the context BaseFee is nil.
func TestSortByGasPriceLegacy(t *testing.T) {
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
		big.NewInt(4),
	)
	if err := gasprice.SortByGasPrice()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 3 {
		t.Fatalf("got length %d, want 3", len(ctx.Batch))
	} else if !testutils.IsTxsEqual(ctx.Batch[0], tx3) {
		t.Fatal("incorrect order: first tx out of place")
	} else if !testutils.IsTxsEqual(ctx.Batch[1], tx2) {
		t.Fatal("incorrect order: second tx out of place")
	} else if !testutils.IsTxsEqual(ctx.Batch[2], tx1) {
		t.Fatal("incorrect order: third tx out of place")
	}
}
