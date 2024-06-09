package lib

import "errors"

type inType int
type adType int
type condType int
type targetType int

const (
	in_Nop inType = iota
	in_Xor
	in_Jp
	in_Add
	in_Di
	in_Ld8
	in_Ld16
)


const (
	target_A targetType = iota
	target_B
	target_C
	target_D
	target_E
	target_F
	target_H
	target_L
	target_nn
	target_None
)

const (
	cond_None condType = iota
	cond_C
	cond_NC
	cond_Z
	cond_NZ
)

func (cpu *CPU) checkCond( ct condType) (bool, error) {
	c := cpu.GetFlag(cf)
	z := cpu.GetFlag(zf)

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

// TODO: fix and destination are the other way around
type Instruction struct {
	InstructionType inType
	Destination		targetType
	Source			targetType
	ConditionType   condType
}

var instructions = map[uint8]Instruction{
	0x00: {in_Nop, target_None, target_None, cond_None},
	0xC3: {in_Jp, target_nn, target_None, cond_None},
	0x80: {in_Add, target_B, target_A, cond_None},
	0xF3: {in_Di, target_None, target_None, cond_None},
	//TODO: More instructions
}
