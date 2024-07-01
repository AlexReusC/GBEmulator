package lib

import "fmt"

type Bus struct {
	cart *Cart
	cpu *CPU
	wram [0x2000]uint8
	hram [0x80]uint8
}

func LoadBus(rb *Cart,  c *CPU) (*Bus, error) {
	b := &Bus{cart: rb, cpu: c}

	return b, nil
}

func (b *Bus) BusRead(a uint16) uint8 {
	// ROM data
	if a < 0x0800 {
		return b.cart.CartRead(a)
	}
	// Video RAM 
	if a < 0xA000 {
		fmt.Println("Bus read not implemented", a)
		panic(0)
		return 0x0000
	}
	// Cartridge/external RAM
	if a < 0xC000 {
		return b.cart.CartRead(a)
	}
	// Working RAM
	if a < 0xE000 {
		return b.WramRead(a)
	}
	// Echo RAM (prohibited)
	if a < 0xFE00 {
		return 0
	}
	// Object attribute memory 
	if a < 0xFEA0 {
		fmt.Println("Bus read not implemented", a)
		panic(0)
		return 0 //TODO
	}
	// Reserved (prohibited)
	if a < 0xFF00 {
		return 0
	}
	// IO registers
	if a < 0xFF80 {
		fmt.Println("Bus read not implemented", a)
		panic(0)
		return 0 //TODO
	}
	// High RAM
	if a < 0xFFFF {
		return b.HramRead(a)
	}
	// CPU enable registerr
	if a == 0xFFFF {
		return b.cpu.GetIeRegister()
	}
	return 0
}

func (b *Bus) BusWrite(a uint16, v uint8) {
	if a < 0x0800 {
		b.cart.CartWrite(a, v)
	}	// Video RAM 
	if a < 0xA000 {
		fmt.Printf("Bus write not implemented %x\n", a)
		panic(0)
		return 
	}
	// Cartridge/external RAM
	if a < 0xC000 {
		b.cart.CartWrite(a, v)
		return
	}
	// Working RAM
	if a < 0xE000 {
		b.WramWrite(a, v)
		return
	}
	// Echo RAM (prohibited)
	if a < 0xFE00 {
		return
	}
	// Object attribute memory 
	if a < 0xFEA0 {
		fmt.Printf("Bus write not implemented %x\n", a)
		panic(0)
		return //TODO
	}
	// Reserved (prohibited)
	if a < 0xFF00 {
		return
	}
	// IO registers
	if a < 0xFF80 {
		fmt.Printf("Bus write not implemented %x\n", a)
		//panic(0)
		return //TODO
	}
	// High RAM
	if a < 0xFFFF {
		b.HramWrite(a, v)
		return
	}
	// CPU enable registerr
	if a == 0xFFFF {
		b.cpu.SetIeRegister(v)
		return
	}

	fmt.Println("Bus write not implemented")
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


func (b *Bus) WramRead(a uint16) uint8 {
	return b.wram[a-0xC000]
}

func (b *Bus) WramWrite(a uint16, v uint8) {
	b.wram[a-0xC000] = v
}

func (b *Bus) HramRead(a uint16) uint8 {
	return b.hram[a-0xFF80]
}

func (b *Bus) HramWrite(a uint16, v uint8) {
	b.hram[a-0xFF80] = v
}