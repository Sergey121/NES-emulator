package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sergey121/nes-emulator/internal/bus"
	"github.com/sergey121/nes-emulator/internal/cpu"
	"github.com/sergey121/nes-emulator/internal/ppu"
	"github.com/sergey121/nes-emulator/internal/rom"
)

type Game struct {
	cpu     *cpu.CPU
	ppu     *ppu.PPU
	ebImage *ebiten.Image
}

func NewGame() *Game {
	// path := "./assets/roms/Tetris.nes"
	path := "./assets/roms/Super Mario Bros (E).nes"
	// path := "./assets/roms/test_cpu_exec_space_apu.nes"

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
	bus.PPU.Reset()

	bus.CPURead(0x2002)

	for i := 0; i < 341; i++ {
		bus.ClockPPU()
	}

	ebImage := ebiten.NewImage(256, 240)

	return &Game{
		cpu:     cpuInstance,
		ppu:     ppu,
		ebImage: ebImage,
	}
}

func (g *Game) Update() error {
	// Один кадр ≈ 29780 PPU-тактов
	for i := 0; i < 29780; i++ {
		// fmt.Println(g.cpu.Trace())
		// g.cpu.Step()
		// g.cpu.Clock()
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
	// game := NewGame()

	// ebiten.SetWindowSize(512, 480)
	// ebiten.SetWindowTitle("NES Emulator")

	// if err := ebiten.RunGame(game); err != nil {
	// 	panic(err)
	// }

	path := "./assets/roms/nestest (1).nes"

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
	bus.PPU.Reset()

	cpuInstance.PC = 0xC000

	// 12000 - is ok
	for i := 0; i < 13000; i++ {
		if cpuInstance.CyclesLeft == 0 {
			fmt.Println(cpuInstance.Trace(ppu.Scanline(), ppu.Cycle()))
		}
		cpuInstance.Clock()
	}
}
