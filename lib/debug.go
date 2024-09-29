package lib

import (
	"fmt"
	"log"
	"os"
)

type Debug struct {
	debugMsg []rune
	msgSize  int
}

func LoadDebug() *Debug {
	d := &Debug{debugMsg: make([]rune, 1024), msgSize: 0}

	return d
}

func (d *Debug) DebugUpdate(m *MMU) {
	if m.Read(0xFF02) == 0x81 {
		c := rune(m.Read(0xFF01))

		d.debugMsg[d.msgSize] = c
		d.msgSize += 1

		m.Write(0xFF02, 0)

		d.DebugPrint()
	}
}

func (d *Debug) DebugPrint() {
	if d.debugMsg[0] != 0 {
		fmt.Printf("DBG: %s\n", string(d.debugMsg))
	}
}

func (d *Debug) GetMsg() string {
	return string(d.debugMsg)
}

func Log(c *CPU, i Instruction, f *os.File) {
	pcData := fmt.Sprintf("Pc: %x, (%02x %02x %02x) -> ", c.Register.pc, c.currentOpcode, c.MMURead(c.Register.pc+1), c.MMURead(c.Register.pc+2))
	flags := fmt.Sprintf("%c%c%c%c", c.FormatFlag(flagZ, 'Z'), c.FormatFlag(flagN, 'N'), c.FormatFlag(flagH, 'H'), c.FormatFlag(flagC, 'C'))
	output := fmt.Sprintf("%s Inst: %-6s Dest: %-6s Src: %-6s A: %02x F: %s BC: %02x%02x DE: %02x%02x  HL: %02x%02x SP: %x \n", pcData, i.InstructionType, i.Destination, i.Source, c.Register.a, flags, c.Register.b, c.Register.c, c.Register.d, c.Register.e, c.Register.h, c.Register.l, c.Register.sp)

	if _, err := f.Write([]byte(output)); err != nil {
        log.Fatal(err)
    }
}

func PrintLog(c *CPU, i Instruction) {
	pcData := fmt.Sprintf("Pc: %x, (%02x %02x %02x) -> ", c.Register.pc, c.currentOpcode, c.MMURead(c.Register.pc+1), c.MMURead(c.Register.pc+2))
	flags := fmt.Sprintf("%c%c%c%c", c.FormatFlag(flagZ, 'Z'), c.FormatFlag(flagN, 'N'), c.FormatFlag(flagH, 'H'), c.FormatFlag(flagC, 'C'))
	output := fmt.Sprintf("%s Inst: %-6s Dest: %-6s Src: %-6s A: %02x F: %s BC: %02x%02x DE: %02x%02x  HL: %02x%02x SP: %x \n", pcData, i.InstructionType, i.Destination, i.Source, c.Register.a, flags, c.Register.b, c.Register.c, c.Register.d, c.Register.e, c.Register.h, c.Register.l, c.Register.sp)

	fmt.Print(output)
}

func DoctorLog(c *CPU, f *os.File) {
	doctor := fmt.Sprintf("A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n", c.Register.a, c.Register.f, c.Register.b, c.Register.c, c.Register.d, c.Register.e, c.Register.h, c.Register.l, c.Register.sp, c.Register.pc, c.MMURead(c.Register.pc), c.MMURead(c.Register.pc+1), c.MMURead(c.Register.pc+2), c.MMURead(c.Register.pc+3))
	
	//fmt.Print(doctor)

	if _, err := f.Write([]byte(doctor)); err != nil {
       log.Fatal(err)
    }
}

func GetInstructionsCycles(c *CPU) {
	for i := 0; i <= 0xFF; i++ {
		v, ok := instructions[uint8(i)]
		if ok {
			c.SourceTarget = v.Source
			c.DestinationTarget = v.Destination
			c.currentOpcode = uint8(i)	
			c.CurrentConditionResult = true
			if v.ConditionType != cond_None{
				cycles, _ := c.ExecuteInstruction(v)
				fmt.Printf("%x True %d\n",i, cycles)
				c.CurrentConditionResult = false 
				cycles, _ = c.ExecuteInstruction(v)
				fmt.Printf("%x False %d\n",i, cycles)
			} else{
				cycles, _ := c.ExecuteInstruction(v)
				fmt.Printf("%x %d\n",i, cycles)
			}
		}

	}
}