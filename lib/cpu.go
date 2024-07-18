package lib

import (
	"errors"
	"fmt"
	"os"
)

type flagRegister = int

const (
	flagZ 	flagRegister 	= 7 	//zero flag 				-> bit 7
	flagN 	flagRegister	= 6 	//substraction flag (BCD) 	-> bit 6
	flagH 	flagRegister	= 5 	//half carry flag (BCD) 	-> bit 5
	flagC	flagRegister	= 4		//carry clag				-> bit 4
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

type Data struct {
	Value uint16
	IsAddr bool
}

type CPU struct {
	Register registers
	Bus *Bus
	Debug *Debug

	Source Data
	SourceTarget target
	DestinationTarget target
	ImmediateData uint16
	CurrentConditionResult bool
	currentOpcode uint8

	InterruptorMasterEnabled bool
}

func LoadCpu(b *Bus, d *Debug) (*CPU, error) {
	c := &CPU{
		Register: registers{
			a: 	0x01,
			b: 	0x00,
			c: 	0x13, 
			d:	0x00, 
			e: 	0xD8, 
			h:	0x01, 
			l:	0x4D, 
			f: 	0xB0, 
			pc: 0x0100, 
			sp: 0xFFFE,
			}, 
		Bus: b, 
		Debug: d,
	}

	return c, nil
}

func (c *CPU) GetIeRegister() uint8 {
	return c.Bus.ieRegister
}

func (c *CPU) SetIeRegister(n uint8) {
	c.Bus.ieRegister = n
}

func (c *CPU) BusRead(a uint16) uint8 {
	return c.Bus.BusRead(a)
}

func (c *CPU) BusRead16(a uint16) uint16 {
	return c.Bus.BusRead16(a)
}

func (c *CPU) BusWrite(a uint16, v uint8) {
	c.Bus.BusWrite(a, v)
}

func (c *CPU) BusWrite16(a uint16, v uint16) {
	c.Bus.BusWrite16(a, v)
}

func (c *CPU) GetImmediateData(destination target, source target) {
	c.Register.pc += 1
	if destination == n || destination == n_M || source == n || source == n_M {
		c.ImmediateData = uint16(c.BusRead(c.Register.pc))
		c.Register.pc += 1
	}
	if destination == nn || destination == nn_M || source == nn || source == nn_M {
		c.ImmediateData = c.BusRead16(c.Register.pc)
		c.Register.pc += 2
	} 
}

func (c *CPU) GetFlag(flag flagRegister) bool {
	return c.Register.f & (0x1 << flag) != 0
}

func (cpu *CPU) Step(f *os.File) error {
	cpu.currentOpcode = cpu.BusRead(cpu.Register.pc)
	instruction, ok := instructions[cpu.currentOpcode]
	if !ok {
		return fmt.Errorf("opcode %x not implemented", cpu.currentOpcode)
	}
	Log(cpu, f)

	cpu.GetImmediateData(instruction.Destination, instruction.Source)

	//Get destination
	cpu.DestinationTarget = instruction.Destination
	
	//Get source
	data, err := cpu.GetTarget(instruction.Source)
	if err != nil{
		return err
	}
	cpu.Source = data
	cpu.SourceTarget = instruction.Source

	//Conditional mode
	currentCondition := instruction.ConditionType
	conditionResult, err := cpu.checkCond(currentCondition)
	if err != nil{
		return err
	}
	cpu.CurrentConditionResult = conditionResult

	//Serial print
	cpu.Debug.DebugUpdate(cpu.Bus)
	cpu.Debug.DebugPrint()

	//Instruction type
	switch instruction.InstructionType {
		case Nop:
			cpu.Nop()
		case Jp:
			cpu.Jp()
		case Jr:
			cpu.Jr()
		case Ld8:
			cpu.Ld8()
		case Ld16:
			cpu.Ld16()
		case Ldh:
			cpu.Ldh()
		case Push:
			cpu.Push()
		case Pop:
			cpu.Pop()
		case Call:
			cpu.Call()
		case Ret:
			cpu.Ret()
		case Reti:
			cpu.Reti()
		case Rst:
			cpu.Rst()
		case Di:
			cpu.Di()
		case Ei:
			cpu.Ei()
		case Daa:
			cpu.Daa()
		case Rlca:
			cpu.Rlca()
		case Rla:
			cpu.Rla()
		case Rrca:
			cpu.Rrca()
		case Rra:
			cpu.Rra()
		case Ccf:
			cpu.Ccf()
		case Cpl:
			cpu.Cpl()
		case Scf:
			cpu.Scf()
		case Inc:
			cpu.Inc()
		case Dec:
			cpu.Dec()
		case Add:
			cpu.Add()
		case AddHl:
			cpu.AddHl()
		case Add16_8:
			cpu.Add16_8()
		case Adc:
			cpu.Adc()
		case Sub:
			cpu.Sub()
		case Sbc:
			cpu.Sbc()	
		case Or:
			cpu.Or()
		case And:
			cpu.And()
		case Xor:
			cpu.Xor()
		case Cp:
			cpu.Cp()
		case Cb:
			err := cpu.Cb()
			if err != nil{
				return err
			}
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
