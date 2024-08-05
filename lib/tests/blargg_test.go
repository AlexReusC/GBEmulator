package lib

import (
	"gbemulator/lib"
	"regexp"
	"testing"
)

func TestBlarggTests(t *testing.T) {
	emu, _ := lib.LoadEmulator(lib.WithCart("../../roms/03-op sp,hl.gb"))
	for cycles := 0; cycles < 2_000_000; cycles++ {
		emu.Run()
	}

	cpu := emu.Cpu

	r, _ := regexp.Compile("Passed")
	if !r.MatchString(cpu.Debug.GetMsg()){
		t.Fatal("->", cpu.Debug.GetMsg()[10:30])
	}
}