package lib

import "fmt"

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
		err := cpu.Step(bus, cart)
		if err != nil {
			fmt.Println(err)
		}
	}

}