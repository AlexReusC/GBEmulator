package lib

import "fmt"

//TODO: Change to class emulator
func Run() {
	cart, err := LoadCart()
	if err != nil {
		return
	}

	bus, err := LoadBus(cart)
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