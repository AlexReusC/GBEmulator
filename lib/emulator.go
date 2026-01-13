package lib

import (
	"errors"
	"fmt"
	"os"
)

type Emulator struct {
	Cpu  *CPU
	cart *Cart
	file *os.File
	ppu  *PPU
	mmu  *MMU

	cpuCycles int
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

// Initialize emulator and main systems
// TODO: still a lot of refactor
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
	emulator.mmu = b

	debug := LoadDebug()

	cpu, err := LoadCpu(emulator.mmu, debug, clock)
	if err != nil {
		return nil, errors.New("cpu failed")
	}
	emulator.Cpu = cpu

	emulator.cpuCycles = 0

	return emulator, nil
}

// Main emulator loop
func (e *Emulator) Run() {
	if e.cpuCycles <= 0 {
		cycles, err := e.Cpu.Step(e.file)
		if err != nil {
			fmt.Println(err)
			return
		}
		e.Cpu.UpdateClock(cycles)
		e.ppu.Update(cycles)
		e.Cpu.HandleInterrupts()
		e.cpuCycles += cycles
	}
	e.cpuCycles--
}
