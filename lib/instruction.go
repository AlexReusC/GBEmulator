package lib

import "errors"

type inType string 
type condType string
type targetType string

const (
	in_Nop inType	 	= "Nop"
	in_Xor inType		= "Xor"
	in_Jp inType		= "Jp"
	in_Jr inType 		= "Jr"
	in_Add inType		= "Add"
	in_Di inType		= "Di"
	in_Ld8 inType		= "Ld8"
	in_Ld16 inType		= "Ld16"
	in_Ldh inType 		= "Ldh"
	in_Push inType 		= "Push"
	in_Pop inType		= "Pop"
	in_Call inType 		= "Call"
	in_Ret inType		= "Ret"
	in_Reti inType		= "Reti"
	in_Rst inType		= "Rst"
)

const (
	target_A targetType 		= "A"
	target_B targetType			= "B"
	target_C targetType			= "C"
	target_D targetType			= "D"
	target_E targetType			= "E"
	target_F targetType			= "F"
	target_H targetType			= "H"
	target_L targetType			= "L"
	target_AF targetType		= "AF"
	target_BC targetType		= "BC"
	target_DE targetType		= "DE"
	target_HL targetType		= "HL"
	target_SP targetType		= "SP"
	target_n targetType			= "n"
	target_nn targetType		= "nn"
	target_C_M targetType		= "(C)"
	target_BC_M targetType		= "(BC)"
	target_DE_M targetType		= "(DE)"
	target_HL_M targetType		= "(HL)"
	target_HLP_M targetType		= "(HL+)"
	target_HLM_M targetType		= "(HL-)"
	target_n_M 	targetType		= "(n)"
	target_nn_M targetType		= "(nn)"
	target_None targetType		= "none"
)

const (
	cond_None condType	= "None"
	cond_C condType		= "C"
	cond_NC condType	= "NC"
	cond_Z condType		= "Z"
	cond_NZ condType	= "NZ"
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
	0x18: {in_Jr, target_None, target_n, cond_None},
	0x1A: {in_Ld8, target_A, target_DE_M, cond_None},
	0x1E: {in_Ld8, target_E, target_n, cond_None},
	// 0x2X
	0x20: {in_Jr, target_None, target_n, cond_NZ},
	0x21: {in_Ld16, target_HL, target_nn, cond_None},
	0x22: {in_Ld8, target_HLP_M, target_A, cond_None},
	0x26: {in_Ld8, target_H, target_n, cond_None},
	0x28: {in_Jr, target_None, target_n, cond_Z},
	0x2A: {in_Ld8, target_A, target_HLP_M, cond_None},
	0x2E: {in_Ld8, target_L, target_n, cond_None},
	// 0x3X
	0x30: {in_Jr, target_None, target_n, cond_NC},
	0x31: {in_Ld16, target_SP, target_nn, cond_None},
	0x32: {in_Ld8, target_HLM_M, target_A, cond_None},
	0x36: {in_Ld8, target_HL_M, target_n, cond_None},
	0x38: {in_Jr, target_None, target_n, cond_C},
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

	//0xCX
	0xC0: {in_Ret, target_None, target_None, cond_NZ},
	0xC1: {in_Pop, target_None, target_BC, cond_None},
	0xC2: {in_Jp, target_None, target_nn, cond_NZ},
	0xC3: {in_Jp, target_None, target_nn, cond_None},
	0xC4: {in_Call, target_None, target_nn, cond_NZ},
	0xC5: {in_Push, target_None, target_BC, cond_None},
	0xC7: {in_Rst, target_None, target_None, cond_None},
	0xC8: {in_Ret, target_None, target_None, cond_Z},
	0xC9: {in_Ret, target_None, target_None, cond_None},
	0xCA: {in_Jp, target_None, target_nn, cond_Z},
	0xCC: {in_Call, target_None, target_nn, cond_Z},
	0xCD: {in_Call, target_None, target_nn, cond_None},
	0xCF: {in_Rst, target_None, target_None, cond_None},
	//0xDX
	0xD0: {in_Ret, target_None, target_None, cond_NC},
	0xD1: {in_Pop, target_None, target_DE, cond_None},
	0xD2: {in_Jp, target_None, target_nn, cond_NC},
	0xD4: {in_Call, target_None, target_nn, cond_NC},
	0xD5: {in_Push, target_None, target_DE, cond_None},
	0xD7: {in_Rst, target_None, target_None, cond_None},
	0xD8: {in_Ret, target_None, target_None, cond_C},
	0xD9: {in_Reti, target_None, target_None, cond_None},
	0xDA: {in_Jp, target_None, target_nn, cond_C},
	0xDC: {in_Jp, target_None, target_nn, cond_C},
	0xDF: {in_Rst, target_None, target_None, cond_None},
	//0xEX
	0xE0: {in_Ldh, target_n_M, target_A, cond_None},
	0xE1: {in_Pop, target_None, target_HL, cond_None},
	0xE2: {in_Ldh, target_C_M, target_A, cond_None},
	0xE5: {in_Push, target_None, target_HL, cond_None},
	0xE7: {in_Push, target_None, target_None, cond_None},
	0xE9: {in_Jp, target_None, target_HL, cond_None},
	0xEA: {in_Ld8, target_nn_M, target_A, cond_None},
	0xEF: {in_Rst, target_None, target_None, cond_None},
	//0xFX
	0xF0: {in_Ldh, target_A, target_n_M, cond_None},
	0xF1: {in_Pop, target_None, target_AF, cond_None},
	0xF2: {in_Ldh, target_A, target_C_M, cond_None},
	0xF3: {in_Di, target_None, target_None, cond_None},
	0xF5: {in_Push, target_None, target_AF, cond_None},
	0xF7: {in_Rst, target_None, target_None, cond_None},
	0xFA: {in_Ld8, target_A, target_nn_M, cond_None},
	0xFF: {in_Rst, target_None, target_None, cond_None},
	//TODO: More instructions
}
