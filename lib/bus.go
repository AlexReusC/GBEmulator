package lib

type Bus struct {
}

func (b *Bus) BusRead(a uint16, c *Cart) uint8 {
	if a < 0x0800 {
		return c.CartRead(a)
	}

	return 0x0000
}
