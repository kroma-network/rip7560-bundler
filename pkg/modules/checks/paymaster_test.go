package checks

import (
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestNilPaymasterAndData calls checks.ValidatePaymasterAndData with no paymaster set. Expects nil.
func TestNilPaymasterAndData(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	*tx.PaymasterData = []byte{}
	err := ValidatePaymasterAndData(tx, testutils.MockGetCodeZero)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

// TestZeroByteCodePaymasterAndData calls checks.ValidatePaymasterAndData with paymaster contract not
// deployed. Expects error.
func TestZeroByteCodePaymasterAndData(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	*tx.PaymasterData = testutils.DummyPaymasterData
	err := ValidatePaymasterAndData(tx, testutils.MockGetCodeZero)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}
