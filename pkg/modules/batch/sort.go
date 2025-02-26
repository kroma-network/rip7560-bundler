package batch

import (
	"sort"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// SortByNonce returns a BatchHandlerFunc that ensures txs with same sender is ordered by ascending nonce
// regardless of gas price.
func SortByNonce() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		sort.SliceStable(ctx.Batch, func(i, j int) bool {
			return ctx.Batch[i].Sender == ctx.Batch[j].Sender &&
				*ctx.Batch[i].Nonce < *ctx.Batch[j].Nonce
		})

		return nil
	}
}
