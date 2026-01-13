package lib

const CLOCKSPEED = 4_194_304

type Clock struct {
	Divider      uint16 //Div
	Counter      uint8  //Tima
	CounterTimer int
	Modulo       uint8 //Tma
	Control      uint8 //Tac/TMC
}

func LoadClock() (*Clock, error) {
	clock := &Clock{Divider: 0xABCC}
	return clock, nil
}

func (c *Clock) SetClockFrequency() {
	f := c.Control & 0x03
	switch f {
	case 0x00:
		c.CounterTimer = 1024
	case 0x01:
		c.CounterTimer = 16
	case 0x10:
		c.CounterTimer = 64
	case 0x11:
		c.CounterTimer = 256
	}
}

func (c *Clock) Update(cycles int) bool {
	setTimerInterruptor := false
	for i := 0; i < cycles*4; i++ {
		prevDiv := c.Divider
		c.Divider++

		timerUpdate := false

		switch c.Control & (0b11) {
		case 0b00:
			bPrevDiv := (prevDiv & (1 << 9)) != 0
			bDiv := (c.Divider & (1 << 9)) != 0
			timerUpdate = bPrevDiv && !bDiv
		case 0b01:
			bPrevDiv := (prevDiv & (1 << 3)) != 0
			bDiv := (c.Divider & (1 << 3)) != 0
			timerUpdate = bPrevDiv && !bDiv
		case 0b10:
			bPrevDiv := (prevDiv & (1 << 5)) != 0
			bDiv := (c.Divider & (1 << 5)) != 0
			timerUpdate = bPrevDiv && !bDiv
		case 0b11:
			bPrevDiv := (prevDiv & (1 << 7)) != 0
			bDiv := (c.Divider & (1 << 7)) != 0
			timerUpdate = bPrevDiv && !bDiv
		}

		if timerUpdate && (c.Control&(1<<2) != 0) {
			c.Counter += 1

			//Update counter
			if c.Counter == 0xFF {
				c.Counter = c.Modulo
				setTimerInterruptor = true
			}
		}
	}
	return setTimerInterruptor
}

func (c *Clock) Write(a uint16, v uint8) {
	switch a {
	case 0xFF04:
		c.Divider = 0
	case 0xFF05:
		c.Counter = v
	case 0xFF06:
		c.Modulo = v
	case 0xFF07:
		c.Control = v

	default:
		panic(0)
	}
}

func (c *Clock) Read(a uint16) uint8 {
	switch a {
	case 0xFF04:
		return uint8(c.Divider >> 8)
	case 0xFF05:
		return c.Counter
	case 0xFF06:
		return c.Modulo
	case 0xFF07:
		return c.Control
	default:
		panic(0)
	}
}
