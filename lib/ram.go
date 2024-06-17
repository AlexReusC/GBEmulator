package lib

type Ram struct {
	Wram [0x2000]uint8
	Hram [0x80]uint8
}

func LoadRam() (*Ram, error) {
	r := &Ram{}

	return r, nil
}

func (r *Ram) WramRead(a uint16) uint8 {
	return r.Wram[a-0xC000]
}

func (r *Ram) WramWrite(a uint16, v uint8) {
	r.Wram[a-0xC000] = v
}

func (r *Ram) HramRead(a uint16) uint8 {
	return r.Wram[a-0xFF80]
}

func (r *Ram) HramWrite(a uint16, v uint8) {
	r.Wram[a-0xFF80] = v
}
