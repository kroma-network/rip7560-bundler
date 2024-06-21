// These are the same BundlerCollectorTracer types from github.com/eth-infinitism/bundler ported for Go.

package tracer

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
)

type HexMap = map[string]string
type Counts = map[string]int

// AccessInfo provides context on read and write counts by storage slots.
type AccessInfo struct {
	Reads  HexMap `json:"reads"`
	Writes Counts `json:"writes"`
}
type AccessMap = map[common.Address]AccessInfo

// ContractSizeInfo provides context on the code size and call type used to access upstream contracts.
type ContractSizeInfo struct {
	ContractSize int    `json:"contractSize"`
	Opcode       string `json:"opcode"`
}
type ContractSizeMap map[common.Address]ContractSizeInfo

// ExtCodeAccessInfoMap provides context on potentially illegal use of EXTCODESIZE.
type ExtCodeAccessInfoMap map[common.Address]string

// CallFromEntryPoint provides context on opcodes and storage access made via the EntryPoint to UserOperation
// entities.
type CallFromEntryPointInfo struct {
	TopLevelMethodSig     hexutil.Bytes        `json:"topLevelMethodSig"`
	TopLevelTargetAddress common.Address       `json:"topLevelTargetAddress"`
	Opcodes               Counts               `json:"opcodes"`
	Access                AccessMap            `json:"access"`
	ContractSize          ContractSizeMap      `json:"contractSize"`
	ExtCodeAccessInfo     ExtCodeAccessInfoMap `json:"extCodeAccessInfo"`
	OOG                   bool                 `json:"oog"`
}

//func (c *CallFromEntryPointInfo) UnmarshalJSON(input []byte) error {
//	err := json.Unmarshal(input, c)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// CallInfo provides context on internal calls made during tracing.
type CallInfo struct {
	// Common
	Type string `json:"type"`

	// Method info
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Method string         `json:"method"`
	Value  string         `json:"value"`
	Gas    float64        `json:"gas"`

	// Exit info
	GasUsed float64 `json:"gasUsed"`
	Data    any     `json:"data"`
}

func (c *CallInfo) UnmarshalJSON(input []byte) error {
	type callFrame struct {
		Type    vm.OpCode       `json:"-"`
		From    common.Address  `json:"from"`
		To      *common.Address `json:"to"`
		Input   hexutil.Bytes   `json:"input"`
		Gas     hexutil.Uint64  `json:"gas"`
		GasUsed hexutil.Uint64  `json:"gasUsed"`
		Value   *hexutil.Big    `json:"value"`
	}
	data := callFrame{}
	err := json.Unmarshal(input, &data)
	if err != nil {
		return err
	}
	c.Type = string(data.Type)
	c.From = data.From
	c.To = *data.To
	c.Data = common.Bytes2Hex(data.Input)
	//TODO : is this needed?
	//c.Method = common.Bytes2Hex(data.Input)
	c.Value = data.Value.String()
	c.Gas = float64(data.Gas)
	c.GasUsed = float64(data.GasUsed)
	return nil
}

// LogInfo provides context from LOG opcodes during each step in the EVM trace.
type LogInfo struct {
	Topics []string `json:"topics"`
	Data   string   `json:"data"`
}

// BundlerCollectorReturn is the return value from performing an EVM trace with BundlerCollectorTracer.js.
type BundlerCollectorReturn struct {
	CallsFromEntryPoint []CallFromEntryPointInfo `json:"callsFromEntryPoint"`
	Keccak              []string                 `json:"keccak"`
	Calls               []CallInfo               `json:"calls"`
	Logs                []LogInfo                `json:"logs"`
	Debug               []any                    `json:"debug"`
}

// BundlerExecutionReturn is the return value from performing an EVM trace with BundlerExecutionTracer.js.
type BundlerExecutionReturn struct {
	Reverts            []string `json:"reverts"`
	ValidationOOG      bool     `json:"validationOOG"`
	ExecutionOOG       bool     `json:"executionOOG"`
	ExecutionGasLimit  float64  `json:"executionGasLimit"`
	UserOperationEvent *LogInfo `json:"userOperationEvent,omitempty"`
	Output             string   `json:"output"`
	Error              string   `json:"error"`
}
