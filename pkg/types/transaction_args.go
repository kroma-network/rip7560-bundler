// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

// TransactionArgs represents the arguments to construct a new transaction
// or a message call.
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`

	// For BlobTxType
	BlobFeeCap *hexutil.Big  `json:"maxFeePerBlobGas"`
	BlobHashes []common.Hash `json:"blobVersionedHashes,omitempty"`

	// For BlobTxType transactions with blob sidecar
	Blobs       []kzg4844.Blob       `json:"blobs"`
	Commitments []kzg4844.Commitment `json:"commitments"`
	Proofs      []kzg4844.Proof      `json:"proofs"`

	// This configures whether blobs are allowed to be passed.
	blobSidecarAllowed bool

	// Introduced by RIP-7560 Transaction
	Sender        *common.Address `json:"sender"`
	Signature     *hexutil.Bytes  `json:"signature"`
	Paymaster     *common.Address `json:"paymaster,omitempty"`
	PaymasterData *hexutil.Bytes  `json:"paymasterData,omitempty"`
	Deployer      *common.Address `json:"deployer,omitempty"`
	DeployerData  *hexutil.Bytes  `json:"deployerData,omitempty"`
	BuilderFee    *hexutil.Big
	ValidationGas *hexutil.Uint64 `json:"verificationGasLimit"`
	PaymasterGas  *hexutil.Uint64 `json:"paymasterVerificationGasLimit"`
	PostOpGas     *hexutil.Uint64 `json:"paymasterPostOpGasLimit"`
	BigNonce      *hexutil.Big    `json:"bigNonce"`
}

// from retrieves the transaction sender address.
func (args *TransactionArgs) from() common.Address {
	if args.From == nil {
		return common.Address{}
	}
	return *args.From
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *TransactionArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

func (args *TransactionArgs) sender() *common.Address {
	if args.Sender != nil {
		return args.Sender
	}
	return nil
}

func (args *TransactionArgs) signature() []byte {
	if args.Signature != nil {
		return *args.Signature
	}
	return nil
}

func (args *TransactionArgs) paymaster() *common.Address {
	if args.Paymaster != nil {
		return args.Paymaster
	}
	return nil
}

func (args *TransactionArgs) paymasterData() []byte {
	if args.PaymasterData != nil {
		return *args.PaymasterData
	}
	return nil
}

func (args *TransactionArgs) paymasterGas() uint64 {
	if args.PaymasterGas != nil {
		return uint64(*args.PaymasterGas)
	}
	return 0
}

func (args *TransactionArgs) validationGas() uint64 {
	if args.ValidationGas != nil {
		return uint64(*args.ValidationGas)
	}
	return 0
}

func (args *TransactionArgs) postOpGas() uint64 {
	if args.PostOpGas != nil {
		return uint64(*args.PostOpGas)
	}
	return 0
}

func (args *TransactionArgs) deployer() *common.Address {
	if args.Deployer != nil {
		return args.Deployer
	}
	return nil
}

func (args *TransactionArgs) deployerData() []byte {
	if args.DeployerData != nil {
		return *args.DeployerData
	}
	return nil
}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) toTransaction() *types.Transaction {
	var data types.TxData
	switch {
	case args.Sender != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		aatx := types.Rip7560AccountAbstractionTx{
			To:         &common.Address{},
			ChainID:    (*big.Int)(args.ChainID),
			Gas:        uint64(*args.Gas),
			GasFeeCap:  (*big.Int)(args.MaxFeePerGas),
			GasTipCap:  (*big.Int)(args.MaxPriorityFeePerGas),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: al,
			// RIP-7560 parameters
			Sender:        args.Sender,
			Signature:     args.signature(),
			Paymaster:     args.paymaster(),
			PaymasterData: args.paymasterData(),
			Deployer:      args.deployer(),
			DeployerData:  args.deployerData(),
			BuilderFee:    (*big.Int)(args.BuilderFee),
			ValidationGas: args.validationGas(),
			PaymasterGas:  args.paymasterGas(),
			PostOpGas:     args.postOpGas(),
			// RIP-7712 parameter
			BigNonce: (*big.Int)(args.BigNonce),
		}
		data = &aatx
		hash := types.NewTx(data).Hash()
		log.Error("RIP-7560 transaction created", "sender", aatx.Sender.Hex(), "hash", hash)
	case args.BlobHashes != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		data = &types.BlobTx{
			To:         *args.To,
			ChainID:    uint256.MustFromBig((*big.Int)(args.ChainID)),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasFeeCap:  uint256.MustFromBig((*big.Int)(args.MaxFeePerGas)),
			GasTipCap:  uint256.MustFromBig((*big.Int)(args.MaxPriorityFeePerGas)),
			Value:      uint256.MustFromBig((*big.Int)(args.Value)),
			Data:       args.data(),
			AccessList: al,
			BlobHashes: args.BlobHashes,
			BlobFeeCap: uint256.MustFromBig((*big.Int)(args.BlobFeeCap)),
		}
		if args.Blobs != nil {
			data.(*types.BlobTx).Sidecar = &types.BlobTxSidecar{
				Blobs:       args.Blobs,
				Commitments: args.Commitments,
				Proofs:      args.Proofs,
			}
		}

	case args.MaxFeePerGas != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		data = &types.DynamicFeeTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasFeeCap:  (*big.Int)(args.MaxFeePerGas),
			GasTipCap:  (*big.Int)(args.MaxPriorityFeePerGas),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: al,
		}

	case args.AccessList != nil:
		data = &types.AccessListTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasPrice:   (*big.Int)(args.GasPrice),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: *args.AccessList,
		}

	default:
		data = &types.LegacyTx{
			To:       args.To,
			Nonce:    uint64(*args.Nonce),
			Gas:      uint64(*args.Gas),
			GasPrice: (*big.Int)(args.GasPrice),
			Value:    (*big.Int)(args.Value),
			Data:     args.data(),
		}
	}
	return types.NewTx(data)
}

// ToTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) ToTransaction() *types.Transaction {
	return args.toTransaction()
}
