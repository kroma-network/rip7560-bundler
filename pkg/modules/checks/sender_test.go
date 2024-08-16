package checks

import (
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestSenderExistAndInitCodeDNE calls checks.ValidateSender where sender contract exist and initCode does
// not. Expect nil.
func TestSenderExistAndInitCodeDNE(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	*tx.DeployerData = []byte{}
	err := ValidateSender(tx, testutils.MockGetCode)

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

// TestSenderAndInitCodeExist calls checks.ValidateSender where sender contract and initCode exist.
// Expect error.
func TestSenderAndInitCodeExist(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	err := ValidateSender(tx, testutils.MockGetCode)

	if err == nil {
		t.Fatalf(`got nil, want err`)
	}
}

// TestSenderDNEAndInitCodeExist calls checks.ValidateSender where sender contract does not exist and
// initCode does. Expect nil.
func TestSenderDNEAndInitCodeExist(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	err := ValidateSender(tx, testutils.MockGetCodeZero)

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

// TestSenderAndInitCodeDNE calls checks.ValidateSender where sender contract and initCode does not exist.
// Expect error.
func TestSenderAndInitCodeDNE(t *testing.T) {
	tx := testutils.MockValidInitRip7560Tx()
	*tx.DeployerData = []byte{}
	err := ValidateSender(tx, testutils.MockGetCodeZero)

	if err == nil {
		t.Fatalf(`got nil, want err`)
	}
}
