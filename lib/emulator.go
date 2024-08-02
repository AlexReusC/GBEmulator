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
		fmt.Println(err)
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

	clock, err := LoadClock(cpu)
	if err != nil {
		return
	}

	bus.BusLoadClock(clock)

	/*
	for i := 0; i <= 0xFF; i++ {
		v, ok := instructions[uint8(i)]
		if ok {
			cpu.SourceTarget = v.Source
			cpu.DestinationTarget = v.Destination
			cpu.currentOpcode = uint8(i)	
			cpu.CurrentConditionResult = true
			if v.ConditionType != cond_None{
				cycles, _ := cpu.ExecuteInstruction(v)
				fmt.Printf("%x True %d\n",i, cycles)
				cpu.CurrentConditionResult = false 
				cycles, _ = cpu.ExecuteInstruction(v)
				fmt.Printf("%x False %d\n",i, cycles)
			} else{
				cycles, _ := cpu.ExecuteInstruction(v)
				fmt.Printf("%x %d\n",i, cycles)
			}

		}
		
	}
	*/	
	for {
		cycles, err := cpu.Step(f)
		if err != nil {
			fmt.Println(err)
			return
		}
			clock.Update(cycles)
		//TODO: GPU
		cpu.HandleInterrupts()
	}

}