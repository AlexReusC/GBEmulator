package lib

// This code are all logic for instructions but independent of opcodes. It tries to group all of the registers and constants
// on the same function (for the most part). Conditional modes are handled from the CPU

func (c *CPU) Nop() int {
	return 1
}

func (c *CPU) Jp() int {
	if c.CurrentConditionResult {
		c.Register.pc = c.Source

		if c.SourceTarget == HL {
			return 1
		}
		return 4
	}
	return 3 //Jp cc untaken
}

func (c *CPU) Jr() int {
	if c.CurrentConditionResult {
		c.Register.pc = uint16(int16(c.Register.pc) + int16(int8(c.Source)))
		return 3
	}
	return 2
}

func (c *CPU) Ld8() int {
	cycles := 1
	var input uint8

	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}

	c.SetTarget(c.DestinationTarget, uint16(input))

	if IsPointer(c.SourceTarget) || IsPointer(c.DestinationTarget) {
		cycles += 1
	}
	if c.SourceTarget == nn_M || c.DestinationTarget == nn_M {
		cycles += 2
	}
	if c.SourceTarget == n {
		cycles += 1
	}
	return cycles
}

func (c *CPU) Ld16() int {
	cycles := 1
	//Ld16 has no addresses in load
	c.SetTarget(c.DestinationTarget, c.Source)

	if c.SourceTarget == nn {
		cycles += 2
	}
	if c.DestinationTarget == SP {
		cycles += 1
	}
	if c.DestinationTarget == nn_M16 {
		cycles += 4
	}
	return cycles
}

// Ld hl, SP+e8
func (c *CPU) LdSPn() int {
	input8 := c.Source
	input16 := c.Register.sp
	result := uint16(int16(int8(uint8(input8))) + int16(input16))

	c.SetTarget(c.DestinationTarget, result)

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (input16&0x0F)+(input8&0x0F) > 0x0F)
	c.SetFlag(flagC, (input16&0xFF)+(input8&0xFF) > 0xFF)
	return 3
}

func (c *CPU) Ldh() int {
	var input uint8

	if IsPointer(c.SourceTarget) {
		input = c.MMURead(0xFF00 | uint16(c.Source))
		c.SetTarget(A, uint16(input)) //If destination is not address is always register A
	} else {
		input = uint8(c.Source)
		destinationData, _ := c.GetTarget(c.DestinationTarget)
		c.MMUWrite(0xFF00|destinationData, input)
	}

	if c.SourceTarget == C_M || c.DestinationTarget == C_M {
		return 2
	}
	return 3
}

func (c *CPU) Push() int {
	c.Register.sp -= 1
	c.MMUWrite(c.Register.sp, uint8((c.Source&0xFF00)>>8))

	c.Register.sp -= 1
	c.MMUWrite(c.Register.sp, uint8(c.Source&0xFF))

	return 4
}

func (c *CPU) Pop() int {
	lo := uint16(c.MMURead(c.Register.sp))
	hi := uint16(c.MMURead(c.Register.sp + 1))
	c.Register.sp += 2
	result := (hi << 8) | lo

	if c.DestinationTarget == AF { //clear lower nibble of F (always have to be zero)
		result = result & 0xFFF0
	}
	c.SetTarget(c.DestinationTarget, result)

	return 3
}

func (c *CPU) Call() int {
	if c.CurrentConditionResult {
		//Push pc
		c.Register.sp -= 1
		c.MMUWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
		c.Register.sp -= 1
		c.MMUWrite(c.Register.sp, uint8(c.Register.pc&0x00FF))

		//Jp nn
		c.Register.pc = c.Source

		return 6
	}

	return 3
}

func (c *CPU) Ret() int {
	if c.CurrentConditionResult {
		//Pop
		lo := uint16(c.MMURead(c.Register.sp))
		c.Register.sp += 1

		hi := uint16(c.MMURead(c.Register.sp))
		c.Register.sp += 1
		//Jp
		c.Register.pc = (hi << 8) | lo

		if c.currentOpcode == 0xc9 {
			return 4
		}
		return 5
	}

	return 2
}

