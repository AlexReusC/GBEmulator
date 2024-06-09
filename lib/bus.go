package lib

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

	return 0x0000
}

func (b *Bus) BusRead16(a uint16) uint16 {
	//TODO
	return 0
}

func (b *Bus) BusWrite(a uint16, v uint8) {
	//TODO
}

func (b *Bus) BusWrite16(a uint16, v uint16) {
	//TODO
}