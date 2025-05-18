package ppu

// https://www.nesdev.org/wiki/PPU_registers
type PPU struct {
	// Control registers
	PPUCTRL   byte // $2000
	PPUMASK   byte // $2001
	PPUStatus byte // $2002
	OAMADDR   byte // $2003
	OAMDATA   byte // $2004
	PPUSCROLL byte // $2005
	PPUADDR   byte // $2006
	PPUDATA   byte // $2007
	OAMDMA    byte // $4014

	VRAM [0x800]byte // 2kb internal RAM

	// Palette RAM
	PaletteTable [0x20]byte // 32 bytes of palette data

	// Object Attribute Memory
	OAM [0x100]byte // 256 bytes of OAM data

	// CHR ROM
	CHR []byte

	// Internal registers
	v uint16 // VRAM address
	t uint16 // Temporary VRAM address
	x byte   // Fine X scroll
	w bool   // Write toggle

	// Timers
	scanline int // Current scanline
	cycle    int // Current cycle
	frame    int // Current frame

	// Flags
	nmiOccurred bool // NMI occurred flag
	nmiOutput   bool // NMI output flag
	nmiPrevious bool // Previous NMI output flag

	bufferedRead byte // Buffered read value
}

func New() *PPU {
	return &PPU{}
}

func (ppu *PPU) Cycle() int {
	return ppu.cycle
}

func (ppu *PPU) Scanline() int {
	return ppu.scanline
}

func (ppu *PPU) NMIOccurred() bool {
	return ppu.nmiOccurred
}

func (ppu *PPU) ClearNMI() {
	ppu.nmiOccurred = false
}

func (ppu *PPU) ReadRegister(addr uint16) byte {
	switch addr {
	case 0x2007: // PPU Data
		value := ppu.read(ppu.v)

		var result byte

		if ppu.v >= 0x3F00 && ppu.v < 0x3FFF {
			// Read from palette directly
			result = value
		} else {
			// Deffered read: return buffered value and then update
			result = ppu.bufferedRead
			ppu.bufferedRead = value
		}

		// Increment the VRAM address
		if ppu.PPUCTRL&0x04 == 0 {
			ppu.v += 1
		} else {
			ppu.v += 32
		}

		return result
	}
	return 0 // Unmapped memory
}

func (ppu *PPU) WriteRegister(addr uint16, value byte) {
	switch addr {
	case 0x2005: // PPU Scroll
		if !ppu.w {
			ppu.t = (ppu.t & 0xFFE0) | uint16(value>>3)
			ppu.x = value & 0x07
			ppu.w = true
		} else {
			ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
			ppu.t = (ppu.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
			ppu.w = false
		}
	case 0x2006: // PPU Address
		if !ppu.w {
			ppu.t = (ppu.t & 0x00FF) | ((uint16(value) & 0x3F) << 8)
			ppu.w = true
		} else {
			ppu.t = (ppu.t & 0xFF00) | uint16(value)
			ppu.v = ppu.t
			ppu.w = false
		}
	case 0x2007: // PPU Data
		ppu.write(ppu.v, value)

		// Increment the VRAM address
		if ppu.PPUCTRL&0x04 == 0 {
			ppu.v += 1
		} else {
			ppu.v += 32
		}
	}
}

func (ppu *PPU) Step() {
	// VBlank start
	if ppu.scanline == 241 && ppu.cycle == 1 {
		// Set the VBlank flag
		ppu.PPUStatus |= 0x80 // Set the VBlank flag (bit 7)
		ppu.nmiOccurred = true

		// If NMI is enabled, trigger NMI
		if ppu.PPUCTRL&0x80 != 0 {
			// Trigger NMI
			ppu.nmiOutput = true
		}
	}

	// VBlank end
	if ppu.scanline == 261 && ppu.cycle == 1 {
		// Clear the VBlank flag
		ppu.PPUStatus &^= 0x80 // Clear the VBlank flag (bit 7)
		ppu.nmiOccurred = false
		ppu.nmiOutput = false
	}

	// Increment the cycle
	ppu.cycle++
	if ppu.cycle >= 341 {
		ppu.cycle = 0
		ppu.scanline++

		if ppu.scanline >= 262 {
			ppu.scanline = 0
			ppu.frame++
		}
	}
}

func (ppu *PPU) read(addr uint16) byte {
	addr &= 0x3FFF // Mask to 14 bits

	switch {
	case addr < 0x2000:
		return ppu.CHR[addr]
	case addr < 0x3F00:
		return ppu.VRAM[addr-0x2000]
	case addr < 0x4000:
		return ppu.PaletteTable[(addr-0x3F00)%32]
	default:
		return 0 // Unmapped memory
	}
}

func (ppu *PPU) write(addr uint16, value byte) {
	addr &= 0x3FFF // Mask to 14 bits

	switch {
	case addr < 0x2000:
		// CHR RAM (not used, not ROM)
	case addr < 0x3F00:
		ppu.VRAM[addr-0x2000] = value
	case addr < 0x4000:
		ppu.PaletteTable[(addr-0x3F00)%32] = value
	default:
		// Unmapped memory
	}
}
