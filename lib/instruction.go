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
	target_BC
	target_DE
	target_HL
	target_SP
	target_n
	target_nn
	target_BC_M
	target_DE_M
	target_HL_M
	target_HLP_M
	target_HLM_M
	target_nn_M
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

type Instruction struct {
	InstructionType inType
	Destination		targetType
	Source			targetType
	ConditionType   condType
}

var instructions = map[uint8]Instruction{
	// 0x0X
	0x00: {in_Nop, target_None, target_None, cond_None},
	0x01: {in_Ld16, target_BC, target_nn, cond_None},
	0x02: {in_Ld8, target_BC_M, target_A, cond_None},
	0x03: {in_Ld8, target_B, target_n, cond_None},
	0x08: {in_Ld16, target_nn_M, target_SP, cond_None},
	0x0A: {in_Ld8, target_A, target_BC_M, cond_None},
	0x0E: {in_Ld8, target_C, target_n, cond_None},
	// 0x1X
	0x11: {in_Ld16, target_DE, target_nn, cond_None},
	0x12: {in_Ld8, target_DE_M, target_A, cond_None},
	0x16: {in_Ld8, target_D, target_n, cond_None},
	0x1A: {in_Ld8, target_A, target_DE_M, cond_None},
	0x1E: {in_Ld8, target_E, target_n, cond_None},
	// 0x2X
	0x21: {in_Ld16, target_HL, target_nn, cond_None},
	0x22: {in_Ld8, target_HLP_M, target_A, cond_None},
	0x26: {in_Ld8, target_H, target_n, cond_None},
	0x2A: {in_Ld8, target_A, target_HLP_M, cond_None},
	0x2E: {in_Ld8, target_L, target_n, cond_None},
	// 0x3X
	0x31: {in_Ld16, target_SP, target_nn, cond_None},
	0x32: {in_Ld8, target_HLM_M, target_A, cond_None},
	0x36: {in_Ld8, target_HL_M, target_n, cond_None},
	0x3A: {in_Ld8, target_A, target_HLM_M, cond_None},
	0x3E: {in_Ld8, target_A, target_n, cond_None},
	//0x4X
	0x41: {in_Ld8, target_B, target_B, cond_None},
	0x42: {in_Ld8, target_B, target_C, cond_None},
	0x43: {in_Ld8, target_B, target_D, cond_None},
	0x44: {in_Ld8, target_B, target_E, cond_None},
	0x45: {in_Ld8, target_B, target_H, cond_None},
	0x46: {in_Ld8, target_B, target_HL_M, cond_None},
	0x47: {in_Ld8, target_B, target_A, cond_None},
	0x48: {in_Ld8, target_C, target_B, cond_None},
	0x49: {in_Ld8, target_C, target_C, cond_None},
	0x4A: {in_Ld8, target_C, target_D, cond_None},
	0x4B: {in_Ld8, target_C, target_E, cond_None},
	0x4C: {in_Ld8, target_C, target_H, cond_None},
	0x4D: {in_Ld8, target_C, target_L, cond_None},
	0x4E: {in_Ld8, target_C, target_HL_M, cond_None},
	0x4F: {in_Ld8, target_C, target_A, cond_None},
	//0x5xX
	0x50: {in_Ld8, target_D, target_B, cond_None},
	0x51: {in_Ld8, target_D, target_C, cond_None},
	0x52: {in_Ld8, target_D, target_D, cond_None},
	0x53: {in_Ld8, target_D, target_E, cond_None},
	0x54: {in_Ld8, target_D, target_H, cond_None},
	0x55: {in_Ld8, target_D, target_L, cond_None},
	0x56: {in_Ld8, target_D, target_HL_M, cond_None},
	0x57: {in_Ld8, target_D, target_A, cond_None},
	0x58: {in_Ld8, target_E, target_B, cond_None},
	0x59: {in_Ld8, target_E, target_C, cond_None},
	0x5A: {in_Ld8, target_E, target_D, cond_None},
	0x5B: {in_Ld8, target_E, target_E, cond_None},
	0x5C: {in_Ld8, target_E, target_H, cond_None},
	0x5D: {in_Ld8, target_E, target_L, cond_None},
	0x5E: {in_Ld8, target_E, target_HL_M, cond_None},
	0x5F: {in_Ld8, target_E, target_A, cond_None},
	//0x6X
	0x60: {in_Ld8, target_H, target_B, cond_None},
	0x61: {in_Ld8, target_H, target_C, cond_None},
	0x62: {in_Ld8, target_H, target_D, cond_None},
	0x63: {in_Ld8, target_H, target_E, cond_None},
	0x64: {in_Ld8, target_H, target_H, cond_None},
	0x65: {in_Ld8, target_H, target_L, cond_None},
	0x66: {in_Ld8, target_H, target_HL_M, cond_None},
	0x67: {in_Ld8, target_H, target_A, cond_None},
	0x68: {in_Ld8, target_L, target_B, cond_None},
	0x69: {in_Ld8, target_L, target_C, cond_None},
	0x6A: {in_Ld8, target_L, target_D, cond_None},
	0x6B: {in_Ld8, target_L, target_E, cond_None},
	0x6C: {in_Ld8, target_L, target_H, cond_None},
	0x6D: {in_Ld8, target_L, target_L, cond_None},
	0x6E: {in_Ld8, target_L, target_HL_M, cond_None},
	0x6F: {in_Ld8, target_L, target_A, cond_None},
	//0x7X
	0x70: {in_Ld8, target_HL_M, target_B, cond_None},
	0x71: {in_Ld8, target_HL_M, target_C, cond_None},
	0x72: {in_Ld8, target_HL_M, target_D, cond_None},
	0x73: {in_Ld8, target_HL_M, target_E, cond_None},
	0x74: {in_Ld8, target_HL_M, target_H, cond_None},
	0x75: {in_Ld8, target_HL_M, target_L, cond_None},
	//HALT instruction
	0x77: {in_Ld8, target_HL_M, target_A, cond_None},
	0x78: {in_Ld8, target_A, target_B, cond_None},
	0x79: {in_Ld8, target_A, target_C, cond_None},
	0x7A: {in_Ld8, target_A, target_D, cond_None},
	0x7B: {in_Ld8, target_A, target_E, cond_None},
	0x7C: {in_Ld8, target_A, target_H, cond_None},
	0x7D: {in_Ld8, target_A, target_L, cond_None},
	0x7E: {in_Ld8, target_A, target_HL_M, cond_None},
	0x7F: {in_Ld8, target_A, target_A, cond_None},

	//0x8X
	0x80: {in_Add, target_B, target_A, cond_None},

	0xC3: {in_Jp, target_None, target_nn, cond_None},
	//0xEX
	0xEA: {in_Ld8, target_nn_M, target_A, cond_None},
	//0xFX
	0xF3: {in_Di, target_None, target_None, cond_None},
	0xFA: {in_Ld8, target_A, target_nn_M, cond_None},
	//TODO: More instructions
}
