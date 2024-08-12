package lib

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	priorityMask = 1 << 7
	yFlip        = 1 << 6
	xFlip        = 1 << 5
	dmgPalette   = 1 << 4
	bank         = 1 << 3
	cbgPalette   = 0b111
)

type Sprite struct {
	yPos       uint8
	xPos       uint8
	tileIdx    uint8
	attributes uint8
}

type PPU struct {
	oam  [40]Sprite
	vram [0x2000]uint8
}

func LoadPpu() (*PPU, error) {
	c := &PPU{}
	return c, nil
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