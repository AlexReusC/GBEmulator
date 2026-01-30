package lib

type Serial struct {
	data    uint8
	control uint8
}

func (s *Serial) SerialRead(a uint16) uint8 {
	if a == 0xFF00 {
		return 0xFF //TODO: this for testing, implement later
	}

	if a == 0xFF01 {
		return s.data
	}

	if a == 0xFF02 {
		return s.control
	}

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

}
