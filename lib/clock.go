package lib

const CLOCKSPEED = 4_194_304

type Clock struct {
	Divider      uint8 //Div
	DividerTimer uint16
	Counter      uint8 //Tima
	CounterTimer uint16
	Modulo       uint8 //Tma
	Control      uint8 //Tac/TMC
	CPU          *CPU
}

func LoadClock(c *CPU) (*Clock, error) {
	clock := &Clock{CPU: c}
	return clock, nil
}

func (c *Clock) Update(cycles int) {
	//update divider timer
	c.DividerTimer += uint16(cycles)

	//check if clock enabled
	if c.Control&0x04 == 0 {
		return
	}

	//Update countertimer

	//Check counter
	//Update counter
	//check overflow
	if c.Counter == 0xFF {
		c.Counter = c.Modulo
		c.CPU.RequestInterrupt(TIMER)
	} else {
		c.Counter += 1
	}
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
		return c.Divider
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