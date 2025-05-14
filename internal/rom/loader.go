package rom

import (
	"fmt"
	"os"
)

func LoadRom(path string) (*Cartridge, error) {
	// Open the ROM file
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return createCartridge(data)
}

func createCartridge(data []byte) (*Cartridge, error) {
	if len(data) < 16 || string(data[0:4]) != "NES\x1A" {
		return nil, fmt.Errorf("invalid NES ROM file")
	}

	prgBanks := int(data[4])
	chrBanks := int(data[5])

	flag6 := data[6] // Mapper and mirroring flags
	flag7 := data[7] // Mapper and mirroring flags

	mapper := (flag7 & 0xF0) | (flag6 >> 4)
	// Only mapper 0 is supported for now
	// TODO: Support other mappers
	if mapper != 0 {
		return nil, fmt.Errorf("unsupported mapper: %d", mapper)
	}

	var mirroring MirroringType
	if flag6&0x01 != 0 {
		mirroring = Vertical
	} else {
		mirroring = Horizontal
	}

	hasTraining := flag6&0x04 != 0 // Trainer present
	offset := 16

	if hasTraining { // Trainer present, skip it
		offset += 512
	}

	prgSize := prgBanks * 16 * 1024
	chrSize := chrBanks * 8 * 1024

	if len(data) < offset+int(prgSize)+int(chrSize) { // Check if the ROM file is complete
		return nil, fmt.Errorf("ROM file is incomplete")
	}

	prg := data[offset : offset+prgSize]
	offset += prgSize

	var chr []byte
	hasCHRRom := chrSize > 0
	if hasCHRRom {
		chr = data[offset : offset+chrSize]
	} else {
		chr = make([]byte, 8*1024)
	}

	cartridge := &Cartridge{
		PRG:       prg,
		CHR:       chr,
		Mapper:    mapper,
		Mirroring: mirroring,
		HasCHRROM: hasCHRRom,
	}

	return cartridge, nil
}
