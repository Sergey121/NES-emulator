package bus

import (
	"github.com/sergey121/nes-emulator/internal/cpu"
	"github.com/sergey121/nes-emulator/internal/ppu"
	"github.com/sergey121/nes-emulator/internal/rom"
)

type Bus struct {
	CPU       *cpu.CPU
	PPU       *ppu.PPU
	Cartridge *rom.Cartridge
	// RAM is the 2KB of RAM in the NES
	RAM [0x800]byte // 2KB of RAM
}

func New(cpu *cpu.CPU, ppu *ppu.PPU, cartridge *rom.Cartridge) *Bus {
	return &Bus{
		CPU:       cpu,
		PPU:       ppu,
		Cartridge: cartridge,
	}
}

func (b *Bus) CPURead(addr uint16) byte {
	// TODO: Add APU/controller read
	switch {
	case addr < 0x2000:
		// Read from RAM and mirrors
		return b.RAM[addr%0x800]
	case addr >= 0x2000 && addr < 0x4000:
		// PPU registers ($2000-$3FFF), mirrors every 8 bytes
		return b.PPU.ReadRegister(0x2000 + (addr % 8))
	case addr >= 0x8000:
		// Cartridge ROM ($8000-$FFFF)
		return b.Cartridge.ReadPRG(addr)
	default:
		// Unmapped memory
		return 0
	}
}
