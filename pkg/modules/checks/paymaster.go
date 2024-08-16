package checks

import (
	"errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

// ValidatePaymasterAndData checks the paymasterAndData is either zero bytes or the first 20 bytes contain an
// address that
//
//	currently has nonempty code on chain
func ValidatePaymasterAndData(
	txArgs *transaction.TransactionArgs,
	gc GetCodeFunc,
) error {
	if len(txArgs.GetPaymasterData()) == 0 {
		return nil
	}

	pm := txArgs.GetPaymaster()
	bytecode, err := gc(pm)
	if err != nil {
		return err
	}
	if len(bytecode) == 0 {
		return errors.New("paymaster: code not deployed")
	}

	return nil
}
