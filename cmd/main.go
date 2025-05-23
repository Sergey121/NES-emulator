package main

import (
	"fmt"

	"github.com/sergey121/nes-emulator/internal/bus"
	"github.com/sergey121/nes-emulator/internal/cpu"
	"github.com/sergey121/nes-emulator/internal/ppu"
	"github.com/sergey121/nes-emulator/internal/rom"
)

func main() {
	// Load cartridge
	// how to load from a file from root folder?
	path := "./assets/roms/Super Mario Bros (E).nes"

	cartridge, err := rom.LoadRom(path)
	if err != nil {
		panic(err)
	}

	ppu := ppu.New()
	bus := bus.New(ppu, cartridge)
	cpuInstance := cpu.New()

	cpuInstance.AttachBus(bus)
	bus.AttachCPU(cpuInstance)

	err = cpuInstance.InsertCartridge(cartridge)
	if err != nil {
		panic(err)
	}

	cpuInstance.Reset()

	for {
		fmt.Println(cpuInstance.Trace())
		cpuInstance.Step()
	}
}
