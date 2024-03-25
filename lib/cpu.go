package lib

import "fmt"

type registers struct {
	a  uint8
	b  uint8
	c  uint8
	d  uint8
	e  uint8
	f  uint8
	h  uint8
	l  uint8
	sp uint16
	pc uint16
}

type CPU struct {
	Register registers
}

func LoadCpu() (*CPU, error) {
	c := &CPU{Register: registers{pc: 0x0100}}

	return c, nil
}

func (cpu *CPU) Step(b *Bus, cart *Cart) {
	opcode := b.BusRead(cpu.Register.pc, cart)
	fmt.Printf("Opcode: %x, Pc: %x\n", opcode, cpu.Register.pc)
	cpu.Register.pc += 1
}
