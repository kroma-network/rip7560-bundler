package client

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

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

// Eth_sendRip7560Transaction routes method calls to *Client.SendRip7560Transaction.
func (r *RpcAdapter) Eth_sendRip7560Transaction(input map[string]interface{}) (string, error) {
	txArgs, err := transaction.New(input)
	if err != nil {
		return "", err
	}
	return r.client.SendRip7560Transaction(txArgs)
}

// Eth_estimateRip7560TransactionGas routes method calls to *Client.EstimateRip7560TransactionGas.
func (r *RpcAdapter) Eth_estimateRip7560TransactionGas(
	input map[string]interface{},
	os optional_stateOverride,
) (*gas.GasEstimates, error) {
	txArgs, err := transaction.New(input)
	if err != nil {
		return nil, err
	}
	return r.client.EstimateRip7560TransactionGas(txArgs, os)
}

// Eth_getRip7560TransactionReceipt routes method calls to *Client.GetRip7560TransactionReceipt.
func (r *RpcAdapter) Eth_getRip7560TransactionReceipt(
	input map[string]interface{},
) (*types.Receipt, error) {
	txArgs, err := transaction.New(input)
	if err != nil {
		return nil, err
	}
	return r.client.GetRip7560TransactionReceipt(txArgs)
}

// Eth_chainId routes method calls to *Client.ChainID.
func (r *RpcAdapter) Eth_chainId() (string, error) {
	return r.client.ChainID()
}

func (r *RpcAdapter) Aa_getRip7560Bundle(input map[string]interface{}) (*transaction.GetRip7560BundleResult, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var args transaction.GetRip7560BundleArgs
	if err := json.Unmarshal(jsonData, &args); err != nil {
		return nil, err
	}

	ret, err := r.bundler.GetRip7560Bundle(args)

	return ret, err
}

// Debug_bundler_clearState routes method calls to *Debug.ClearState.
func (r *RpcAdapter) Debug_bundler_clearState() (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.ClearState()
}

// Debug_bundler_dumpMempool routes method calls to *Debug.DumpMempool.
func (r *RpcAdapter) Debug_bundler_dumpMempool() ([]*transaction.TransactionArgs, error) {
	if r.debug == nil {
		return []*transaction.TransactionArgs{}, errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.DumpMempool()
}

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
