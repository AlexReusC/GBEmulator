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
)

const (
	am_D16 adType = iota
	am_Imp
	am_Mr_R
	am_R_Mr
	am_R_HlI
	am_R_HlD
	am_HlI_R
	am_HlD_R

)

const (
	target_A targetType = iota
)

const (
	cond_None condType = iota
	cond_C
	cond_NC
	cond_Z
	cond_NZ
)

func checkCond(cpu *CPU, ct condType) (bool, error) {
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

type Instruction struct {
	InstructionType inType
	AddressMode     adType
	Source			*targetType
	Destination		*targetType
	ConditionType   condType
}

var instructions = map[uint8]Instruction{
	0x00: {in_Nop, am_Imp, nil, nil, cond_None},
	0xC3: {in_Jp, am_D16, nil, nil, cond_None},
}
