package testutils

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
)

var (
	MockRip7560TxData = `{
		"from":                          "0x0000000000000000000000000000000000007560",
		"to":                            null,
		"value":						 "0x0",
		"gasPrice":                      "0x1",
		"maxFeePerGas":                  "0x2",
		"maxPriorityFeePerGas":          "0x1",
		"data":                          "0xb61d27f6000000000000000000000000e7f1725e7734ce288f8367e1bb143e90bb3f05120000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000246057361d000000000000000000000000000000000000000000000000000000000000004d00000000000000000000000000000000000000000000000000000000",
		"sender":                        "0xAEd18Be2bc078345cd0c2ae5f4DA16230E5ac0d4",
		"paymaster":                     "0x7560000000000000000000000000000000000000",
		"paymasterData":                 "0x",
		"builderFee":                    "0x186A0",
		"gas":                           "0xAAE60",
		"verificationGasLimit":          "0xAAE60",
		"paymasterVerificationGasLimit": "0xAAE60",
		"paymasterPostOpGasLimit":       "0xAAE60",
		"subType":                       "0x1",
		"signature":                     "0x00",
		"bigNonce":  					 "0x10000000000000000"
	}`
	MockByteCode = common.Hex2Bytes("6080604052")
)

// Returns a valid initial RIP-7560 transaction
func MockValidInitRip7560Tx() *transaction.TransactionArgs {
	var txArgs = new(transaction.TransactionArgs)
	err := json.Unmarshal([]byte(MockRip7560TxData), txArgs)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}
	return txArgs
}

func IsTxsEqual(tx1 *transaction.TransactionArgs, tx2 *transaction.TransactionArgs) bool {
	return cmp.Equal(
		tx1,
		tx2,
		cmp.Comparer(func(a *big.Int, b *big.Int) bool {
			if a == nil || b == nil {
				return a == nil && b == nil
			}
			return a.Cmp(b) == 0
		}),
		cmpopts.IgnoreUnexported(transaction.TransactionArgs{}, hexutil.Big{}),
	)
}

func GetTxsDiff(tx1 *transaction.TransactionArgs, tx2 *transaction.TransactionArgs) string {
	return cmp.Diff(
		tx1,
		tx2,
		cmp.Comparer(func(a *big.Int, b *big.Int) bool {
			if a == nil || b == nil {
				return a == nil && b == nil
			}
			return a.Cmp(b) == 0
		}),
		cmpopts.IgnoreUnexported(transaction.TransactionArgs{}, hexutil.Big{}),
	)
}
