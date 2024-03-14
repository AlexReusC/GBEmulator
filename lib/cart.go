package lib

import (
	"encoding/binary"
	"fmt"
	"os"
)

type header struct {
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
	Header header
	Length int64
	Rom []byte
}

func (c *Cart) LoadCart() {
	file, err := os.Open("./roms/dmg-acid2.gb")
	if err != nil{
		fmt.Println("Failed to open")
		return
	}
	defer file.Close()

	myCart := Cart{}

	fi, err := file.Stat()
	if err != nil {
		fmt.Println("Couldn't get info")
		return
	}
	myCart.Length = fi.Size()


	file.Seek(0x0100, 0)
	cartHeader := header{}
	if err := binary.Read(file, binary.LittleEndian, &cartHeader); err != nil {
		fmt.Println("Invalid header")
		return
	}
	myCart.Header = cartHeader
	file.Seek(0, 0)

	cartRom := make([]byte, myCart.Length)	
	if err := binary.Read(file, binary.LittleEndian, &cartRom); err != nil {
		fmt.Println("Invalid rom", err)
		return
	}
	myCart.Rom = cartRom

	fmt.Printf("Title: %s\n", myCart.Header.Title)
	fmt.Printf("Type: % x\n", myCart.Header.CartridgeType)	
	fmt.Printf("Nintendo logo: % x\n", myCart.Header.Logo)
	fmt.Printf("Rom: %d KB\n", 32 <<  myCart.Header.RomSize)
	fmt.Printf("Ram: %x\n", myCart.Header.RamSize)
	fmt.Printf("Lic Code: %x\n", myCart.Header.OldLicenseeCode)
	fmt.Printf("Rom Version: %x\n", myCart.Header.MaskRomVersion)

}