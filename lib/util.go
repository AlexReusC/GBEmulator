package lib

func BoolToUint(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func Union16(h, l uint8) uint16 {
	return uint16(h)<<8 | uint16(l)
}