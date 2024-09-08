package lib

import (
	"fmt"
	"image"
	"image/color"
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
	HBlank PPUMode = iota
	VBlank
	OamSearch
	PixelTransfer
)

type Sprite struct {
	yPos       uint8
	xPos       uint8
	tileIdx    uint8
	attributes uint8
}

type PPU struct {
	dots uint16
	pixels uint16
	mode PPUMode
	pixelFetcher *PixelFetcher
	Image *image.RGBA
	MMU *MMU
	//registers

	oam  [40]Sprite
	vram [0x2000]uint8
	//lcd
	lcdControl, stat uint8
	scy, scx uint8
	ly uint8 //scan line
	lyc uint8
	wy, wx uint8
	bgp, obp0, obp1 uint8
}

func LoadPpu() (*PPU, error) {
	//ppu initial values
	p := new(PPU)
	p.Image = image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{160+128, 192}})
	p.pixelFetcher = LoadPixelFetcher(p)

	p.lcdControl = 0x91
	return p, nil
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
		case a == 0xFF48: p.obp0 = v & 0b11111100
		case a == 0xFF49: p.obp1 = v & 0b11111100
		case a == 0xFF4A: p.wy = v
		case a == 0xFF4B: p.wx = v
		default:
			panic(0)
	}
}  

func (p *PPU) GetColor() {
	//TODO: palettes
}

func (p *PPU) Update(cycles int) {
	switch p.mode {
		case HBlank:
			p.dots++
			//HBlank work
			if p.dots == 456 {
				p.ly++
				if p.ly == p.lyc {
					p.MMU.RequestInterrupt(LCDSATUS)
				}
				p.dots = 0
				if p.ly < 144 {
					p.mode = OamSearch
				} else {
					p.mode = VBlank
				}
			}
		case VBlank:
			p.dots++
			//VBlank work here
			if p.dots == 456 {
				p.ly++
				if p.ly == p.lyc {
					p.MMU.RequestInterrupt(LCDSATUS)
				}
				if p.ly == 154 { //ppu has visited last line (153)
					p.ly = 0
					p.mode = OamSearch
				}
				p.dots = 0
			}
		case OamSearch:
			//OamSearch work here
			p.dots++
			if p.dots == 80 {
				p.mode = PixelTransfer
			}
		case PixelTransfer: 
			//PixelTransfer work here
			colors := []color.RGBA{{0xFF, 0xFF, 0xFF, 1}, {0xC0, 0xC0, 0xC0, 1}, {40, 40, 40, 1}, {0, 0, 0, 1}}
			p.dots++

			index := (((p.pixels+uint16(p.scx)) / 8) + ((uint16(p.ly+p.scy) / 8) * 32))

			p.pixelFetcher.Update(p.dots, index)

			if len(p.pixelFetcher.queue) <= 8 {
				return
			}

			pixelData := p.pixelFetcher.queue[0]
			p.pixelFetcher.queue = p.pixelFetcher.queue[1:]

			p.Image.SetRGBA(int(p.pixels), int(p.ly), colors[int(pixelData)])
			p.pixels++

			if p.pixels == 160 {
				p.pixels = 0
				p.mode = HBlank
				p.pixelFetcher.queue = []uint8{}
			}
		default:
			panic(fmt.Sprintf("unexpected ppu mode %d", p.mode))
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

func (p *PPU) DisplayTile(tile int, x int, y int) {
	colors := []color.RGBA{{0xFF, 0xFF, 0xFF, 1}, {0xC0, 0xC0, 0xC0, 1}, {40, 40, 40, 1}, {0, 0, 0, 1}}

	var tileY int
	for tileY = 0; tileY < 16; tileY += 2 {
		byte1 := p.vram[(tile * 16) + tileY]
		byte2 := p.vram[(tile * 16) + tileY + 1]

		var tileX int
		for tileX = 0; tileX < 8; tileX++ {
			bit1 := (byte1 & (1 << tileX)) >> (tileX)
			bit2 := (byte2 & (1 << tileX)) >> (tileX)

			color := colors[(bit2<<1)|bit1]
			//image.Set(int(7-tileX+(x*8)), int((tileY/2))+(y*8), color)		
			//p.Image.SetRGBA(int(p.pixels), int(p.ly), colors[int(pixelData)])
			//fmt.Println(int(7-tileX+(x*8)), int((tileY/2))+(y*8))
			p.Image.SetRGBA(int(7-tileX+(x*8)), int((tileY/2))+(y*8), color)
		}
	}
}