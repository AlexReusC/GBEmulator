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

	return 0x0000
}
