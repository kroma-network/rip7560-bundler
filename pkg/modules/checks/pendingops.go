package checks

import (
	"fmt"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"
)

var (
	minPriceBump                = int64(10)
	ErrReplacementTxUnderpriced = fmt.Errorf(
		"pending txs: replacement tx must increase maxFeePerGas and MaxPriorityFeePerGas by >= %d%%",
		minPriceBump,
	)
)

// calcNewThresholds returns new threshold values where newFee = oldFee  * (100 + minPriceBump) / 100.
func calcNewThresholds(cap *big.Int, tip *big.Int) (newCap *big.Int, newTip *big.Int) {
	a := big.NewInt(100 + minPriceBump)
	aFeeCap := big.NewInt(0).Mul(a, cap)
	aTip := big.NewInt(0).Mul(a, tip)

	b := big.NewInt(100)
	newCap = aFeeCap.Div(aFeeCap, b)
	newTip = aTip.Div(aTip, b)

	return newCap, newTip
}

// ValidatePendingTxs checks the pending Transactions by the same sender and only passes if:
//
//  1. Sender doesn't have another Transactions already present in the pool.
//  2. It replaces an existing Transactions with same nonce and higher fee.
func ValidatePendingTxs(
	tx *transaction.TransactionArgs,
	penTxs []*transaction.TransactionArgs,
) error {
	if len(penTxs) > 0 {
		var oldTx *transaction.TransactionArgs
		for _, penTx := range penTxs {
			if tx.Nonce == penTx.Nonce {
				oldTx = penTx
			}
		}

		if oldTx != nil {
			newMf, newMpf := calcNewThresholds(oldTx.MaxFeePerGas.ToInt(), oldTx.MaxPriorityFeePerGas.ToInt())

			if tx.MaxFeePerGas.ToInt().Cmp(newMf) < 0 || tx.MaxPriorityFeePerGas.ToInt().Cmp(newMpf) < 0 {
				return ErrReplacementTxUnderpriced
			}
		}
	}
	return nil
}
