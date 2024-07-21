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

//TODO: Separate both functions, this is temporal

func SetBit(b uint8, n int, c bool) uint8 {
	if c {
		b |= (1 << n)
	} else {
		b &= ^(1 << n)
	}
	return b
}

func UnsetBit(b uint8, n int) uint8 {
	return SetBit(b, n, false)
}