package main

import (
	"fmt"
	"gbemulator/lib"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("no file passed")
		return
	}
	p := os.Args[1]
	//logging
	f, err := os.Create("../gameboy-doctor/debug.txt")
    if err != nil {
        log.Fatal(err)
    }
	defer f.Close()
	e, err := lib.LoadEmulator(lib.WithFile(f), lib.WithCart(p))
	if err != nil {
		fmt.Println(err)
		return
	}

	lib.RunGame(e)
}