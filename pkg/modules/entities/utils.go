package entities

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
)

type addressCounter map[common.Address]int

type status int64

const (
	ok status = iota
	throttled
	banned
)

var (
	emaHours       = 24
	txsCountPrefix = dbutils.JoinValues("entity", "txsCount")
)

func getTxsCountKey(entity common.Address) []byte {
	return []byte(dbutils.JoinValues(txsCountPrefix, entity.String()))
}

func getTxsCountValue(txsSeen int, txsIncluded int) []byte {
	return []byte(
		dbutils.JoinValues(strconv.Itoa(txsSeen), strconv.Itoa(txsIncluded), fmt.Sprint(time.Now().Unix())),
	)
}

func applyExpWeights(txn *badger.Txn, key []byte, value []byte) (txsSeen int, txsIncluded int, err error) {
	counts := dbutils.SplitValues(string(value))
	txsSeen, err = strconv.Atoi(counts[0])
	if err != nil {
		return 0, 0, err
	}
	txsIncluded, err = strconv.Atoi(counts[1])
	if err != nil {
		return 0, 0, err
	}
	lastUpdated, err := strconv.ParseInt(counts[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	dur := time.Since(time.Unix(lastUpdated, 0))
	for i := int(dur.Hours()); i > 0; i-- {
		if txsSeen < 24 && txsIncluded < 24 {
			break
		}

		txsSeen -= txsSeen / emaHours
		txsIncluded -= txsIncluded / emaHours
	}

	e := badger.NewEntry(key, getTxsCountValue(txsSeen, txsIncluded))
	err = txn.SetEntry(e)

	return txsSeen, txsIncluded, err
}

func getTxsCountByEntity(
	txn *badger.Txn,
	entity common.Address,
) (txsSeen int, txsIncluded int, err error) {
	key := getTxsCountKey(entity)
	item, err := txn.Get(key)
	if err != nil && err == badger.ErrKeyNotFound {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var value []byte
	err = item.Value(func(val []byte) error {
		value = append([]byte{}, val...)
		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return applyExpWeights(txn, key, value)
}

func incrementTxsSeenByEntity(txn *badger.Txn, entity common.Address) error {
	txsSeen, txsIncluded, err := getTxsCountByEntity(txn, entity)
	if err != nil {
		return err
	}

	e := badger.NewEntry(getTxsCountKey(entity), getTxsCountValue(txsSeen+1, txsIncluded))
	return txn.SetEntry(e)
}

func incrementTxsIncludedByEntity(txn *badger.Txn, count addressCounter) error {
	for entity, n := range count {
		txsSeen, txsIncluded, err := getTxsCountByEntity(txn, entity)
		if err != nil {
			return err
		}

		e := badger.NewEntry(
			getTxsCountKey(entity),
			getTxsCountValue(txsSeen, txsIncluded+n),
		)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
	}

	return nil
}

func getStatus(txn *badger.Txn, entity common.Address, repConst *ReputationConstants) (status, error) {
	txsSeen, txsIncluded, err := getTxsCountByEntity(txn, entity)
	if err != nil {
		return ok, err
	}
	if txsSeen == 0 {
		return ok, nil
	}

	minExpectedIncluded := txsSeen / repConst.MinInclusionRateDenominator
	if minExpectedIncluded <= txsIncluded+repConst.ThrottlingSlack {
		return ok, nil
	} else if minExpectedIncluded <= txsIncluded+repConst.BanSlack {
		return throttled, nil
	} else {
		return banned, nil
	}
}

func overrideEntity(txn *badger.Txn, entry *ReputationOverride) error {
	return txn.SetEntry(
		badger.NewEntry(getTxsCountKey(entry.Address), getTxsCountValue(entry.TxsSeen, entry.TxsIncluded)),
	)
}
