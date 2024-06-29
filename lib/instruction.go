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
	A targetType 		= "A"
	B targetType		= "B"
	C targetType		= "C"
	D targetType		= "D"
	E targetType		= "E"
	F targetType		= "F"
	H targetType		= "H"
	L targetType		= "L"
	AF targetType		= "AF"
	BC targetType		= "BC"
	DE targetType		= "DE"
	HL targetType		= "HL"
	SP targetType		= "SP"
	n targetType		= "n"
	nn targetType		= "nn"
	C_M targetType		= "(C)"
	BC_M targetType		= "(BC)"
	DE_M targetType		= "(DE)"
	HL_M targetType		= "(HL)"
	HLP_M targetType	= "(HL+)"
	HLM_M targetType	= "(HL-)"
	n_M targetType		= "(n)"
	nn_M targetType		= "(nn)"
	None targetType		= "none"
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
	0x00: {in_Nop,  None, 	None, 	cond_None},
	0x01: {in_Ld16, BC, 	nn, 	cond_None},
	0x02: {in_Ld8, 	BC_M, 	A, 		cond_None},
	0x03: {in_Ld8, 	B, 		n, 		cond_None},
	0x08: {in_Ld16, nn_M,	SP, 	cond_None},
	0x0A: {in_Ld8, 	A,		BC_M, 	cond_None},
	0x0E: {in_Ld8, 	C, 		n, 		cond_None},
	// 0x1X
	0x11: {in_Ld16, DE,		nn, 	cond_None},
	0x12: {in_Ld8, 	DE_M, 	A, 		cond_None},
	0x16: {in_Ld8, 	D, 		n, 		cond_None},
	0x18: {in_Jr, 	None, 	n, 		cond_None},
	0x1A: {in_Ld8,	A, 		DE_M, 	cond_None},
	0x1E: {in_Ld8,	E, 		n, 		cond_None},
	// 0x2X
	0x20: {in_Jr, 	None, 	n,		cond_NZ},
	0x21: {in_Ld16, HL, 	nn, 	cond_None},
	0x22: {in_Ld8, 	HLP_M, 	A, 		cond_None},
	0x26: {in_Ld8, 	H, 		n, 		cond_None},
	0x28: {in_Jr, 	None, 	n, 		cond_Z},
	0x2A: {in_Ld8, 	A, 		HLP_M,	cond_None},
	0x2E: {in_Ld8, 	L, 		n, 		cond_None},
	// 0x3X
	0x30: {in_Jr, 	None, 	n,		cond_NC},
	0x31: {in_Ld16, SP, 	nn, 	cond_None},
	0x32: {in_Ld8, 	HLM_M, 	A, 		cond_None},
	0x36: {in_Ld8, 	HL_M, 	n, 		cond_None},
	0x38: {in_Jr, 	None, 	n, 		cond_C},
	0x3A: {in_Ld8, 	A, 		HLM_M, 	cond_None},
	0x3E: {in_Ld8, 	A, 		n, 		cond_None},
	//0x4X
	0x41: {in_Ld8, 	B, 		B, 		cond_None},
	0x42: {in_Ld8, 	B, 		C, 		cond_None},
	0x43: {in_Ld8, 	B, 		D, 		cond_None},
	0x44: {in_Ld8, 	B, 		E, 		cond_None},
	0x45: {in_Ld8, 	B, 		H, 		cond_None},
	0x46: {in_Ld8, 	B, 		HL_M, 	cond_None},
	0x47: {in_Ld8, 	B, 		A, 		cond_None},
	0x48: {in_Ld8, 	C, 		B, 		cond_None},
	0x49: {in_Ld8, 	C, 		C, 		cond_None},
	0x4A: {in_Ld8, 	C, 		D, 		cond_None},
	0x4B: {in_Ld8, 	C, 		E, 		cond_None},
	0x4C: {in_Ld8, 	C, 		H, 		cond_None},
	0x4D: {in_Ld8, 	C, 		L, 		cond_None},
	0x4E: {in_Ld8, 	C, 		HL_M, 	cond_None},
	0x4F: {in_Ld8, 	C,		A, 		cond_None},
	//0x5xX
	0x50: {in_Ld8, 	D, 		B, 		cond_None},
	0x51: {in_Ld8, 	D, 		C, 		cond_None},
	0x52: {in_Ld8, 	D, 		D, 		cond_None},
	0x53: {in_Ld8, 	D, 		E, 		cond_None},
	0x54: {in_Ld8, 	D, 		H, 		cond_None},
	0x55: {in_Ld8, 	D, 		L, 		cond_None},
	0x56: {in_Ld8, 	D, 		HL_M, 	cond_None},
	0x57: {in_Ld8, 	D, 		A, 		cond_None},
	0x58: {in_Ld8, 	E, 		B, 		cond_None},
	0x59: {in_Ld8, 	E, 		C, 		cond_None},
	0x5A: {in_Ld8, 	E, 		D, 		cond_None},
	0x5B: {in_Ld8, 	E, 		E, 		cond_None},
	0x5C: {in_Ld8, 	E, 		H, 		cond_None},
	0x5D: {in_Ld8, 	E, 		L, 		cond_None},
	0x5E: {in_Ld8, 	E, 		HL_M, 	cond_None},
	0x5F: {in_Ld8, 	E, 		A, 		cond_None},
	//0x6X
	0x60: {in_Ld8, 	H, 		B, 		cond_None},
	0x61: {in_Ld8, 	H, 		C, 		cond_None},
	0x62: {in_Ld8, 	H, 		D, 		cond_None},
	0x63: {in_Ld8, 	H, 		E, 		cond_None},
	0x64: {in_Ld8, 	H, 		H, 		cond_None},
	0x65: {in_Ld8, 	H, 		L, 		cond_None},
	0x66: {in_Ld8, 	H, 		HL_M, 	cond_None},
	0x67: {in_Ld8, 	H, 		A, 		cond_None},
	0x68: {in_Ld8, 	L, 		B, 		cond_None},
	0x69: {in_Ld8, 	L, 		C, 		cond_None},
	0x6A: {in_Ld8, 	L, 		D, 		cond_None},
	0x6B: {in_Ld8, 	L, 		E, 		cond_None},
	0x6C: {in_Ld8, 	L, 		H, 		cond_None},
	0x6D: {in_Ld8, 	L, 		L, 		cond_None},
	0x6E: {in_Ld8, 	L, 		HL_M, 	cond_None},
	0x6F: {in_Ld8, 	L, 		A, 		cond_None},
	//0x7X
	0x70: {in_Ld8, 	HL_M, 	B, 		cond_None},
	0x71: {in_Ld8, 	HL_M, 	C, 		cond_None},
	0x72: {in_Ld8, 	HL_M, 	D, 		cond_None},
	0x73: {in_Ld8, 	HL_M, 	E, 		cond_None},
	0x74: {in_Ld8, 	HL_M, 	H, 		cond_None},
	0x75: {in_Ld8, 	HL_M, 	L, 		cond_None},
	//HALT instruction
	0x77: {in_Ld8, 	HL_M, 	A, 		cond_None},
	0x78: {in_Ld8, 	A, 		B, 		cond_None},
	0x79: {in_Ld8, 	A, 		C, 		cond_None},
	0x7A: {in_Ld8, 	A, 		D, 		cond_None},
	0x7B: {in_Ld8, 	A, 		E, 		cond_None},
	0x7C: {in_Ld8, 	A, 		H, 		cond_None},
	0x7D: {in_Ld8, 	A, 		L, 		cond_None},
	0x7E: {in_Ld8, 	A, 		HL_M, 	cond_None},
	0x7F: {in_Ld8, 	A, 		A, 		cond_None},

	//0x8X
	0x80: {in_Add, 	B, 		A, 		cond_None},

	//0xCX
	0xC0: {in_Ret, 	None, 	None, 	cond_NZ},
	0xC1: {in_Pop, 	BC,		None,  	cond_None},
	0xC2: {in_Jp, 	None, 	nn, 	cond_NZ},
	0xC3: {in_Jp, 	None, 	nn, 	cond_None},
	0xC4: {in_Call, None, 	nn, 	cond_NZ},
	0xC5: {in_Push, None, 	BC, 	cond_None},
	0xC7: {in_Rst, 	None, 	None, 	cond_None},
	0xC8: {in_Ret, 	None, 	None, 	cond_Z},
	0xC9: {in_Ret, 	None, 	None, 	cond_None},
	0xCA: {in_Jp, 	None, 	nn, 	cond_Z},
	0xCC: {in_Call, None, 	nn, 	cond_Z},
	0xCD: {in_Call, None, 	nn, 	cond_None},
	0xCF: {in_Rst, 	None, 	None, 	cond_None},
	//0xDX
	0xD0: {in_Ret, 	None, 	None, 	cond_NC},
	0xD1: {in_Pop, 	DE,  	None,  	cond_None},
	0xD2: {in_Jp, 	None, 	nn, 	cond_NC},
	0xD4: {in_Call, None, 	nn, 	cond_NC},
	0xD5: {in_Push, None, 	DE, 	cond_None},
	0xD7: {in_Rst, 	None, 	None, 	cond_None},
	0xD8: {in_Ret, 	None, 	None, 	cond_C},
	0xD9: {in_Reti, None, 	None, 	cond_None},
	0xDA: {in_Jp, 	None, 	nn, 	cond_C},
	0xDC: {in_Jp, 	None, 	nn, 	cond_C},
	0xDF: {in_Rst, 	None, 	None,	cond_None},
	//0xEX
	0xE0: {in_Ldh, 	n_M, 	A, 		cond_None},
	0xE1: {in_Pop, 	HL, 	None, 	cond_None},
	0xE2: {in_Ldh, 	C_M, 	A, 		cond_None},
	0xE5: {in_Push, None, 	HL, 	cond_None},
	0xE7: {in_Push, None, 	None, 	cond_None},
	0xE9: {in_Jp, 	None, 	HL, 	cond_None},
	0xEA: {in_Ld8, 	nn_M, 	A, 		cond_None},
	0xEF: {in_Rst, 	None, 	None, 	cond_None},
	//0xFX
	0xF0: {in_Ldh, 	A, 		n_M, 	cond_None},
	0xF1: {in_Pop, 	AF, 	None,  	cond_None},
	0xF2: {in_Ldh, 	A, 		C_M, 	cond_None},
	0xF3: {in_Di, 	None, 	None, 	cond_None},
	0xF5: {in_Push, None, 	AF, 	cond_None},
	0xF7: {in_Rst, 	None, 	None, 	cond_None},
	0xFA: {in_Ld8, 	A, 		nn_M, 	cond_None},
	0xFF: {in_Rst, 	None, 	None, 	cond_None},
	//TODO: More instructions
}
