package lib

import (
	"fmt"
	"image"
	"image/color"
	"sort"
)

type PPUMode uint8
type Palette uint8

const (
	priorityMaskBit = 1 << 7
	yFlipBit        = 1 << 6
	xFlipBit        = 1 << 5
	dmgPaletteBit   = 1 << 4
	bankBit         = 1 << 3
	cbgPaletteBit   = 0b111
)

const (
	HBlank PPUMode = iota
	VBlank
	OamSearch
	PixelTransfer
)

type Sprite struct {
	yPos      uint8
	xPos      uint8
	tileIndex uint8
	palette   bool
	xFlipped  bool
	yFlipped  bool
	priority  bool
}

const (
	Bgp Palette = iota
	Obp0
	Obp1
)

type PPU struct {
	dots          uint16
	pixels        uint16 // x pos in screen
	Image         *image.RGBA
	MMU           *MMU
	sprites       []Sprite
	spritesInLine int

	//memory
	oam  [0xA0]uint8 // 160 = 40 * 4
	vram [0x2000]uint8
	//lcd
	lcdControl, stat uint8
	scy, scx         uint8
	ly               uint8 //scan line
	lyc              uint8
	wy, wx           uint8
	bgp, obp0, obp1  uint8
}

func LoadPpu() (*PPU, error) {
	//ppu initial values
	p := new(PPU)
	p.Image = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{160 + 128, 192}})

	p.lcdControl = 0x91
	p.bgp = 0xFC
	p.obp0 = 0xFF
	p.obp1 = 0xFF
	return p, nil
}

func (p *PPU) LcdRead(a uint16) uint8 {
	switch {
	case a == 0xFF40:
		return p.lcdControl
	case a == 0xFF41:
		return p.stat
	case a == 0xFF42:
		return p.scy
	case a == 0xFF43:
		return p.scx
	case a == 0xFF44:
		return p.ly
	case a == 0xFF45:
		return p.lyc
	case a == 0xFF46:
		panic(0) //dma has no read
	case a == 0xFF47:
		return p.bgp
	case a == 0xFF48:
		return p.obp0
	case a == 0xFF49:
		return p.obp1
	case a == 0xFF4A:
		return p.wy
	case a == 0xFF4B:
		return p.wx
	default:
		panic(0)
	}
}

func (p *PPU) LcdWrite(a uint16, v uint8) {
	switch {
	case a == 0xFF40:
		p.lcdControl = v
	case a == 0xFF41:
		p.stat = v
	case a == 0xFF42:
		p.scy = v
	case a == 0xFF43:
		p.scx = v
	case a == 0xFF44:
		p.ly = v
	case a == 0xFF45:
		p.lyc = v
	case a == 0xFF46:
		panic(0) //dma write is processed elsewhere
	case a == 0xFF47:
		p.bgp = v
	case a == 0xFF48:
		p.obp0 = v & 0b11111100
	case a == 0xFF49:
		p.obp1 = v & 0b11111100
	case a == 0xFF4A:
		p.wy = v
	case a == 0xFF4B:
		p.wx = v
	default:
		panic(0)
	}
}

func (p *PPU) GetLcdPpuEnable() bool { return BitIsSet(p.lcdControl, 7) }
func (p *PPU) GetWindowMapArea() uint16 {
	if BitIsSet(p.lcdControl, 6) {
		return 0x1C00
	} else {
		return 0x1800
	}
}
func (p *PPU) GetWindowEnable() bool { return BitIsSet(p.lcdControl, 5) }
func (p *PPU) GetBGWindowTileArea() uint16 {
	if BitIsSet(p.lcdControl, 4) {
		return 0
	} else {
		return 0x0800
	}
}
func (p *PPU) GetBGTileArea() uint16 {
	if BitIsSet(p.lcdControl, 3) {
		return 0x1C00
	} else {
		return 0x1800
	}
}
func (p *PPU) GetObjHeight() uint8 {
	if BitIsSet(p.lcdControl, 2) {
		return 8
	} else {
		return 16
	}
}
func (p *PPU) GetObjEnable() bool      { return BitIsSet(p.lcdControl, 1) }
func (p *PPU) GetBGWindowEnable() bool { return BitIsSet(p.lcdControl, 0) }

func (p *PPU) LycSourceSelected() bool       { return BitIsSet(p.stat, 6) }
func (p *PPU) OamSearchSourceSelected() bool { return BitIsSet(p.stat, 5) }
func (p *PPU) VBlankSourceSelected() bool    { return BitIsSet(p.stat, 4) }
func (p *PPU) HBlankSourceSelected() bool    { return BitIsSet(p.stat, 3) }
func (p *PPU) SetLycLcEqual(c bool)          { p.stat = SetBitWithCond(p.stat, 2, c) }
func (p *PPU) GetMode() PPUMode              { return PPUMode(p.stat & 0b11) }
func (p *PPU) SetMode(m PPUMode) {
	p.stat &= 0b11111100
	p.stat |= uint8(m)
}

