package simulation

import (
	"context"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

// SimulateValidation makes a static call to eth_callRip7560Validation and returns the
// results without any state changes.
func SimulateValidation(
	rpc *rpc.Client,
	tx *transaction.TransactionArgs,
) (*core.ValidationPhaseResult, error) {
	var res core.ValidationPhaseResult
	if err := rpc.CallContext(context.Background(), &res, "eth_callRip7560Validation", tx, "latest"); err != nil {
		return nil, err
	}

	return &res, nil
}
