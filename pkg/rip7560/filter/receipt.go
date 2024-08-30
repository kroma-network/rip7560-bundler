package filter

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetRip7560TransactionReceipt(
	eth *ethclient.Client,
	txHash string) (*types.Receipt, error) {

	hash := common.HexToHash(txHash)
	receipt, err := eth.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
