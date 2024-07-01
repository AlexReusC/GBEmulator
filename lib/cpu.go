package lib

import (
	"errors"
	"fmt"
)

type flagRegister = int

const (
	zf 	flagRegister 	= 7 	//zero flag 				-> bit 7
	nf 	flagRegister	= 6 	//substraction flag (BCD) 	-> bit 6
	hf 	flagRegister	= 5 	//half carry flag (BCD) 	-> bit 5
	cf	flagRegister	= 4		//carry clag				-> bit 4
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
	SourceTarget target
	DestinationTarget target
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

func (cpu *CPU) checkCond( ct conditional) (bool, error) {
	c := cpu.GetFlag(cf)
	z := cpu.GetFlag(zf)

	switch ct {
	case cond_None:
		return true, nil
	case cond_C:
		return c, nil
	case cond_NC:
		return !c, nil
	case cond_Z:
		return z, nil
	case cond_NZ:
		return !z, nil
	default:
		return false, errors.New("invalid condition")
	}
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

func (c *CPU) GetTarget(t target, b *Bus) (Data, error) {
	switch t {
		case A:
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
			val := c.GetTargetHL()
			c.SetRegister(HL, val + 1)
			return Data{val, true}, nil 
		case HLM_M:
			val := c.GetTargetHL()
			c.SetRegister(HL, val - 1)
			return Data{val, true}, nil 
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

func (c *CPU) SetRegister(t target, v uint16)  {
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
	//TODO: logging for flags
	fmt.Printf("Inst: %-6s Dest: %-6s Src: %-6s A: %02x BC: %02x%02x DE: %02x%02x  HL: %02x%02x\n", instruction.InstructionType, instruction.Destination, instruction.Source, cpu.Register.a, cpu.Register.b, cpu.Register.c, cpu.Register.d, cpu.Register.e, cpu.Register.h, cpu.Register.l)
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
		case Nop:
			cpu.Nop()
		case Jp:
			cpu.Jp()
		case Jr:
			cpu.Jr()
		case Di:
			cpu.Di()
		case Ld8:
			cpu.Ld8(b)
		case Ld16:
			cpu.Ld16(b)
		case Ldh:
			cpu.Ldh(b)
		case Push:
			cpu.Push(b)
		case Pop:
			cpu.Pop(b)
		case Call:
			cpu.Call(b)
		case Ret:
			cpu.Ret(b)
		case Reti:
			cpu.Reti(b)
		case Rst:
			cpu.Rst(b)
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
		case Xor:
			cpu.Xor()
		default:
			return errors.New("invalid instruction")
	}
	return nil
}
