package lib

import (
	"encoding/binary"
	"errors"
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
	Rom []uint8
}

func LoadCart() (*Cart, error) {
	if len(os.Args) <= 1 {
		return nil, errors.New("no file passed")
	}
	
	file, err := os.Open(os.Args[1])
	if err != nil{
		fmt.Println("Failed to open")
		return nil, err
	}
	defer file.Close()

	myCart := &Cart{}

	fi, err := file.Stat()
	if err != nil {
		fmt.Println("Couldn't get info")
		return nil, err
	}
	myCart.Length = fi.Size()


	file.Seek(0x0100, 0)
	cartHeader := header{}
	if err := binary.Read(file, binary.LittleEndian, &cartHeader); err != nil {
		fmt.Println("Invalid header")
		return nil, err
	}
	myCart.Header = cartHeader
	file.Seek(0, 0)

	cartRom := make([]uint8, myCart.Length)	
	if err := binary.Read(file, binary.LittleEndian, &cartRom); err != nil {
		fmt.Println("Invalid rom", err)
		return nil, err
	}
	myCart.Rom = cartRom

	fmt.Printf("Title: %s\n", myCart.Header.Title)
	fmt.Printf("Type: % x\n", myCart.Header.CartridgeType)	
	fmt.Printf("Nintendo logo: % x\n", myCart.Header.Logo)
	fmt.Printf("Rom: %d KB\n", 32 <<  myCart.Header.RomSize)
	fmt.Printf("Ram: %x\n", myCart.Header.RamSize)
	fmt.Printf("Lic Code: %x\n", myCart.Header.OldLicenseeCode)
	fmt.Printf("Rom Version: %x\n", myCart.Header.MaskRomVersion)
	//fmt.Printf("Rom: %x\n", myCart.Rom)

	return myCart, nil

}

func (c *Cart) CartRead(a uint16) uint8 {
	return c.Rom[a]
}

func (c *Cart) CartWrite(a uint16, val uint8) {
	fmt.Printf("Cart write not implemented, %x\n", a)
	c.Rom[a] = val
}