package bundler

import (
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

func adjustBatchSize(max int, batch []*transaction.TransactionArgs) []*transaction.TransactionArgs {
	if len(batch) > max && max > 0 {
		return batch[:max]
	}
	return batch
}
