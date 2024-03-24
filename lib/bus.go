package lib

type Bus struct {
}

func BusRead(a uint16) uint8 {
	if a < 0x0800 {
		//return CartRead(a)
	}

	return 0x0000
}
