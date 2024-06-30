package lib

func (c *CPU) Nop() {
}

func (c *CPU) Xor() {
	c.Register.a ^= uint8(c.Destination.Value & 0xFF)
	c.SetFlags(int(c.Register.a), -1, -1, -1)
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

func (c *CPU) Add() {

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