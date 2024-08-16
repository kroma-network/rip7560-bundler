package simulation

import (
	mapset "github.com/deckarep/golang-set/v2"
)

var (
	// Only one create2 opcode is allowed if these two conditions are met:
	// 	1. op.initcode.length != 0
	// 	2. During account simulation (i.e. before markerOpCode)
	create2OpCode = "CREATE2"

	// List of opcodes not allowed during simulation for depth > 1 (i.e. account, paymaster, or contracts
	// called by them).
	bannedOpCodes = mapset.NewSet(
		"GASPRICE",
		"GASLIMIT",
		"DIFFICULTY",
		"TIMESTAMP",
		"BASEFEE",
		"BLOCKHASH",
		"NUMBER",
		"ORIGIN",
		"GAS",
		"CREATE",
		"COINBASE",
		"SELFDESTRUCT",
	)

	// List of opcodes not allowed during validation for unstaked entities.
	bannedUnstakedOpCodes = mapset.NewSet(
		"SELFBALANCE",
		"BALANCE",
	)

	revertOpCode = "REVERT"
	returnOpCode = "RETURN"
)
