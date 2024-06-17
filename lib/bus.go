package lib

import "fmt"

type Bus struct {
	cart *Cart
	ram *Ram
}

func LoadBus(c *Cart, r *Ram) (*Bus, error) {
	b := &Bus{cart: c, ram: r}

	return b, nil
}

func (b *Bus) BusRead(a uint16) uint8 {
	// ROM data
	if a < 0x0800 {
		return b.cart.CartRead(a)
	}
	// Video RAM 
	if a < 0xA000 {
		fmt.Println("Bus read not implemented")
		return 0x0000
	}
	// Cartridge/external RAM
	if a < 0xC000 {
		return b.cart.CartRead(a)
	}
	// Working RAM
	if a < 0xE000 {
		return b.ram.WramRead(a)
	}
	// Echo RAM (prohibited)
	if a < 0xFE00 {
		return 0
	}
	// Object attribute memory 
	if a < 0xFEA0 {
		fmt.Println("Bus read not implemented")
		return 0 //TODO
	}
	// Reserved (prohibited)
	if a < 0xFF00 {
		return 0
	}
	// IO registers
	if a < 0xFF80 {
		fmt.Println("Bus read not implemented")
		return 0 //TODO
	}
	// High RAM
	if a < 0xFFFF {
		return b.ram.HramRead(a)
	}
	// CPU enable registerr
	if a == 0xFFFF {
		fmt.Println("Bus read not implemented")
		return 0	//TODO
	}
	return 0
}

func (b *Bus) BusWrite(a uint16, v uint8) {
	if a < 0x0800 {
		b.cart.CartWrite(a, v)
	}	// Video RAM 
	if a < 0xA000 {
		fmt.Println("Bus read not implemented")
		return 
	}
	// Cartridge/external RAM
	if a < 0xC000 {
		b.cart.CartWrite(a, v)
		return
	}
	// Working RAM
	if a < 0xE000 {
		b.ram.WramWrite(a, v)
		return
	}
	// Echo RAM (prohibited)
	if a < 0xFE00 {
		return
	}
	// Object attribute memory 
	if a < 0xFEA0 {
		fmt.Println("Bus read not implemented")
		return //TODO
	}
	// Reserved (prohibited)
	if a < 0xFF00 {
		return
	}
	// IO registers
	if a < 0xFF80 {
		fmt.Println("Bus read not implemented")
		return //TODO
	}
	// High RAM
	if a < 0xFFFF {
		b.ram.HramWrite(a, v)
		return
	}
	// CPU enable registerr
	if a == 0xFFFF {
		fmt.Println("Bus read not implemented")
		return	//TODO
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