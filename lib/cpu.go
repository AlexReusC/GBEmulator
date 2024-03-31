package lib

import (
	"errors"
	"fmt"
)

type flagRegister = int

const (
	zf flagRegister 	= 7 	//zero flag 				-> bit 7
	nf 					= 6 	//substraction flag (BCD) 	-> bit 6
	hf 					= 5 	//half carry flag (BCD) 	-> bit 5
	cf					= 4		//carry clag				-> bit 4
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
	CurrentData uint16
}

func LoadCpu() (*CPU, error) {
	c := &CPU{Register: registers{pc: 0x0100}}

	return c, nil
}

func (c *CPU) GetFlag(flag flagRegister) bool {
	return c.Register.f & (0x1 << flag) != 0
}

func SetBit(b uint8, n int, c bool) uint8 {
	if c {
		b |= (1 << n)
	}else{
		b &= ^(1 << n)
	}
	return b
}

func (c *CPU) SetFlags(flagZ int, flagN int, flagH int, flagC int) {
	if flagZ != -1{
		c.Register.f = SetBit(c.Register.f, zf, flagZ > 0)
	}
	if flagN != -1{
		c.Register.f = SetBit(c.Register.f, nf, flagN > 0)
	}
	if flagH != -1{
		c.Register.f = SetBit(c.Register.f, hf, flagH > 0)
	}
	if flagC != -1{
		c.Register.f = SetBit(c.Register.f, cf, flagC > 0)
	}
} 

func (c *CPU) Nop(){
}

func (c *CPU) Xor(){
	c.Register.a ^= uint8(c.CurrentData & 0xFF)
	c.SetFlags(int(c.Register.a), -1, -1, -1)
}

func (c *CPU) Jp(){
	c.Register.pc = c.CurrentData
}

func (cpu *CPU) Step(b *Bus) error {
	opcode := b.BusRead(cpu.Register.pc)
	fmt.Printf("Opcode: %x, Pc: %x\n", opcode, cpu.Register.pc)
	instruction, ok := instructions[opcode]
	if !ok {
		return errors.New("opcode not implemented")
	}
	fmt.Printf("Instruction: %x, Address mode: %x\n", instruction.InstructionType, instruction.AddressMode)
	cpu.Register.pc += 1

	//Address mode
	switch instruction.AddressMode {
	case am_Imp:
		break
	case am_N16:
		lo := uint16(b.BusRead(cpu.Register.pc))
		hi := uint16(b.BusRead(cpu.Register.pc+1))
		cpu.CurrentData = hi << 8 | lo
	}

	//Conditional mode

	//Instruction type
	switch instruction.InstructionType {
		case in_Nop:
			cpu.Nop()
		case in_Xor:
			cpu.Xor()
		case in_Jp:
			cpu.Jp()
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
