// Package mempool provides a local representation of all the Rip-7560 transactions that are known to the bundler
// which have passed all Client checks and pending action by the Bundler.
package mempool

import (
	"bytes"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

// Mempool provides read and write access to a pool of pending AA Transactions which have passed all Client
// checks.
type Mempool struct {
	db    *badger.DB
	queue *rip7560TxQueues
}

// New creates an instance of a mempool that uses an embedded DB to persist and load AA Transactions from disk
// incase of a reset.
func New(db *badger.DB) (*Mempool, error) {
	queue := newRip7560TxQueue()
	err := loadFromDisk(db, queue)
	if err != nil {
		return nil, err
	}

	return &Mempool{db, queue}, nil
}

// GetTxs returns all the AA Transactions associated with Sender address.
func (m *Mempool) GetTxs(sender common.Address) ([]*transaction.TransactionArgs, error) {
	if sender == (common.Address{}) {
		return []*transaction.TransactionArgs{}, nil
	}
	txs := m.queue.GetTxs(sender)
	return txs, nil
}

// AddTx adds a AA Transaction to the mempool or replace an existing one with the Sender, and
// Nonce values.
func (m *Mempool) AddTx(tx *transaction.TransactionArgs) error {
	var buf bytes.Buffer
	err := rlp.Encode(&buf, tx)
	if err != nil {
		return fmt.Errorf("failed to RLP encode transaction: %v", err)
	}
	err = m.db.Update(func(txn *badger.Txn) error {
		return txn.Set(getUniqueKey(tx.GetSender(), tx.Nonce, tx.BigNonce), buf.Bytes())
	})
	if err != nil {
		return err
	}

	m.queue.AddTx(tx)
	return nil
}

// RemoveTxs removes a list of AA Transactions from the mempool by Sender, and Nonce values.
func (m *Mempool) RemoveTxs(txs ...*transaction.TransactionArgs) error {
	err := m.db.Update(func(txn *badger.Txn) error {
		for _, tx := range txs {
			err := txn.Delete(getUniqueKey(tx.GetSender(), tx.Nonce, tx.BigNonce))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	m.queue.RemoveTxs(txs...)
	return nil
}

// Dump will return a list of AA Transactions from the mempool by EntryPoint in the order it arrived.
func (m *Mempool) Dump() ([]*transaction.TransactionArgs, error) {
	return m.queue.All(), nil
}

// Clear will clear the entire embedded db and reset it to a clean state.
func (m *Mempool) Clear() error {
	if err := m.db.DropAll(); err != nil {
		return err
	}
	m.queue = newRip7560TxQueue()

	return nil
}
