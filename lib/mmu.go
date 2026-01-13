package lib

import "fmt"

type MMU struct {
	cart             *Cart
	wram             [0x2000]uint8
	hram             [0x80]uint8
	serial           *Serial
	ieRegister       uint8
	interruptorFlags uint8
	clock            *Clock
	ppu              *PPU
}

func LoadBus(rb *Cart, s *Serial, c *Clock, p *PPU) (*MMU, error) {
	b := &MMU{cart: rb, serial: s, clock: c, ppu: p}

	p.MMU = b
	return b, nil
}

func (m *MMU) Read(a uint16) uint8 {
	switch {
	case a < 0x8000: // ROM data
		return m.cart.CartRead(a)
	case a < 0xA000: // Video RAM
		return m.ppu.VramRead(a)
	case a < 0xC000: // Cartridge/external RAM
		return m.cart.CartRead(a)
	case a < 0xE000: // Working RAM
		return m.WramRead(a)
	case a < 0xFE00: // Echo RAM (prohibited)
		return 0
	case a < 0xFEA0: // Object attribute memory
		return m.ppu.oamRead(a)
	case a < 0xFF00: // Reserved (prohibited)
		return 0
	case a < 0xFF03: // IO registers
		return m.serial.SerialRead(a)
	case a >= 0xFF04 && a <= 0xFF07:
		return m.clock.Read(a)
	case a == 0xFF0F:
		return m.interruptorFlags
	case a >= 0xFF40 && a <= 0xFF4B:
		return m.ppu.LcdRead(a)
	case a == 0xFF4D:
		return 0xFF
	case a < 0xFF80:
		return 0
	case a < 0xFFFF: // High RAM
		return m.HramRead(a)
	case a == 0xFFFF: // CPU enable registerr
		return m.GetIeRegister()
	default:
		return 0
	}
}

func (m *MMU) Write(a uint16, v uint8) {
	switch {
	case a < 0x8000:
		m.cart.CartWrite(a, v)
	case a < 0xA000: // Video RAM
		m.ppu.VramWrite(a, v)
	case a < 0xC000: // Cartridge/external RAM
		m.cart.CartWrite(a, v)
	case a < 0xE000: // Working RAM
		m.WramWrite(a, v)
	case a < 0xFE00: // Echo RAM (prohibited)
		return
	case a < 0xFEA0: // Object attribute memory
		m.ppu.oamwrite(a, v)
	case a < 0xFF00: // Reserved (prohibited)
		return
	case a < 0xFF03: // IO registers
		m.serial.SerialWrite(a, v)
	case a >= 0xFF04 && a <= 0xFF07:
		m.clock.Write(a, v)
	case a == 0xFF0F:
		m.interruptorFlags = v
	case a >= 0xFF40 && a <= 0xFF4B:
		if a == 0xFF46 {
			m.DmaTransfer(v)
		}
		m.ppu.LcdWrite(a, v)
	case a < 0xFF80:
	case a >= 0xFF80 && a < 0xFFFF: // High RAM
		m.HramWrite(a, v)
	case a == 0xFFFF: // CPU enable registerr
		m.SetIeRegister(v)
	default:
		fmt.Println("Bus write unavailable", a)
	}
}

func (m *MMU) DmaTransfer(a uint8) {
	realAddress := uint16(a) << 8
	for i := uint16(0); i < 0xA0; i++ {
		v := m.Read(realAddress + i)
		//fmt.Printf("memory %x %x\n", realAddress + i, v)
		m.ppu.oamwrite(0xFE00+i, v)
	}
}

func (m *MMU) Read16(a uint16) uint16 {
	lo := uint16(m.Read(a))
	hi := uint16(m.Read(a + 1))
	return (hi << 8) | lo
}

func (m *MMU) Write16(a uint16, v uint16) {
	m.Write(a+1, uint8((v>>8)&0xFF))
	m.Write(a, uint8(v&0xFF))
}

func (m *MMU) WramRead(a uint16) uint8     { return m.wram[a-0xC000] }
func (m *MMU) WramWrite(a uint16, v uint8) { m.wram[a-0xC000] = v }
func (m *MMU) HramRead(a uint16) uint8     { return m.hram[a-0xFF80] }
func (m *MMU) HramWrite(a uint16, v uint8) { m.hram[a-0xFF80] = v }
func (m *MMU) GetIeRegister() uint8        { return m.ieRegister }
func (m *MMU) SetIeRegister(ir uint8)      { m.ieRegister = ir }

func (m *MMU) RequestInterrupt(i InterruptorBit) {
	m.interruptorFlags = SetBit(m.interruptorFlags, int(i))
}
