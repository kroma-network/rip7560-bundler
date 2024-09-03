package mempool

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

var (
	keyPrefix = dbutils.JoinValues("mempool")
)

func getUniqueKey(sender common.Address, nonce *hexutil.Uint64, bigNonce *hexutil.Big) []byte {
	if bigNonce == nil {
		bigNonce = new(hexutil.Big)
		bigNonce.ToInt().SetInt64(-1)
	}
	return []byte(
		dbutils.JoinValues(keyPrefix, sender.String(), nonce.String(), bigNonce.String()),
	)
}

func getTransactionsFromDBValue(serializedTx []byte) (*transaction.TransactionArgs, error) {
	var decodedTx *transaction.TransactionArgs
	err := rlp.DecodeBytes(serializedTx, &decodedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP decode transactions: %v", err)
	}
	return decodedTx, nil
}

func loadFromDisk(db *badger.DB, q *rip7560TxQueues) error {
	return db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		prefix := []byte(keyPrefix)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			err := item.Value(func(v []byte) error {
				tx, err := getTransactionsFromDBValue(v)
				if err != nil {
					return err
				}

				q.AddTx(tx)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})
}
