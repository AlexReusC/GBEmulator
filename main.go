package main

import (
	"fmt"
	"gbemulator/lib"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("no file passed")
		return
	}
	file := os.Args[1]
	//logging for gb doctor
	//f, err := os.Create("../gameboy-doctor/debug.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	e, err := lib.LoadEmulator(lib.WithCart(file))
	if err != nil {
		fmt.Println(err)
		return
	}

	lib.RunGame(e)
}
