package checks

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

// ValidateSender accepts a transaction and a generic function that can retrieve the bytecode of the sender.
// Either the sender is deployed (non-zero length bytecode) or the initCode is not empty (but not both).
func ValidateSender(tx *transaction.TransactionArgs, gc GetCodeFunc) error {
	if tx.GetSender() == (common.Address{}) {
		return errors.New("sender is required")
	}
	bytecode, err := gc(tx.GetSender())
	if err != nil {
		return err
	}

	if len(bytecode) == 0 && tx.DeployerData != nil && len(tx.GetDeployerData()) == 0 {
		return errors.New("sender: not deployed, initCode must be set")
	}
	if len(bytecode) > 0 && tx.DeployerData != nil && len(tx.GetDeployerData()) > 0 {
		return errors.New("sender: already deployed, initCode must be empty")
	}

	return nil
}
