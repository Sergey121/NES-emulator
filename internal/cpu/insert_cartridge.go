package cpu

import (
	"fmt"

	"github.com/sergey121/nes-emulator/internal/rom"
)

func (cpu *CPU) InsertCartridge(cartridge *rom.Cartridge) error {
	prgLen := len(cartridge.PRG)

	switch prgLen {
	case 16 * 1024:
		copy(cpu.Memory[0x8000:0xC000], cartridge.PRG)
		copy(cpu.Memory[0xC000:0x10000], cartridge.PRG)

	case 32 * 1024:
		copy(cpu.Memory[0x8000:0x10000], cartridge.PRG)

	default:
		return fmt.Errorf("unsupported PRG size: %d", prgLen)
	}

	// Reset Vector should be set to the start of the PRG ROM
	// If not (for testing purposes), set it to 0x8000
	if cpu.Memory[ResetVector] == 0x00 && cpu.Memory[ResetVector+1] == 0x00 {
		cpu.Memory[ResetVector] = 0x00
		cpu.Memory[ResetVector+1] = 0x80
	}

	return nil
}
