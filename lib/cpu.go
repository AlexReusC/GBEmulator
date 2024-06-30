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

type Data struct {
	Value uint16
	IsAddr bool
}

type CPU struct {
	Register registers

	Source Data
	Destination Data
	SourceTarget targetType
	DestinationTarget targetType
	CurrentConditionResult bool
	currentOpcode uint8

	InterruptorMasterEnabled bool
	IeRegister uint8
}

func LoadCpu() (*CPU, error) {
	c := &CPU{Register: registers{pc: 0x0100, a: 0x01}}

	return c, nil
}

func (c *CPU) GetIeRegister() uint8 {
	return c.IeRegister
}

func (c *CPU) SetIeRegister(n uint8) {
	c.IeRegister = n
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

func (c *CPU) GetTargetAF() uint16{
	hi := uint16(c.Register.a)
	lo := uint16(c.Register.f)
	return (hi << 8 | lo)
}

func (c *CPU) GetTargetBC() uint16{
	hi := uint16(c.Register.b)
	lo := uint16(c.Register.c)
	return (hi << 8 | lo)
}

func (c *CPU) GetTargetDE() uint16{
	hi := uint16(c.Register.d)
	lo := uint16(c.Register.e)
	return (hi << 8 | lo)
}

func (c *CPU) GetTargetHL() uint16{
	hi := uint16(c.Register.h)
	lo := uint16(c.Register.l)
	return (hi << 8 | lo)
}

func (c *CPU) GetTarget(t targetType, b *Bus) (Data, error) {
	switch t {
		case  A:
			return Data{uint16(c.Register.a), false}, nil
		case B:
			return Data{uint16(c.Register.b), false}, nil
		case C:
			return Data{uint16(c.Register.c), false}, nil
		case D:
			return Data{uint16(c.Register.d), false}, nil
		case E:
			return Data{uint16(c.Register.e), false}, nil
		case F:
			return Data{uint16(c.Register.f), false}, nil
		case H:
			return Data{uint16(c.Register.h), false}, nil
		case L:
			return Data{uint16(c.Register.l), false}, nil
		case AF:
			return Data{c.GetTargetAF(), false}, nil
		case BC:
			return Data{c.GetTargetBC(), false}, nil
		case DE:
			return Data{c.GetTargetDE(), false}, nil
		case HL:
			return Data{c.GetTargetHL(), false}, nil
		case SP:
			return Data{c.Register.sp, false}, nil
		case n:
			n := uint16(b.BusRead(c.Register.pc))
			c.Register.pc += 1
			return Data{n, false}, nil
		case nn:		
			nn := b.BusRead16(c.Register.pc)
			c.Register.pc += 2
			return Data{nn, false}, nil
		case C_M:
			return Data{uint16(c.Register.c), true}, nil
		case BC_M:
			return Data{c.GetTargetBC(), true}, nil
		case DE_M:
			return Data{c.GetTargetDE(), true}, nil
		case HL_M:
			return Data{c.GetTargetHL(), true}, nil
		case HLP_M:
			return Data{0, false}, errors.New("target not implemented: (HL+)")
		case HLM_M:
			return Data{0, false}, errors.New("target not implemented: (HL-)")
		case n_M:
			n := uint16(b.BusRead(c.Register.pc))
			c.Register.pc += 1
			return Data{n, true}, nil
		case nn_M:
			nn := b.BusRead16(c.Register.pc)
			c.Register.pc += 2
			return Data{nn, true}, nil
		case None:
			return Data{0, false}, nil
		// TODO: Other targets
		default:
			return Data{0, false}, errors.New("unknown target type")
	}
} 

func (c *CPU) SetRegister(t targetType, v uint16)  {
	switch t {
		case A:
			c.Register.a = uint8(v)
		case B:
			c.Register.b = uint8(v)
		case C:
			c.Register.c = uint8(v)
		case D:
			c.Register.d = uint8(v)
		case E:
			c.Register.e = uint8(v)
		case F:
			c.Register.f = uint8(v)	
		case H:
			c.Register.h = uint8(v)	
		case L:
			c.Register.l = uint8(v)	
		case AF:
			c.Register.a = uint8((v & 0xFF00) >> 8)
			c.Register.f = uint8(v & 0xFF)
		case BC:
			c.Register.b = uint8((v & 0xFF00) >> 8)
			c.Register.c = uint8(v & 0xFF)
		case DE:
			c.Register.d = uint8((v & 0xFF00) >> 8)
			c.Register.e = uint8(v & 0xFF)
		case HL:
			c.Register.h = uint8((v & 0xFF00) >> 8)
			c.Register.l = uint8(v & 0xFF)
		case SP:
			c.Register.sp = v
		default:
			fmt.Printf("Unknown register %x for setting\n", t)
			panic(0)
	}
}

func (cpu *CPU) Step(b *Bus) error {
	cpu.currentOpcode = b.BusRead(cpu.Register.pc)
	fmt.Printf("Pc: %x, (%02x %02x %02x) -> ", cpu.Register.pc, cpu.currentOpcode, b.BusRead(cpu.Register.pc+1), b.BusRead(cpu.Register.pc+2))
	instruction, ok := instructions[cpu.currentOpcode]
	if !ok {
		return errors.New("opcode not implemented")
	}
	fmt.Printf("Instruction: %-6s Destination: %-6s Source: %-6s A: %02x BC: %02x%02x DE: %02x%02x  HL: %02x%02x\n", instruction.InstructionType, instruction.Destination, instruction.Source, cpu.Register.a, cpu.Register.b, cpu.Register.c, cpu.Register.d, cpu.Register.e, cpu.Register.h, cpu.Register.l)
	cpu.Register.pc += 1

	//Get destination, including inmediate
	data, err := cpu.GetTarget(instruction.Destination, b)
	if err != nil{
		return err
	}
	cpu.Destination = data
	cpu.DestinationTarget = instruction.Destination
	

	//Get source, including inmediate
	data, err = cpu.GetTarget(instruction.Source, b)
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

	//Instruction type
	switch instruction.InstructionType {
		case in_Nop:
			cpu.Nop()
		case in_Xor:
			cpu.Xor()
		case in_Jp:
			cpu.Jp()
		case in_Jr:
			cpu.Jr()
		case in_Di:
			cpu.Di()
		case in_Ld8:
			cpu.Ld8(b)
		case in_Ld16:
			cpu.Ld16(b)
		case in_Ldh:
			cpu.Ldh(b)
		case in_Push:
			cpu.Push(b)
		case in_Pop:
			cpu.Pop(b)
		case in_Call:
			cpu.Call(b)
		case in_Ret:
			cpu.Ret(b)
		case in_Reti:
			cpu.Reti(b)
		case in_Rst:
			cpu.Rst(b)
		case in_Inc:
			cpu.Inc()
		case in_Dec:
			cpu.Dec()
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
