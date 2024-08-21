package mempool

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestAddTxToMempool verifies that a UserOperation can be added to the mempool and later retrieved without
// any changes.
func TestAddTxToMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	txArgs := testutils.MockValidInitRip7560Tx()

	if err := mem.AddTx(txArgs); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memTxs, err := mem.GetTxs(txArgs.GetSender())
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memTxs) != 1 {
		t.Fatalf("got length %d, want 1", len(memTxs))
	}

	if !testutils.IsTxsEqual(txArgs, memTxs[0]) {
		t.Fatalf("txs not equal: %s", testutils.GetTxsDiff(txArgs, memTxs[0]))
	}
}

// TestReplaceTxInMempool verifies that a RIP-7560 transaction with same Sender and Nonce can replace another
// UserOperation already in the mempool.
func TestReplaceTxInMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	tx1 := testutils.MockValidInitRip7560Tx()
	tx2 := testutils.MockValidInitRip7560Tx()
	tx2.MaxPriorityFeePerGas = (*hexutil.Big)(big.NewInt(0).Add((*big.Int)(tx1.MaxPriorityFeePerGas), common.Big1))

	if err := mem.AddTx(tx1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddTx(tx2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memTxs, err := mem.GetTxs(tx2.GetSender())
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memTxs) != 1 {
		t.Fatalf("got length %d, want 1", len(memTxs))
	}

	if !testutils.IsTxsEqual(tx2, memTxs[0]) {
		t.Fatalf("txs not equal: %s", testutils.GetTxsDiff(tx2, memTxs[0]))
	}
}

// TestRemoveTxsFromMempool verifies that a UserOperation can be added to the mempool and later removed.
func TestRemoveTxsFromMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	txArgs := testutils.MockValidInitRip7560Tx()

	if err := mem.AddTx(txArgs); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if err := mem.RemoveTxs(txArgs); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memTxs, err := mem.GetTxs(txArgs.GetSender())
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memTxs) != 0 {
		t.Fatalf("got length %d, want 0", len(memTxs))
	}
}

// TestDumpFromMempool verifies that bundles are being built with UserOperations in the mempool. Ordering is
// FIFO and more specific sorting and filtering is left up to downstream modules to implement.
func TestDumpFromMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)

	tx1 := testutils.MockValidInitRip7560Tx()
	tx1.MaxFeePerGas = (*hexutil.Big)(big.NewInt(4))
	tx1.MaxPriorityFeePerGas = (*hexutil.Big)(big.NewInt(3))

	tx2 := testutils.MockValidInitRip7560Tx()
	tx2.Sender = &testutils.ValidAddress2
	tx2.MaxFeePerGas = (*hexutil.Big)(big.NewInt(5))
	tx2.MaxPriorityFeePerGas = (*hexutil.Big)(big.NewInt(2))

	tx3 := testutils.MockValidInitRip7560Tx()
	tx3.Sender = &testutils.ValidAddress3
	tx3.MaxFeePerGas = (*hexutil.Big)(big.NewInt(6))
	tx3.MaxPriorityFeePerGas = (*hexutil.Big)(big.NewInt(1))

	if err := mem.AddTx(tx1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddTx(tx2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddTx(tx3); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if memTxs, err := mem.Dump(); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(memTxs) != 3 {
		t.Fatalf("got length %d, want 3", len(memTxs))
	} else if !testutils.IsTxsEqual(memTxs[0], tx1) {
		t.Fatal("incorrect order: first tx out of place")
	} else if !testutils.IsTxsEqual(memTxs[1], tx2) {
		t.Fatal("incorrect order: second tx out of place")
	} else if !testutils.IsTxsEqual(memTxs[2], tx3) {
		t.Fatal("incorrect order: third tx out of place")
	}
}

// TestNewMempoolLoadsFromDisk verifies that a new Mempool instance is built from txs saved in the DB without
// including txs previously removed.
func TestNewMempoolLoadsFromDisk(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem1, _ := New(db)
	tx1 := testutils.MockValidInitRip7560Tx()
	tx2 := testutils.MockValidInitRip7560Tx()
	tx2.BigNonce = (*hexutil.Big)(big.NewInt(0).Add((*big.Int)(tx1.BigNonce), common.Big1))
	tx2.MaxPriorityFeePerGas = (*hexutil.Big)(big.NewInt(0).Add((*big.Int)(tx1.MaxPriorityFeePerGas), common.Big1))

	if err := mem1.AddTx(tx1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem1.AddTx(tx2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem1.RemoveTxs(tx1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	mem2, _ := New(db)
	memTxs, err := mem2.GetTxs(tx2.GetSender())
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memTxs) != 1 {
		t.Fatalf("got length %d, want 1", len(memTxs))
	}

	if !testutils.IsTxsEqual(tx2, memTxs[0]) {
		t.Fatalf("txsS not equal: %s", testutils.GetTxsDiff(tx2, memTxs[0]))
	}
}
