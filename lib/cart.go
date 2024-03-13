package lib

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Header struct {
	_ [0x0100]uint8
	EntryPoint [0x4]uint8
	Logo [0x30]uint8
	Title [16]uint8
	NewLicenseeCode uint16
	SgbFlag uint8
	CartridgeType uint8
	RomSize uint8
	RamSize uint8
	DestinationCode uint8
	OldLicenseeCode uint8
	MaskRomVersion uint8
	HeaderChecksum uint8
	GlobalChecksum uint16
}

type Cart struct{

}

func (c *Cart) LoadCart() {
	file, err := os.Open("./roms/dmg-acid2.gb")
	if err != nil{
		fmt.Println("Failed to open")
	}
	defer file.Close()


	myHeader := Header{}
	if err := binary.Read(file, binary.LittleEndian, &myHeader); err != nil {
		fmt.Println("Invalid file")
	}


	fmt.Printf("Title: %s\n", myHeader.Title)
	fmt.Printf("Type: % x\n", myHeader.CartridgeType)	
	fmt.Printf("Nintendo logo: % x\n", myHeader.Logo)
	fmt.Printf("Rom: %d KB\n", 32 << myHeader.RomSize)
	fmt.Printf("Ram: %x\n", myHeader.RamSize)
	fmt.Printf("Lic Code: %x\n", myHeader.OldLicenseeCode)
	fmt.Printf("Rom Version: %x\n", myHeader.MaskRomVersion)

}