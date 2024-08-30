package lib

import (
	"errors"
	"fmt"
	"os"
)

type Emulator struct {
	Cpu *CPU
	cart *Cart
	file *os.File
	ppu *PPU
	bus *Bus
}

func WithFile(f *os.File) func(e *Emulator) {
	return func(e *Emulator) {
		e.file = f
	}
}

func WithCart(p string) func(e *Emulator) {
	return func(e *Emulator) {
		cart, err := LoadCart(p)
		if err != nil {
			panic(err)
		}
		e.cart = cart
	}
}

//TODO: still a lot of refactor
func LoadEmulator(options ...func(*Emulator)) (*Emulator, error) {
	emulator := new(Emulator)

	for _, o := range options {
		o(emulator)
	}

	clock, err := LoadClock()
	if err != nil {
		return nil, errors.New("clock failed")
	}

	ppu, err := LoadPpu()
	if err != nil {
		return nil, errors.New("ppu failed")
	}
	emulator.ppu = ppu 

	serial := &Serial{data: 0, control: 0}
	b, err := LoadBus(emulator.cart, serial, clock, emulator.ppu)
	if err != nil {
		return nil, errors.New("bus failed")
	}
	emulator.bus = b

	debug := LoadDebug()

	cpu, err := LoadCpu(emulator.bus, debug, clock)
	if err != nil {
		return nil, errors.New("cpu failed")
	}
	emulator.Cpu = cpu

	return emulator, nil 
}

func (e *Emulator) Run() {
	cycles, err := e.Cpu.Step(e.file)
	if err != nil {
		fmt.Println(err)
		return
	}
	e.Cpu.UpdateClock(cycles)
	e.ppu.Update(cycles)
	e.Cpu.HandleInterrupts()
}