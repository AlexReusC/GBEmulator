package lib

import "fmt"

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