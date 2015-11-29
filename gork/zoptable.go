package gork

import (
	"fmt"
)

type ZeroOpFunc func(*ZMachine)
type OneOpFunc func(*ZMachine, uint16)
type TwoOpFunc func(*ZMachine, uint16, uint16)
type VarOpFunc func(*ZMachine, []uint16)

var zeroOpFuncs = []ZeroOpFunc{
	ZReturnTrue,
	ZReturnFalse,
	ZPrint,
	ZPrintRet,
	nil,
	nil,
	nil,
	nil,
	nil,
	ZPop,
}

var oneOpFuncs = []OneOpFunc{
	ZJ0,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	ZPrintAt,
	nil,
	nil,
	ZPrintObject,
	ZReturn,
	nil,
	ZPrintAtPacked,
	ZLoad,
	ZNot,
}

var twoOpFuncs = []TwoOpFunc{
	ZNOOP,
	nil, // ZJe is a two op func but it accepts VAR count of args
	ZJl,
	ZJg,
	nil,
	nil,
	nil,
	nil,
	ZOr,
	ZAnd,
	nil,
	nil,
	nil,
	ZStore,
	nil,
	ZLoadW,
	ZLoadB,
	nil,
	nil,
	nil,
	ZAdd,
	ZSub,
	ZMul,
	ZDiv,
	ZMod,
}

var varOpFuncs = []VarOpFunc{
	ZCall,
	ZStoreW,
	ZStoreB,
	nil,
	nil,
	ZPrintChar,
	ZPrintNum,
	nil,
	ZPush,
	ZPull,
}

func ZCall(zm *ZMachine, operands []uint16) {
	routineAddr := PackedAddress(operands[0])

	retAddr := zm.seq.pos
	zm.seq.pos = routineAddr
	routine := NewZRoutine(zm.seq, retAddr)

	zm.stack.Push(routine)
	// fmt.Println(routine)

	if routineAddr == 0 {
		ZReturnFalse(zm)
		return
	}

	if len(operands) > 1 {
		// copy operands to locals
		for i, v := range operands[1:] {
			routine.locals[i] = v
		}
	}
}

func ZReturn(zm *ZMachine, retValue uint16) {
	zm.seq.pos = zm.stack.Pop().retAddr
	// fmt.Printf("Returning to 0x%X\n", zm.seq.pos)
	zm.StoreReturn(retValue)
}

func ZReturnFalse(zm *ZMachine) {
	ZReturn(zm, uint16(0))
}

func ZReturnTrue(zm *ZMachine) {
	ZReturn(zm, uint16(1))
}

func ZJe(zm *ZMachine, args []uint16) {
	conditionOk := false
	for _, v := range args[1:] {
		if v == args[0] {
			conditionOk = true
			break
		}
	}
	zm.Branch(conditionOk)
}

func ZJl(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.Branch(lhs < rhs)
}

func ZJg(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.Branch(lhs > rhs)
}

func ZJ0(zm *ZMachine, op uint16) {
	zm.Branch(op == 0)
}

func ZPrint(zm *ZMachine) {
	str := zm.seq.DecodeZString(zm.header)
	fmt.Print(str)
}

func ZPrintRet(zm *ZMachine) {
	ZPrint(zm)
	fmt.Println("")
	ZReturnTrue(zm)
}

func ZPrintObject(zm *ZMachine, obj uint16) {
	// objects are 1-based
	fmt.Print(zm.objects[obj-1].name)
}

func ZPrintAt(zm *ZMachine, addr uint16) {
	str := zm.seq.mem.DecodeZStringAt(addr, zm.header)
	fmt.Print(str)
}

func ZPrintAtPacked(zm *ZMachine, paddr uint16) {
	ZPrintAt(zm, PackedAddress(paddr))
}

func ZPrintNum(zm *ZMachine, args []uint16) {
	fmt.Print(args[0])
}

func ZPrintChar(zm *ZMachine, args []uint16) {
	// print only ASCII
	if args[0] == 13 {
		fmt.Println("")
	} else if args[0] >= 32 && args[0] <= 126 {
		fmt.Printf("%c", args[0])
	} // ignore everything else
}

func ZAdd(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.StoreReturn(lhs + rhs)
}

func ZSub(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.StoreReturn(lhs - rhs)
}

func ZMul(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.StoreReturn(lhs * rhs)
}

func ZDiv(zm *ZMachine, lhs uint16, rhs uint16) {
	if rhs == 0 {
		panic("division by zero error")
	}
	zm.StoreReturn(lhs / rhs)
}

func ZMod(zm *ZMachine, lhs uint16, rhs uint16) {
	if rhs == 0 {
		panic("mod by zero error")
	}
	zm.StoreReturn(lhs % rhs)
}

func ZOr(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.StoreReturn(lhs | rhs)
}

func ZAnd(zm *ZMachine, lhs uint16, rhs uint16) {
	zm.StoreReturn(lhs & rhs)
}

// v3
func ZNot(zm *ZMachine, arg uint16) {
	zm.StoreReturn(^arg)
}

func ZNOOP(_ *ZMachine, _ uint16, _ uint16) {
	panic("NO OP 2OP")
}

func ZLoad(zm *ZMachine, varnum uint16) {
	zm.StoreReturn(zm.GetVarAt(byte(varnum)))
}

func ZLoadB(zm *ZMachine, array uint16, bidx uint16) {
	// TODO access violation
	zm.StoreReturn(uint16(zm.seq.mem.ByteAt(array + bidx)))
}

func ZLoadW(zm *ZMachine, array uint16, widx uint16) {
	// TODO access violation
	// index is the index of the nth word
	zm.StoreReturn(zm.seq.mem.WordAt(array + widx*2))
}

func ZStore(zm *ZMachine, varnum uint16, value uint16) {
	zm.StoreVarAt(byte(varnum), value)
}

func ZStoreB(zm *ZMachine, args []uint16) {
	// TODO access violation
	addr := args[0] + args[1]
	zm.seq.mem.WriteByteAt(addr, byte(args[2]))
}

func ZStoreW(zm *ZMachine, args []uint16) {
	// TODO access violation
	// index is the index of the nth word
	addr := args[0] + args[1]*2
	zm.seq.mem.WriteWordAt(addr, args[2])
}

func ZPush(zm *ZMachine, args []uint16) {
	zm.stack.Top().locals = append(zm.stack.Top().locals, args[0])
}

func ZPull(zm *ZMachine, args []uint16) {
	varnum := byte(args[0]) - 1
	topLocals := &zm.stack.Top().locals
	*topLocals = append((*topLocals)[:varnum], (*topLocals)[varnum:]...)
	// should not zm.StoreReturn popped value
}

func ZPop(zm *ZMachine) {
	topLocals := &zm.stack.Top().locals
	*topLocals = append((*topLocals)[:1], (*topLocals)[1:]...)
}