package lib

//TODO: clock behavior

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

func (c *CPU) Nop() {
}

func (c *CPU) Jp() {
	if c.CurrentConditionResult {
		c.Register.pc = c.Source.Value
	}
}

func (c *CPU) Jr() {
	if c.CurrentConditionResult {
		c.Register.pc = c.Register.pc + c.Source.Value
	}
}

func (c *CPU) Ld8() {
	var input uint8

	if c.Source.IsAddr {
		input = c.BusRead(c.Source.Value)
	} else {
		input = uint8(c.Source.Value)
	}

	if c.Destination.IsAddr {
		c.BusWrite(c.Destination.Value, input)
	} else {
		c.SetRegister(c.DestinationTarget, uint16(input))
	}

	//TODO: (HL) SP+e instruction
}

func (c *CPU) Ld16() {
	//Ld16 has no addresses in load
	if c.Destination.IsAddr {
		c.BusWrite16(c.Destination.Value, c.Source.Value)
	} else {
		c.SetRegister(c.DestinationTarget, c.Source.Value)
	}
}

func (c *CPU) Ldh() {
	var input uint8

	if c.Source.IsAddr {
		input = c.BusRead(0xFF00 | uint16(c.BusRead(c.Source.Value)))
	} else {
		input = uint8(c.Source.Value)
	}

	if c.Destination.IsAddr {
		c.BusWrite(0xFF00|c.Destination.Value, input)
	} else {
		c.SetRegister(A, uint16(input)) //If destination is not address is always register A
	}
}

func (c *CPU) Push() {
	c.Register.sp -= 1
	c.BusWrite(c.Register.sp, uint8((c.Source.Value&0xFF00)>>8))

	c.Register.sp -= 1
	c.BusWrite(c.Register.sp, uint8(c.Source.Value&0xFF))
}

func (c *CPU) Pop() {
	lo := uint16(c.BusRead(c.Register.sp))
	hi := uint16(c.BusRead(c.Register.sp+1))
	c.Register.sp += 2
	result := (hi<<8)|lo

	c.SetRegister(c.DestinationTarget, result)
}

func (c *CPU) Call() {
	if c.CurrentConditionResult {
		//Push pc
		c.Register.sp -= 1
		c.BusWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
		c.Register.sp -= 1
		c.BusWrite(c.Register.sp, uint8(c.Register.pc&0x00FF))

		//Jp nn
		c.Register.pc = c.Source.Value
	}
}

func (c *CPU) Ret() {
	if c.CurrentConditionResult {
		//Pop
		lo := uint16(c.BusRead(c.Register.sp))
		c.Register.sp += 1

		hi := uint16(c.BusRead(c.Register.sp))
		c.Register.sp += 1
		//Jp
		c.Register.pc = (hi << 8) | lo
	}
}

func (c *CPU) Reti() {
	//Pop
	lo := uint16(c.BusRead(c.Register.sp))
	c.Register.sp += 1

	hi := uint16(c.BusRead(c.Register.sp))
	c.Register.sp += 1
	//Jp
	c.Register.pc = (hi << 8) | lo

	c.InterruptorMasterEnabled = true
}

func (c *CPU) Rst() {
	c.Register.sp -= 1
	c.BusWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
	c.Register.sp -= 1
	c.BusWrite(c.Register.sp, uint8(c.Register.pc&0xFF))

	c.Register.pc = (0x00 << 8) | rstAddress[c.currentOpcode]
}

func (c *CPU) Di() {
	c.InterruptorMasterEnabled = false
}

func (c *CPU) Ei() {
	//TODO: Implement this when adding threads
}

//Decimal Adjust Accumulator
//https://blog.ollien.com/posts/gb-daa/
func (c *CPU) Daa() {
	modifiedVal := c.Register.a
	if (c.GetFlag(flagN)){ //after substraction
		if (c.GetFlag(flagC) || modifiedVal > 0x99){
			modifiedVal += 0x60
			c.SetFlag(flagC, true)
		}
		if (c.GetFlag(flagH) || (modifiedVal & 0x0F) > 0x09){
			modifiedVal += 0x06
		}
	} else { //after addition
		if (c.GetFlag(flagC)){
			modifiedVal -= 0x60
		}
		if (c.GetFlag(flagH)){
			modifiedVal -= 0x6
		}
	}
	c.SetRegister(A, uint16(modifiedVal))

	c.SetFlag(flagZ, modifiedVal == 0)
	c.SetFlag(flagH, false)
}

