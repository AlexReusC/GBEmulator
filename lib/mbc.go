package lib

import "fmt"

type MBC1 struct {
	rom []uint8

	ramEnabled bool
	romBank    uint8
}

func loadMBC1(rom []uint8) *MBC1 {
	mbc1 := &MBC1{
		rom:        rom,
		romBank:    1,
		ramEnabled: false,
	}
	return mbc1
}

func (m *MBC1) read(a uint16) uint8 {
	switch {
	case a < 0x4000: //ROM bank 00
		return m.rom[a]
	case a < 0x8000: //ROM bank 01-7F
		offset := (a - 0x4000) + (uint16(m.romBank) * 0x4000)
		return m.rom[offset]
	case a >= 0xA000 && a < 0xC000: //RAM bank 00-03
		//TODO
		return m.rom[a]
	default:
		return m.rom[a]
	}
}

func (m *MBC1) write(a uint16, v uint8) {
	fmt.Printf("rom write -> a: %x, v: %x\n", a, v)
	switch {
	case a < 0x2000: //switch ram, value "intercepted"
		if v == 0x0A {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}
	case a < 0x4000: //rom bank
		v = v & 0x1F
		if v == 0x00 || v == 0x20 || v == 0x40 || v == 0x61 {
			v++
		}
		m.romBank = v
	case a < 0x6000: //upper bits bank number
	case a < 0x8000: //rom/ram mode
	case a >= 0xA000 && a < 0xC000: //ram banks
	}
}
