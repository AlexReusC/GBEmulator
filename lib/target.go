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
	c.Register.f = SetBitWithCond(c.Register.f, flag, cond)
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

//TODO: remove isAddress info
/*Function for getting the target value dynamically
*/
func (c *CPU) GetTarget(t target) (uint16, error) {
	switch t {
	case A:
		return uint16(c.Register.a), nil
	case B:
		return uint16(c.Register.b), nil
	case C:
		return uint16(c.Register.c), nil
	case D:
		return uint16(c.Register.d), nil
	case E:
		return uint16(c.Register.e), nil
	case F:
		return uint16(c.Register.f), nil
	case H:
		return uint16(c.Register.h), nil
	case L:
		return uint16(c.Register.l), nil
	case SPe8:
		val := uint16(int16(c.Register.sp) + int16(int8(uint8(c.Immediate)))) 
		return val, nil
	case AF:
		return c.GetTargetAF(), nil
	case BC:
		return c.GetTargetBC(), nil
	case DE:
		return c.GetTargetDE(), nil
	case HL:
		return  c.GetTargetHL(), nil
	case SP:
		return c.Register.sp, nil
	case n:
		return c.Immediate, nil
	case nn:
		return c.Immediate, nil
	case C_M:
		return uint16(c.Register.c), nil
	case BC_M:
		return c.GetTargetBC(), nil
	case DE_M:
		return c.GetTargetDE(), nil
	case HL_M:
		return c.GetTargetHL(), nil
	case HLP_M:
		val := c.GetTargetHL()
		c.SetTarget(HL, val+1)
		return val, nil
	case HLM_M:
		val := c.GetTargetHL()
		c.SetTarget(HL, val-1)
		return val, nil
	case n_M:
		return c.Immediate, nil
	case nn_M:
		return c.Immediate, nil
	case None:
		return 0, nil
	default:
		return 0, errors.New("unknown target type")
	}
}

func (c *CPU) SetTarget(t target, v uint16) {
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
	case BC_M:
		c.MMUWrite(c.GetTargetBC(), uint8(v))
	case DE_M:
		c.MMUWrite(c.GetTargetDE(), uint8(v))
	case HL_M:
		c.MMUWrite(c.GetTargetHL(), uint8(v))
	case HLP_M:
		c.MMUWrite(c.GetTargetHL(), uint8(v))
		c.SetTarget(HL, c.GetTargetHL()+1)
	case HLM_M:
		c.MMUWrite(c.GetTargetHL(), uint8(v))
		c.SetTarget(HL, c.GetTargetHL()-1)
	case SP:
		c.Register.sp = v
	case nn_M:
		c.MMUWrite(c.Immediate, uint8(v))
	case nn_M16:
		c.MMUWrite16(c.Immediate, v)
	default:
		fmt.Printf("Unknown register %s for setting\n", t)
		panic(0)
	}
}