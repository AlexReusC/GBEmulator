package lib

import (
	"errors"
	"fmt"
)

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

func (c *CPU) proc_nop(){

}

func (cpu *CPU) Step(b *Bus, cart *Cart) error {
	opcode := b.BusRead(cpu.Register.pc, cart)
	fmt.Printf("Opcode: %x, Pc: %x\n", opcode, cpu.Register.pc)
	cpu.Register.pc += 1
	switch opcode {
		case 0x00:
			cpu.proc_nop()
		case 0xAF:
			proc_xor()
		case 0xC3:
			proc_jp()
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
