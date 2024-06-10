package lib

import "fmt"

type Bus struct {
	cart *Cart
}

func LoadBus(c *Cart) (*Bus, error) {
	b := &Bus{cart: c}

	return b, nil
}

func (b *Bus) BusRead(a uint16) uint8 {
	if a < 0x0800 {
		return b.cart.CartRead(a)
	}

	//TODO: more reads
	fmt.Println("Bus read ot implemented")
	return 0x0000
}

func (b *Bus) BusWrite(a uint16, v uint8) {
	if a < 0x0800 {
		b.cart.CartWrite(a, v)
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