package lib

import "fmt"

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
	c.Bus.interruptorFlags = UnsetBit(c.Bus.interruptorFlags, int(b))
	c.MasterInterruptEnabled = true

	c.Register.sp -= 1
	c.BusWrite(c.Register.sp, uint8((c.Register.pc&0xFF00)>>8))
	c.Register.sp -= 1
	c.BusWrite(c.Register.sp, uint8(c.Register.pc&0xFF))

	fmt.Printf("Process %b \n", b)
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
		if (c.Bus.ieRegister&flag != 0) && (c.Bus.interruptorFlags&flag != 0) {
			c.ProcessInterrupt(interruptor.Bit, interruptor.Address)
			return
		}
	}
}

func (c *CPU) RequestInterrupt(i InterruptorBit) {
	c.Bus.interruptorFlags |= (1 << i)
}
