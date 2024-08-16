package gas

import "math/big"

// GasEstimates provides estimate values for all gas fields in a UserOperation.
type GasEstimates struct {
	VerificationGasLimit *big.Int `json:"verificationGasLimit"`
	CallGasLimit         *big.Int `json:"callGasLimit"`
}
