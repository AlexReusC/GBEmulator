package lib

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var DebugScreenOffset = 22

type Screen struct {
	emulator  *Emulator
	debugging bool
}

func (s *Screen) Draw(screen *ebiten.Image) {
	//debug
	var tileNum int = 0
	for y := 0; y < 24; y++ {
		for x := 0; x < 16; x++ {
			s.emulator.ppu.DebugDisplayTile(tileNum, x+DebugScreenOffset, y)
			tileNum++
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	image := ebiten.NewImageFromImage(s.emulator.ppu.Image)
	screen.DrawImage(image, op)
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 650, 400
}

func (s *Screen) Update() error {
	//TODO: probably this is unstable
	s.emulator.Run()
	return nil
}

func RunGame(e *Emulator) {
	screen := &Screen{emulator: e}
	ebiten.SetWindowSize(650, 400)
	ebiten.SetWindowTitle("GBEmulator")

	ebiten.SetTPS(240 * 80 * 4)
	if err := ebiten.RunGame(screen); err != nil {
		log.Fatal(err)
	}
}
