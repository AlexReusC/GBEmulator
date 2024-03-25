package lib

func Run() {
	cart, err := LoadCart()
	if err != nil {
		return
	}

	bus := &Bus{}

	cpu, err := LoadCpu()
	if err != nil {
		return
	}

	for {
		cpu.Step(bus, cart)
	}

}