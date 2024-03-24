package lib

func Run() {
	cart := &Cart{}
	cart.LoadCart()

	cpu, err := LoadCpu()
	if err != nil {
		return
	}

	for {
		cpu.Step()
	}

}