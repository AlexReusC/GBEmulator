package lib

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Screen struct {
	emulator *Emulator
}

func (s *Screen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 100, 50, 50})
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

func (s *Screen) Update() error {
	// todo
	s.emulator.Run()
	return nil
}

func RunGame(e *Emulator) {
	screen := &Screen{emulator: e}
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("GBEmulator")
	ebiten.SetTPS(60*154 * 114)
	if err := ebiten.RunGame(screen); err != nil {
		log.Fatal(err)
	}
}