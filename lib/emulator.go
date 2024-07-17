package lib

import (
	"fmt"
	"log"
	"os"
)

//TODO: Change to class emulator
func Run() {
	cart, err := LoadCart()
	if err != nil {
		return
	}

	serial := &Serial{data: 0, control: 0}

	bus, err := LoadBus(cart, serial)
	if err != nil {
		return
	}

	//logging
	f, err := os.Create("../gameboy-doctor/debug.txt")
    if err != nil {
        log.Fatal(err)
    }
	defer f.Close()

	debug := LoadDebug()

	cpu, err := LoadCpu(bus, debug)
	if err != nil {
		return
	}

	for {
		err := cpu.Step(f)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}