package simulation

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SimulateValidation makes a static call to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func SimulateValidation(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
) (*reverts.ValidationResultRevert, error) {
	_, err := entrypoint.NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return nil, err
	}

	//var res []interface{}
	//rawCaller := &entrypoint.EntrypointRaw{Contract: ep}
	//err = rawCaller.Call(nil, &res, "simulateValidation", entrypoint.UserOperation(*op))
	//if err == nil {
	//	return nil, stdError.New("unexpected result from simulateValidation")
	//}

	sim, simErr := reverts.NewValidationResult(err)
	if simErr != nil {
		fo, foErr := reverts.NewFailedOp(err)
		if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}

	return sim, nil
}

// SimulateRIP7560Validation makes a static call to eth_callRip7560Validation and returns the
// results without any state changes.
func SimulateRIP7560Validation(
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
