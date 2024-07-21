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

	Halted bool

	Source Data
	SourceTarget target
	DestinationTarget target
	ImmediateData uint16
	CurrentConditionResult bool
	currentOpcode uint8

	MasterInterruptEnabled bool
	EnableMasterInterruptAfter int
	Interrupts uint8
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

func (c *CPU) GetFlag(flag flagRegister) bool {
	return c.Register.f & (0x1 << flag) != 0
}

func (c *CPU) Step(f *os.File) error {
	if !c.Halted {
		instruction, err := c.FetchInstruction(f)
		if err != nil{
			return err
		}

		err = c.DecodeInstruction(instruction)
		if err != nil{
			return err
		}

		//Serial print
		c.Debug.DebugUpdate(c.Bus)
		c.Debug.DebugPrint()

		err = c.ExecuteInstruction(instruction)
		if err != nil{
			return err
		}

	}else{
		//cycle()
		if c.Interrupts != 0 {
			c.Halted = true
		}
	}

	//Manage EI instruction
	if c.EnableMasterInterruptAfter > 0 {
		c.EnableMasterInterruptAfter -= 1
		if c.EnableMasterInterruptAfter == 0 {
			c.MasterInterruptEnabled = true
		}
	}
	c.HandleInterrupts()
	return nil
}

func (c *CPU) FetchInstruction(f *os.File) (Instruction, error) {
	c.currentOpcode = c.BusRead(c.Register.pc)
	instruction, ok := instructions[c.currentOpcode]
	if !ok {
		return instruction, fmt.Errorf("opcode %x not implemented", c.currentOpcode)
	}
	Log(c, f)
	c.Register.pc += 1
	if instruction.Destination == n || instruction.Destination == n_M || instruction.Source == n || instruction.Source == n_M {
		c.ImmediateData = uint16(c.BusRead(c.Register.pc))
		c.Register.pc += 1
	}
	if instruction.Destination == nn || instruction.Destination == nn_M || instruction.Source == nn || instruction.Source == nn_M {
		c.ImmediateData = c.BusRead16(c.Register.pc)
		c.Register.pc += 2
	}
	return instruction, nil
}

func (c *CPU) DecodeInstruction(instruction Instruction) error{
		//Get destination
		c.DestinationTarget = instruction.Destination
		
		//Get source
		data, err := c.GetTarget(instruction.Source)
		if err != nil{
			return err
		}
		c.Source = data
		c.SourceTarget = instruction.Source

		//Conditional mode
		currentCondition := instruction.ConditionType
		conditionResult, err := c.checkCond(currentCondition)
		if err != nil{
			return err
		}
		c.CurrentConditionResult = conditionResult

		return nil
}

func (c *CPU) ExecuteInstruction(i Instruction) error {
		//Instruction type
		switch i.InstructionType {
			case Nop:
				c.Nop()
			case Jp:
				c.Jp()
			case Jr:
				c.Jr()
			case Ld8:
				c.Ld8()
			case Ld16:
				c.Ld16()
			case Ldh:
				c.Ldh()
			case Push:
				c.Push()
			case Pop:
				c.Pop()
			case Call:
				c.Call()
			case Ret:
				c.Ret()
			case Reti:
				c.Reti()
			case Rst:
				c.Rst()
			case Di:
				c.Di()
			case Ei:
				c.Ei()
			case Daa:
				c.Daa()
			case Rlca:
				c.Rlca()
			case Rla:
				c.Rla()
			case Rrca:
				c.Rrca()
			case Rra:
				c.Rra()
			case Ccf:
				c.Ccf()
			case Cpl:
				c.Cpl()
			case Scf:
				c.Scf()
			case Inc:
				c.Inc()
			case Dec:
				c.Dec()
			case Add:
				c.Add()
			case AddHl:
				c.AddHl()
			case Add16_8:
				c.Add16_8()
			case Adc:
				c.Adc()
			case Sub:
				c.Sub()
			case Sbc:
				c.Sbc()	
			case Or:
				c.Or()
			case And:
				c.And()
			case Xor:
				c.Xor()
			case Cp:
				c.Cp()
			case Cb:
				err := c.Cb()
				if err != nil{
					return err
				}
			default:
				return errors.New("invalid instruction")
		}
		return nil
}