package lib

import "fmt"

func Run() {
	cart, err := LoadCart()
	if err != nil {
		return
	}

	ram, err := LoadRam()
	if err != nil {
		return
	}

	bus, err := LoadBus(cart, ram)
	if err != nil {
		return
	}

	cpu, err := LoadCpu()
	if err != nil {
		return
	}

	for {
		err := cpu.Step(bus)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}