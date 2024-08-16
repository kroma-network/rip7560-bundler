// Package bundler provides the mediator for processing outgoing UserOperation batches to the EntryPoint.
package bundler

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"math"
	"math/big"
	"time"

	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// Bundler controls the end to end process of creating a batch of RIP-7560 transactions from the mempool and serveing
// it to the sequencer.
type Bundler struct {
	mempool      *mempool.Mempool
	chainID      *big.Int
	batchHandler modules.BatchHandlerFunc
	logger       logr.Logger
	meter        metric.Meter
	gbf          gasprice.GetBaseFeeFunc
	ggt          gasprice.GetGasTipFunc
	// TODO : is this needed?
	ggp gasprice.GetLegacyGasPriceFunc
}

// New initializes a new RIP-7560 bundler which can be extended with modules for validating batches and
// excluding RIP-7560 transaction that should not be sent to the sequencer and/or dropped from the mempool.
func New(mempool *mempool.Mempool, chainID *big.Int) *Bundler {
	return &Bundler{
		mempool:      mempool,
		chainID:      chainID,
		batchHandler: noop.BatchHandler,
		logger:       logger.NewZeroLogr().WithName("bundler"),
		meter:        otel.GetMeterProvider().Meter("bundler"),
		gbf:          gasprice.NoopGetBaseFeeFunc(),
		ggt:          gasprice.NoopGetGasTipFunc(),
		ggp:          gasprice.NoopGetLegacyGasPriceFunc(),
	}
}

// SetGetBaseFeeFunc defines the function used to retrieve an estimate for basefee during each bundler run.
func (i *Bundler) SetGetBaseFeeFunc(gbf gasprice.GetBaseFeeFunc) {
	i.gbf = gbf
}

// SetGetGasTipFunc defines the function used to retrieve an estimate for gas tip during each bundler run.
func (i *Bundler) SetGetGasTipFunc(ggt gasprice.GetGasTipFunc) {
	i.ggt = ggt
}

// SetGetLegacyGasPriceFunc defines the function used to retrieve an estimate for gas price during each
// bundler run.
func (i *Bundler) SetGetLegacyGasPriceFunc(ggp gasprice.GetLegacyGasPriceFunc) {
	i.ggp = ggp
}

// UseLogger defines the logger object used by the Bundler instance based on the go-logr/logr interface.
func (i *Bundler) UseLogger(logger logr.Logger) {
	i.logger = logger.WithName("bundler")
}

// UserMeter defines an opentelemetry meter object used by the Bundler instance to capture metrics during each
// run.
func (i *Bundler) UserMeter(meter metric.Meter) error {
	i.meter = meter
	_, err := i.meter.Int64ObservableGauge(
		"bundler_mempool_size",
		metric.WithInt64Callback(func(ctx context.Context, io metric.Int64Observer) error {
			size := 0
			batch, err := i.mempool.Dump()
			if err != nil {
				return err
			}
			size += len(batch)
			io.Observe(int64(size))
			return nil
		}),
	)
	return err
}

// UseModules defines the BatchHandlers to process batches after it has gone through the standard checks.
func (i *Bundler) UseModules(handlers ...modules.BatchHandlerFunc) {
	i.batchHandler = modules.ComposeBatchHandlerFunc(handlers...)
}

func (i *Bundler) GetRip7560Bundle(args transaction.GetRip7560BundleArgs) (*transaction.GetRip7560BundleResult, error) {
	// Init logger
	start := time.Now()
	l := i.logger.
		WithName("run").
		WithValues("chain_id", i.chainID.String())

	result := &transaction.GetRip7560BundleResult{
		Bundle: make([]transaction.TransactionArgs, 0),
		// TODO : is this needed?
		ValidForBlock: (*hexutil.Big)(big.NewInt(math.MaxInt64)),
	}

	// Get all pending RIP-7560 transactions from the mempool. This will be in FIFO order. Downstream modules should sort it
	// based on more specific strategies.
	batch, err := i.mempool.Dump()
	if err != nil {
		l.Error(err, "bundler run error")
		return result, err
	}
	if len(batch) == 0 {
		return result, nil
	}
	batch = adjustBatchSize(int(args.MaxBundleSize), batch)

	// Get current block basefee
	bf, err := i.gbf()
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Get suggested gas tip
	var gt *big.Int
	if bf != nil {
		gt, err = i.ggt()
		if err != nil {
			l.Error(err, "bundler run error")
			return nil, err
		}
	}

	// Get suggested gas price (for networks that don't support EIP-1559)
	gp, err := i.ggp()
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Create context and execute modules.
	ctx := modules.NewBatchHandlerContext(batch, i.chainID, bf, gt, gp)
	if err := i.batchHandler(ctx); err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Remove RIP-7560 transactions that remain in the context from mempool.
	rmTxs := append([]*transaction.TransactionArgs{}, ctx.Batch...)
	dh := []string{}
	dr := []string{}
	for _, item := range ctx.PendingRemoval {
		rmTxs = append(rmTxs, item.Tx)
		dh = append(dh, item.Tx.ToTransaction().Hash().String())
		dr = append(dr, item.Reason)
	}
	if err := i.mempool.RemoveTxs(rmTxs...); err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Add tx to result && Update logs for the current run.
	bat := []string{}
	for _, txArgs := range ctx.Batch {
		result.Bundle = append(result.Bundle, *txArgs)
		bat = append(bat, txArgs.ToTransaction().Hash().String())
	}
	l = l.WithValues("batch_aatx_hashes", bat)
	l = l.WithValues("dropped_aatx_hashes", dh)
	l = l.WithValues("dropped_aatx_reasons", dr)

	for k, v := range ctx.Data {
		l = l.WithValues(k, v)
	}
	l = l.WithValues("duration", time.Since(start))
	l.Info("bundler run ok")
	return result, nil
}
