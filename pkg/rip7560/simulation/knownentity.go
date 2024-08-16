package simulation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/tracers/native"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

type knownEntity map[string]struct {
	Address common.Address
	Info    *native.Level
}

func newKnownEntity(
	tx *transaction.TransactionArgs,
	res *native.Rip7560ValidationResult,
) (knownEntity, error) {
	var si, fi, pi *native.Level
	for _, c := range res.CallsFromEntryPoint {
		switch c.TopLevelTargetAddress {
		case tx.GetSender():
			si = c
		case tx.GetPaymaster():
			pi = c
		case tx.GetDeployer():
			if c.TopLevelMethodSig.String() == methods.CreateAccountSelector {
				fi = c
			}
		default:
		}
	}

	return knownEntity{
		"account": {
			Address: tx.GetSender(),
			Info:    si,
		},
		"deployer": {
			Address: tx.GetDeployer(),
			Info:    fi,
		},
		"paymaster": {
			Address: tx.GetPaymaster(),
			Info:    pi,
		},
	}, nil
}

func addr2KnownEntity(tx *transaction.TransactionArgs, addr common.Address) string {
	if addr == tx.GetDeployer() {
		return "deployer"
	} else if addr == tx.GetSender() {
		return "account"
	} else if addr == tx.GetPaymaster() {
		return "paymaster"
	} else {
		return addr.String()
	}
}
