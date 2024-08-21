package checks

import (
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

func TestNoPendingTxs(t *testing.T) {
	var penTxs []*transaction.TransactionArgs
	tx := testutils.MockValidInitRip7560Tx()
	err := ValidatePendingTxs(
		tx,
		penTxs,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestPendingTxsWithNewTx(t *testing.T) {
	penTx := testutils.MockValidInitRip7560Tx()
	penTxs := []*transaction.TransactionArgs{penTx}
	tx := testutils.MockValidInitRip7560Tx()
	tx.Nonce = (*hexutil.Uint64)(&testutils.DummyNonce1)
	err := ValidatePendingTxs(
		tx,
		penTxs,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestPendingTxsWithNoGasFeeReplacement(t *testing.T) {
	penTx := testutils.MockValidInitRip7560Tx()
	penTxs := []*transaction.TransactionArgs{penTx}
	tx := testutils.MockValidInitRip7560Tx()
	err := ValidatePendingTxs(
		tx,
		penTxs,
	)

	if !errors.Is(err, ErrReplacementTxUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementTxUnderpriced", err)
	}
}

func TestPendingTxsWithOnlyMaxFeeReplacement(t *testing.T) {
	penTx := testutils.MockValidInitRip7560Tx()
	penTxs := []*transaction.TransactionArgs{penTx}
	tx := testutils.MockValidInitRip7560Tx()
	maxFeePerGas, _ := calcNewThresholds((*big.Int)(tx.MaxFeePerGas), (*big.Int)(tx.MaxPriorityFeePerGas))
	tx.MaxFeePerGas = (*hexutil.Big)(maxFeePerGas)
	err := ValidatePendingTxs(
		tx,
		penTxs,
	)

	if !errors.Is(err, ErrReplacementTxUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementTxUnderpriced", err)
	}
}

func TestPendingTxsWithOnlyMaxPriorityFeeReplacement(t *testing.T) {
	penTx := testutils.MockValidInitRip7560Tx()
	penTxs := []*transaction.TransactionArgs{penTx}
	tx := testutils.MockValidInitRip7560Tx()
	maxFeePerGas, _ := calcNewThresholds((*big.Int)(tx.MaxFeePerGas), (*big.Int)(tx.MaxPriorityFeePerGas))
	tx.MaxFeePerGas = (*hexutil.Big)(maxFeePerGas)
	err := ValidatePendingTxs(
		tx,
		penTxs,
	)

	if !errors.Is(err, ErrReplacementTxUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementTxUnderpriced", err)
	}
}

func TestPendingTxsWithOkGasFeeReplacement(t *testing.T) {
	penTx := testutils.MockValidInitRip7560Tx()
	penTxs := []*transaction.TransactionArgs{penTx}
	tx := testutils.MockValidInitRip7560Tx()
	maxFeePerGas, _ := calcNewThresholds((*big.Int)(tx.MaxFeePerGas), (*big.Int)(tx.MaxPriorityFeePerGas))
	tx.MaxFeePerGas = (*hexutil.Big)(maxFeePerGas)
	err := ValidatePendingTxs(
		tx,
		penTxs,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}
