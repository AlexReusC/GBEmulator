package lib

import "errors"

type inType int
type adType int
type condType int

const (
	in_Nop inType = iota
	in_Xor
	in_Jp
)

const (
	am_N16 adType = iota
	am_Imp
)

const (
	cond_None condType = iota
	cond_C
	cond_NC
	cond_Z
	cond_NZ
)

func checkCond(cpu *CPU, ct condType) (bool, error) {
	c := cpu.getFlag(c)
	z := cpu.getFlag(z)

	switch ct {
	case cond_None:
		return true, nil
	case cond_C:
		return c, nil
	case cond_NC:
		return !c, nil
	case cond_Z:
		return z, nil
	case cond_NZ:
		return !z, nil
	default:
		return false, errors.New("invalid condition")
	}
} 

type Instruction struct {
	InstructionType inType
	AddressMode     adType
	ConditionType   condType
}

var instructions = map[uint8]Instruction{
	0x00: {in_Nop, am_Imp, cond_None},
	0xC3: {in_Jp, am_N16, cond_None},
}
