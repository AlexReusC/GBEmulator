package lib

import (
	"fmt"
)

//TODO: flags behavior
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
		//cycles(1)
	}
}

func (c *CPU) Jr() {
	if c.CurrentConditionResult {
		c.Register.pc = c.Register.pc + c.Source.Value
		//cycles(1)
	}
}

func (c *CPU) Di() {
	c.InterruptorMasterEnabled = false
}

func (c *CPU) Ld8(b *Bus) {
	var input uint8

	if c.Source.IsAddr {
		input = b.BusRead(c.Source.Value)
	} else {
		input = uint8(c.Source.Value)
	}

	if c.Destination.IsAddr {
		b.BusWrite(c.Destination.Value, input)
	} else {
		c.SetRegister(c.DestinationTarget, uint16(input))
	}

	//TODO: (HL) SP+e instruction
}

func (c *CPU) Ld16(b *Bus) {
	//Ld16 has no addresses in load
	if c.Destination.IsAddr {
		b.BusWrite16(c.Destination.Value, c.Source.Value)
	} else {
		c.SetRegister(c.DestinationTarget, c.Source.Value)
	}
}

func (c *CPU) Ldh(b *Bus) {
	var input uint8

	if c.Source.IsAddr {
		input = b.BusRead(0xFF00 | uint16(b.BusRead(c.Source.Value)))
	} else {
		input = uint8(c.Source.Value)
	}

	if c.Destination.IsAddr {
		b.BusWrite(0xFF00|c.Destination.Value, input)
	} else {
		c.SetRegister(A, uint16(input)) //If destination is not address is always register A
	}
}

func (c *CPU) Push(b *Bus) {
	c.Register.sp -= 1
	b.BusWrite(c.Register.sp, uint8((c.Source.Value&0xFF00)>>8))

	c.Register.sp -= 1
	b.BusWrite(c.Register.sp, uint8(c.Source.Value&0xFF))
}

func (c *CPU) Pop(b *Bus) {
	lo := uint16(b.BusRead(c.Register.sp))
	c.Register.sp += 1

	hi := uint16(b.BusRead(c.Register.sp))
	c.Register.sp += 1

	c.SetRegister(c.DestinationTarget, (hi<<8)|lo)
}

func (c *CPU) Call(b *Bus) {
	if c.CurrentConditionResult {
		//Push pc
		c.Register.sp -= 1
		b.BusWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
		c.Register.sp -= 1
		b.BusWrite(c.Register.sp, uint8(c.Register.pc&0x00FF))

		//Jp nn
		c.Register.pc = c.Source.Value
	}
}

func (c *CPU) Ret(b *Bus) {
	if c.CurrentConditionResult {
		//Pop
		lo := uint16(b.BusRead(c.Register.sp))
		c.Register.sp += 1

		hi := uint16(b.BusRead(c.Register.sp))
		c.Register.sp += 1
		//Jp
		c.Register.pc = (hi << 8) | lo
	}
}

func (c *CPU) Reti(b *Bus) {
	//Pop
	lo := uint16(b.BusRead(c.Register.sp))
	c.Register.sp += 1

	hi := uint16(b.BusRead(c.Register.sp))
	c.Register.sp += 1
	//Jp
	c.Register.pc = (hi << 8) | lo

	c.InterruptorMasterEnabled = true
}

func (c *CPU) Rst(b *Bus) {
	c.Register.sp -= 1
	b.BusWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
	c.Register.sp -= 1
	b.BusWrite(c.Register.sp, uint8(c.Register.pc&0xFF))

	c.Register.pc = (0x00 << 8) | rstAddress[c.currentOpcode]
}

func (c *CPU) Inc() {
	newVal := c.Source.Value + 1
	c.SetRegister(c.SourceTarget, newVal)
}

func (c *CPU) Dec() {
	newVal := c.Source.Value - 1
	c.SetRegister(c.SourceTarget, newVal)
}

func (c *CPU) Add() {
	input := c.Source.Value
	result := c.Register.a + uint8(input)
	c.SetRegister(A, uint16(result))
}

func (c *CPU) AddHl() {
	input := c.Source.Value
	result := c.GetTargetHL() + input
	c.SetRegister(HL, result)
}

func (c *CPU) Add16_8() {
	input16 := c.Destination.Value
	input8 := c.Source.Value
	result := input8 + input16

	c.SetRegister(c.DestinationTarget, result)
}

