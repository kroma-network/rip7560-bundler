package checks

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestMFLessThanBF calls checks.ValidateFeePerGas with a MaxFeePerGas < base fee. Expect error.
func TestMFLessThanBF(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	gbf := testutils.GetMockBaseFeeFunc(common.Big2)
	*tx.MaxFeePerGas = hexutil.Big(*common.Big1)
	*tx.MaxPriorityFeePerGas = hexutil.Big(*common.Big0)
	err := ValidateFeePerGas(tx, gbf)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestMFEqualBF calls checks.ValidateFeePerGas with a MaxFeePerGas == base fee. Expect nil.
func TestMFEqualBF(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	*tx.MaxFeePerGas = hexutil.Big(*common.Big1)
	*tx.MaxPriorityFeePerGas = hexutil.Big(*common.Big0)
	err := ValidateFeePerGas(tx, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMFMoreThanBF calls checks.ValidateFeePerGas with a MaxFeePerGas > base fee. Expect nil.
func TestMFMoreThanBF(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	*tx.MaxFeePerGas = hexutil.Big(*common.Big2)
	*tx.MaxPriorityFeePerGas = hexutil.Big(*common.Big0)
	err := ValidateFeePerGas(tx, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMPFMoreThanMF calls checks.ValidateFeePerGas with a MaxPriorityFeePerGas > MaxFeePerGas. Expect error.
func TestMPFMoreThanMF(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	*tx.MaxFeePerGas = hexutil.Big(*common.Big2)
	*tx.MaxPriorityFeePerGas = hexutil.Big(*common.Big3)
	err := ValidateFeePerGas(tx, gbf)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestMPFEqualMF calls checks.ValidateFeePerGas with a MaxPriorityFeePerGas == MaxFeePerGas. Expect nil.
func TestMPFEqualMF(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	*tx.MaxFeePerGas = hexutil.Big(*common.Big2)
	*tx.MaxPriorityFeePerGas = hexutil.Big(*common.Big2)
	err := ValidateFeePerGas(tx, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMPFLessThanMF calls checks.ValidateFeePerGas with a MaxPriorityFeePerGas < MaxFeePerGas. Expect nil.
func TestMPFLessThanMF(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	*tx.MaxFeePerGas = hexutil.Big(*common.Big2)
	*tx.MaxPriorityFeePerGas = hexutil.Big(*common.Big1)
	err := ValidateFeePerGas(tx, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}
