package lib

type procedure string
type procedureCb string
type conditional string
type target string

const (
	Nop     procedure = "Nop"
	Jp      procedure = "Jp"
	Jr      procedure = "Jr"
	Di      procedure = "Di"
	Ld8     procedure = "Ld8"
	Ld16    procedure = "Ld16"
	Ldh     procedure = "Ldh"
	Push    procedure = "Push"
	Pop     procedure = "Pop"
	Call    procedure = "Call"
	Ret     procedure = "Ret"
	Reti    procedure = "Reti"
	Rst     procedure = "Rst"
	Inc     procedure = "Inc"
	Dec     procedure = "Dec"
	Add     procedure = "Add"
	AddHl   procedure = "AddHl"
	Add16_8 procedure = "Add16_8"
	Adc     procedure = "Adc"
	Sub     procedure = "Sub"
	Sbc     procedure = "Sbc"
	And     procedure = "And"
	Xor     procedure = "Xor"
	Or      procedure = "Or"
	Cp      procedure = "Cp"
	Cb      procedure = "Cb"
)

const (
	Rlc  procedureCb = "Rlc"
	Rrc  procedureCb = "Rrc"
	Rl   procedureCb = "Rl"
	Rr   procedureCb = "Rr"
	Sla  procedureCb = "Sla"
	Sra  procedureCb = "Sra"
	Swap procedureCb = "Swap"
	Srl  procedureCb = "Srl"
	Bit  procedureCb = "Bit"
	Res  procedureCb = "Res"
	Set  procedureCb = "Set"
)

const (
	A     target = "A"
	B     target = "B"
	C     target = "C"
	D     target = "D"
	E     target = "E"
	F     target = "F"
	H     target = "H"
	L     target = "L"
	AF    target = "AF"
	BC    target = "BC"
	DE    target = "DE"
	HL    target = "HL"
	SP    target = "SP"
	e8    target = "e8"
	n     target = "n"
	nn    target = "nn"
	C_M   target = "(C)"
	BC_M  target = "(BC)"
	DE_M  target = "(DE)"
	HL_M  target = "(HL)"
	HLP_M target = "(HL+)"
	HLM_M target = "(HL-)"
	n_M   target = "(n)"
	nn_M  target = "(nn)"
	None  target = "none"
)

const (
	cond_None conditional = "None"
	cond_C    conditional = "C"
	cond_NC   conditional = "NC"
	cond_Z    conditional = "Z"
	cond_NZ   conditional = "NZ"
)

type Instruction struct {
	InstructionType procedure
	Destination     target
	Source          target
	ConditionType   conditional
}

type CbOpcode struct {
	Instruction procedureCb
	Register    target
	Bit         uint8
}

