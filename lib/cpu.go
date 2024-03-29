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

func (c *CPU) nop_00(){
	fmt.Println("NOP INSTRUCTION")
}

func (c *CPU) xor_AF(){
	fmt.Println("NOP INSTRUCTION")
}

func (c *CPU) jp_C3(b *Bus){
	lo := uint16(b.BusRead(c.Register.pc))
	hi := uint16(b.BusRead(c.Register.pc+1))
	a := hi << 8 | lo
	c.Register.pc = a
}

func (cpu *CPU) Step(b *Bus) error {
	opcode := b.BusRead(cpu.Register.pc)
	fmt.Printf("Opcode: %x, Pc: %x\n", opcode, cpu.Register.pc)
	cpu.Register.pc += 1
	switch opcode {
		case 0x00:
			cpu.nop_00()
		case 0xAF:
			cpu.xor_AF()
		case 0xC3:
			cpu.jp_C3(b)
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
