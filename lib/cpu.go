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

type size = int

const (
	u16 size = iota
	u8
	undefined
)

type Data struct {
	Value uint16
	IsAddr bool
	Length size
}

type CPU struct {
	Register registers

	Source Data
	Destination Data
	DestinationTarget targetType

	CurrentConditionResult bool
	InterruptorMasterEnabled bool
}

func LoadCpu() (*CPU, error) {
	c := &CPU{Register: registers{pc: 0x0100, a: 0x01}}

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

func (c *CPU) GetTarget(t targetType, b *Bus) (Data, error) {
	switch t {
		case  target_A:
			return Data{uint16(c.Register.a), false, u8}, nil
		case target_SP:
			return Data{c.Register.sp, false, u16}, nil
		case target_nn:
			lo := uint16(b.BusRead(c.Register.pc))
			hi := uint16(b.BusRead(c.Register.pc+1))
			c.Register.pc += 2
			return Data{(hi << 8 | lo), false, u16}, nil
		case target_None:
			return Data{0, false, undefined}, nil
		// TODO: Other targets
		default:
			return Data{0, false, undefined}, errors.New("unknown target type")
	}
} 

func (c *CPU) SetRegister(t targetType, v uint16)  {
	switch t {
		case target_SP:
			c.Register.sp = v

		default:
			fmt.Println("Unknown register for setting")
			panic(0)
	}
}

func (c *CPU) Nop(){
}

func (c *CPU) Xor(){
	c.Register.a ^= uint8(c.Destination.Value & 0xFF)
	c.SetFlags(int(c.Register.a), -1, -1, -1)
}

func (c *CPU) Jp(){
	if c.CurrentConditionResult{
		c.Register.pc = c.Source.Value
		//cycles(1)
	}
}

func (c *CPU) Add() {

}

func (c *CPU) Di() {
	c.InterruptorMasterEnabled = false
}

func (c *CPU) Ld8(b *Bus) {
	var sourceVal uint8

	if c.Source.IsAddr{
		sourceVal = b.BusRead(c.Source.Value)
	}else{
		sourceVal = uint8(c.Source.Value)
	}

	if c.Destination.IsAddr{
		b.BusWrite(c.Destination.Value, sourceVal)
	} else{
		c.SetRegister(c.DestinationTarget, uint16(sourceVal))
	}

	//TODO: (HL) SP+e instruction
}

func (c *CPU) Ld16(b *Bus){
 	//Ld16 has no addresses in load
	if c.Destination.IsAddr{
		b.BusWrite16(c.Destination.Value,  c.Source.Value)
	}else{
		c.SetRegister(c.DestinationTarget, c.Source.Value)
	}
}

func (cpu *CPU) Step(b *Bus) error {
	opcode := b.BusRead(cpu.Register.pc)
	fmt.Printf("Pc: %x, (%02x %02x %02x) -> ", cpu.Register.pc, opcode, b.BusRead(cpu.Register.pc+1), b.BusRead(cpu.Register.pc+2))
	instruction, ok := instructions[opcode]
	if !ok {
		return errors.New("opcode not implemented")
	}
	fmt.Printf("Instruction: %s, Destination: %s, Source: %s\n", instruction.InstructionType, instruction.Destination, instruction.Source)
	cpu.Register.pc += 1

	//Get source
	data, err := cpu.GetTarget(instruction.Source, b)
	if err != nil{
		return err
	}
	cpu.Source = data

	//Get destination
	data, err = cpu.GetTarget(instruction.Destination, b)
	if err != nil{
		return err
	}
	cpu.Destination = data
	if !cpu.Destination.IsAddr {
		cpu.DestinationTarget = instruction.Destination
	}

	//Conditional mode
	currentCondition := instruction.ConditionType
	conditionResult, err := cpu.checkCond(currentCondition)
	if err != nil{
		return err
	}
	cpu.CurrentConditionResult = conditionResult

	//Instruction type
	switch instruction.InstructionType {
		case in_Nop:
			cpu.Nop()
		case in_Xor:
			cpu.Xor()
		case in_Jp:
			cpu.Jp()
		case in_Di:
			cpu.Di()
		case in_Ld8:
			cpu.Ld8(b)
		case in_Ld16:
			cpu.Ld16(b)
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
