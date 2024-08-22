package simulation

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/tracers/native"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
)

type TraceInput struct {
	Rpc     *rpc.Client
	Tx      *transaction.TransactionArgs
	ChainID *big.Int
}

type TraceOutput struct {
	TouchedContracts []common.Address
}

// TraceSimulateValidation makes call to debug_traceRip7560Validation to geth and returns
// information related to the validation phase of a RIP-7560 transaction.
func TraceSimulateValidation(in *TraceInput) (*TraceOutput, error) {
	var res native.Rip7560ValidationResult
	req := in.Tx
	if err := in.Rpc.CallContext(context.Background(), &res, "debug_traceRip7560Validation", &req, "latest"); err != nil {
		return nil, err
	}

	knownEntity, err := newKnownEntity(in.Tx, &res)
	if err != nil {
		return nil, err
	}

	ic := mapset.NewSet[common.Address]()
	for title, entity := range knownEntity {
		if entity.Info == nil {
			continue
		}
		if entity.Info.Oog {
			return nil, fmt.Errorf("%s OOG", title)
		}
		if _, ok := entity.Info.ExtCodeAccessInfo[config.EntryPointAddress]; ok {
			return nil, fmt.Errorf("%s has forbidden EXTCODE* access to the EntryPoint", title)
		}
		for opcode := range entity.Info.Opcodes {
			if bannedOpCodes.Contains(opcode) {
				return nil, fmt.Errorf("%s uses banned opcode: %s", title, opcode)
			}

			if bannedUnstakedOpCodes.Contains(opcode) {
				return nil, fmt.Errorf("unstaked %s uses banned opcode: %s", title, opcode)
			}
		}

		ic.Add(entity.Address)
		for addr := range entity.Info.ContractSize {
			ic.Add(addr)
		}
	}
	if knownEntity["deployer"].Info != nil {
		create2Count, ok := knownEntity["deployer"].Info.Opcodes[create2OpCode]
		if ok && (create2Count > 1 || len(in.Tx.GetDeployerData()) == 0) {
			return nil, fmt.Errorf("deployer with too many %s", create2OpCode)
		}
	}
	if knownEntity["account"].Info != nil {
		_, ok := knownEntity["account"].Info.Opcodes[create2OpCode]
		if ok {
			return nil, fmt.Errorf("account uses banned opcode: %s", create2OpCode)
		}
	}
	if knownEntity["paymaster"].Info != nil {
		_, ok := knownEntity["paymaster"].Info.Opcodes[create2OpCode]
		if ok {
			return nil, fmt.Errorf("paymaster uses banned opcode: %s", create2OpCode)
		}
	}

	// TODO : is this needed?
	//targetAddresses := []common.Address{in.Tx.GetSender()}
	//for _, entity := range knownEntity {
	//	targetAddresses = append(targetAddresses, entity.Address)
	//}
	//slotsByEntity := newStorageSlotsByEntity(res.Keccak, targetAddresses)
	//for title, entity := range knownEntity {
	//	v := &storageSlotsValidator{
	//		Tx:                    in.Tx,
	//		SenderSlots:           slotsByEntity[in.Tx.GetSender()],
	//		EntityName:            title,
	//		EntityAddr:            entity.Address,
	//		EntityAccessMap:       entity.Info.Access,
	//		EntityContractSizeMap: entity.Info.ContractSize,
	//		EntitySlots:           slotsByEntity[entity.Address],
	//	}
	//	if _, err := v.Process(); err != nil {
	//		return nil, err
	//	}
	//}

	callStack := newCallStack(res.Calls)
	for _, call := range callStack {
		if call.Method == methods.ValidatePaymasterTransactionSelector {
			out, err := methods.DecodevalidatePaymasterTransactionOutputOutput(call.Return)
			if err != nil {
				return nil, fmt.Errorf(
					"unexpected tracing result for tx: %s, %s",
					in.Tx.ToTransaction().Hash(),
					err,
				)
			}

			if len(out.Context) != 0 {
				return nil, errors.New("unstaked paymaster must not return context")
			}
		} else if call.Value.Cmp(common.Big0) == 1 {
			return nil, fmt.Errorf(
				"%s has a forbidden value transfer to %s",
				addr2KnownEntity(in.Tx, call.From),
				addr2KnownEntity(in.Tx, call.To),
			)
		}
	}

	return &TraceOutput{
		TouchedContracts: ic.ToSlice(),
	}, nil
}