func (c *CPU) Adc() {
	input := uint8(c.Source.Value)
	carryBit := BoolToUint(c.GetFlag(flagC))

	result := uint16(c.Register.a + input + carryBit)
	c.SetRegister(A, result)
}

func (c *CPU) Sub() {
	input := c.Source.Value
	result := uint16(c.Register.a - uint8(input))
	c.SetRegister(A, result)
}

func (c *CPU) Sbc() {
	input := uint8(c.Source.Value)
	carryBit := BoolToUint(c.GetFlag(flagC))

	result := uint16(c.Register.a - input - carryBit)
	c.SetRegister(A, result)
}

func (c *CPU) And() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a & input) 
	c.SetRegister(A, result)
}


func (c *CPU) Xor() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a & input) 
	c.SetRegister(A, result)
	c.SetFlags(int(c.Register.a), -1, -1, -1)
}

func (c *CPU) Or() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a | input) 
	c.SetRegister(A, result)
}


func (c *CPU) Cp() {
	input := uint8(c.Source.Value)
	result := uint16(c.Register.a - input) 
	//TODO: flags
	fmt.Println("Cp ins: ",result)
}

func (c *CPU) Rlc(input uint16, t target) {
	msbOn := input & 0x80 //128
	modifiedVal := input << 1
	if(msbOn != 0){
		 modifiedVal |= 0x1
	}
	c.SetRegister(t, modifiedVal)
	//set flags
}

func (c *CPU) Rrc(input uint16, t target) {
	lsbOn := input & 0x01
	modifiedVal := input >> 1
	if(lsbOn != 0){
		modifiedVal |= 0x80
	}
	c.SetRegister(t, modifiedVal)
	//set flags
}

func (c *CPU) Rl(input uint16, t target) {
	//TODO: remove when implementing flags
	//msbOn := input & 0x80
	modifiedVal := input << 1
	if (c.GetFlag(flagC)){
		modifiedVal |= 0x01
	}
	c.SetRegister(t, modifiedVal)
	//Set flags
}

func (c *CPU) Rr(input uint16, t target) {
	//TODO: remove when implementing flags
	//lsbOn := input & 0x01
	modifiedVal := input >> 1
	if (c.GetFlag(flagC)){
		modifiedVal |= 0x80
	}
	c.SetRegister(t, modifiedVal)
	//Set flags
}

func (c *CPU) Sla(input uint16, t target) {
	//TODO
	//msbOn := input & 0x80
	modifiedVal := input << 1
	c.SetRegister(t, modifiedVal)
	//Set flags
}

func (c *CPU) Sra(input uint16, t target) {
	//TODO: remove when implementing flags
	msbOn := input & 0x80
	//lsbOn := input & 0x01
	modifiedVal := msbOn | (input >> 1)
	c.SetRegister(t, modifiedVal)
	//Set flags
}

func (c *CPU) Swap(input uint16, t target) {
	modifiedVal := ((input & 0x0F) << 4) | ((input & 0xF0) >> 4)
	c.SetRegister(t, modifiedVal)
	//Set flags
}

func (c *CPU) Srl(input uint16, t target) {
	//TODO: remove when implementing flags
	//lsbOn := input & 0x01
	modifiedVal := input >> 0x01
	c.SetRegister(t, modifiedVal)
}

func (c *CPU) Bit(input uint16, t target, b uint8) {
	//TODO: flags
}

func (c *CPU) Res(input uint16, t target, b uint8) {
	modifiedVal := input & ^(0x01 << b)
	c.SetRegister(t, modifiedVal)
	//TODO: flags
}

func (c *CPU) Set(input uint16, t target, b uint8) {
	modifiedVal := input | (0x01 << b)
	c.SetRegister(t, modifiedVal)
	//TODO: flags
}


func (c *CPU) Cb(b *Bus) error {
	op := uint8(c.Source.Value)
	instruction := cbOpcodes[op]

	data, err := c.GetTarget(instruction.Register, b)
	if err != nil{
		return err
	}

	input := data.Value
	if data.IsAddr {
		input = uint16(b.BusRead(data.Value))
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
			c.Bit(input, instruction.Register, instruction.Bit)
		case Res:
			c.Res(input, instruction.Register, instruction.Bit)
		case Set:
			c.Set(input, instruction.Register, instruction.Bit)
	}

	return nil
}