package lib

import (
	"fmt"
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

func (d *Debug) DebugUpdate(b *Bus) {
	if b.BusRead(0xFF02) == 0x81 {
		c := rune(b.BusRead(0xFF01))

		d.debugMsg[d.msgSize] = c
		d.msgSize += 1

		b.BusWrite(0xFF02, 0)
	}
}

func (d *Debug) DebugPrint() {
	if d.debugMsg[0] != 0 {
		fmt.Printf("DBG: %s\n", string(d.debugMsg))
	}
}

func Log(c *CPU, f *os.File) {
	//pcData := fmt.Sprintf("Pc: %x, (%02x %02x %02x) -> ", c.Register.pc, c.currentOpcode, c.BusRead(c.Register.pc+1), c.BusRead(c.Register.pc+2))

	doctor := fmt.Sprintf("A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n", c.Register.a, c.Register.f, c.Register.b, c.Register.c, c.Register.d, c.Register.e, c.Register.h, c.Register.l, c.Register.sp, c.Register.pc, c.BusRead(c.Register.pc), c.BusRead(c.Register.pc+1), c.BusRead(c.Register.pc+2), c.BusRead(c.Register.pc+3))
	
	//flags := fmt.Sprintf("%c%c%c%c", c.FormatFlag(flagZ, 'Z'), c.FormatFlag(flagN, 'N'), c.FormatFlag(flagH, 'H'), c.FormatFlag(flagC, 'C'))
	//output := fmt.Sprintf("%s Inst: %-6s Dest: %-6s Src: %-6s A: %02x F: %s BC: %02x%02x DE: %02x%02x  HL: %02x%02x SP: %x \n", pcData, instruction.InstructionType, instruction.Destination, instruction.Source, c.Register.a, flags, c.Register.b, c.Register.c, c.Register.d, c.Register.e, c.Register.h, c.Register.l, c.Register.sp)
	fmt.Print(doctor)

	// if _, err := f.Write([]byte(doctor)); err != nil {
    //     log.Fatal(err)
    // }

	// if c.printTimers {
	// flags := fmt.Sprintf("Div: %x, Div tim: %x, counter: %x, counter tim: %x, modulo: %x, control: %x\n", 
	// c.Bus.clock.Divider,
	// c.Bus.clock.Counter,
	// c.Bus.clock.CounterTimer,
	// c.Bus.clock.Modulo,
	// c.Bus.clock.Control)

	// if _, err := f.Write([]byte(flags)); err != nil {
    //     log.Fatal(err)
    // }
	// }
	
}