package lib

import (
	"fmt"
)

type PixelFetcherMode uint8

const (
	ReadTileId PixelFetcherMode = iota
	ReadTileData0
	ReadTileData1
	Idle
	PushPixelToQueue
)

type PixelFetcher struct {
	mode PixelFetcherMode
	ppu  *PPU
	tileId uint8
	PixelQueue []PixelData
	SpritesInTile []Sprite
	PalettesInTile []Palette
	lineX uint8
	fetchX uint8
	tileData [2]uint8
	spriteXPositions []uint8
	tileSpriteData0 []uint8
	tileSpriteData1 []uint8
}

type PixelData struct {
	color uint8
	palette Palette
}

func LoadPixelFetcher(ppu *PPU) *PixelFetcher {
	pf := new(PixelFetcher)
	pf.ppu = ppu
	return pf
}

func (pf *PixelFetcher) CleanPF() {
	pf.PixelQueue = []PixelData{}
}

func (pf *PixelFetcher) GetSpritePixelData(x, bit int) (uint8, Palette) {
	var pixelColor uint8
	var pixelPalette Palette

	for i := 0; i < len(pf.PalettesInTile); i++ {
		offset := int(pf.spriteXPositions[i]) - x
		if offset < 0 || offset > 7 {
			continue
		}
		bit += offset

		lo := (pf.tileSpriteData0[i] & (1 << bit)) >> (bit)
		hi := (pf.tileSpriteData1[i] & (1 << bit)) >> (bit)

		if pf.SpritesInTile[i].priority {
			break
		}
		pixelColor = (hi<<1)|lo
		pixelPalette = pf.PalettesInTile[i]

		if pixelColor != 0x00 {
			break
		}
	}
	return pixelColor, pixelPalette
}

func (pf *PixelFetcher) Update(d, index uint16) {
	p := pf.ppu
	mapX := ( p.pixelFetcher.fetchX + p.scx ) / 8
	mapY := ( p.ly + p.scy ) / 8


	if d&0x1 == 0 {
		switch pf.mode {
			case ReadTileId:
				//load background data
				pf.tileId = p.vram[p.GetBGTileArea() + uint16(mapX) + (uint16(mapY) * 32)]
				if p.GetBGWindowTileArea() == 0x0800 {
					pf.tileId += 128
				}
				//load sprites data
				pf.SpritesInTile = []Sprite{}
				pf.PalettesInTile =[]Palette{}
				pf.tileSpriteData0 = []uint8{}
				pf.tileSpriteData1 = []uint8{}
				pf.spriteXPositions = []uint8{}
				for i := 0; i < len(p.sprites) && len(pf.SpritesInTile) < 3; i++ {
					spriteXPos := (p.sprites[i].xPos - 8) + (p.scx % 8)
					spriteLo, spriteHi := spriteXPos, spriteXPos + 8 //beginning and end of sprite
					//check if sprite is inside this tile
					if (spriteLo >= pf.fetchX && spriteLo < pf.fetchX + 8) || (spriteHi >= pf.fetchX && spriteHi < pf.fetchX + 8){
						pf.SpritesInTile = append(pf.SpritesInTile, p.sprites[i])
						pf.spriteXPositions = append(pf.spriteXPositions, spriteXPos)
					} 
				}

				pf.fetchX += 8
				pf.mode = ReadTileData0
			case ReadTileData0: 
				//background data
				pf.tileData[0] = pf.ppu.vram[p.GetBGWindowTileArea() + (uint16(pf.tileId) * 16) + ((uint16(pf.ppu.ly+pf.ppu.scy) % 8) * 2)]
				//sprite data
				for i := 0; i < len(pf.SpritesInTile); i++ {
					spritePixels := p.vram[(uint16(pf.SpritesInTile[i].tileIndex) * 16) + ((uint16(p.ly) + 16) - uint16(pf.SpritesInTile[i].yPos)) * 2 ]
					pf.tileSpriteData0 = append(pf.tileSpriteData0, spritePixels)
				}

				pf.mode = ReadTileData1
			case ReadTileData1:
				//background data
				pf.tileData[1] = pf.ppu.vram[p.GetBGWindowTileArea() + (uint16(pf.tileId) * 16) + ((uint16(pf.ppu.ly+pf.ppu.scy) % 8) * 2) + 1]

				//sprite and palette data
				for i := 0; i < len(pf.SpritesInTile); i++ {
					spritePixels := p.vram[(uint16(pf.SpritesInTile[i].tileIndex) * 16) + ((uint16(p.ly) + 16) - uint16(pf.SpritesInTile[i].yPos)) * 2 + 1]
					pf.tileSpriteData1 = append(pf.tileSpriteData1, spritePixels)
					pf.PalettesInTile = append(pf.PalettesInTile, p.GetPaletteSprite( pf.SpritesInTile[i].palette ))

				}

				pf.mode = Idle
			case Idle:
				pf.mode = PushPixelToQueue
			case PushPixelToQueue:
				if len(pf.PixelQueue) <= 8 {
					x := int(pf.fetchX) - (8 - (int(p.scx) % 8))
					for bit := 7; bit >= 0; bit-- {
						pixelColor, pixelPalette := uint8(0x00), Bgp 

						if p.GetBGWindowEnable() {
							lo := (pf.tileData[0] & (1 << bit)) >> (bit)
							hi := (pf.tileData[1] & (1 << bit)) >> (bit)
							pixelColor = (hi<<1)|lo
						}

						if p.GetObjEnable() {
							if color, palette := pf.GetSpritePixelData(x, bit); color != 0x00 {
								pixelColor = color
								pixelPalette = palette
							}
						}

						if x >= 0 {
							pf.PixelQueue = append(pf.PixelQueue, PixelData{pixelColor, pixelPalette})
						}
					}
					pf.mode = ReadTileId
				}

			default:
				panic(fmt.Sprintf("unexpected pixel fetcher mode %d", pf.mode))
		}
	}

}