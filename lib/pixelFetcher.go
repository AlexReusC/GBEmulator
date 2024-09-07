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
	tileData [2]uint8
}

func LoadPixelFetcher(ppu *PPU) *PixelFetcher {
	pf := new(PixelFetcher)
	pf.ppu = ppu
	return pf
}

func (pf *PixelFetcher) Update(d, index uint16) {
	if d&0x1 == 0 {
		switch pf.mode {
		case ReadTileId:
			var offsetIndex uint16
			if BitIsSet(pf.ppu.lcdControl, 3){
				offsetIndex = 0x1C00
			} else {
				offsetIndex = 0x1800
			}
			pf.tileId = pf.ppu.vram[offsetIndex + index]
			if !BitIsSet(pf.ppu.lcdControl, 4){
				pf.tileId += 128
			}
			pf.mode = ReadTileData0
		case ReadTileData0: 
			var offsetIndex uint16
			if !BitIsSet(pf.ppu.lcdControl, 4) {
				offsetIndex = 0x800
			}
			pf.tileData[0] = pf.ppu.vram[offsetIndex + (uint16(pf.tileId) * 16) + ((uint16(pf.ppu.ly+pf.ppu.scy) % 8) * 2)]
			pf.mode = ReadTileData1
		case ReadTileData1:
			var offsetIndex uint16
			if !BitIsSet(pf.ppu.lcdControl, 4){
				offsetIndex = 0x800
			} 
			pf.tileData[1] = pf.ppu.vram[offsetIndex + (uint16(pf.tileId) * 16) + ((uint16(pf.ppu.ly+pf.ppu.scy) % 8) * 2) + 1]
			pf.mode = Idle
		case Idle:
			pf.mode = PushPixelToQueue
		case PushPixelToQueue:
			if len(pf.queue) <= 8 {
			for bit := 7; bit >= 0; bit-- {
				lo := (pf.tileData[0] & (1 << bit)) >> (bit)
				hi := (pf.tileData[1] & (1 << bit)) >> (bit)

				pixelData := (hi<<1)|lo

				pf.queue = append(pf.queue, pixelData)
			}
			pf.mode = ReadTileId
			}

		default:
		panic(fmt.Sprintf("unexpected pixel fetcher mode %d", pf.mode))
		}
	}

}