func (c *CPU) Reti() int {
	//Pop
	lo := uint16(c.MMURead(c.Register.sp))
	c.Register.sp += 1

	hi := uint16(c.MMURead(c.Register.sp))
	c.Register.sp += 1
	//Jp
	c.Register.pc = (hi << 8) | lo

	c.MasterInterruptEnabled = true

	return 4
}

func (c *CPU) Rst() int {

	var rstAddress = map[uint8]uint16{
		0xC7: 0x00,
		0xD7: 0x10,
		0xE7: 0x20,
		0xF7: 0x30,
		0xCF: 0x08,
		0xDF: 0x18,
		0xEF: 0x28,
		0xFF: 0x38,
	}

	c.Register.sp -= 1
	c.MMUWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
	c.Register.sp -= 1
	c.MMUWrite(c.Register.sp, uint8(c.Register.pc&0xFF))

	c.Register.pc = (0x00 << 8) | rstAddress[c.currentOpcode]

	return 4
}

func (c *CPU) Di() int {
	c.MasterInterruptEnabled = false

	return 1
}

func (c *CPU) Ei() int {
	c.EnableMasterInterruptAfter = 2

	return 1
}

func (c *CPU) Halt() int {
	c.Halted = true
	return 0
}

// Decimal Adjust Accumulator
// https://blog.ollien.com/posts/gb-daa/
func (c *CPU) Daa() int {
	result := c.Register.a
	var offset uint8 = 0

	if c.GetFlag(flagH) || (!c.GetFlag(flagN) && (result&0x0F) > 0x09) {
		offset |= 0x06
	}

	if c.GetFlag(flagC) || (!c.GetFlag(flagN) && result > 0x99) {
		offset |= 0x60
		c.SetFlag(flagC, true)
	} else {
		c.SetFlag(flagC, false)
	}

	if c.GetFlag(flagN) {
		result -= offset
	} else {
		result += offset
	}

	c.SetTarget(A, uint16(result))
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagH, false)
	return 1
}

// Rotate left A reg
func (c *CPU) Rlca() int {
	msbOn := c.Register.a & 0x80 //128
	modifiedVal := c.Register.a << 1
	if msbOn != 0 {
		modifiedVal |= 0x1
	}
	c.SetTarget(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
	return 1
}

// Rotate right A reg
func (c *CPU) Rrca() int {
	lsbOn := c.Register.a & 0x01
	modifiedVal := c.Register.a >> 1
	if lsbOn != 0 {
		modifiedVal |= 0x80
	}
	c.SetTarget(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
	return 1
}

// Rotate left A reg, through carry
func (c *CPU) Rla() int {
	msbOn := c.Register.a & 0x80
	modifiedVal := c.Register.a << 1
	if c.GetFlag(flagC) {
		modifiedVal |= 0x01
	}
	c.SetTarget(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
	return 1
}

// Rotate right A reg, through carry
func (c *CPU) Rra() int {
	lsbOn := c.Register.a & 0x01
	modifiedVal := c.Register.a >> 1
	if c.GetFlag(flagC) {
		modifiedVal |= 0x80
	}
	c.SetTarget(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
	return 1
}

// Complement accumulator
func (c *CPU) Cpl() int {
	modifiedVal := ^c.Register.a
	c.SetTarget(A, uint16(modifiedVal))

	c.SetFlag(flagN, true)
	c.SetFlag(flagH, true)
	return 1
}

// Complement carry flag
func (c *CPU) Ccf() int {
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, !c.GetFlag(flagC))
	return 1
}

// Set Carry Flag
func (c *CPU) Scf() int {
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, true)
	return 1
}

// add uin16 or uint8
func (c *CPU) Inc() int {
	input := c.Source
	if IsPointer(c.SourceTarget) {
		input = uint16(c.MMURead(input))
	}
	result := input + 1
	c.SetTarget(c.SourceTarget, result)

	if (c.currentOpcode & 0x03) != 0x03 {
		c.SetFlag(flagZ, 0xFF&result == 0)
		c.SetFlag(flagN, false)
		c.SetFlag(flagH, (input&0x0F)+0x01 == 0x10)
	}

	if Isr8(c.SourceTarget) {
		return 1
	} else if c.SourceTarget == HL_M {
		return 3
	}
	return 2 //r16
}

func (c *CPU) Dec() int {
	input := c.Source
	if IsPointer(c.SourceTarget) {
		input = uint16(c.MMURead(input))
	}
	result := input - 1
	c.SetTarget(c.SourceTarget, result)

	if (c.currentOpcode & 0x0B) != 0x0B {
		c.SetFlag(flagZ, 0xFF&result == 0)
		c.SetFlag(flagN, true)
		c.SetFlag(flagH, (input&0x0F) == 0x00) //4 trailing zeroes
	}

	if Isr8(c.SourceTarget) {
		return 1
	} else if c.SourceTarget == HL_M {
		return 3
	}
	return 2 //r16
}

func (c *CPU) Add() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	a := c.Register.a
	result := a + input
	c.SetTarget(A, uint16(result))

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (a&0x0F)+(input&0x0F) > 0x0F)
	c.SetFlag(flagC, input > 0xFF-a)

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Add a,[HL] & Add A, n8
}

