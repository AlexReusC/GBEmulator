package lib

import (
	"errors"
	"fmt"
)

func (cpu *CPU) checkCond(ct conditional) (bool, error) {
	c := cpu.GetFlag(flagC)
	z := cpu.GetFlag(flagZ)

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

func (c *CPU) SetFlag(flag flagRegister, cond bool) {
	c.Register.f = SetBit(c.Register.f, flag, cond)
}

func (c *CPU) FormatFlag(f flagRegister, s rune) rune {
	if c.GetFlag(f){
		return s
	} else {
		return '-'
	}
}

func (c *CPU) GetTargetAF() uint16 {
	hi := uint16(c.Register.a)
	lo := uint16(c.Register.f)
	return (hi<<8 | lo)
}

func (c *CPU) GetTargetBC() uint16 {
	hi := uint16(c.Register.b)
	lo := uint16(c.Register.c)
	return (hi<<8 | lo)
}

func (c *CPU) GetTargetDE() uint16 {
	hi := uint16(c.Register.d)
	lo := uint16(c.Register.e)
	return (hi<<8 | lo)
}

func (c *CPU) GetTargetHL() uint16 {
	hi := uint16(c.Register.h)
	lo := uint16(c.Register.l)
	return (hi<<8 | lo)
}

//TODO: move pc logic
func (c *CPU) GetTarget(t target) (Data, error) {
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
		n := uint16(c.BusRead(c.Register.pc))
		c.Register.pc += 1
		return Data{n, false}, nil
	case nn:
		nn := c.BusRead16(c.Register.pc)
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
		c.SetRegister(HL, c.GetTargetHL()+1)
		return Data{val, true}, nil
	case HLM_M:
		val := c.GetTargetHL()
		c.SetRegister(HL, c.GetTargetHL()-1)
		return Data{val, true}, nil
	case n_M:
		n := uint16(c.BusRead(c.Register.pc))
		c.Register.pc += 1
		return Data{n, true}, nil
	case nn_M:
		nn := c.BusRead16(c.Register.pc)
		c.Register.pc += 2
		return Data{nn, true}, nil
	case None:
		return Data{0, false}, nil
	// TODO: Other targets
	default:
		return Data{0, false}, errors.New("unknown target type")
	}
}

func (c *CPU) SetRegister(t target, v uint16) {
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
		c.Register.f = uint8(v & 0x00FF)
	case BC:
		c.Register.b = uint8((v & 0xFF00) >> 8)
		c.Register.c = uint8(v & 0xFF)
	case DE:
		c.Register.d = uint8((v & 0xFF00) >> 8)
		c.Register.e = uint8(v & 0xFF)
	case HL:
		c.Register.h = uint8((v & 0xFF00) >> 8)
		c.Register.l = uint8(v & 0xFF)
	case HL_M:
		c.BusWrite(c.GetTargetHL(), uint8(v))
	case SP:
		c.Register.sp = v
	default:
		fmt.Printf("Unknown register %x for setting\n", t)
		panic(0)
	}
}