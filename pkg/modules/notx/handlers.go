// Package noop implements basic no-operation modules which are used by default for both Client and Bundler.
package notx

import "github.com/stackup-wallet/stackup-bundler/pkg/modules"

// BatchHandler takes a BatchHandlerCtx and returns nil error.
func BatchHandler(ctx *modules.BatchHandlerCtx) error {
	return nil
}

// Rip7560TxHandler takes a TxHandlerCtx and returns nil error.
func Rip7560TxHandler(ctx *modules.TxHandlerCtx) error {
	return nil
}