//Rotate left A reg
func (c *CPU) Rlca() {
	msbOn := c.Register.a & 0x80 //128
	modifiedVal := c.Register.a << 1
	if(msbOn != 0){
		 modifiedVal |= 0x1
	}
	c.SetRegister(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
}

//Rotate right A reg
func (c *CPU) Rrca() {
	lsbOn := c.Register.a & 0x01
	modifiedVal := c.Register.a >> 1
	if(lsbOn != 0){
		modifiedVal |= 0x80
	}
	c.SetRegister(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
}

//Rotate left A reg, through carry
func (c *CPU) Rla() {
	msbOn := c.Register.a & 0x80
	modifiedVal := c.Register.a << 1
	if (c.GetFlag(flagC)){
		modifiedVal |= 0x01
	}
	c.SetRegister(A, uint16(modifiedVal))

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
}

//Rotate right A reg, through carry
func (c *CPU) Rra() {
	lsbOn := c.Register.a & 0x01
	modifiedVal := c.Register.a >> 1
	if (c.GetFlag(flagC)){
		modifiedVal |= 0x80
	}
	c.SetRegister(A, uint16(modifiedVal))
	
	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
}

//Complement accumulator
func (c *CPU) Cpl() {
	modifiedVal := ^c.Register.a
	c.SetRegister(A, uint16(modifiedVal))
	
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, true)
}

//Complement carry flag
func (c *CPU) Ccf() {
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, !c.GetFlag(flagC))
}

func (c *CPU) Scf() {
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, true)
}

func (c *CPU) Inc() {
	input := c.Source.Value
	newVal := input + 1
	c.SetRegister(c.SourceTarget, newVal)

	if ((c.currentOpcode & 0x03) == 0){
		c.SetFlag(flagZ, newVal == 0)
		c.SetFlag(flagN, false)
		c.SetFlag(flagH, (input & 0x0F) + 0x01 == 0x10 )
	}
}

func (c *CPU) Dec() {
	input := c.Source.Value
	newVal := input - 1
	c.SetRegister(c.SourceTarget, newVal)

	if ((c.currentOpcode & 0x0B) == 0){
		c.SetFlag(flagZ, newVal == 0)
		c.SetFlag(flagN, true)
		c.SetFlag(flagH, ^(input & 0x0F) == 0x0F ) //4 trailing zeroes
	}
}

func (c *CPU) Add() {
	input := uint8(c.Source.Value)
	result := c.Register.a + input
	c.SetRegister(A, uint16(result))

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (c.Register.a & 0x0F) + (input & 0x0F) > 0x0F)
	c.SetFlag(flagC, input > 0xFF - c.Register.a)
}

func (c *CPU) AddHl() {
	input := c.Source.Value	
	result := c.GetTargetHL() + input
	c.SetRegister(HL, result)

	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (c.GetTargetHL() & 0x0FFF) + (input & 0x0FFF) > 0x0FFF)
	c.SetFlag(flagC, input > 0xFFFF - c.GetTargetHL())
}

func (c *CPU) Add16_8() {
	input16 := c.Destination.Value
	input8 := c.Source.Value
	result := input8 + input16
	c.SetRegister(c.DestinationTarget, result)

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (input16 & 0x0F) + (input8 & 0x0F) > 0x0F)
	c.SetFlag(flagC, (input16 & 0xFF) + (input8 & 0xFF) > 0xFF)
}

func (c *CPU) Adc() {
	input := uint8(c.Source.Value)
	carryBit := BoolToUint(c.GetFlag(flagC))

	result := uint16(c.Register.a + input + carryBit)
	c.SetRegister(A, result)

	c.SetFlag(flagZ, false)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, (c.Register.a & 0x0F) + (input & 0x0F) > 0x0F)
	c.SetFlag(flagC, input > 0xFF - c.Register.a)
}

func (c *CPU) Sub() {
	input := c.Source.Value
	result := uint16(c.Register.a - uint8(input))
	c.SetRegister(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, ^(input & 0x0F) == 0x0F)
	c.SetFlag(flagC, input > uint16(c.Register.a))
}

