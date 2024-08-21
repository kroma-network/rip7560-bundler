package gas

import "math/big"

// GasEstimates provides estimate values for all gas fields in a Rip-7560 transactions.
type GasEstimates struct {
	VerificationGasLimit *big.Int `json:"verificationGasLimit"`
	CallGasLimit         *big.Int `json:"callGasLimit"`
}
