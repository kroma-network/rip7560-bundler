package client

import (
	"encoding/json"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/entities"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

// Debug exposes methods used for testing the bundler. These should not be made available in production.
type Debug struct {
	eoa     *signer.EOA
	eth     *ethclient.Client
	mempool *mempool.Mempool
	rep     *entities.Reputation
	bundler *bundler.Bundler
	chainID *big.Int
}

func NewDebug(
	eoa *signer.EOA,
	eth *ethclient.Client,
	mempool *mempool.Mempool,
	rep *entities.Reputation,
	bundler *bundler.Bundler,
	chainID *big.Int,
) *Debug {
	return &Debug{eoa, eth, mempool, rep, bundler, chainID}
}

// ClearState clears the bundler mempool and reputation data of paymasters/accounts/factories/aggregators.
func (d *Debug) ClearState() (string, error) {
	if err := d.mempool.Clear(); err != nil {
		return "", err
	}

	return "ok", nil
}

// DumpMempool dumps the current RIP-7560 transactions mempool in order of arrival.
func (d *Debug) DumpMempool() ([]*transaction.TransactionArgs, error) {
	txs, err := d.mempool.Dump()
	if err != nil {
		return []*transaction.TransactionArgs{}, err
	}

	return txs, nil

	//res := []map[string]any{}
	//for _, tx := range txs {
	//	data, err := tx.MarshalJSON()
	//	if err != nil {
	//		return []map[string]any{}, err
	//	}
	//
	//	item := make(map[string]any)
	//	if err := json.Unmarshal(data, &item); err != nil {
	//		return []map[string]any{}, err
	//	}
	//
	//	res = append(res, item)
	//}
	//
	//return res, nil
}

// SetReputation allows the bundler to set the reputation of given addresses.
func (d *Debug) SetReputation(entries []any, ep string) (string, error) {
	roArr := []*entities.ReputationOverride{}
	for _, entry := range entries {
		b, err := json.Marshal(entry)
		if err != nil {
			return "", err
		}

		ro := &entities.ReputationOverride{}
		if err := json.Unmarshal(b, ro); err != nil {
			return "", err
		}

		roArr = append(roArr, ro)
	}
	if err := d.rep.Override(roArr); err != nil {
		return "", err
	}

	return "ok", nil
}

// DumpReputation returns the reputation data of all known addresses.
// TODO: Implement
func (d *Debug) DumpReputation(ep string) ([]map[string]any, error) {
	return []map[string]any{}, nil
}
