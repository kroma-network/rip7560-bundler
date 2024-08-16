// Package checks implements modules for running an array of standard validations for both the Client and
// Bundler.
package checks

import (
	"math/big"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/entities"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/simulation"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"golang.org/x/sync/errgroup"
)

// Standalone exposes modules to perform basic Client and Bundler checks as specified in EIP-4337. It is
// intended for bundlers that are independent of an Ethereum node and hence relies on a given ethClient to
// query blockchain state.
type Standalone struct {
	db                 *badger.DB
	rpc                *rpc.Client
	eth                *ethclient.Client
	maxVerificationGas *big.Int
	maxBatchGasLimit   *big.Int
	repConst           *entities.ReputationConstants
}

// New returns a Standalone instance with methods that can be used in Client and Bundler modules to perform
// standard checks as specified in EIP-4337.
func New(
	db *badger.DB,
	rpc *rpc.Client,
	maxVerificationGas *big.Int,
	maxBatchGasLimit *big.Int,
	repConst *entities.ReputationConstants,
) *Standalone {
	eth := ethclient.NewClient(rpc)
	return &Standalone{
		db,
		rpc,
		eth,
		maxVerificationGas,
		maxBatchGasLimit,
		repConst,
	}
}

// ValidateTxValues returns a UserOpHandler that runs through some first line sanity checks for new UserOps
// received by the Client. This should be one of the first modules executed by the Client.
func (s *Standalone) ValidateTxValues() modules.Rip7560TxHandlerFunc {
	return func(ctx *modules.TxHandlerCtx) error {
		gc := getCodeWithEthClient(s.eth)

		g := new(errgroup.Group)
		g.Go(func() error { return ValidateSender(ctx.Tx, gc) })
		g.Go(func() error { return ValidatePaymasterAndData(ctx.Tx, gc) })
		g.Go(func() error { return ValidateFeePerGas(ctx.Tx, gasprice.GetBaseFeeWithEthClient(s.eth)) })
		g.Go(func() error { return ValidatePendingTxs(ctx.Tx, ctx.GetPendingSenderTxs()) })

		if err := g.Wait(); err != nil {
			return errors.NewRPCError(errors.INVALID_FIELDS, err.Error(), err.Error())
		}
		return nil
	}
}

// TODO : implement scale-out structure
func (s *Standalone) SimulateTx() modules.Rip7560TxHandlerFunc {
	return func(ctx *modules.TxHandlerCtx) error {
		gc := getCodeWithEthClient(s.eth)
		g := new(errgroup.Group)
		g.Go(func() error {
			sim, err := simulation.SimulateValidation(s.rpc, ctx.Tx)

			if err != nil {
				return errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, err.Error(), err.Error())
			}
			if sim.SenderValidUntil != 0 &&
				uint64(time.Now().Unix()) >= sim.SenderValidUntil-30 {
				return errors.NewRPCError(
					errors.SHORT_DEADLINE,
					"expires too soon",
					nil,
				)
			}
			return nil
		})
		g.Go(func() error {
			out, err := simulation.TraceSimulateValidation(&simulation.TraceInput{
				Rpc:     s.rpc,
				Tx:      ctx.Tx,
				ChainID: ctx.ChainID,
			})
			if err != nil {
				return errors.NewRPCError(errors.BANNED_OPCODE, err.Error(), err.Error())
			}

			ch, err := getCodeHashes(out.TouchedContracts, gc)
			if err != nil {
				return errors.NewRPCError(errors.BANNED_OPCODE, err.Error(), err.Error())
			}
			return saveCodeHashes(s.db, ctx.Tx.ToTransaction().Hash(), ch)
		})

		return g.Wait()
	}
}

// CodeHashes returns a BatchHandler that verifies the code for any interacted contracts has not changed since
// the first simulation.
func (s *Standalone) CodeHashes() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		gc := getCodeWithEthClient(s.eth)

		end := len(ctx.Batch) - 1
		for i := end; i >= 0; i-- {
			aaTxArgs := ctx.Batch[i]
			chs, err := getSavedCodeHashes(s.db, aaTxArgs.ToTransaction().Hash())
			if err != nil {
				return err
			}

			changed, err := hasCodeHashChanges(chs, gc)
			if err != nil {
				return err
			}
			if changed {
				ctx.MarkOpIndexForRemoval(i, "code hash changed")
			}
		}
		return nil
	}
}

// Clean returns a BatchHandler that clears the DB of data that is no longer required. This should be one of
// the last modules executed by the Bundler.
func (s *Standalone) Clean() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		all := append([]*transaction.TransactionArgs{}, ctx.Batch...)
		for _, item := range ctx.PendingRemoval {
			all = append(all, item.Tx)
		}
		var hashes []common.Hash
		for _, aaTxArgs := range all {
			hashes = append(hashes, aaTxArgs.ToTransaction().Hash())
		}

		return removeSavedCodeHashes(s.db, hashes...)
	}
}
