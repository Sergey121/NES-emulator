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
}

func (ppu *PPU) ReadRegister(addr uint16) byte {
	panic("ReadRegister not implemented")
}
