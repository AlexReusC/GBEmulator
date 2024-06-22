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

	cpu, err := LoadCpu()
	if err != nil {
		return
	}

	//Probably should make that the cpu creates the bus
	bus, err := LoadBus(cart, ram, cpu)
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