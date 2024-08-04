package lib

import "fmt"

type Bus struct {
	cart *Cart
	wram [0x2000]uint8
	hram [0x80]uint8
	serial *Serial
	ieRegister uint8
	interruptorFlags uint8
	clock *Clock
}

func LoadBus(rb *Cart, s *Serial, c *Clock) (*Bus, error) {
	b := &Bus{cart: rb, serial: s, clock: c}

	return b, nil
}


func (b *Bus) BusRead(a uint16) uint8 {
	switch {
	case a < 0x8000: // ROM data
		return b.cart.CartRead(a)
	case a < 0xA000: // Video RAM
		//fmt.Printf("Bus read not implemented %x\n", a)
		return 0x0000
	case a < 0xC000: // Cartridge/external RAM
		return b.cart.CartRead(a)
	case a < 0xE000: // Working RAM
		return b.WramRead(a)
	case a < 0xFE00: // Echo RAM (prohibited)
		return 0
	case a < 0xFEA0: // Object attribute memory
		//fmt.Printf("Bus read not implemented %x\n", a)
		return 0 //TODO
	case a < 0xFF00: // Reserved (prohibited)
		return 0
	case a < 0xFF03: // IO registers
		return b.serial.SerialRead(a, b.clock)
	case a >= 0xFF04 && a <= 0xFF07:
		return b.clock.Read(a)
	case a == 0xFF0F:
		return b.interruptorFlags
	case a < 0xFF44:
		//fmt.Println("address not implemented")
		return 0
	case a == 0xFF44: //GPU
		return 0x90
	case a == 0xFF4D:
		return 0xFF
	case a < 0xFF80:
		//fmt.Println("address not implemented")
		return 0
	case a < 0xFFFF: // High RAM
		return b.HramRead(a)
	case a == 0xFFFF: // CPU enable registerr
		return b.GetIeRegister()
	default:
		return 0
	}
}

func (b *Bus) BusWrite(a uint16, v uint8) {
	switch {
	case a < 0x8000:
		b.cart.CartWrite(a, v)
	case a < 0xA000: // Video RAM 
		//fmt.Printf("Bus write not implemented %x\n", a)
	case a < 0xC000: // Cartridge/external RAM
		b.cart.CartWrite(a, v)
	case a < 0xE000: // Working RAM
		b.WramWrite(a, v)
	case a < 0xFE00: // Echo RAM (prohibited)
		return
	case a < 0xFEA0: // Object attribute memory
		//fmt.Printf("Bus write not implemented %x\n", a)
		return //TODO
	case a < 0xFF00: // Reserved (prohibited)
		return 
	case a < 0xFF03:// IO registers
		b.serial.SerialWrite(a, v, b.clock)
	case a >= 0xFF04 && a <= 0xFF07:
		b.clock.Write(a, v)
	case a == 0xFF0F:
		b.interruptorFlags = v
	case a < 0xFF80:
		//fmt.Println("address not implemented")
	case a >= 0xFF80 && a < 0xFFFF: // High RAM
		b.HramWrite(a, v)
	case a == 0xFFFF: // CPU enable registerr
		b.SetIeRegister(v)
	default:
		fmt.Println("Bus write unavailable", a)
	}
}

func (b *Bus) BusRead16(a uint16) uint16 {
	lo := uint16(b.BusRead(a))
	hi := uint16(b.BusRead(a + 1))
	return (hi << 8) | lo
}

func (b *Bus) BusWrite16(a uint16, v uint16) {
	b.BusWrite(a+1, uint8((v>>8)&0xFF))
	b.BusWrite(a, uint8(v&0xFF))
}


func (b *Bus) WramRead(a uint16) uint8 { return b.wram[a-0xC000] }
func (b *Bus) WramWrite(a uint16, v uint8) { b.wram[a-0xC000] = v }
func (b *Bus) HramRead(a uint16) uint8 { return b.hram[a-0xFF80] }
func (b *Bus) HramWrite(a uint16, v uint8) { b.hram[a-0xFF80] = v }
func (b *Bus) GetIeRegister() uint8 { return b.ieRegister }
func (b *Bus) SetIeRegister(ir uint8) { b.ieRegister = ir }