// Add Hl, SP & Add Hl, r16
func (c *CPU) AddHl() int {
	hl := c.GetTargetHL()
	input := c.Source
	result := c.GetTargetHL() + input
	c.SetTarget(HL, result)

	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (hl&0x0FFF)+(input&0x0FFF) > 0x0FFF)
	c.SetFlag(flagC, input > 0xFFFF-hl)
	return 2
}

// Add Sp, e8
func (c *CPU) Add16_8() int {
	input16 := c.Register.sp
	input8 := c.Source
	result := uint16(int16(int8(uint8(input8))) + int16(input16))
	c.SetTarget(c.DestinationTarget, result)

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (input16&0x0F)+(input8&0x0F) > 0x0F)
	c.SetFlag(flagC, (input16&0xFF)+(input8&0xFF) > 0xFF)
	return 4
}

func (c *CPU) Adc() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	a := c.Register.a
	carryBit := BoolToUint(c.GetFlag(flagC))

	result := uint16(a + input + carryBit)
	c.SetTarget(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (a&0x0F)+(input&0x0F)+carryBit > 0x0F)
	c.SetFlag(flagC, int(input) > 0xFF-int(a)-int(carryBit))
	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Adc a,[HL] & Adc A, n8
}

func (c *CPU) Sub() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	a := c.Register.a
	result := uint16(c.Register.a - input)
	c.SetTarget(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, (a&0x0F) < (input&0x0F))
	c.SetFlag(flagC, input > a)

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Sub a,[HL] & Sub A, n8
}

func (c *CPU) Sbc() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	a := c.Register.a
	carryBit := BoolToUint(c.GetFlag(flagC))

	result := uint16(a - input - carryBit)
	c.SetTarget(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, (a&0x0F) < ((input&0x0F)+carryBit))
	c.SetFlag(flagC, int(input)+int(carryBit) > int(a))

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Sbc a,[HL] & Sbc A, n8
}

func (c *CPU) And() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	result := uint16(c.Register.a & input)
	c.SetTarget(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, true)
	c.SetFlag(flagC, false)

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //And a,[HL] & And A, n8
}

func (c *CPU) Xor() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	result := uint16(c.Register.a ^ input)
	c.SetTarget(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, false)

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Xor a,[HL] & Xor A, n8
}

func (c *CPU) Or() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	result := uint16(c.Register.a | input)
	c.SetTarget(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, false)

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Or a,[HL] & Or A, n8
}

func (c *CPU) Cp() int {
	var input uint8
	if IsPointer(c.SourceTarget) {
		input = c.MMURead(c.Source)
	} else {
		input = uint8(c.Source)
	}
	a := c.Register.a
	result := uint16(a - input)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, (a&0x0F) < (input&0x0F))
	c.SetFlag(flagC, input > a)

	if Isr8(c.SourceTarget) {
		return 1
	}
	return 2 //Cp a,[HL] & Cp A, n8
}