func (c *CPU) Sbc() {
	input := uint8(c.Source.Value)
	carryBit := BoolToUint(c.GetFlag(flagC))

	result := uint16(c.Register.a - input - carryBit)
	c.SetRegister(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, ^(input & 0x0F) == 0x0F)
	c.SetFlag(flagC, input + carryBit > c.Register.a)
}

func (c *CPU) And() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a & input) 
	c.SetRegister(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, true)
	c.SetFlag(flagC, false)
}


func (c *CPU) Xor() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a & input) 
	c.SetRegister(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, false)
}

func (c *CPU) Or() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a | input) 
	c.SetRegister(A, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false)
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, false)
}


func (c *CPU) Cp() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a - input) 
	
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, true)
	c.SetFlag(flagH, ^(input & 0x0F) == 0x0F)
	c.SetFlag(flagC, input > c.Register.a)
}

func (c *CPU) Rlc(input uint16, t target) {
	msbOn := input & 0x80 //128
	result := input << 1
	if(msbOn != 0){
		 result |= 0x1
	}
	c.SetRegister(t, result)
	
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
}

func (c *CPU) Rrc(input uint16, t target) {
	lsbOn := input & 0x01
	result := input >> 1
	if(lsbOn != 0){
		result |= 0x80
	}
	c.SetRegister(t, result)
	
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
}

func (c *CPU) Rl(input uint16, t target) {
	msbOn := input & 0x80
	result := input << 1
	if (c.GetFlag(flagC)){
		result |= 0x01
	}
	c.SetRegister(t, result)
	
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
}

func (c *CPU) Rr(input uint16, t target) {
	lsbOn := input & 0x01
	result := input >> 1
	if (c.GetFlag(flagC)){
		result |= 0x80
	}
	c.SetRegister(t, result)
	
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
}

func (c *CPU) Sla(input uint16, t target) {
	msbOn := input & 0x80
	result := input << 1
	c.SetRegister(t, result)
	
	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, msbOn != 0)
}

func (c *CPU) Sra(input uint16, t target) {
	msbOn := input & 0x80
	lsbOn := input & 0x01
	result := msbOn | (input >> 1)
	c.SetRegister(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
}

func (c *CPU) Swap(input uint16, t target) {
	result := ((input & 0x0F) << 4) | ((input & 0xF0) >> 4)
	c.SetRegister(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, false)
}

func (c *CPU) Srl(input uint16, t target) {
	lsbOn := input & 0x01
	result := input >> 0x01
	c.SetRegister(t, result)

	c.SetFlag(flagZ, result == 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, false)
	c.SetFlag(flagC, lsbOn != 0)
}

func (c *CPU) Bit(input uint16, b uint8) {
	c.SetFlag(flagZ, input & (1 << b) != 0)
	c.SetFlag(flagN, false) 
	c.SetFlag(flagH, true)
}

func (c *CPU) Res(input uint16, t target, b uint8) {
	result := input & ^(0x01 << b)
	c.SetRegister(t, result)
}

func (c *CPU) Set(input uint16, t target, b uint8) {
	result := input | (0x01 << b)
	c.SetRegister(t, result)
}

func (c *CPU) Cb() error {
	op := uint8(c.Source.Value)
	instruction := cbOpcodes[op]

	data, err := c.GetTarget(instruction.Register)
	if err != nil{
		return err
	}

	input := data.Value
	if data.IsAddr {
		input = uint16(c.BusRead(data.Value))
	}

	switch instruction.Instruction{
		case Rlc:
			c.Rlc(input, instruction.Register)
		case Rrc:
			c.Rrc(input, instruction.Register)
		case Rl:
			c.Rl(input, instruction.Register)
		case Rr:
			c.Rr(input, instruction.Register)
		case Sla:
			c.Sla(input, instruction.Register)
		case Sra:
			c.Sra(input, instruction.Register)	
		case Swap:
			c.Swap(input, instruction.Register)
		case Srl:
			c.Srl(input, instruction.Register)
		case Bit:
			c.Bit(input, instruction.Bit)
		case Res:
			c.Res(input, instruction.Register, instruction.Bit)
		case Set:
			c.Set(input, instruction.Register, instruction.Bit)
	}

	return nil
}