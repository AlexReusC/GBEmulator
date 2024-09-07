package lib

import (
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
	Clock *Clock

	Halted bool

	Source Data
	SourceTarget target
	DestinationTarget target
	Immediate uint16
	CurrentConditionResult bool
	currentOpcode uint8

	MasterInterruptEnabled bool
	EnableMasterInterruptAfter int
	Interrupts uint8
}

func LoadCpu(b *Bus, d *Debug, cl *Clock) (*CPU, error) {
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
		Clock: cl,
	}

	return c, nil
}


func (c *CPU) BusRead(a uint16) uint8 { return c.Bus.BusRead(a) }
func (c *CPU) BusRead16(a uint16) uint16 { return c.Bus.BusRead16(a) }
func (c *CPU) BusWrite(a uint16, v uint8) { c.Bus.BusWrite(a, v) }
func (c *CPU) BusWrite16(a uint16, v uint16) { c.Bus.BusWrite16(a, v) }
func (c *CPU) GetFlag(flag flagRegister) bool { return c.Register.f & (0x1 << flag) != 0 }
func (c *CPU) UpdateClock(cycles int) {
	changeTimer := c.Clock.Update(cycles)
	if changeTimer {
		c.Bus.RequestInterrupt(TIMER)
	}
}

func (c *CPU) Step(f *os.File) (int, error) {
	cycles := 0
	if !c.Halted {
		instruction, err := c.FetchInstruction(f)
		if err != nil{
			return 0, err
		}

		err = c.DecodeInstruction(instruction)
		if err != nil{
			return 0, err
		}

		//Serial print
		c.Debug.DebugUpdate(c.Bus)

		instructionCycles, err := c.ExecuteInstruction(instruction)
		if err != nil{
			return 0, err
		}
		cycles += instructionCycles
	}else{
		cycles += 1
		if c.Bus.interruptorFlags != 0 {
			c.Halted = false
		}
	}

	//Manage EI instruction
	if c.EnableMasterInterruptAfter > 0 {
		c.EnableMasterInterruptAfter -= 1
		if c.EnableMasterInterruptAfter == 0 {
			c.MasterInterruptEnabled = true
		}
	}
	return cycles, nil
}

func (c *CPU) FetchInstruction(f *os.File) (Instruction, error) {
	c.currentOpcode = c.BusRead(c.Register.pc)
	instruction, ok := instructions[c.currentOpcode]
	if !ok {
		return instruction, fmt.Errorf("opcode %x not implemented", c.currentOpcode)
	}
	if f != nil {
		DoctorLog(c, f)
	}
	c.Register.pc += 1
	if IsImmediateTarget8(instruction.Source) || IsImmediateTarget8(instruction.Destination) {
		c.Immediate = uint16(c.BusRead(c.Register.pc))
		c.Register.pc += 1
	} else if IsImmediateTarget16(instruction.Source) || IsImmediateTarget16(instruction.Destination) {
		c.Immediate = c.BusRead16(c.Register.pc)
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

func (c *CPU) ExecuteInstruction(i Instruction) (int, error) {
		cycles := 0
		switch i.InstructionType {
			case Nop:
				cycles += c.Nop()
			case Jp:
				cycles += c.Jp()
			case Jr:
				cycles += c.Jr()
			case Ld8:
				cycles += c.Ld8()
			case Ld16:
				cycles += c.Ld16()
			case Ldh:
				cycles += c.Ldh()
			case LdSPn:
				cycles += c.LdSPn()
			case Push:
				cycles += c.Push()
			case Pop:
				cycles += c.Pop()
			case Call:
				cycles += c.Call()
			case Ret:
				cycles += c.Ret()
			case Reti:
				cycles += c.Reti()
			case Rst:
				cycles += c.Rst()
			case Di:
				cycles += c.Di()
			case Ei:
				cycles += c.Ei()
			case Halt:
				cycles += c.Halt()
			case Daa:
				cycles += c.Daa()
			case Rlca:
				cycles += c.Rlca()
			case Rla:
				cycles += c.Rla()
			case Rrca:
				cycles += c.Rrca()
			case Rra:
				cycles += c.Rra()
			case Ccf:
				cycles += c.Ccf()
			case Cpl:
				cycles += c.Cpl()
			case Scf:
				cycles += c.Scf()
			case Inc:
				cycles += c.Inc()
			case Dec:
				cycles += c.Dec()
			case Add:
				cycles += c.Add()
			case AddHl:
				cycles += c.AddHl()
			case Add16_8:
				cycles += c.Add16_8()
			case Adc:
				cycles += c.Adc()
			case Sub:
				cycles += c.Sub()
			case Sbc:
				cycles += c.Sbc()	
			case Or:
				cycles += c.Or()
			case And:
				cycles += c.And()
			case Xor:
				cycles += c.Xor()
			case Cp:
				cycles += c.Cp()
			case Cb:
				instructionCycles, err := c.Cb()
				if err != nil{
					return 0, err
				}
				cycles += instructionCycles
			default:
				return 0, fmt.Errorf("invalid instruction %s", i.InstructionType)
		}
		return cycles, nil
}