package modules

// ComposeBatchHandlerFunc combines many BatchHandlers into one.
func ComposeBatchHandlerFunc(fns ...BatchHandlerFunc) BatchHandlerFunc {
	return func(ctx *BatchHandlerCtx) error {
		for _, fn := range fns {
			err := fn(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// ComposeTransactionHandlerFunc combines many Rip7560TxHandlers into one.
func ComposeTransactionHandlerFunc(fns ...Rip7560TxHandlerFunc) Rip7560TxHandlerFunc {
	return func(ctx *TxHandlerCtx) error {
		for _, fn := range fns {
			err := fn(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