func (p *PPU) GetColor(colorPixel uint8, palettePixel Palette) color.RGBA {
	colors := []color.RGBA{{0xFF, 0xFF, 0xFF, 1}, {0xC0, 0xC0, 0xC0, 1}, {40, 40, 40, 1}, {0, 0, 0, 1}}
	var paletteData uint8

	if palettePixel == Bgp {
		paletteData = p.bgp
	} else if palettePixel == Obp0 {
		paletteData = p.obp0
	} else if palettePixel == Obp1 {
		paletteData = p.obp1
	}

	id := (paletteData & (0b11 << (2 * colorPixel))) >> (2 * colorPixel)
	return colors[id]
}

func (p *PPU) GetPaletteSprite(dmgP bool) Palette {
	if dmgP {
		return Obp1
	}
	return Obp0
}

func (p *PPU) GetSpritesInLine() {
	p.sprites = []Sprite{}
	p.spritesInLine = 0
	spriteHeight := p.GetObjHeight()

	//get 10 sprites in this line
	for i := 0; i < 40 && p.spritesInLine <= 10; i++ {
		spriteY, spriteX, spriteIndex, spriteFlags := p.oam[i*4], p.oam[i*4+1], p.oam[i*4+2], p.oam[i*4+3]
		//sprite is touching y
		if (spriteY > p.ly+16) || (spriteY+spriteHeight <= p.ly+16) {
			continue
		}
		p.spritesInLine++ //sprite is counted even if next condition is true
		//sprite is invisible because of x
		if spriteX == 0 {
			continue
		}

		palette := (spriteFlags & dmgPaletteBit) != 0
		priority := (spriteFlags & priorityMaskBit) != 0
		yFlipped := (spriteFlags & yFlipBit) != 0
		xFlipped := (spriteFlags & xFlipBit) != 0

		p.sprites = append(p.sprites, Sprite{spriteY, spriteX, spriteIndex, palette, xFlipped, yFlipped, priority})
	}

	sort.SliceStable(p.sprites, func(i, j int) bool {
		return p.sprites[i].xPos < p.sprites[j].xPos
	})

}

func (p *PPU) UpdateLy() {
	p.ly++

	if p.ly == p.lyc {
		p.SetLycLcEqual(true)
		if p.LycSourceSelected() {
			p.MMU.RequestInterrupt(LCDSATUS)
		}
	} else {
		p.SetLycLcEqual(false)
	}
}

func (p *PPU) Update(cycles int) {
	switch p.GetMode() {
	case HBlank:
		p.dots++
		if p.dots == 456 {
			p.UpdateLy()
			p.dots = 0
			if p.ly < 144 { //rendered line
				p.SetMode(OamSearch)
				if p.OamSearchSourceSelected() {
					p.MMU.RequestInterrupt(LCDSATUS)
				}
			} else { //not rendered line
				p.SetMode(VBlank)
				p.MMU.RequestInterrupt(VBLANK)
				if p.VBlankSourceSelected() {
					p.MMU.RequestInterrupt(LCDSATUS)
				}
			}
		}
	case VBlank:
		p.dots++
		if p.dots == 456 {
			p.UpdateLy()
			if p.ly == 154 { //ppu has visited last line (153)
				p.ly = 0
				p.SetMode(OamSearch)
			}
			p.dots = 0
		}
	case OamSearch:
		p.dots++
		if p.dots == 80 {

			p.pixels = 0
			p.SetMode(PixelTransfer)
			p.GetSpritesInLine()
		}
	case PixelTransfer:
		p.SetMode(HBlank)
	default:
		panic(fmt.Sprintf("unexpected ppu mode %d", p.GetMode()))
	}
}

func (p *PPU) VramRead(a uint16) uint8     { return p.vram[a-0x8000] }
func (p *PPU) VramWrite(a uint16, v uint8) { p.vram[a-0x8000] = v }

func (p *PPU) oamRead(a uint16) uint8 {
	return p.oam[a-0xFE00]
}

func (p *PPU) oamwrite(a uint16, v uint8) {
	p.oam[a-0xFE00] = v
}

func (p *PPU) DebugDisplayTile(tile int, x int, y int) {
	colors := []color.RGBA{{0xFF, 0xFF, 0xFF, 1}, {0xC0, 0xC0, 0xC0, 1}, {40, 40, 40, 1}, {0, 0, 0, 1}, {255, 0, 0, 1}}

	var tileY int
	for tileY = 0; tileY < 16; tileY += 2 {
		byte1 := p.vram[(tile*16)+tileY]
		byte2 := p.vram[(tile*16)+tileY+1]

		var tileX int
		for tileX = 0; tileX < 8; tileX++ {
			bit1 := (byte1 & (1 << tileX)) >> (tileX)
			bit2 := (byte2 & (1 << tileX)) >> (tileX)

			color := colors[(bit2<<1)|bit1]
			p.Image.SetRGBA(int(7-tileX+(x*8)), int((tileY/2))+(y*8), color)
		}
	}
}
