package lib

import (
	"errors"
	"fmt"
)

type flagRegister = int

const (
	z flagRegister 	= 7 	//zero flag 				-> bit 7
	n 				= 6 	//substraction flag (BCD) 	-> bit 6
	h 				= 5 	//half carry flag (BCD) 	-> bit 5
	c				= 4		//carry clag				-> bit 4
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

func (c *CPU) getFlag(flag flagRegister) bool {
	return c.Register.f & (0x1 << flag) != 0
}


func (c *CPU) nop(){
	fmt.Println("NOP INSTRUCTION")
}

func (c *CPU) xor(){
	fmt.Println("Xor INSTRUCTION")
}

func (c *CPU) jp(){
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
			cpu.nop()
		case in_Xor:
			cpu.xor()
		case in_Jp:
			cpu.jp()
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
