package lib

type Serial struct {
	data    uint8
	control uint8
}

// TODO: remove clock from here
func (s *Serial) SerialRead(a uint16) uint8 {
	if a == 0xFF01 {
		return s.data
	}

	if a == 0xFF02 {
		return s.control
	}

	//fmt.Printf("Unsupported serial read %x\n", a)
	return 0
}

func (s *Serial) SerialWrite(a uint16, v uint8) {
	if a == 0xFF01 {
		s.data = v
		return
	}
	if a == 0xFF02 {
		s.control = v
		return
	}

	//fmt.Printf("Unsupported serial read %x\n", a)
}