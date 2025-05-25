package main

import (
	"testing"

	"github.com/sergey121/nes-emulator/internal/bus"
	"github.com/sergey121/nes-emulator/internal/cpu"
	"github.com/sergey121/nes-emulator/internal/ppu"
	"github.com/sergey121/nes-emulator/internal/rom"
)

func TestNestestROM(t *testing.T) {
	cartridge, err := rom.LoadRom("../assets/roms/nestest (1).nes")
	if err != nil {
		t.Fatal(err)
	}

	ppuInstance := ppu.New(cartridge.CHR)
	busInstance := bus.New(ppuInstance, cartridge)
	cpuInstance := cpu.New()

	cpuInstance.AttachBus(busInstance)
	busInstance.AttachCPU(cpuInstance)

	cpuInstance.Reset()
	busInstance.PPU.Reset()
	cpuInstance.PC = 0xC000

	testCases := TestCases{
		7:     {A: 0x00, X: 0x00, Y: 0x00, SP: 0xFD, P: 0x24, ppuScanline: 0, ppuCycle: 21},
		101:   {A: 0xFF, X: 0x00, Y: 0x00, SP: 0xFB, P: 0xE4, ppuScanline: 0, ppuCycle: 303},
		501:   {A: 0xF5, X: 0x00, Y: 0x00, SP: 0xFB, P: 0x6F, ppuScanline: 4, ppuCycle: 139},
		1000:  {A: 0x80, X: 0x80, Y: 0x80, SP: 0xFB, P: 0x25, ppuScanline: 8, ppuCycle: 272},
		1501:  {A: 0x96, X: 0x69, Y: 0x69, SP: 0xFB, P: 0x25, ppuScanline: 13, ppuCycle: 70},
		2000:  {A: 0x33, X: 0x80, Y: 0x01, SP: 0x80, P: 0x27, ppuScanline: 17, ppuCycle: 203},
		2500:  {A: 0x00, X: 0x55, Y: 0x69, SP: 0xFB, P: 0x27, ppuScanline: 21, ppuCycle: 339},
		3000:  {A: 0x00, X: 0x00, Y: 0x5F, SP: 0xF9, P: 0x27, ppuScanline: 26, ppuCycle: 134},
		3502:  {A: 0x80, X: 0x00, Y: 0x69, SP: 0xFB, P: 0x27, ppuScanline: 30, ppuCycle: 276},
		4000:  {A: 0x23, X: 0x55, Y: 0x11, SP: 0xFB, P: 0x65, ppuScanline: 35, ppuCycle: 65},
		4505:  {A: 0x7F, X: 0x33, Y: 0x86, SP: 0xF9, P: 0x25, ppuScanline: 39, ppuCycle: 216},
		5001:  {A: 0x40, X: 0x33, Y: 0x91, SP: 0xFB, P: 0x25, ppuScanline: 43, ppuCycle: 340},
		5501:  {A: 0x3F, X: 0x9D, Y: 0x40, SP: 0xF9, P: 0x25, ppuScanline: 48, ppuCycle: 135},
		6001:  {A: 0x80, X: 0x55, Y: 0xA9, SP: 0xFB, P: 0xE5, ppuScanline: 52, ppuCycle: 271},
		6501:  {A: 0x55, X: 0x33, Y: 0xB8, SP: 0xF9, P: 0x64, ppuScanline: 57, ppuCycle: 66},
		7002:  {A: 0xFF, X: 0x33, Y: 0xC2, SP: 0xF9, P: 0xA5, ppuScanline: 61, ppuCycle: 205},
		7503:  {A: 0xFF, X: 0x33, Y: 0xCC, SP: 0xFB, P: 0x27, ppuScanline: 66, ppuCycle: 3},
		8003:  {A: 0x00, X: 0xD9, Y: 0x40, SP: 0xF9, P: 0x26, ppuScanline: 70, ppuCycle: 139},
		10000: {A: 0xFF, X: 0x42, Y: 0x00, SP: 0xFB, P: 0x27, ppuScanline: 87, ppuCycle: 333},
		12001: {A: 0x80, X: 0x55, Y: 0x29, SP: 0xF9, P: 0xE5, ppuScanline: 105, ppuCycle: 198},
		14500: {A: 0x97, X: 0x33, Y: 0x78, SP: 0xFB, P: 0xE5, ppuScanline: 127, ppuCycle: 193},
		17002: {A: 0x37, X: 0x02, Y: 0x9A, SP: 0xF9, P: 0x25, ppuScanline: 149, ppuCycle: 197},
		20002: {A: 0x4A, X: 0x02, Y: 0xBF, SP: 0xFB, P: 0x67, ppuScanline: 175, ppuCycle: 331},
		24002: {A: 0x29, X: 0x02, Y: 0xFF, SP: 0xF9, P: 0x65, ppuScanline: 211, ppuCycle: 55},
		26554: {A: 0x00, X: 0xFF, Y: 0x15, SP: 0xFD, P: 0x27, ppuScanline: 233, ppuCycle: 209},
	}

	for i := 0; i < 26550; i++ {
		cycle := cpuInstance.Cycles

		tc, ok := testCases[cycle]

		if ok {
			if cpuInstance.A != tc.A {
				t.Errorf("Cycle %d: Expected A = %02X, got %02X", cycle, tc.A, cpuInstance.A)
				return
			}
			if cpuInstance.X != tc.X {
				t.Errorf("Cycle %d: Expected X = %02X, got %02X", cycle, tc.X, cpuInstance.X)
				return
			}
			if cpuInstance.Y != tc.Y {
				t.Errorf("Cycle %d: Expected Y = %02X, got %02X", cycle, tc.Y, cpuInstance.Y)
				return
			}
			if cpuInstance.SP != tc.SP {
				t.Errorf("Cycle %d: Expected SP = %02X, got %02X", cycle, tc.SP, cpuInstance.SP)
				return
			}
			if cpuInstance.P != tc.P {
				t.Errorf("Cycle %d: Expected P = %02X, got %02X", cycle, tc.P, cpuInstance.P)
				return
			}
			if ppuInstance.Scanline() != tc.ppuScanline {
				t.Errorf("Cycle %d: Expected PPU Scanline = %d, got %d", cycle, tc.ppuScanline, ppuInstance.Scanline())
				return
			}
			if ppuInstance.Cycle() != tc.ppuCycle {
				t.Errorf("Cycle %d: Expected PPU Cycle = %d, got %d", cycle, tc.ppuCycle, ppuInstance.Cycle())
				return
			}
		}

		cpuInstance.Clock()
	}
}

type TestCase struct {
	// cycle       int
	A           byte
	X           byte
	Y           byte
	SP          byte
	P           byte
	ppuScanline int
	ppuCycle    int
}

type TestCases map[int]TestCase
