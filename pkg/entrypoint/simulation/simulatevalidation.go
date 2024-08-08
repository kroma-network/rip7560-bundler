package simulation

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SimulateValidation makes a static call to eth_callRip7560Validation and returns the
// results without any state changes.
func SimulateValidation(
	rpc *rpc.Client,
	op *userop.UserOperation,
) (*reverts.ValidationPhaseResult, error) {
	var res reverts.ValidationPhaseResult
	req := client.CreateUserOperationArgs(op)
	if err := rpc.CallContext(context.Background(), &res, "eth_callRip7560Validation", &req, "latest"); err != nil {
		return nil, err
	}

	return &res, nil
}
