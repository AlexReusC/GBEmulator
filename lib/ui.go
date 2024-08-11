package lib

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Screen struct {
	emulator *Emulator
	debugging bool
}


func (s *Screen) DisplayTile(image *ebiten.Image, add uint16, tileNum uint16, xOffset int, yOffset int) {
	colors := []color.Color{color.White, color.Gray16{0xFF55}, color.Gray16{0xFFAA}, color.Black}

	var tileY uint16
	for tileY = 0; tileY < 16; tileY += 2 {
		byte1 := s.emulator.bus.BusRead(add + (tileNum * 16) + tileY)
		byte2 := s.emulator.bus.BusRead(add + (tileNum * 16) + tileY + 1)

		var tileX int
		for tileX = 0; tileX < 8; tileX++{
			bit1 := (byte1 & (1 << tileX)) >> (tileX)
			bit2 := (byte2 & (1 << tileX)) >> (tileX)

			color := colors[(bit2 << 1) | bit1]
			image.Set(int(7-tileX+(xOffset*8)), int((tileY/2)) + (yOffset * 8), color)
		}
	}
}

func (s *Screen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 100, 50, 50})

	image := ebiten.NewImage(620, 620)
	var tileNum int = 0

	for y := 0; y < 24; y++ {
		for x := 0; x < 16; x++{
			s.emulator.ppu.DisplayTile(tileNum, image, x, y)
			tileNum++
		}
	}
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	screen.DrawImage(image, op)
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 384
}

func (s *Screen) Update() error {
	s.emulator.Run()
	return nil
}

func RunGame(e *Emulator) {
	screen := &Screen{emulator: e}
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("GBEmulator")
	ebiten.SetTPS(60*40)
	if err := ebiten.RunGame(screen); err != nil {
		log.Fatal(err)
	}
}