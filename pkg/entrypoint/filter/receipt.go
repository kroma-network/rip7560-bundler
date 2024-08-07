package filter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type parsedTransaction struct {
	BlockHash         common.Hash    `json:"blockHash"`
	BlockNumber       string         `json:"blockNumber"`
	From              common.Address `json:"from"`
	CumulativeGasUsed string         `json:"cumulativeGasUsed"`
	GasUsed           string         `json:"gasUsed"`
	Logs              []*types.Log   `json:"logs"`
	LogsBloom         types.Bloom    `json:"logsBloom"`
	TransactionHash   common.Hash    `json:"transactionHash"`
	TransactionIndex  string         `json:"transactionIndex"`
	EffectiveGasPrice string         `json:"effectiveGasPrice"`
}

type UserOperationReceipt struct {
	UserOpHash    common.Hash        `json:"userOpHash"`
	Sender        common.Address     `json:"sender"`
	Paymaster     common.Address     `json:"paymaster"`
	Nonce         string             `json:"nonce"`
	Success       bool               `json:"success"`
	ActualGasCost string             `json:"actualGasCost"`
	ActualGasUsed string             `json:"actualGasUsed"`
	From          common.Address     `json:"from"`
	Receipt       *parsedTransaction `json:"receipt"`
	Logs          []*types.Log       `json:"logs"`
}

func GetUserOperationReceipt(
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
