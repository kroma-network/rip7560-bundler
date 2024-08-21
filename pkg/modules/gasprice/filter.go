package gasprice

import (
	"github.com/stackup-wallet/stackup-bundler/internal/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// FilterUnderpriced returns a BatchHandlerFunc that will filter out all the Rip7560Txs that are below either the
// dynamic or legacy GasPrice set in the context.
func FilterUnderpriced() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		var b []*transaction.TransactionArgs
		for _, txArgs := range ctx.Batch {
			if ctx.BaseFee != nil && ctx.BaseFee.Cmp(common.Big0) != 0 && ctx.Tip != nil {
				gp := big.NewInt(0).Add(ctx.BaseFee, ctx.Tip)
				dgp := txArgs.GetDynamicGasPrice(ctx.BaseFee)
				if dgp.Cmp(gp) >= 0 {
					b = append(b, txArgs)
				}
			} else if ctx.GasPrice != nil {
				if utils.CompareHexBigWithBig(txArgs.MaxFeePerGas, ctx.GasPrice) >= 0 {
					b = append(b, txArgs)
				}
			}
		}

		ctx.Batch = b
		return nil
	}
}