var instructions = map[uint8]Instruction{
	// 0x0X
	0x00: {Nop, None, None, cond_None},
	0x01: {Ld16, BC, nn, cond_None},
	0x02: {Ld8, BC_M, A, cond_None},
	0x03: {Inc, None, BC, cond_None},
	0x04: {Inc, None, B, cond_None},
	0x05: {Dec, None, B, cond_None},
	0x06: {Ld8, B, n, cond_None},
	0x08: {Ld16, nn_M, SP, cond_None},
	0x09: {AddHl, HL, BC, cond_None},
	0x0A: {Ld8, A, BC_M, cond_None},
	0x0B: {Dec, None, BC, cond_None},
	0x0C: {Inc, None, C, cond_None},
	0x0D: {Dec, None, C, cond_None},
	0x0E: {Ld8, C, n, cond_None},
	// 0x1X
	0x11: {Ld16, DE, nn, cond_None},
	0x12: {Ld8, DE_M, A, cond_None},
	0x13: {Inc, None, DE, cond_None},
	0x14: {Inc, None, D, cond_None},
	0x15: {Dec, None, D, cond_None},
	0x16: {Ld8, D, n, cond_None},
	0x18: {Jr, None, n, cond_None},
	0x19: {AddHl, HL, DE, cond_None},
	0x1A: {Ld8, A, DE_M, cond_None},
	0x1B: {Dec, None, DE, cond_None},
	0x1C: {Inc, None, E, cond_None},
	0x1D: {Dec, None, E, cond_None},
	0x1E: {Ld8, E, n, cond_None},
	// 0x2X
	0x20: {Jr, None, n, cond_NZ},
	0x21: {Ld16, HL, nn, cond_None},
	0x22: {Ld8, HLP_M, A, cond_None},
	0x23: {Inc, None, HL, cond_None},
	0x24: {Inc, None, H, cond_None},
	0x25: {Dec, None, H, cond_None},
	0x26: {Ld8, H, n, cond_None},
	0x28: {Jr, None, n, cond_Z},
	0x29: {AddHl, HL, HL, cond_None},
	0x2A: {Ld8, A, HLP_M, cond_None},
	0x2B: {Dec, None, HL, cond_None},
	0x2C: {Inc, None, L, cond_None},
	0x2D: {Dec, None, L, cond_None},
	0x2E: {Ld8, L, n, cond_None},
	// 0x3X
	0x30: {Jr, None, n, cond_NC},
	0x31: {Ld16, SP, nn, cond_None},
	0x32: {Ld8, HLM_M, A, cond_None},
	0x33: {Inc, None, SP, cond_None},
	0x34: {Inc, None, HL_M, cond_None},
	0x35: {Dec, None, HL_M, cond_None},
	0x36: {Ld8, HL_M, n, cond_None},
	0x38: {Jr, None, n, cond_C},
	0x39: {AddHl, HL, SP, cond_None},
	0x3A: {Ld8, A, HLM_M, cond_None},
	0x3B: {Dec, None, SP, cond_None},
	0x3C: {Inc, None, A, cond_None},
	0x3D: {Dec, None, A, cond_None},
	0x3E: {Ld8, A, n, cond_None},
	//0x4X
	0x41: {Ld8, B, B, cond_None},
	0x42: {Ld8, B, C, cond_None},
	0x43: {Ld8, B, D, cond_None},
	0x44: {Ld8, B, E, cond_None},
	0x45: {Ld8, B, H, cond_None},
	0x46: {Ld8, B, HL_M, cond_None},
	0x47: {Ld8, B, A, cond_None},
	0x48: {Ld8, C, B, cond_None},
	0x49: {Ld8, C, C, cond_None},
	0x4A: {Ld8, C, D, cond_None},
	0x4B: {Ld8, C, E, cond_None},
	0x4C: {Ld8, C, H, cond_None},
	0x4D: {Ld8, C, L, cond_None},
	0x4E: {Ld8, C, HL_M, cond_None},
	0x4F: {Ld8, C, A, cond_None},
	//0x5xX
	0x50: {Ld8, D, B, cond_None},
	0x51: {Ld8, D, C, cond_None},
	0x52: {Ld8, D, D, cond_None},
	0x53: {Ld8, D, E, cond_None},
	0x54: {Ld8, D, H, cond_None},
	0x55: {Ld8, D, L, cond_None},
	0x56: {Ld8, D, HL_M, cond_None},
	0x57: {Ld8, D, A, cond_None},
	0x58: {Ld8, E, B, cond_None},
	0x59: {Ld8, E, C, cond_None},
	0x5A: {Ld8, E, D, cond_None},
	0x5B: {Ld8, E, E, cond_None},
	0x5C: {Ld8, E, H, cond_None},
	0x5D: {Ld8, E, L, cond_None},
	0x5E: {Ld8, E, HL_M, cond_None},
	0x5F: {Ld8, E, A, cond_None},
	//0x6X
	0x60: {Ld8, H, B, cond_None},
	0x61: {Ld8, H, C, cond_None},
	0x62: {Ld8, H, D, cond_None},
	0x63: {Ld8, H, E, cond_None},
	0x64: {Ld8, H, H, cond_None},
	0x65: {Ld8, H, L, cond_None},
	0x66: {Ld8, H, HL_M, cond_None},
	0x67: {Ld8, H, A, cond_None},
	0x68: {Ld8, L, B, cond_None},
	0x69: {Ld8, L, C, cond_None},
	0x6A: {Ld8, L, D, cond_None},
	0x6B: {Ld8, L, E, cond_None},
	0x6C: {Ld8, L, H, cond_None},
	0x6D: {Ld8, L, L, cond_None},
	0x6E: {Ld8, L, HL_M, cond_None},
	0x6F: {Ld8, L, A, cond_None},
	//0x7X
	0x70: {Ld8, HL_M, B, cond_None},
	0x71: {Ld8, HL_M, C, cond_None},
	0x72: {Ld8, HL_M, D, cond_None},
	0x73: {Ld8, HL_M, E, cond_None},
	0x74: {Ld8, HL_M, H, cond_None},
	0x75: {Ld8, HL_M, L, cond_None},
	//HALT instruction
	0x77: {Ld8, HL_M, A, cond_None},
	0x78: {Ld8, A, B, cond_None},
	0x79: {Ld8, A, C, cond_None},
	0x7A: {Ld8, A, D, cond_None},
	0x7B: {Ld8, A, E, cond_None},
	0x7C: {Ld8, A, H, cond_None},
	0x7D: {Ld8, A, L, cond_None},
	0x7E: {Ld8, A, HL_M, cond_None},
	0x7F: {Ld8, A, A, cond_None},

	//0x8X
	0x80: {Add, None, B, cond_None},
	0x81: {Add, None, C, cond_None},
	0x82: {Add, None, D, cond_None},
	0x83: {Add, None, E, cond_None},
	0x84: {Add, None, H, cond_None},
	0x85: {Add, None, L, cond_None},
	0x86: {Add, None, HL_M, cond_None},
	0x87: {Add, None, A, cond_None},
	0x88: {Adc, None, B, cond_None},
	0x89: {Adc, None, C, cond_None},
	0x8A: {Adc, None, D, cond_None},
	0x8B: {Adc, None, E, cond_None},
	0x8C: {Adc, None, H, cond_None},
	0x8D: {Adc, None, L, cond_None},
	0x8E: {Adc, None, HL_M, cond_None},
	0x8F: {Adc, None, A, cond_None},
	//0x9X
	0x90: {Sub, None, B, cond_None},
	0x91: {Sub, None, C, cond_None},
	0x92: {Sub, None, D, cond_None},
	0x93: {Sub, None, E, cond_None},
	0x94: {Sub, None, H, cond_None},
	0x95: {Sub, None, L, cond_None},
	0x96: {Sub, None, HL_M, cond_None},
	0x97: {Sub, None, A, cond_None},
	0x98: {Sbc, None, B, cond_None},
	0x99: {Sbc, None, C, cond_None},
	0x9A: {Sbc, None, D, cond_None},
	0x9B: {Sbc, None, E, cond_None},
	0x9C: {Sbc, None, H, cond_None},
	0x9D: {Sbc, None, L, cond_None},
	0x9E: {Sbc, None, HL_M, cond_None},
	0x9F: {Sbc, None, A, cond_None},
	//0xAX
	0xA0: {And, None, B, cond_None},
	0xA1: {And, None, C, cond_None},
	0xA2: {And, None, D, cond_None},
	0xA3: {And, None, E, cond_None},
	0xA4: {And, None, H, cond_None},
	0xA5: {And, None, L, cond_None},
	0xA6: {And, None, HL_M, cond_None},
	0xA7: {And, None, A, cond_None},
	0xA8: {Xor, None, B, cond_None},
	0xA9: {Xor, None, C, cond_None},
	0xAA: {Xor, None, D, cond_None},
	0xAB: {Xor, None, E, cond_None},
	0xAC: {Xor, None, H, cond_None},
	0xAD: {Xor, None, L, cond_None},
	0xAE: {Xor, None, HL_M, cond_None},
	0xAF: {Xor, None, A, cond_None},
	//0xBX
	0xB0: {Or, None, B, cond_None},
	0xB1: {Or, None, C, cond_None},
	0xB2: {Or, None, D, cond_None},
	0xB3: {Or, None, E, cond_None},
	0xB4: {Or, None, H, cond_None},
	0xB5: {Or, None, L, cond_None},
	0xB6: {Or, None, HL, cond_None},
	0xB7: {Or, None, A, cond_None},
	0xB8: {Cp, None, B, cond_None},
	0xB9: {Cp, None, C, cond_None},
	0xBA: {Cp, None, D, cond_None},
	0xBB: {Cp, None, E, cond_None},
	0xBC: {Cp, None, H, cond_None},
	0xBD: {Cp, None, L, cond_None},
	0xBE: {Cp, None, HL, cond_None},
	0xBF: {Cp, None, A, cond_None},

	//0xCX
	0xC0: {Ret, None, None, cond_NZ},
	0xC1: {Pop, BC, None, cond_None},
	0xC2: {Jp, None, nn, cond_NZ},
	0xC3: {Jp, None, nn, cond_None},
	0xC4: {Call, None, nn, cond_NZ},
	0xC5: {Push, None, BC, cond_None},
	0xC6: {Add, None, n, cond_None},
	0xC7: {Rst, None, None, cond_None},
	0xC8: {Ret, None, None, cond_Z},
	0xC9: {Ret, None, None, cond_None},
	0xCA: {Jp, None, nn, cond_Z},
	0xCC: {Call, None, nn, cond_Z},
	0xCD: {Call, None, nn, cond_None},
	0xCE: {Adc, None, n, cond_None},
	0xCF: {Rst, None, None, cond_None},
	//0xDX
	0xD0: {Ret, None, None, cond_NC},
	0xD1: {Pop, DE, None, cond_None},
	0xD2: {Jp, None, nn, cond_NC},
	0xD4: {Call, None, nn, cond_NC},
	0xD5: {Push, None, DE, cond_None},
	0xD6: {Sub, None, n, cond_None},
	0xD7: {Rst, None, None, cond_None},
	0xD8: {Ret, None, None, cond_C},
	0xD9: {Reti, None, None, cond_None},
	0xDA: {Jp, None, nn, cond_C},
	0xDC: {Jp, None, nn, cond_C},
	0xDE: {Sbc, None, n, cond_C},
	0xDF: {Rst, None, None, cond_None},
	//0xEX
	0xE0: {Ldh, n_M, A, cond_None},
	0xE1: {Pop, HL, None, cond_None},
	0xE2: {Ldh, C_M, A, cond_None},
	0xE5: {Push, None, HL, cond_None},
	0xE6: {And, None, n, cond_None},
	0xE7: {Push, None, None, cond_None},
	0xE8: {Add16_8, SP, e8, cond_None},
	0xE9: {Jp, None, HL, cond_None},
	0xEA: {Ld8, nn_M, A, cond_None},
	0xEE: {Xor, None, n, cond_None},
	0xEF: {Rst, None, None, cond_None},
	//0xFX
	0xF0: {Ldh, A, n_M, cond_None},
	0xF1: {Pop, AF, None, cond_None},
	0xF2: {Ldh, A, C_M, cond_None},
	0xF3: {Di, None, None, cond_None},
	0xF5: {Push, None, AF, cond_None},
	0xF6: {Or, None, n, cond_None},
	0xF7: {Rst, None, None, cond_None},
	0xFA: {Ld8, A, nn_M, cond_None},
	0xFE: {Cp, None, n, cond_None},
	0xFF: {Rst, None, None, cond_None},
	//TODO: More instructions
}

