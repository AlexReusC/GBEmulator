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
	queue []uint8 //TODO: create custom queue
	lineX uint8
	fetchX uint8
	tileData [2]uint8
}

func LoadPixelFetcher(ppu *PPU) *PixelFetcher {
	pf := new(PixelFetcher)
	pf.ppu = ppu
	return pf
}

func (pf *PixelFetcher) CleanPF() {
	pf.queue = []uint8{}
}

func (pf *PixelFetcher) Update(d, index uint16) {
	p := pf.ppu
	mapX := ( p.pixelFetcher.fetchX + p.scx ) / 8
	mapY := ( p.ly + p.scy ) / 8


	if d&0x1 == 0 {
		switch pf.mode {
		case ReadTileId:
			if p.GetBwObjEnable(){
				pf.tileId = p.vram[p.GetBGTileArea() + uint16(mapX) + (uint16(mapY) * 32)]
				if p.GetBGWindowTileArea() == 0x0800 {
					pf.tileId += 128
				}
			}
			pf.fetchX += 8
			pf.mode = ReadTileData0
		case ReadTileData0: 
			pf.tileData[0] = pf.ppu.vram[p.GetBGWindowTileArea() + (uint16(pf.tileId) * 16) + ((uint16(pf.ppu.ly+pf.ppu.scy) % 8) * 2)]
			pf.mode = ReadTileData1
		case ReadTileData1:
			pf.tileData[1] = pf.ppu.vram[p.GetBGWindowTileArea() + (uint16(pf.tileId) * 16) + ((uint16(pf.ppu.ly+pf.ppu.scy) % 8) * 2) + 1]
			pf.mode = Idle
		case Idle:
			pf.mode = PushPixelToQueue
		case PushPixelToQueue:
			if len(pf.queue) <= 8 {
				x := int(pf.fetchX) - (8 - (int(p.scx) % 8))
				for bit := 7; bit >= 0; bit-- {
					lo := (pf.tileData[0] & (1 << bit)) >> (bit)
					hi := (pf.tileData[1] & (1 << bit)) >> (bit)

					pixelData := (hi<<1)|lo

					if x >= 0 {
						pf.queue = append(pf.queue, pixelData)
					}
				}
				pf.mode = ReadTileId
			}

		default:
			panic(fmt.Sprintf("unexpected pixel fetcher mode %d", pf.mode))
		}
	}

}