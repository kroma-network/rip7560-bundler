package mempool

import (
	"encoding/json"
	"math/big"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	keyPrefix = dbutils.JoinValues("mempool")
)

func getUniqueKey(sender common.Address, nonce *big.Int) []byte {
	return []byte(
		dbutils.JoinValues(keyPrefix, sender.String(), nonce.String()),
	)
}

func getUserOpFromDBValue(value []byte) (*userop.UserOperation, error) {
	data := make(map[string]any)
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, err
	}

	op, err := userop.New(data)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func loadFromDisk(db *badger.DB, q *userOpQueues) error {
	return db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		prefix := []byte(keyPrefix)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			err := item.Value(func(v []byte) error {
				op, err := getUserOpFromDBValue(v)
				if err != nil {
					return err
				}

				q.AddOp(op)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})
}
