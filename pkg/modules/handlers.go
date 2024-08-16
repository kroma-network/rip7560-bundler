// Package modules provides standard interfaces for extending the Client and Bundler packages with
// middleware.
package modules

// BatchHandlerFunc is an interface to support modular processing of UserOperation batches by the Bundler.
type BatchHandlerFunc func(ctx *BatchHandlerCtx) error

// Rip7560TxHandlerFunc is an interface to support modular processing of single Rip7560Tx by the Client.
type Rip7560TxHandlerFunc func(ctx *TxHandlerCtx) error
