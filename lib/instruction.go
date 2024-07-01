package lib

import "errors"

type procedure string 
type conditional string
type target string

const (
	Nop 		procedure = "Nop"
	Jp 			procedure = "Jp"
	Jr 			procedure = "Jr"
	Di 			procedure = "Di"
	Ld8 		procedure = "Ld8"
	Ld16 		procedure = "Ld16"
	Ldh 		procedure = "Ldh"
	Push 		procedure = "Push"
	Pop 		procedure = "Pop"
	Call 		procedure = "Call"
	Ret 		procedure = "Ret"
	Reti 		procedure = "Reti"
	Rst 		procedure = "Rst"
	Inc 		procedure = "Inc"
	Dec 		procedure = "Dec"
	Add 		procedure = "Add"
	AddHl 		procedure = "AddHl"
	Add16_8 	procedure = "Add16_8"
	Adc 		procedure = "Adc"
	Sub 		procedure = "Sub"
	Sbc 		procedure = "Sbc"
	And 		procedure = "And"
	Xor 		procedure = "Xor"
	Or 			procedure = "Or"
	Cp 			procedure = "Cp"

)

const (
	A 		target	= "A"
	B 		target	= "B"
	C 		target	= "C"
	D 		target	= "D"
	E 		target	= "E"
	F 		target	= "F"
	H 		target	= "H"
	L 		target	= "L"
	AF 		target	= "AF"
	BC 		target	= "BC"
	DE 		target	= "DE"
	HL 		target	= "HL"
	SP 		target	= "SP"
	e8 		target	= "e8"
	n 		target	= "n"
	nn 		target	= "nn"
	C_M		target	= "(C)"
	BC_M	target	= "(BC)"
	DE_M	target	= "(DE)"
	HL_M	target	= "(HL)"
	HLP_M	target	= "(HL+)"
	HLM_M	target	= "(HL-)"
	n_M		target	= "(n)"
	nn_M	target	= "(nn)"
	None	target	= "none"
)

const (
	cond_None 	conditional	= "None"
	cond_C 		conditional	= "C"
	cond_NC 	conditional	= "NC"
	cond_Z 		conditional	= "Z"
	cond_NZ 	conditional	= "NZ"
)

func (cpu *CPU) checkCond( ct conditional) (bool, error) {
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
	InstructionType procedure
	Destination		target
	Source			target
	ConditionType   conditional
}

