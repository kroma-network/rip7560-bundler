package filter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetRip7560TransactionReceipt(
	eth *ethclient.Client,
	txHash string,
	blkRange uint64) (*types.Receipt, error) {

	header, err := eth.HeaderByNumber(context.Background(), nil)
	if err != nil {

		return nil, fmt.Errorf("failed to retrieve latest block header: %v", err)
	}

	latestBlock := header.Number
	startBlock := new(big.Int).Sub(latestBlock, new(big.Int).SetUint64(blkRange))

	for blockNumber := new(big.Int).Set(startBlock); blockNumber.Cmp(latestBlock) <= 0; blockNumber.Add(blockNumber, big.NewInt(1)) {
		block, err := eth.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve block %v: %v", blockNumber, err)
		}

		for _, tx := range block.Transactions() {
			if tx.Hash().String() == txHash {
				receipt, err := eth.TransactionReceipt(context.Background(), tx.Hash())
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve transaction receipt: %v", err)
				}

				return receipt, nil
			}
		}
	}
	return nil, fmt.Errorf("transaction not found in the specified block range")
}