var cbOpcodes = map[uint8]CbOpcode{
	//0x0X
	0x00: {Rlc, B, 0},
	0x01: {Rlc, C, 0},
	0x02: {Rlc, D, 0},
	0x03: {Rlc, E, 0},
	0x04: {Rlc, H, 0},
	0x05: {Rlc, L, 0},
	0x06: {Rlc, HL_M, 0},
	0x07: {Rlc, A, 0},
	0x08: {Rrc, B, 0},
	0x09: {Rrc, C, 0},
	0x0A: {Rrc, D, 0},
	0x0B: {Rrc, E, 0},
	0x0C: {Rrc, H, 0},
	0x0D: {Rrc, L, 0},
	0x0E: {Rrc, HL_M, 0},
	0x0F: {Rrc, A, 0},
	// 0x1X
	0x10: {Rl, B, 0},
	0x11: {Rl, C, 0},
	0x12: {Rl, D, 0},
	0x13: {Rl, E, 0},
	0x14: {Rl, H, 0},
	0x15: {Rl, L, 0},
	0x16: {Rl, HL_M, 0},
	0x17: {Rl, A, 0},
	0x18: {Rr, B, 0},
	0x19: {Rr, C, 0},
	0x1A: {Rr, D, 0},
	0x1B: {Rr, E, 0},
	0x1C: {Rr, H, 0},
	0x1D: {Rr, L, 0},
	0x1F: {Rr, HL_M, 0},
	0x1E: {Rr, A, 0},
	// 0x
	0x20: {Sla, B, 0},
	0x21: {Sla, C, 0},
	0x22: {Sla, D, 0},
	0x23: {Sla, E, 0},
	0x24: {Sla, H, 0},
	0x25: {Sla, L, 0},
	0x26: {Sla, HL_M, 0},
	0x27: {Sla, A, 0},
	0x28: {Sra, B, 0},
	0x29: {Sra, C, 0},
	0x2A: {Sra, D, 0},
	0x2B: {Sra, E, 0},
	0x2C: {Sra, H, 0},
	0x2D: {Sra, L, 0},
	0x2E: {Sra, HL_M, 0},
	0x2F: {Sra, A, 0},
	// 0x
	0x30: {Swap, B, 0},
	0x31: {Swap, C, 0},
	0x32: {Swap, D, 0},
	0x33: {Swap, E, 0},
	0x34: {Swap, H, 0},
	0x35: {Swap, L, 0},
	0x36: {Swap, HL_M, 0},
	0x37: {Swap, A, 0},
	0x38: {Srl, B, 0},
	0x39: {Srl, C, 0},
	0x3A: {Srl, D, 0},
	0x3B: {Srl, E, 0},
	0x3C: {Srl, H, 0},
	0x3D: {Srl, L, 0},
	0x3E: {Srl, HL_M, 0},
	0x3F: {Srl, A, 0},
	//0x4
	0x40: {Bit, B, 0},
	0x41: {Bit, C, 0},
	0x42: {Bit, D, 0},
	0x43: {Bit, E, 0},
	0x44: {Bit, H, 0},
	0x45: {Bit, L, 0},
	0x46: {Bit, HL_M, 0},
	0x47: {Bit, A, 0},
	0x48: {Bit, B, 1},
	0x49: {Bit, C, 1},
	0x4A: {Bit, D, 1},
	0x4B: {Bit, E, 1},
	0x4C: {Bit, H, 1},
	0x4D: {Bit, L, 1},
	0x4E: {Bit, HL_M, 1},
	0x4F: {Bit, A, 1},
	//0x5
	0x50: {Bit, B, 2},
	0x51: {Bit, C, 2},
	0x52: {Bit, D, 2},
	0x53: {Bit, E, 2},
	0x54: {Bit, H, 2},
	0x55: {Bit, L, 2},
	0x56: {Bit, HL_M, 2},
	0x57: {Bit, A, 2},
	0x58: {Bit, B, 3},
	0x59: {Bit, C, 3},
	0x5A: {Bit, D, 3},
	0x5B: {Bit, E, 3},
	0x5C: {Bit, H, 3},
	0x5D: {Bit, L, 3},
	0x5E: {Bit, HL_M, 3},
	0x5F: {Bit, A, 3},
	//0x6
	0x60: {Bit, B, 4},
	0x61: {Bit, C, 4},
	0x62: {Bit, D, 4},
	0x63: {Bit, E, 4},
	0x64: {Bit, H, 4},
	0x65: {Bit, L, 4},
	0x66: {Bit, HL_M, 4},
	0x67: {Bit, A, 4},
	0x68: {Bit, B, 5},
	0x69: {Bit, C, 5},
	0x6A: {Bit, D, 5},
	0x6B: {Bit, E, 5},
	0x6C: {Bit, H, 5},
	0x6D: {Bit, L, 5},
	0x6E: {Bit, HL_M, 5},
	0x6F: {Bit, A, 5},
	//0x7
	0x70: {Bit, B, 6},
	0x71: {Bit, C, 6},
	0x72: {Bit, D, 6},
	0x73: {Bit, E, 6},
	0x74: {Bit, H, 6},
	0x75: {Bit, L, 6},
	0x76: {Bit, HL_M, 6},
	0x77: {Bit, A, 6},
	0x78: {Bit, B, 7},
	0x79: {Bit, C, 7},
	0x7A: {Bit, D, 7},
	0x7B: {Bit, E, 7},
	0x7C: {Bit, H, 7},
	0x7D: {Bit, L, 7},
	0x7E: {Bit, HL_M, 7},
	0x7F: {Bit, A, 7},

	//0x8
	0x80: {Res, B, 0},
	0x81: {Res, C, 0},
	0x82: {Res, D, 0},
	0x83: {Res, E, 0},
	0x84: {Res, H, 0},
	0x85: {Res, L, 0},
	0x86: {Res, HL_M, 0},
	0x87: {Res, A, 0},
	0x88: {Res, B, 1},
	0x89: {Res, C, 1},
	0x8A: {Res, D, 1},
	0x8B: {Res, E, 1},
	0x8C: {Res, H, 1},
	0x8D: {Res, L, 1},
	0x8E: {Res, HL_M, 1},
	0x8F: {Res, A, 1},
	//0x9
	0x90: {Res, B, 2},
	0x91: {Res, C, 2},
	0x92: {Res, D, 2},
	0x93: {Res, E, 2},
	0x94: {Res, H, 2},
	0x95: {Res, L, 2},
	0x96: {Res, HL_M, 2},
	0x97: {Res, A, 2},
	0x98: {Res, B, 3},
	0x99: {Res, C, 3},
	0x9A: {Res, D, 3},
	0x9B: {Res, E, 3},
	0x9C: {Res, H, 3},
	0x9D: {Res, L, 3},
	0x9E: {Res, HL_M, 3},
	0x9F: {Res, A, 3},
	//0xA
	0xA0: {Res, B, 4},
	0xA1: {Res, C, 4},
	0xA2: {Res, D, 4},
	0xA3: {Res, E, 4},
	0xA4: {Res, H, 4},
	0xA5: {Res, L, 4},
	0xA6: {Res, HL_M, 4},
	0xA7: {Res, A, 4},
	0xA8: {Res, B, 5},
	0xA9: {Res, C, 5},
	0xAA: {Res, D, 5},
	0xAB: {Res, E, 5},
	0xAC: {Res, H, 5},
	0xAD: {Res, L, 5},
	0xAE: {Res, HL_M, 5},
	0xAF: {Res, A, 5},
	//0xB
	0xB0: {Res, B, 6},
	0xB1: {Res, C, 6},
	0xB2: {Res, D, 6},
	0xB3: {Res, E, 6},
	0xB4: {Res, H, 6},
	0xB5: {Res, L, 6},
	0xB6: {Res, HL_M, 6},
	0xB7: {Res, A, 6},
	0xB8: {Res, B, 7},
	0xB9: {Res, C, 7},
	0xBA: {Res, D, 7},
	0xBB: {Res, E, 7},
	0xBC: {Res, H, 7},
	0xBD: {Res, L, 7},
	0xBE: {Res, HL_M, 7},
	0xBF: {Res, A, 7},

	//0xC
	0xC0: {Set, B, 0},
	0xC1: {Set, C, 0},
	0xC2: {Set, D, 0},
	0xC3: {Set, E, 0},
	0xC4: {Set, H, 0},
	0xC5: {Set, L, 0},
	0xC6: {Set, HL_M, 0},
	0xC7: {Set, A, 0},
	0xC8: {Set, B, 1},
	0xC9: {Set, C, 1},
	0xCA: {Set, D, 1},
	0xCB: {Set, E, 1},
	0xCC: {Set, H, 1},
	0xCD: {Set, L, 1},
	0xCE: {Set, HL_M, 1},
	0xCF: {Set, A, 1},

	//0xD
	0xD0: {Set, B, 2},
	0xD1: {Set, C, 2},
	0xD2: {Set, D, 2},
	0xD3: {Set, E, 2},
	0xD4: {Set, H, 2},
	0xD5: {Set, L, 2},
	0xD6: {Set, HL_M, 2},
	0xD7: {Set, A, 2},
	0xD8: {Set, B, 3},
	0xD9: {Set, C, 3},
	0xDA: {Set, D, 3},
	0xDB: {Set, E, 3},
	0xDC: {Set, H, 3},
	0xDD: {Set, L, 3},
	0xDE: {Set, HL_M, 3},
	0xDF: {Set, A, 3},
	//0xE
	0xE0: {Set, B, 4},
	0xE1: {Set, C, 4},
	0xE2: {Set, D, 4},
	0xE3: {Set, E, 4},
	0xE4: {Set, H, 4},
	0xE5: {Set, L, 4},
	0xE6: {Set, HL_M, 4},
	0xE7: {Set, A, 4},
	0xE8: {Set, B, 5},
	0xE9: {Set, C, 5},
	0xEA: {Set, D, 5},
	0xEB: {Set, E, 5},
	0xEC: {Set, H, 5},
	0xED: {Set, L, 5},
	0xEE: {Set, HL_M, 5},
	0xEF: {Set, A, 5},
	//0xF
	0xF0: {Set, B, 6},
	0xF1: {Set, C, 6},
	0xF2: {Set, D, 6},
	0xF3: {Set, E, 6},
	0xF4: {Set, H, 6},
	0xF5: {Set, L, 6},
	0xF6: {Set, HL_M, 6},
	0xF7: {Set, A, 6},
	0xF8: {Set, B, 7},
	0xF9: {Set, C, 7},
	0xFA: {Set, D, 7},
	0xFB: {Set, E, 7},
	0xFC: {Set, H, 7},
	0xFD: {Set, L, 7},
	0xFE: {Set, HL_M, 7},
	0xFF: {Set, A, 7},
}