// Rotate Regiser Left
func (c *CPU) Rlc(input uint16, t target) int {
	msbOn := input & 0x80 //128
	result := input << 1
	if msbOn != 0 {
		result |= 0x1
	}
	c.SetTarget(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Rlc [hl]
}

// Rotate Regiser Right
func (c *CPU) Rrc(input uint16, t target) int {
	lsbOn := input & 0x01
	result := input >> 1
	if lsbOn != 0 {
		result |= 0x80
	}
	c.SetTarget(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Rrc [hl]
}

// Rotate r Left, through carry
func (c *CPU) Rl(input uint16, t target) int {
	msbOn := input & 0x80
	result := input << 1
	if c.GetFlag(flagC) {
		result |= 0x01
	}
	c.SetTarget(t, result)

	c.SetFlag(flagZ, (result&0xFF) == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Rl [hl]
}

// Rotate r right, through carry
func (c *CPU) Rr(input uint16, t target) int {
	lsbOn := input & 0x01
	result := input >> 1
	if c.GetFlag(flagC) {
		result |= 0x80
	}
	c.SetTarget(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Rr [hl]
}

// Shift Left Arithmetically
func (c *CPU) Sla(input uint16, t target) int {
	msbOn := input & 0x80
	result := input << 1
	c.SetTarget(t, result)

	c.SetFlag(flagZ, (result&0xFF) == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Sla [hl]
}

// Shift Right Arithmetically
func (c *CPU) Sra(input uint16, t target) int {
	msbOn := input & 0x80
	lsbOn := input & 0x01
	result := msbOn | (input >> 1)
	c.SetTarget(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Rrc [hl]
}

func (c *CPU) Swap(input uint16, t target) int {
	result := ((input & 0x0F) << 4) | ((input & 0xF0) >> 4)
	c.SetTarget(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, false)

	if Isr8(t) {
		return 2
	}
	return 4 //Swap [hl]
}

// Shift Right Logically
func (c *CPU) Srl(input uint16, t target) int {
	lsbOn := input & 0x01
	result := input >> 1
	c.SetTarget(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)

	if Isr8(t) {
		return 2
	}
	return 4 //Rrc [hl]
}

func (c *CPU) Bit(input uint16, t target, b uint8) int {
	c.SetFlag(flagZ, input&(1<<b) == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, true)

	if Isr8(t) {
		return 2
	}
	return 3 //Bit [hl]
}

func (c *CPU) Res(input uint16, t target, b uint8) int {
	result := input & ^(0x01 << b)
	c.SetTarget(t, result)

	if Isr8(t) {
		return 2
	}
	return 4 //Res [hl]
}

func (c *CPU) Set(input uint16, t target, b uint8) int {
	result := input | (0x01 << b)
	c.SetTarget(t, result)

	if Isr8(t) {
		return 2
	}
	return 4 //Rrc [hl]
}

func (c *CPU) Cb() (int, error) {
	cycles := 0
	op := uint8(c.Source)
	instruction := cbOpcodes[op]

	data, err := c.GetTarget(instruction.Register)
	if err != nil {
		return 0, err
	}

	input := data
	if IsPointer(c.SourceTarget) {
		input = uint16(c.MMURead(data))
	}

	switch instruction.Instruction {
	case Rlc:
		cycles += c.Rlc(input, instruction.Register)
	case Rrc:
		cycles += c.Rrc(input, instruction.Register)
	case Rl:
		cycles += c.Rl(input, instruction.Register)
	case Rr:
		cycles += c.Rr(input, instruction.Register)
	case Sla:
		cycles += c.Sla(input, instruction.Register)
	case Sra:
		cycles += c.Sra(input, instruction.Register)
	case Swap:
		cycles += c.Swap(input, instruction.Register)
	case Srl:
		cycles += c.Srl(input, instruction.Register)
	case Bit:
		cycles += c.Bit(input, instruction.Register, instruction.Bit)
	case Res:
		cycles += c.Res(input, instruction.Register, instruction.Bit)
	case Set:
		cycles += c.Set(input, instruction.Register, instruction.Bit)
	}

	return cycles, nil
}
