# GBEmulator
Gameboy Emulator written in Golang

## About
GBEmulator was created as a project that could teach more about how the internals of a computer works. It's not intended to be the most precise emulation but a way to create a system that
allows me to understand the core systems and be able to play some games from the base I created.

## How to run
GBEmulator uses the [ebiten library](https://github.com/hajimehoshi/ebiten) for graphical rendering. For running use: 
```
go run main [location of ROM]
```
## Features
- [x] CPU
  - [x] All instructions
  - [x] Interrupts
  - [x] Clock
- [x] Memory
- [ ] Graphics
  - [x] Pixel Pipeline 
  - [x] Tiles
  - [x] OAM
  - [ ] Window
- [ ] Sound

## Resources
- [The Ultimate Game Boy Talk (33c3)](https://www.youtube.com/watch?v=HyzD8pNlpwI&ab_channel=media.ccc.de)
- https://gbdev.io/pandocs/
- https://rgbds.gbdev.io/docs/v0.7.0/gbz80.7
- https://emudev.de/gameboy-emulator/overview/
