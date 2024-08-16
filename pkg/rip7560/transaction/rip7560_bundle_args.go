package transaction

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TODO : remove & update ref to geth package
type GetRip7560BundleArgs struct {
	MinBaseFee    uint64 `json:"minBaseFee"`
	MaxBundleGas  uint64 `json:"maxBundleGas"`
	MaxBundleSize uint64 `json:"maxBundleSize"`
}

// TODO : remove & update ref to geth package
type GetRip7560BundleResult struct {
	Bundle        []TransactionArgs
	ValidForBlock *hexutil.Big
}
