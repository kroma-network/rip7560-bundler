package types

import "github.com/ethereum/go-ethereum/common/hexutil"

type GetRip7560BundleArgs struct {
	MinBaseFee    uint64
	MaxBundleGas  uint64
	MaxBundleSize uint64
}

type GetRip7560BundleResult struct {
	Bundle        []TransactionArgs
	ValidForBlock *hexutil.Big
}
