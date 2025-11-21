package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sergey121/nes-emulator/internal/bus"
	"github.com/sergey121/nes-emulator/internal/cpu"
	"github.com/sergey121/nes-emulator/internal/input"
	"github.com/sergey121/nes-emulator/internal/ppu"
	"github.com/sergey121/nes-emulator/internal/rom"
)

type Game struct {
	cpu     *cpu.CPU
	ppu     *ppu.PPU
	bus     *bus.Bus
	ebImage *ebiten.Image
}

func getTestPath(part string) string {
	return "./assets/tests/" + part + ".nes"
}

func NewGame() *Game {
	// path := "./assets/roms/Tetris.nes"
	path := "./assets/roms/Super Mario Bros (E).nes"
	// path := "./assets/roms/test_cpu_exec_space_apu.nes"

	// path := getTestPath("palette")

	cartridge, err := rom.LoadRom(path)
	if err != nil {
		panic(err)
	}

	ppu := ppu.New(cartridge.CHR)

	bus := bus.New(ppu, cartridge)
	cpuInstance := cpu.New()

	cpuInstance.AttachBus(bus)
	bus.AttachCPU(cpuInstance)

	cpuInstance.Reset()

	ebImage := ebiten.NewImage(256, 240)

	return &Game{
		cpu:     cpuInstance,
		ppu:     ppu,
		bus:     bus,
		ebImage: ebImage,
	}
}

func (g *Game) Update() error {
	// Update controller state
	var buttons byte
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		buttons |= input.ButtonA
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		buttons |= input.ButtonB
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		buttons |= input.ButtonSelect
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		buttons |= input.ButtonStart
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		buttons |= input.ButtonUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		buttons |= input.ButtonDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		buttons |= input.ButtonLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		buttons |= input.ButtonRight
	}
	g.bus.Controller1.SetButtons(buttons)

	// Один кадр ≈ 29780 PPU-тактов
	for i := 0; i < 29780; i++ {
		g.cpu.Clock()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ppu.DrawToImage(g.ebImage)     // твой метод отрисовки framebuffer в ebiten.Image
	screen.DrawImage(g.ebImage, nil) // вывод на экран
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return 256, 240
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(512, 480)
	ebiten.SetWindowTitle("NES Emulator")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
