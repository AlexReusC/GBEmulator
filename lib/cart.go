package lib

import (
	"encoding/binary"
	"fmt"
	"os"
)

type header struct {
	EntryPoint      [0x4]uint8
	Logo            [0x30]uint8
	Title           [16]uint8
	NewLicenseeCode uint16
	SgbFlag         uint8
	CartridgeType   uint8
	RomSize         uint8
	RamSize         uint8
	DestinationCode uint8
	OldLicenseeCode uint8
	MaskRomVersion  uint8
	HeaderChecksum  uint8
	GlobalChecksum  uint16
}

type Cart struct {
	Header header
	Length int64
	Rom    []uint8
	mbc1   *MBC1 //TODO: make this generic
}

func LoadCart(p string) (*Cart, error) {
	file, err := os.Open(p)
	if err != nil {
		fmt.Println("Failed to open")
		return nil, err
	}
	defer file.Close()

	cart := &Cart{}

	fi, err := file.Stat()
	if err != nil {
		fmt.Println("Couldn't get info")
		return nil, err
	}
	cart.Length = fi.Size()

	file.Seek(0x0100, 0)
	cartHeader := header{}
	if err := binary.Read(file, binary.LittleEndian, &cartHeader); err != nil {
		fmt.Println("Invalid header")
		return nil, err
	}
	cart.Header = cartHeader
	file.Seek(0, 0)

	cartRom := make([]uint8, cart.Length)
	if err := binary.Read(file, binary.LittleEndian, &cartRom); err != nil {
		fmt.Println("Invalid rom", err)
		return nil, err
	}
	cart.Rom = cartRom

	if err := cart.initMBC(cartRom); err != nil {
		return nil, err
	}

	fmt.Printf("Title: %s\n", cart.Header.Title)
	fmt.Printf("Type: % x\n", cart.Header.CartridgeType)
	fmt.Printf("Nintendo logo: % x\n", cart.Header.Logo)
	fmt.Printf("Rom: %d KB\n", 32<<cart.Header.RomSize)
	fmt.Printf("Ram: %x\n", cart.Header.RamSize)
	fmt.Printf("Lic Code: %x\n", cart.Header.OldLicenseeCode)
	fmt.Printf("Rom Version: %x\n", cart.Header.MaskRomVersion)
	fmt.Printf("Cart length: %x\n", cart.Length)
	fmt.Printf("Cartridge type: 0x%x\n", cart.Header.CartridgeType)

	return cart, nil

}

func (c *Cart) initMBC(rom []uint8) error {
	switch c.Header.CartridgeType {
	case 0x00:
		c.Rom = rom
		return nil
	case 0x01:
		c.mbc1 = loadMBC1(rom)
	default:
		return fmt.Errorf("Unsuported mbc")
	}

	return nil
}

func (c *Cart) CartRead(a uint16) uint8 {
	if c.Header.CartridgeType == 0x01 {
		return c.mbc1.read(a)
	}
	return c.Rom[a]
}

func (c *Cart) CartWrite(a uint16, v uint8) {
	if c.Header.CartridgeType == 0x01 {
		c.mbc1.write(a, v)
	} else {
		c.Rom[a] = v
	}
}
