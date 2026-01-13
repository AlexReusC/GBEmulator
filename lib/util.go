package lib

import "slices"

func BoolToUint(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func Union16(h, l uint8) uint16 {
	return uint16(h)<<8 | uint16(l)
}

func BitIsSet(n uint8, b uint8) bool {
	if n&(1<<b) != 0 {
		return true
	}
	return false
}

func SetBitWithCond(b uint8, n int, c bool) uint8 {
	if c {
		return (b | (1 << n))
	} else {
		return (b & ^(1 << n))
	}
}

func SetBit(b uint8, n int) uint8 {
	return b | (1 << n)
}

func UnsetBit(b uint8, n int) uint8 {
	return b & ^(1 << n)
}

func IsImmediateTarget8(t target) bool {
	var pt = []target{n, n_M, SPe8}
	return slices.Contains(pt, t)
}

func IsImmediateTarget16(t target) bool {
	var pt = []target{nn, nn_M, nn_M16}
	return slices.Contains(pt, t)
}

func Isr8(t target) bool {
	var r8 = []target{A, B, C, D, E, F, H, L}
	return slices.Contains(r8, t)
}

func IsPointer(t target) bool {
	var pt = []target{C_M, BC_M, DE_M, HL_M, HLP_M, HLM_M, n_M, nn_M, nn_M16}
	return slices.Contains(pt, t)
}
