package checks

import (
	"fmt"
	"github.com/stackup-wallet/stackup-bundler/internal/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
)

// ValidateFeePerGas checks the maxFeePerGas is sufficiently high to be included with the current
// block.basefee. Alternatively, if basefee is not supported, then check that maxPriorityFeePerGas is equal to
// maxFeePerGas as a fallback.
func ValidateFeePerGas(txArgs *transaction.TransactionArgs, gbf gasprice.GetBaseFeeFunc) error {
	bf, err := gbf()
	if err != nil {
		return err
	}

	if bf == nil {
		if utils.CompareHexBigWithHexBig(txArgs.MaxPriorityFeePerGas, txArgs.MaxFeePerGas) != 0 {
			return fmt.Errorf("legacy fee mode: maxPriorityFeePerGas must equal maxFeePerGas")
		}

		return nil
	}

	if utils.CompareHexBigWithHexBig(txArgs.MaxPriorityFeePerGas, txArgs.MaxFeePerGas) == 1 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than maxPriorityFeePerGas")
	}

	if utils.CompareHexBigWithBig(txArgs.MaxFeePerGas, bf) < 0 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than current block.basefee(%s)", bf.String())
	}

	return nil
}
