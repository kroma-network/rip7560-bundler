package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"github.com/wangjia184/sortedset"
)

type rip7560TxQueues struct {
	all      *sortedset.SortedSet
	entities map[common.Address]*sortedset.SortedSet
}

func (q *rip7560TxQueues) getEntitiesSortedSet(entity common.Address) *sortedset.SortedSet {
	if _, ok := q.entities[entity]; !ok {
		q.entities[entity] = sortedset.New()
	}

	return q.entities[entity]
}

func (q *rip7560TxQueues) AddTx(tx *transaction.TransactionArgs) {
	rip7560Tx := tx.ToTransaction().Rip7560TransactionData()
	key := string(getUniqueKey(tx.GetSender(), tx.BigNonce))

	q.all.AddOrUpdate(key, sortedset.SCORE(q.all.GetCount()), tx)
	q.getEntitiesSortedSet(*rip7560Tx.Sender).AddOrUpdate(key, sortedset.SCORE(int64(rip7560Tx.Nonce)), tx)
	if deployer := tx.GetDeployer(); deployer != common.HexToAddress("0x") {
		fss := q.getEntitiesSortedSet(deployer)
		fss.AddOrUpdate(key, sortedset.SCORE(fss.GetCount()), tx)
	}
	if paymaster := tx.GetPaymaster(); paymaster != common.HexToAddress("0x") {
		pss := q.getEntitiesSortedSet(paymaster)
		pss.AddOrUpdate(key, sortedset.SCORE(pss.GetCount()), tx)
	}
}

func (q *rip7560TxQueues) GetTxs(entity common.Address) []*transaction.TransactionArgs {
	ess := q.getEntitiesSortedSet(entity)
	nodes := ess.GetByRankRange(-1, -ess.GetCount(), false)
	var batch []*transaction.TransactionArgs
	for _, n := range nodes {
		batch = append(batch, n.Value.(*transaction.TransactionArgs))
	}

	return batch
}

func (q *rip7560TxQueues) All() []*transaction.TransactionArgs {
	nodes := q.all.GetByRankRange(1, -1, false)
	batch := []*transaction.TransactionArgs{}
	for _, n := range nodes {
		batch = append(batch, n.Value.(*transaction.TransactionArgs))
	}

	return batch
}

func (q *rip7560TxQueues) RemoveTxs(txArgsList ...*transaction.TransactionArgs) {

	for _, txArgs := range txArgsList {
		key := string(getUniqueKey(txArgs.GetSender(), txArgs.BigNonce))
		q.all.Remove(key)
		q.getEntitiesSortedSet(txArgs.GetSender()).Remove(key)
		q.getEntitiesSortedSet(txArgs.GetDeployer()).Remove(key)
		q.getEntitiesSortedSet(txArgs.GetPaymaster()).Remove(key)
	}
}

func newRip7560TxQueue() *rip7560TxQueues {
	return &rip7560TxQueues{
		all:      sortedset.New(),
		entities: make(map[common.Address]*sortedset.SortedSet),
	}
}
