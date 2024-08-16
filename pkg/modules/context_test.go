package modules

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
)

func TestNoPendingTxs(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := mempool.New(db)
	tx := testutils.MockValidInitRip7560Tx()
	*tx.DeployerData = []byte{}
	*tx.PaymasterData = []byte{}

	ctx, err := NewTxHandlerContext(
		tx,
		testutils.ChainID,
		mem,
	)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	} else if pso := ctx.GetPendingSenderTxs(); len(pso) != 0 {
		t.Fatalf("pending sender txs: want 0, got %d", len(pso))
	} else if pfo := ctx.GetPendingFactoryTxs(); len(pfo) != 0 {
		t.Fatalf("pending deployer txs: want 0, got %d", len(pfo))
	} else if ppo := ctx.GetPendingPaymasterTxs(); len(ppo) != 0 {
		t.Fatalf("pending paymaster ops: want 0, got %d", len(ppo))
	}
}

func TestGetPendingSenderTxs(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := mempool.New(db)
	tx := testutils.MockValidInitRip7560Tx()
	*tx.DeployerData = []byte{}
	*tx.PaymasterData = []byte{}

	penTx1 := testutils.MockValidInitRip7560Tx()
	_ = mem.AddTx(penTx1)

	penTx2 := testutils.MockValidInitRip7560Tx()
	*penTx2.BigNonce = (hexutil.Big)(*big.NewInt(0).Add((*big.Int)(penTx1.BigNonce), common.Big1))
	_ = mem.AddTx(penTx2)

	penTx3 := testutils.MockValidInitRip7560Tx()
	*penTx3.BigNonce = (hexutil.Big)(*big.NewInt(0).Add((*big.Int)(penTx2.BigNonce), common.Big1))
	_ = mem.AddTx(penTx3)

	ctx, err := NewTxHandlerContext(
		tx,
		testutils.ChainID,
		mem,
	)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	expectedPenTxs := []*transaction.TransactionArgs{penTx3, penTx2, penTx1}
	penTxs := ctx.GetPendingSenderTxs()
	if len(penTxs) != len(expectedPenTxs) {
		t.Fatalf("got length %d, want %d", len(penTxs), len(expectedPenTxs))
	}

	for i, penOp := range penTxs {
		if !testutils.IsTxsEqual(penOp, expectedPenTxs[i]) {
			t.Fatalf("ops not equal: %s", testutils.GetTxsDiff(penOp, expectedPenTxs[i]))
		}
	}
}

func TestGetPendingFactoryTxs(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := mempool.New(db)
	tx := testutils.MockValidInitRip7560Tx()
	*tx.DeployerData = []byte{}
	*tx.PaymasterData = []byte{}

	penTx1 := testutils.MockValidInitRip7560Tx()
	*penTx1.Sender = testutils.ValidAddress1
	*penTx1.DeployerData = testutils.DummyDeployerData
	_ = mem.AddTx(penTx1)

	penTx2 := testutils.MockValidInitRip7560Tx()
	*penTx2.Sender = testutils.ValidAddress2
	*penTx2.DeployerData = testutils.DummyDeployerData
	_ = mem.AddTx(penTx2)

	penTx3 := testutils.MockValidInitRip7560Tx()
	*penTx3.Sender = testutils.ValidAddress3
	*penTx3.DeployerData = testutils.DummyDeployerData
	_ = mem.AddTx(penTx3)

	ctx, err := NewTxHandlerContext(
		tx,
		testutils.ChainID,
		mem,
	)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	expectedPenTxs := []*transaction.TransactionArgs{penTx3, penTx2, penTx1}
	penTxs := ctx.GetPendingSenderTxs()
	if len(penTxs) != len(expectedPenTxs) {
		t.Fatalf("got length %d, want %d", len(penTxs), len(expectedPenTxs))
	}

	for i, penTx := range penTxs {
		if !testutils.IsTxsEqual(penTx, expectedPenTxs[i]) {
			t.Fatalf("ops not equal: %s", testutils.GetTxsDiff(penTx, expectedPenTxs[i]))
		}
	}
}

func TestGetPendingPaymasterTxs(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := mempool.New(db)
	tx := testutils.MockValidInitRip7560Tx()
	*tx.DeployerData = []byte{}
	*tx.PaymasterData = []byte{}

	penTx1 := testutils.MockValidInitRip7560Tx()
	*penTx1.Sender = testutils.ValidAddress1
	*penTx1.PaymasterData = testutils.DummyPaymasterData
	_ = mem.AddTx(penTx1)

	penTx2 := testutils.MockValidInitRip7560Tx()
	*penTx2.Sender = testutils.ValidAddress2
	*penTx2.PaymasterData = testutils.DummyPaymasterData
	_ = mem.AddTx(penTx2)

	penTx3 := testutils.MockValidInitRip7560Tx()
	*penTx3.Sender = testutils.ValidAddress3
	*penTx3.PaymasterData = testutils.DummyPaymasterData
	_ = mem.AddTx(penTx3)

	ctx, err := NewTxHandlerContext(
		tx,
		testutils.ChainID,
		mem,
	)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	expectedPenTxs := []*transaction.TransactionArgs{penTx3, penTx2, penTx1}
	penTxs := ctx.GetPendingSenderTxs()
	if len(penTxs) != len(expectedPenTxs) {
		t.Fatalf("got length %d, want %d", len(penTxs), len(expectedPenTxs))
	}

	for i, penTx := range penTxs {
		if !testutils.IsTxsEqual(penTx, expectedPenTxs[i]) {
			t.Fatalf("txs not equal: %s", testutils.GetTxsDiff(penTx, expectedPenTxs[i]))
		}
	}
}
