package bus

import (
	"github.com/sergey121/nes-emulator/internal/cpu"
	"github.com/sergey121/nes-emulator/internal/input"
	"github.com/sergey121/nes-emulator/internal/ppu"
	"github.com/sergey121/nes-emulator/internal/rom"
)

type Bus struct {
	CPU         *cpu.CPU
	PPU         *ppu.PPU
	Cartridge   *rom.Cartridge
	Controller1 *input.Controller
	// RAM is the 2KB of RAM in the NES
	RAM [0x800]byte // 2KB of RAM
}

func New(ppu *ppu.PPU, cartridge *rom.Cartridge) *Bus {
	return &Bus{
		PPU:         ppu,
		Cartridge:   cartridge,
		Controller1: input.NewController(),
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
	case addr == 0x4016:
		return b.Controller1.Read()
	case addr >= 0x8000:
		// Cartridge ROM ($8000-$FFFF)
		return b.Cartridge.ReadPRG(addr)
	default:
		// Unmapped memory returns open bus (usually high byte of address)
		return byte(addr >> 8)
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
		// Writing $XX will upload 256 bytes of data from CPU page $XX00-$XXFF to the internal PPU OAM.
		startAddr := uint16(value) << 8
		for i := uint16(0); i < 256; i++ {
			data := b.CPURead(startAddr + i)
			b.PPU.OAMDATA = data // This updates OAMDATA register but usually DMA writes directly to OAM
			// Direct write to OAM is better/simpler here as OAMADDR might auto-increment on OAMDATA write
			// but DMA usually writes to OAMADDR, OAMADDR+1, etc.
			// Let's write directly to OAM array to be safe and simple, respecting OAMADDR wrapping if needed.
			// Standard DMA copies to OAM[OAMADDR] and increments OAMADDR?
			// "The CPU is suspended... 256 bytes are read... and written to OAMDATA"
			// So it effectively writes to OAMDATA 256 times.
			// Let's use WriteRegister(0x2004, data) to reuse logic?
			// Or just direct copy if we assume OAMADDR is 0 (usually is).
			// Let's use WriteRegister to be safe with auto-increment logic.
			b.PPU.WriteRegister(0x2004, data)
		}
		// CPU suspension is not implemented yet (requires cycle counting accuracy),
		// but immediate copy works for functionality.

	case addr == 0x4016:
		b.Controller1.Write(value)

	case addr == 0x4017:
		// Controller 2 (not implemented yet)

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

func (b *Bus) StepPPU() {
	b.PPU.Step()
}
