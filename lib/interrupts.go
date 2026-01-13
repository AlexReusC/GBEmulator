package lib

type InterruptorBit int

const (
	VBLANK InterruptorBit = iota
	LCDSATUS
	TIMER
	SERIAL
	JOYPAD
)

type Interruptor struct {
	Bit     InterruptorBit
	Address uint16
}

func (c *CPU) ProcessInterrupt(b InterruptorBit, address uint16) {
	c.Halted = false
	c.MMU.interruptorFlags = UnsetBit(c.MMU.interruptorFlags, int(b))
	c.MasterInterruptEnabled = true

	c.Register.sp -= 1
	c.MMUWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
	c.Register.sp -= 1
	c.MMUWrite(c.Register.sp, uint8(c.Register.pc&0xFF))

	c.Register.pc = address
}

func (c *CPU) HandleInterrupts() {
	if !c.MasterInterruptEnabled {
		return
	}

	interruptors := []Interruptor{
		{VBLANK, 0x40},
		{LCDSATUS, 0x48},
		{TIMER, 0x50},
		{SERIAL, 0x58},
		{JOYPAD, 0x60},
	}

	for _, interruptor := range interruptors {
		var flag uint8 = 1 << interruptor.Bit
		if (c.MMU.ieRegister&flag != 0) && (c.MMU.interruptorFlags&flag != 0) {
			c.ProcessInterrupt(interruptor.Bit, interruptor.Address)
			return
		}
	}
}
