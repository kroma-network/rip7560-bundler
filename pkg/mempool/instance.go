// Package mempool provides a local representation of all the UserOperations that are known to the bundler
// which have passed all Client checks and pending action by the Bundler.
package mempool

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Mempool provides read and write access to a pool of pending UserOperations which have passed all Client
// checks.
type Mempool struct {
	db    *badger.DB
	queue *userOpQueues
}

// New creates an instance of a mempool that uses an embedded DB to persist and load UserOperations from disk
// incase of a reset.
func New(db *badger.DB) (*Mempool, error) {
	queue := newUserOpQueue()
	err := loadFromDisk(db, queue)
	if err != nil {
		return nil, err
	}

	return &Mempool{db, queue}, nil
}

// GetOps returns all the UserOperations associated with an EntryPoint and Sender address.
func (m *Mempool) GetOps(sender common.Address) ([]*userop.UserOperation, error) {
	ops := m.queue.GetOps(sender)
	return ops, nil
}

// AddOp adds a UserOperation to the mempool or replace an existing one with the same EntryPoint, Sender, and
// Nonce values.
func (m *Mempool) AddOp(op *userop.UserOperation) error {
	data, err := op.MarshalJSON()
	if err != nil {
		return err
	}

	err = m.db.Update(func(txn *badger.Txn) error {
		return txn.Set(getUniqueKey(op.Sender, op.Nonce), data)
	})
	if err != nil {
		return err
	}

	m.queue.AddOp(op)
	return nil
}

// RemoveOps removes a list of UserOperations from the mempool by EntryPoint, Sender, and Nonce values.
func (m *Mempool) RemoveOps(ops ...*userop.UserOperation) error {
	err := m.db.Update(func(txn *badger.Txn) error {
		for _, op := range ops {
			err := txn.Delete(getUniqueKey(op.Sender, op.Nonce))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	m.queue.RemoveOps(ops...)
	return nil
}

// Dump will return a list of UserOperations from the mempool by EntryPoint in the order it arrived.
func (m *Mempool) Dump() ([]*userop.UserOperation, error) {
	return m.queue.All(), nil
}

// Clear will clear the entire embedded db and reset it to a clean state.
func (m *Mempool) Clear() error {
	if err := m.db.DropAll(); err != nil {
		return err
	}
	m.queue = newUserOpQueue()

	return nil
}
