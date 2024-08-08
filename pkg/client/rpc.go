package client

import (
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	types2 "github.com/stackup-wallet/stackup-bundler/pkg/types"

	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
)

// Named UserOperation type for jsonrpc package.
type userOperation map[string]any

// Named StateOverride type for jsonrpc package.
type optional_stateOverride map[string]any

// RpcAdapter is an adapter for routing JSON-RPC method calls to the correct client functions.
type RpcAdapter struct {
	client  *Client
	bundler *bundler.Bundler
	debug   *Debug
}

// NewRpcAdapter initializes a new RpcAdapter which can be used with a JSON-RPC server.
func NewRpcAdapter(client *Client, bundler *bundler.Bundler, debug *Debug) *RpcAdapter {
	return &RpcAdapter{client, bundler, debug}
}

// Eth_sendUserOperation routes method calls to *Client.SendUserOperation.
func (r *RpcAdapter) Eth_sendUserOperation(op userOperation) (string, error) {
	return r.client.SendUserOperation(op)
}

// Eth_estimateUserOperationGas routes method calls to *Client.EstimateUserOperationGas.
func (r *RpcAdapter) Eth_estimateUserOperationGas(
	op userOperation,
	os optional_stateOverride,
) (*gas.GasEstimates, error) {
	return r.client.EstimateUserOperationGas(op, os)
}

// Eth_getUserOperationReceipt routes method calls to *Client.GetUserOperationReceipt.
func (r *RpcAdapter) Eth_getUserOperationReceipt(
	op userOperation,
) (*types.Receipt, error) {
	return r.client.GetUserOperationReceipt(op)
}

// Eth_supportedEntryPoints routes method calls to *Client.SupportedEntryPoints.
func (r *RpcAdapter) Eth_supportedEntryPoints() ([]string, error) {
	return r.client.SupportedEntryPoints()
}

// Eth_chainId routes method calls to *Client.ChainID.
func (r *RpcAdapter) Eth_chainId() (string, error) {
	return r.client.ChainID()
}

func (r *RpcAdapter) Aa_getRip7560Bundle(args types2.GetRip7560BundleArgs) (*types2.GetRip7560BundleResult, error) {
	return r.bundler.GetRip7560Bundle(args)
}

// Debug_bundler_clearState routes method calls to *Debug.ClearState.
func (r *RpcAdapter) Debug_bundler_clearState() (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.ClearState()
}

// Debug_bundler_dumpMempool routes method calls to *Debug.DumpMempool.
func (r *RpcAdapter) Debug_bundler_dumpMempool(ep string) ([]map[string]any, error) {
	if r.debug == nil {
		return []map[string]any{}, errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.DumpMempool(ep)
}

// Debug_bundler_sendBundleNow routes method calls to *Debug.SendBundleNow.
//func (r *RpcAdapter) Debug_bundler_sendBundleNow() (string, error) {
//	if r.debug == nil {
//		return "", errors.New("rpc: debug mode is not enabled")
//	}
//
//	return r.debug.SendBundleNow()
//}

// Debug_bundler_setReputation routes method calls to *Debug.SetReputation.
func (r *RpcAdapter) Debug_bundler_setReputation(entries []any, ep string) (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.SetReputation(entries, ep)
}

// Debug_bundler_dumpReputation routes method calls to *Debug.DumpReputation.
func (r *RpcAdapter) Debug_bundler_dumpReputation(ep string) ([]map[string]any, error) {
	if r.debug == nil {
		return []map[string]any{}, errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.DumpReputation(ep)
}
