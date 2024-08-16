package utils

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

func CompareHexBigWithBig(a *hexutil.Big, b *big.Int) int {
	if a == nil && b == nil {
		return 0
	} else if a == nil {
		return -1
	} else if b == nil {
		return 1
	}

	return (*big.Int)(a).Cmp(b)
}
func CompareHexBigWithHexBig(a, b *hexutil.Big) int {
	if a == nil && b == nil {
		return 0
	} else if a == nil {
		return -1
	} else if b == nil {
		return 1
	}

	return (*big.Int)(a).Cmp((*big.Int)(b))
}
