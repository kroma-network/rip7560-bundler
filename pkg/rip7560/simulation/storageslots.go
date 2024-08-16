package simulation

import (
	"fmt"
	"math/big"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers/native"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/pkg/rip7560/transaction"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
)

var (
	accessModeRead       = "read"
	accessModeWrite      = "write"
	associatedSlotOffset = big.NewInt(128)
)

type storageSlots mapset.Set[string]

type storageSlotsByEntity map[common.Address]storageSlots

func newStorageSlotsByEntity(keccak []string, targetAddresses []common.Address) storageSlotsByEntity {
	storageSlotsByEntity := make(storageSlotsByEntity)

	for _, k := range keccak {
		value := hexutil.Encode(crypto.Keccak256(common.Hex2Bytes(k[2:])))

		for _, addr := range targetAddresses {
			if addr == common.HexToAddress("0x") {
				continue
			}
			if _, ok := storageSlotsByEntity[addr]; !ok {
				storageSlotsByEntity[addr] = mapset.NewSet[string]()
			}

			addrPadded := hexutil.Encode(common.LeftPadBytes(addr.Bytes(), 32))
			if strings.HasPrefix(k, addrPadded) {
				storageSlotsByEntity[addr].Add(value)
			}
		}
	}

	return storageSlotsByEntity
}

type storageSlotsValidator struct {
	// Global parameters
	Tx *transaction.TransactionArgs

	// Parameters of specific entities required for all validation
	SenderSlots storageSlots

	// Parameters of the entity under validation
	EntityName            string
	EntityAddr            common.Address
	EntityAccessMap       native.AccessMap
	EntityContractSizeMap native.ContractSizeMap
	EntitySlots           storageSlots
}

func isAssociatedWith(entitySlots storageSlots, slot string) bool {
	slotBN, _ := big.NewInt(0).SetString(slot, 0)
	for _, entitySlot := range entitySlots.ToSlice() {
		entitySlotBN, _ := big.NewInt(0).SetString(entitySlot, 0)
		maxAssocSlotBN := big.NewInt(0).Add(entitySlotBN, associatedSlotOffset)
		if slotBN.Cmp(entitySlotBN) >= 0 && slotBN.Cmp(maxAssocSlotBN) <= 0 {
			return true
		}
	}

	return false
}

func (v *storageSlotsValidator) Process() error {
	senderSlots := v.SenderSlots
	if senderSlots == nil {
		senderSlots = mapset.NewSet[string]()
	}
	entitySlots := v.EntitySlots
	if entitySlots == nil {
		entitySlots = mapset.NewSet[string]()
	}

	for addr, access := range v.EntityAccessMap {
		if addr == v.Tx.GetSender() || addr == config.EntryPointAddress {
			continue
		}

		var mustStakeSlot string
		accessTypes := map[string]any{
			accessModeRead:  access.Reads,
			accessModeWrite: access.Writes,
		}
		for mode, val := range accessTypes {
			slots := []string{}
			if readMap, ok := val.(tracer.HexMap); ok {
				for slot := range readMap {
					slots = append(slots, slot)
				}
			} else if writeMap, ok := val.(tracer.Counts); ok {
				for slot := range writeMap {
					slots = append(slots, slot)
				}
			} else {
				return fmt.Errorf("cannot decode %s access type: %+v", mode, val)
			}

			// FIXME : remove or uncommnet
			//for _, slot := range slots {
			//	if isAssociatedWith(senderSlots, slot) {
			//		if (len(v.Tx.GetDeployerData()) > 0) ||
			//			(len(v.Tx.GetDeployerData()) > 0 && v.EntityAddr != v.Tx.GetSender()) {
			//			mustStakeSlot = slot
			//		} else {
			//			continue
			//		}
			//	} else if amIds := v.AltMempools.HasInvalidStorageAccessException(
			//		v.EntityName,
			//		addr2KnownEntity(v.Tx, addr),
			//		slot,
			//	); (isAssociatedWith(entitySlots, slot) || mode == accessModeRead) && len(amIds) == 0 {
			//		mustStakeSlot = slot
			//	} else if len(amIds) > 0 {
			//		altMempoolIds = append(altMempoolIds, amIds...)
			//	} else {
			//		return fmt.Errorf(
			//			"%s has forbidden %s to %s slot %s",
			//			v.EntityName,
			//			mode,
			//			addr2KnownEntity(v.Tx, addr),
			//			slot,
			//		)
			//	}
			//}
		}

		if mustStakeSlot != "" {
			return fmt.Errorf(
				"unstaked %s accessed %s slot %s",
				v.EntityName,
				addr2KnownEntity(v.Tx, addr),
				mustStakeSlot,
			)
		}
	}

	return nil
}