var instructions = map[uint8]Instruction{
	// 0x0X
	0x00: {Nop,  	None, 	None, 	cond_None},
	0x01: {Ld16, 	BC, 	nn, 	cond_None},
	0x02: {Ld8, 	BC_M, 	A, 		cond_None},
	0x03: {Inc, 	None, 	BC, 	cond_None},
	0x04: {Inc, 	None, 	B, 		cond_None},
	0x05: {Dec, 	None, 	B, 		cond_None},
	0x06: {Ld8, 	B, 		n, 		cond_None},
	0x08: {Ld16, 	nn_M,	SP, 	cond_None},
	0x09: {AddHl, 	HL, 	BC,		cond_None},
	0x0A: {Ld8, 	A,		BC_M, 	cond_None},
	0x0B: {Dec, 	None, 	BC, 	cond_None},
	0x0C: {Inc, 	None, 	C, 		cond_None},
	0x0D: {Dec, 	None, 	C, 		cond_None},
	0x0E: {Ld8, 	C, 		n, 		cond_None},
	// 0x1X
	0x11: {Ld16, 	DE,		nn, 	cond_None},
	0x12: {Ld8, 	DE_M, 	A, 		cond_None},
	0x13: {Inc, 	None, 	DE, 	cond_None},
	0x14: {Inc, 	None, 	D, 		cond_None},
	0x15: {Dec, 	None, 	D, 		cond_None},
	0x16: {Ld8, 	D, 		n, 		cond_None},
	0x18: {Jr, 		None, 	n, 		cond_None},
	0x19: {AddHl, 	HL, 	DE,		cond_None},
	0x1A: {Ld8,		A, 		DE_M, 	cond_None},
	0x1B: {Dec, 	None, 	DE, 	cond_None},
	0x1C: {Inc, 	None, 	E, 		cond_None},
	0x1D: {Dec, 	None, 	E, 		cond_None},
	0x1E: {Ld8,		E, 		n, 		cond_None},
	// 0x2X
	0x20: {Jr, 		None, 	n,		cond_NZ},
	0x21: {Ld16, 	HL, 	nn, 	cond_None},
	0x22: {Ld8, 	HLP_M, 	A, 		cond_None},
	0x23: {Inc, 	None, 	HL, 	cond_None},
	0x24: {Inc, 	None, 	H, 		cond_None},
	0x25: {Dec, 	None, 	H, 		cond_None},
	0x26: {Ld8, 	H, 		n, 		cond_None},
	0x28: {Jr, 		None, 	n, 		cond_Z},
	0x29: {AddHl, 	HL, 	HL,		cond_None},
	0x2A: {Ld8, 	A, 		HLP_M,	cond_None},
	0x2B: {Dec, 	None, 	HL, 	cond_None},
	0x2C: {Inc, 	None, 	L, 		cond_None},
	0x2D: {Dec, 	None, 	L, 		cond_None},
	0x2E: {Ld8, 	L, 		n, 		cond_None},
	// 0x3X
	0x30: {Jr, 		None, 	n,		cond_NC},
	0x31: {Ld16, 	SP, 	nn, 	cond_None},
	0x32: {Ld8, 	HLM_M, 	A, 		cond_None},
	0x33: {Inc, 	None, 	SP, 	cond_None},
	0x34: {Inc, 	None, 	HL_M, 	cond_None},
	0x35: {Dec, 	None, 	HL_M, 	cond_None},
	0x36: {Ld8, 	HL_M, 	n, 		cond_None},
	0x38: {Jr, 		None, 	n, 		cond_C},
	0x39: {AddHl, 	HL, 	SP,		cond_None},
	0x3A: {Ld8, 	A, 		HLM_M, 	cond_None},
	0x3B: {Dec, 	None, 	SP, 	cond_None},
	0x3C: {Inc, 	None, 	A, 		cond_None},
	0x3D: {Dec, 	None, 	A, 		cond_None},
	0x3E: {Ld8, 	A, 		n, 		cond_None},
	//0x4X
	0x41: {Ld8, 	B, 		B, 		cond_None},
	0x42: {Ld8, 	B, 		C, 		cond_None},
	0x43: {Ld8, 	B, 		D, 		cond_None},
	0x44: {Ld8, 	B, 		E, 		cond_None},
	0x45: {Ld8, 	B, 		H, 		cond_None},
	0x46: {Ld8, 	B, 		HL_M, 	cond_None},
	0x47: {Ld8, 	B, 		A, 		cond_None},
	0x48: {Ld8, 	C, 		B, 		cond_None},
	0x49: {Ld8, 	C, 		C, 		cond_None},
	0x4A: {Ld8, 	C, 		D, 		cond_None},
	0x4B: {Ld8, 	C, 		E, 		cond_None},
	0x4C: {Ld8, 	C, 		H, 		cond_None},
	0x4D: {Ld8, 	C, 		L, 		cond_None},
	0x4E: {Ld8, 	C, 		HL_M, 	cond_None},
	0x4F: {Ld8, 	C,		A, 		cond_None},
	//0x5xX
	0x50: {Ld8, 	D, 		B, 		cond_None},
	0x51: {Ld8, 	D, 		C, 		cond_None},
	0x52: {Ld8, 	D, 		D, 		cond_None},
	0x53: {Ld8, 	D, 		E, 		cond_None},
	0x54: {Ld8, 	D, 		H, 		cond_None},
	0x55: {Ld8, 	D, 		L, 		cond_None},
	0x56: {Ld8, 	D, 		HL_M, 	cond_None},
	0x57: {Ld8, 	D, 		A, 		cond_None},
	0x58: {Ld8, 	E, 		B, 		cond_None},
	0x59: {Ld8, 	E, 		C, 		cond_None},
	0x5A: {Ld8, 	E, 		D, 		cond_None},
	0x5B: {Ld8, 	E, 		E, 		cond_None},
	0x5C: {Ld8, 	E, 		H, 		cond_None},
	0x5D: {Ld8, 	E, 		L, 		cond_None},
	0x5E: {Ld8, 	E, 		HL_M, 	cond_None},
	0x5F: {Ld8, 	E, 		A, 		cond_None},
	//0x6X
	0x60: {Ld8, 	H, 		B, 		cond_None},
	0x61: {Ld8, 	H, 		C, 		cond_None},
	0x62: {Ld8, 	H, 		D, 		cond_None},
	0x63: {Ld8, 	H, 		E, 		cond_None},
	0x64: {Ld8, 	H, 		H, 		cond_None},
	0x65: {Ld8, 	H, 		L, 		cond_None},
	0x66: {Ld8, 	H, 		HL_M, 	cond_None},
	0x67: {Ld8, 	H, 		A, 		cond_None},
	0x68: {Ld8, 	L, 		B, 		cond_None},
	0x69: {Ld8, 	L, 		C, 		cond_None},
	0x6A: {Ld8, 	L, 		D, 		cond_None},
	0x6B: {Ld8, 	L, 		E, 		cond_None},
	0x6C: {Ld8, 	L, 		H, 		cond_None},
	0x6D: {Ld8, 	L, 		L, 		cond_None},
	0x6E: {Ld8, 	L, 		HL_M, 	cond_None},
	0x6F: {Ld8, 	L, 		A, 		cond_None},
	//0x7X
	0x70: {Ld8, 	HL_M, 	B, 		cond_None},
	0x71: {Ld8, 	HL_M, 	C, 		cond_None},
	0x72: {Ld8, 	HL_M, 	D, 		cond_None},
	0x73: {Ld8, 	HL_M, 	E, 		cond_None},
	0x74: {Ld8, 	HL_M, 	H, 		cond_None},
	0x75: {Ld8, 	HL_M, 	L, 		cond_None},
	//HALT instruction
	0x77: {Ld8, 	HL_M, 	A, 		cond_None},
	0x78: {Ld8, 	A, 		B, 		cond_None},
	0x79: {Ld8, 	A, 		C, 		cond_None},
	0x7A: {Ld8, 	A, 		D, 		cond_None},
	0x7B: {Ld8, 	A, 		E, 		cond_None},
	0x7C: {Ld8, 	A, 		H, 		cond_None},
	0x7D: {Ld8, 	A, 		L, 		cond_None},
	0x7E: {Ld8, 	A, 		HL_M, 	cond_None},
	0x7F: {Ld8, 	A, 		A, 		cond_None},

	//0x8X
	0x80: {Add, 	None, 	B, 		cond_None},
	0x81: {Add, 	None, 	C, 		cond_None},
	0x82: {Add, 	None, 	D, 		cond_None},
	0x83: {Add, 	None, 	E, 		cond_None},
	0x84: {Add, 	None, 	H, 		cond_None},
	0x85: {Add, 	None, 	L, 		cond_None},
	0x86: {Add, 	None, 	HL_M, 	cond_None},
	0x87: {Add, 	None, 	A, 		cond_None},
	0x88: {Adc, 	None, 	B, 		cond_None},
	0x89: {Adc, 	None, 	C, 		cond_None},
	0x8A: {Adc, 	None, 	D, 		cond_None},
	0x8B: {Adc, 	None, 	E, 		cond_None},
	0x8C: {Adc, 	None, 	H, 		cond_None},
	0x8D: {Adc, 	None, 	L, 		cond_None},
	0x8E: {Adc, 	None, 	HL_M, 	cond_None},
	0x8F: {Adc, 	None, 	A, 		cond_None},
	//0x9X
	0x90: {Sub, 	None, 	B, 		cond_None},
	0x91: {Sub, 	None, 	C, 		cond_None},
	0x92: {Sub, 	None, 	D, 		cond_None},
	0x93: {Sub, 	None, 	E, 		cond_None},
	0x94: {Sub, 	None, 	H, 		cond_None},
	0x95: {Sub, 	None, 	L, 		cond_None},
	0x96: {Sub, 	None, 	HL_M, 	cond_None},
	0x97: {Sub, 	None, 	A, 		cond_None},
	0x98: {Sbc, 	None, 	B, 		cond_None},
	0x99: {Sbc, 	None, 	C, 		cond_None},
	0x9A: {Sbc, 	None, 	D, 		cond_None},
	0x9B: {Sbc, 	None, 	E, 		cond_None},
	0x9C: {Sbc, 	None, 	H, 		cond_None},
	0x9D: {Sbc, 	None, 	L, 		cond_None},
	0x9E: {Sbc, 	None, 	HL_M, 	cond_None},
	0x9F: {Sbc, 	None, 	A, 		cond_None},
	//0xAX
	0xA0: {And, 	None, 	B, 		cond_None},
	0xA1: {And, 	None, 	C, 		cond_None},
	0xA2: {And, 	None, 	D, 		cond_None},
	0xA3: {And, 	None, 	E, 		cond_None},
	0xA4: {And, 	None, 	H, 		cond_None},
	0xA5: {And, 	None, 	L, 		cond_None},
	0xA6: {And, 	None, 	HL_M, 	cond_None},
	0xA7: {And, 	None, 	A, 		cond_None},
	0xA8: {Xor, 	None, 	B, 		cond_None},
	0xA9: {Xor, 	None, 	C, 		cond_None},
	0xAA: {Xor, 	None, 	D, 		cond_None},
	0xAB: {Xor, 	None, 	E, 		cond_None},
	0xAC: {Xor, 	None, 	H, 		cond_None},
	0xAD: {Xor, 	None, 	L, 		cond_None},
	0xAE: {Xor, 	None, 	HL_M, 	cond_None},
	0xAF: {Xor, 	None, 	A, 		cond_None},
	//0xBX
	0xB0: {Or, 		None, 	B, 		cond_None},
	0xB1: {Or, 		None, 	C, 		cond_None},
	0xB2: {Or, 		None, 	D, 		cond_None},
	0xB3: {Or, 		None, 	E, 		cond_None},
	0xB4: {Or, 		None, 	H, 		cond_None},
	0xB5: {Or, 		None, 	L, 		cond_None},
	0xB6: {Or, 		None, 	HL, 	cond_None},
	0xB7: {Or, 		None, 	A, 		cond_None},
	0xB8: {Cp, 		None, 	B, 		cond_None},
	0xB9: {Cp, 		None, 	C, 		cond_None},
	0xBA: {Cp, 		None, 	D, 		cond_None},
	0xBB: {Cp, 		None, 	E, 		cond_None},
	0xBC: {Cp, 		None, 	H, 		cond_None},
	0xBD: {Cp, 		None, 	L, 		cond_None},
	0xBE: {Cp, 		None, 	HL, 	cond_None},
	0xBF: {Cp, 		None, 	A, 		cond_None},
	
	//0xCX
	0xC0: {Ret, 	None, 	None, 	cond_NZ},
	0xC1: {Pop, 	BC,		None,  	cond_None},
	0xC2: {Jp, 		None, 	nn, 	cond_NZ},
	0xC3: {Jp, 		None, 	nn, 	cond_None},
	0xC4: {Call, 	None, 	nn, 	cond_NZ},
	0xC5: {Push, 	None, 	BC, 	cond_None},
	0xC6: {Add, 	None, 	n, 		cond_None},
	0xC7: {Rst, 	None, 	None, 	cond_None},
	0xC8: {Ret, 	None, 	None, 	cond_Z},
	0xC9: {Ret, 	None, 	None, 	cond_None},
	0xCA: {Jp, 		None, 	nn, 	cond_Z},
	0xCC: {Call, 	None, 	nn, 	cond_Z},
	0xCD: {Call, 	None, 	nn, 	cond_None},
	0xCE: {Adc, 	None, 	n, 		cond_None},
	0xCF: {Rst, 	None, 	None, 	cond_None},
	//0xDX
	0xD0: {Ret, 	None, 	None, 	cond_NC},
	0xD1: {Pop, 	DE,  	None,  	cond_None},
	0xD2: {Jp, 		None, 	nn, 	cond_NC},
	0xD4: {Call, 	None, 	nn, 	cond_NC},
	0xD5: {Push, 	None, 	DE, 	cond_None},
	0xD6: {Sub, 	None, 	n, 		cond_None},
	0xD7: {Rst, 	None, 	None, 	cond_None},
	0xD8: {Ret, 	None, 	None, 	cond_C},
	0xD9: {Reti, 	None, 	None, 	cond_None},
	0xDA: {Jp, 		None, 	nn, 	cond_C},
	0xDC: {Jp, 		None, 	nn, 	cond_C},
	0xDE: {Sbc, 	None, 	n, 	cond_C},
	0xDF: {Rst, 	None, 	None,	cond_None},
	//0xEX
	0xE0: {Ldh, 	n_M, 	A, 		cond_None},
	0xE1: {Pop, 	HL, 	None, 	cond_None},
	0xE2: {Ldh, 	C_M, 	A, 		cond_None},
	0xE5: {Push, 	None, 	HL, 	cond_None},
	0xE6: {And, 	None, 	n, 		cond_None},
	0xE7: {Push, 	None, 	None, 	cond_None},
	0xE8: {Add16_8, SP, 	e8,		cond_None},
	0xE9: {Jp, 		None, 	HL, 	cond_None},
	0xEA: {Ld8, 	nn_M, 	A, 		cond_None},
	0xEE: {Xor, 	None, 	n, 		cond_None},
	0xEF: {Rst, 	None, 	None, 	cond_None},
	//0xFX
	0xF0: {Ldh, 	A, 		n_M, 	cond_None},
	0xF1: {Pop, 	AF, 	None,  	cond_None},
	0xF2: {Ldh, 	A, 		C_M, 	cond_None},
	0xF3: {Di, 		None, 	None, 	cond_None},
	0xF5: {Push, 	None, 	AF, 	cond_None},
	0xF6: {Or, 		None, 	n, 		cond_None},
	0xF7: {Rst, 	None, 	None, 	cond_None},
	0xFA: {Ld8, 	A, 		nn_M, 	cond_None},
	0xFE: {Cp, 		None, 	n, 		cond_None},
	0xFF: {Rst, 	None, 	None, 	cond_None},
	//TODO: More instructions
}
