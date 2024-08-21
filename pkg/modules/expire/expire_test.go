package expire

import (
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// TestDropExpired calls (*ExpireHandler).DropExpired and verifies that it marks old Rip-7560 transactions for
// pending removal.
func TestDropExpired(t *testing.T) {
	exp := New(time.Second * 30)
	tx1 := testutils.MockValidInitRip7560Tx()
	tx2 := testutils.MockValidInitRip7560Tx()
	*tx2.Data = common.Hex2Bytes("0xdead")
	exp.seenAt = map[common.Hash]time.Time{
		tx1.ToTransaction().Hash(): time.Now().Add(time.Second * -45),
		tx2.ToTransaction().Hash(): time.Now().Add(time.Second * -15),
	}

	ctx := modules.NewBatchHandlerContext(
		[]*transaction.TransactionArgs{tx1, tx2},
		testutils.ChainID,
		nil,
		nil,
		nil,
	)
	if err := exp.DropExpired()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 1 {
		t.Fatalf("got batch length %d, want 1", len(ctx.Batch))
	} else if len(ctx.PendingRemoval) != 1 {
		t.Fatalf("got pending removal length %d, want 1", len(ctx.Batch))
	} else if !testutils.IsTxsEqual(ctx.Batch[0], tx2) {
		t.Fatal("incorrect batch: Dropped legit tx")
	} else if !testutils.IsTxsEqual(ctx.PendingRemoval[0].Tx, tx1) {
		t.Fatal("incorrect pending removal: Didn't drop bad tx")
	}

}
