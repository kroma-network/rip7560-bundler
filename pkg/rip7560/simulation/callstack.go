package simulation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

type callEntry struct {
	From   common.Address
	To     common.Address
	Value  *big.Int
	Type   vm.OpCode
	Method string
	Revert any
	Return any
}

//func newCallStack(calls []native.CallFrame) []*callEntry {
//	var out []*callEntry
//	stack := utils.NewStack[native.CallFrame]()
//	for _, call := range calls {
//		if call.Type == vm.REVERT || call.Type == vm.RETURN {
//			top, _ := stack.Pop()
//
//			if top.Type == vm.CREATE {
//				// TODO: implement
//			} else if top.Type == vm.CREATE2 {
//				// TODO: implement
//			} else if call.Type == vm.REVERT {
//				// TODO: implement
//			} else {
//				out = append(out, &callEntry{
//					From:   top.From,
//					To:     *top.To,
//					Value:  top.Value,
//					Type:   top.Type,
//					Method: hex.EncodeToString(top.Input),
//					Return: call.Output,
//				})
//			}
//		} else {
//			stack.Push(call)
//		}
//	}
//
//	return out
//}
