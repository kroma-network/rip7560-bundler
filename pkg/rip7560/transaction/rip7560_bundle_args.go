package transaction

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// geth/rip7560pool
type GetRip7560BundleArgs struct {
	MinBaseFee    uint64 `json:"minBaseFee"`
	MaxBundleGas  uint64 `json:"maxBundleGas"`
	MaxBundleSize uint64 `json:"maxBundleSize"`
}

// geth/rip7560pool
type GetRip7560BundleResult struct {
	Bundle        []TransactionArgs
	ValidForBlock *hexutil.Big
}
