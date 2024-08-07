package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/wangjia184/sortedset"
)

// TODO : renaming userOp to aaTx
type userOpQueues struct {
	all      *sortedset.SortedSet
	entities map[common.Address]*sortedset.SortedSet
}

func (q *userOpQueues) getEntitiesSortedSet(entity common.Address) *sortedset.SortedSet {
	if _, ok := q.entities[entity]; !ok {
		q.entities[entity] = sortedset.New()
	}

	return q.entities[entity]
}

func (q *userOpQueues) AddOp(op *userop.UserOperation) {
	key := string(getUniqueKey(op.Sender, op.Nonce))

	q.all.AddOrUpdate(key, sortedset.SCORE(q.all.GetCount()), op)
	q.getEntitiesSortedSet(op.Sender).AddOrUpdate(key, sortedset.SCORE(op.Nonce.Int64()), op)
	if factory := op.GetFactory(); factory != common.HexToAddress("0x") {
		fss := q.getEntitiesSortedSet(factory)
		fss.AddOrUpdate(key, sortedset.SCORE(fss.GetCount()), op)
	}
	if paymaster := op.GetPaymaster(); paymaster != common.HexToAddress("0x") {
		pss := q.getEntitiesSortedSet(paymaster)
		pss.AddOrUpdate(key, sortedset.SCORE(pss.GetCount()), op)
	}
}

func (q *userOpQueues) GetOps(entity common.Address) []*userop.UserOperation {
	ess := q.getEntitiesSortedSet(entity)
	nodes := ess.GetByRankRange(-1, -ess.GetCount(), false)
	var batch []*userop.UserOperation
	for _, n := range nodes {
		batch = append(batch, n.Value.(*userop.UserOperation))
	}

	return batch
}

func (q *userOpQueues) All() []*userop.UserOperation {
	nodes := q.all.GetByRankRange(1, -1, false)
	var batch []*userop.UserOperation
	for _, n := range nodes {
		batch = append(batch, n.Value.(*userop.UserOperation))
	}

	return batch
}

func (q *userOpQueues) RemoveOps(ops ...*userop.UserOperation) {
	for _, op := range ops {
		key := string(getUniqueKey(op.Sender, op.Nonce))
		q.all.Remove(key)
		q.getEntitiesSortedSet(op.Sender).Remove(key)
		q.getEntitiesSortedSet(op.GetFactory()).Remove(key)
		q.getEntitiesSortedSet(op.GetPaymaster()).Remove(key)
	}
}

func newUserOpQueue() *userOpQueues {
	return &userOpQueues{}
}
