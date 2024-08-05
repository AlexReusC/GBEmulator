package main

import (
	"fmt"
	"gbemulator/lib"
	"log"
	"os"
)

func main() {
	//logging
	f, err := os.Create("../gameboy-doctor/debug.txt")
    if err != nil {
        log.Fatal(err)
    }
	defer f.Close()

	if len(os.Args) <= 1 {
		fmt.Println("no file passed")
	}
	p := os.Args[1]

	emu, err := lib.LoadEmulator(lib.WithFile(f), lib.WithCart(p))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		emu.Run()
	}
}