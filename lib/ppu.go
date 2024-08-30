package lib

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type PPUMode  int

const (
	priorityMask = 1 << 7
	yFlip        = 1 << 6
	xFlip        = 1 << 5
	dmgPalette   = 1 << 4
	bank         = 1 << 3
	cbgPalette   = 0b111
)

const (
	hBlankMode PPUMode = iota + 1
	vBlankMode
	oamMode
	PixelTransferModel
)

type Sprite struct {
	yPos       uint8
	xPos       uint8
	tileIdx    uint8
	attributes uint8
}

type PPU struct {
	lineDots uint16
	//registers

	oam  [40]Sprite
	vram [0x2000]uint8
	//lcd
	lcdControl, stat uint8
	scy, scx uint8
	ly, lyc uint8
	wy, wx uint8
	bgp, obp0, obp1 uint8
}

func LoadPpu() (*PPU, error) {
	//ppu initial values
	c := new(PPU)
	return c, nil
}

func (p *PPU) LcdRead(a uint16) uint8 {
	switch {
		case a == 0xFF40: return p.lcdControl 
		case a == 0xFF41: return p.stat
		case a == 0xFF42: return p.scy
		case a == 0xFF43: return p.scx
		case a == 0xFF44: return p.ly 
		case a == 0xFF45: return p.lyc
		case a == 0xFF46: 
			panic(0) //dma has no read
		case a == 0xFF47: return p.bgp
		case a == 0xFF48: return p.obp0
		case a == 0xFF49: return p.obp1
		case a == 0xFF4A: return p.wy
		case a == 0xFF4B: return p.wx
		default:
			panic(0)
	}
}  

func (p *PPU) LcdWrite(a uint16, v uint8) {
	switch {
		case a == 0xFF40: p.lcdControl = v
		case a == 0xFF41: p.stat = v
		case a == 0xFF42: p.scy = v
		case a == 0xFF43: p.scx = v
		case a == 0xFF44: p.ly = v
		case a == 0xFF45: p.lyc = v
		case a == 0xFF46: 
			panic(0) //dma write is processed elsewhere
		case a == 0xFF47: p.bgp = v 
		case a == 0xFF48: p.obp0 = v & 0xFC
		case a == 0xFF49: p.obp1 = v & 0xFC
		case a == 0xFF4A: p.wy = v
		case a == 0xFF4B: p.wx = v
		default:
			panic(0)
	}
}  

func (p *PPU) GetColor() {
	//TODO
}

func (p *PPU) SetMode(m PPUMode) {
	switch m {
		case hBlankMode:
		case vBlankMode:
		case oamMode:
		case PixelTransferModel: 
		default:
			panic(fmt.Sprintf("unexpected ppu mode %d", m))
	}
}

func (p *PPU) Update(cycles int) {

	for i := 0; i < cycles*4; i++ {
		//scanlines
		p.lineDots++
		if p.lineDots > 456 {
			p.lineDots = 0
			p.ly++

			if p.ly < 144 {

			} else if p.ly == 144 { // vertical blank
				
			} else if p.ly == 153 { // new frame
				p.ly = 0
			}
		}
	}
}


func (p *PPU) VramRead(a uint16) uint8     { return p.vram[a-0x8000] }
func (p *PPU) VramWrite(a uint16, v uint8) { p.vram[a-0x8000] = v }

func (p *PPU) oamRead(a uint16) uint8 {
	offsetByte := a - 0xFE00
	sprite := p.oam[offsetByte/4]
	switch offsetByte % 4 {
	case 0:
		return sprite.yPos
	case 1:
		return sprite.xPos
	case 2:
		return sprite.tileIdx
	case 3:
		return sprite.attributes
	default:
		panic(0)
	}
}

func (p *PPU) oamwrite(a uint16, v uint8) {
	offsetByte := a - 0xFE00
	sprite := p.oam[offsetByte/4]
	switch offsetByte % 4 {
	case 0:
		sprite.yPos = v
	case 1:
		sprite.xPos = v
	case 2:
		sprite.tileIdx = v
	case 3:
		sprite.attributes = v
	default:
		panic(0)
	}
}

func (p *PPU) DisplayTile(tile int, image *ebiten.Image, x int, y int) {
	colors := []color.Color{color.White, color.Gray{0xAA}, color.Gray{0x55}, color.Black}

	var tileY int
	for tileY = 0; tileY < 16; tileY += 2 {
		byte1 := p.vram[(tile * 16) + tileY]
		byte2 := p.vram[(tile * 16) + tileY + 1]

		var tileX int
		for tileX = 0; tileX < 8; tileX++ {
			bit1 := (byte1 & (1 << tileX)) >> (tileX)
			bit2 := (byte2 & (1 << tileX)) >> (tileX)

			color := colors[(bit2<<1)|bit1]
			image.Set(int(7-tileX+(x*8)), int((tileY/2))+(y*8), color)
		}
	}
}