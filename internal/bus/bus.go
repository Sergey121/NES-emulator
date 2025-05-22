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

func New(ppu *ppu.PPU, cartridge *rom.Cartridge) *Bus {
	return &Bus{
		PPU:       ppu,
		Cartridge: cartridge,
	}
}

func (b *Bus) AttachCPU(cpu *cpu.CPU) {
	b.CPU = cpu
}

func (b *Bus) ShouldTriggerNMI() bool {
	return b.PPU.NMIOccurred()
}

func (b *Bus) AcknowledgeNMI() {
	b.PPU.ClearNMI()
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

func (b *Bus) CPUWrite(addr uint16, value byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF:
		// 2KB internal RAM, mirrored every 0x800
		b.RAM[addr%0x0800] = value

	case addr >= 0x2000 && addr <= 0x3FFF:
		// PPU registers (mirrored every 8 bytes)
		// fmt.Printf("PPU register write: %04X = %02X\n", addr, value)
		b.PPU.WriteRegister(0x2000+(addr%8), value)

	case addr >= 0x4000 && addr <= 0x4013, addr == 0x4015:
		// APU and I/O (если нужно)

	case addr == 0x4014:
		// OAM DMA
		b.PPU.OAMDMA = value

	case addr == 0x4016 || addr == 0x4017:
		// Controller

	case addr >= 0x8000:
		// PRG ROM — обычно нельзя писать, но можно обработать
		b.Cartridge.WritePRG(addr, value)
	}
}

func (b *Bus) ClockPPU() {
	for i := 0; i < 3; i++ {
		b.PPU.Step()
	}
}
