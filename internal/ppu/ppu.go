package ppu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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

	framebuffer [240][256]byte // 240x256 framebuffer

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

	// Background shift registers
	bgPatternLow    uint16 // битовая плоскость 0
	bgPatternHigh   uint16 // битовая плоскость 1
	bgAttributeLow  uint16 // палитра: младший бит
	bgAttributeHigh uint16 // палитра: старший бит

	// Latches / временные буферы
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
}

func (ppu *PPU) Reset() {
	// Сразу переход на пред-рендер-линию
	ppu.scanline = 0
	ppu.cycle = 21
	ppu.frame = 0
	// Обнулить адреса и лэтчи
	ppu.v = 0
	ppu.t = 0
	ppu.x = 0
	ppu.w = false
	ppu.PPUStatus = 0
	ppu.bufferedRead = 0
	ppu.nmiOccurred = false
	ppu.nmiOutput = false
	// Сбросить шифтеры
	ppu.bgPatternLow = 0
	ppu.bgPatternHigh = 0
	ppu.bgAttributeLow = 0
	ppu.bgAttributeHigh = 0
}

func New(chr []byte) *PPU {
	return &PPU{
		CHR: chr,
	}
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
	case 0x2002: // PPUSTATUS
		value := ppu.PPUStatus
		ppu.PPUStatus &^= 0x80 // сбрасываем флаг VBlank
		ppu.w = false          // сбрасываем toggle для SCROLL/ADDR
		return value

	case 0x2007: // PPUDATA
		value := ppu.read(ppu.v)

		var result byte
		if ppu.v >= 0x3F00 && ppu.v < 0x3FFF {
			result = value
		} else {
			result = ppu.bufferedRead
			ppu.bufferedRead = value
		}

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
	case 0x2000:
		ppu.PPUCTRL = value

	case 0x2001:
		ppu.PPUMASK = value

	case 0x2003:
		ppu.OAMADDR = value

	case 0x2004:
		ppu.OAM[ppu.OAMADDR] = value
		ppu.OAMADDR++

	case 0x2005:
		if !ppu.w {
			ppu.t = (ppu.t & 0xFFE0) | uint16(value>>3)
			ppu.x = value & 0x07
			ppu.w = true
		} else {
			ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
			ppu.t = (ppu.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
			ppu.w = false
		}

	case 0x2006:
		if !ppu.w {
			ppu.t = (ppu.t & 0x00FF) | ((uint16(value) & 0x3F) << 8)
			ppu.w = true
		} else {
			ppu.t = (ppu.t & 0xFF00) | uint16(value)
			ppu.v = ppu.t
			ppu.w = false
		}

	case 0x2007:
		ppu.write(ppu.v, value)
		if ppu.PPUCTRL&0x04 == 0 {
			ppu.v += 1
		} else {
			ppu.v += 32
		}
	}
}

func (ppu *PPU) Step() {
	if ppu.scanline == 241 && ppu.cycle == 1 {
		ppu.PPUStatus |= 0x80

		if ppu.PPUCTRL&0x80 != 0 {
			ppu.nmiOccurred = true
		}
	}

	// VBlank end
	if ppu.scanline == 261 && ppu.cycle == 1 {
		// Clear the VBlank flag
		ppu.PPUStatus &^= 0x80 // Clear the VBlank flag (bit 7)
		ppu.nmiOccurred = false
		ppu.nmiOutput = false
	}

	// Background rendering
	if ppu.scanline >= 0 && ppu.scanline < 240 && ppu.cycle >= 1 && ppu.cycle <= 256 {
		ppu.renderPixel()

		switch ppu.cycle % 8 {
		case 0:
			ppu.loadBackgroundShifters()
			ppu.fetchHighTileByte()
		case 1:
			// idle — ничего
		case 2:
			ppu.fetchNameTableByte()
		case 4:
			ppu.fetchAttributeTableByte()
		case 6:
			ppu.fetchLowTileByte()
		}

		ppu.bgPatternLow <<= 1
		ppu.bgPatternHigh <<= 1
		ppu.bgAttributeLow <<= 1
		ppu.bgAttributeHigh <<= 1

		if ppu.cycle%8 == 0 {
			ppu.incrementCoarseX()
		}
		if ppu.cycle == 256 {
			ppu.incrementY()
		}
	}

	if (ppu.scanline < 240 || ppu.scanline == 261) && ppu.cycle == 257 {
		ppu.transferAddressX()
	}

	if ppu.scanline == 261 && ppu.cycle >= 280 && ppu.cycle <= 304 {
		ppu.transferAddressY()
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
		return ppu.VRAM[(addr-0x2000)%0x800]
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
		ppu.VRAM[(addr-0x2000)%0x800] = value
	case addr < 0x4000:
		if addr >= 0x3F00 && addr < 0x3F20 {
			fmt.Printf("Palette write: %04X = %02X\n", addr, value)
		}
		ppu.PaletteTable[(addr-0x3F00)%32] = value
	default:
		// Unmapped memory
	}
}

func (ppu *PPU) fetchNameTableByte() {
	addr := 0x2000 | (ppu.v & 0x0FFF)
	tile := ppu.read(addr)
	// fmt.Printf(
	// 	"fetch NT: scanline=%3d cycle=%3d  v=%04X  tile=%02X\n",
	// 	ppu.scanline, ppu.cycle, ppu.v&0x0FFF, tile,
	// )
	ppu.nameTableByte = tile
}

func (ppu *PPU) fetchAttributeTableByte() {
	// Attribute for current scquqare 2x2 tile
	coarseX := (ppu.v >> 0) & 0x1F
	coarseY := (ppu.v >> 5) & 0x1F
	nameTable := (ppu.v >> 10) & 0x03

	attrAddr := 0x23C0 | (nameTable << 10) | ((coarseY >> 2) << 3) | (coarseX >> 2)
	ppu.attributeTableByte = ppu.read(attrAddr)
}

func (ppu *PPU) fetchLowTileByte() {
	fineY := (ppu.v >> 12) & 0x07
	tile := ppu.nameTableByte
	table := (ppu.PPUCTRL >> 4) & 0x01
	base := uint16(table) * 0x1000
	addr := base + uint16(tile)*16 + fineY
	ppu.lowTileByte = ppu.read(addr)
}

func (ppu *PPU) fetchHighTileByte() {
	fineY := (ppu.v >> 12) & 0x07
	tile := ppu.nameTableByte
	table := (ppu.PPUCTRL >> 4) & 1
	base := uint16(table) * 0x1000
	addr := base + uint16(tile)*16 + fineY + 8
	ppu.highTileByte = ppu.read(addr)
}

func (ppu *PPU) loadBackgroundShifters() {
	ppu.bgPatternLow = (ppu.bgPatternLow << 8) | uint16(ppu.lowTileByte)
	ppu.bgPatternHigh = (ppu.bgPatternHigh << 8) | uint16(ppu.highTileByte)

	// Палитра из attribute byte
	coarseX := (ppu.v >> 0) & 0x1F
	coarseY := (ppu.v >> 5) & 0x1F
	shift := ((coarseY & 0x02) << 1) | (coarseX & 0x02)

	attr := (ppu.attributeTableByte >> shift) & 0x03
	attrLow := byte(attr & 1)
	attrHigh := byte((attr >> 1) & 1)

	// Расширяем 2-битную палитру на 8 бит
	var repeatedLow, repeatedHigh byte
	if attrLow != 0 {
		repeatedLow = 0xFF
	}
	if attrHigh != 0 {
		repeatedHigh = 0xFF
	}

	ppu.bgAttributeLow = (ppu.bgAttributeLow << 8) | uint16(repeatedLow)
	ppu.bgAttributeHigh = (ppu.bgAttributeHigh << 8) | uint16(repeatedHigh)
}

func (ppu *PPU) renderPixel() {
	x := ppu.cycle - 1
	y := ppu.scanline

	bit0 := (ppu.bgPatternLow >> (15 - ppu.x)) & 1
	bit1 := (ppu.bgPatternHigh >> (15 - ppu.x)) & 1

	paletteIndex := (bit1 << 1) | bit0

	attr0 := (ppu.bgAttributeLow >> 15) & 1
	attr1 := (ppu.bgAttributeHigh >> 15) & 1
	paletteAttribute := (attr1 << 1) | attr0

	color := byte((paletteAttribute << 2) | paletteIndex)

	ppu.framebuffer[y][x] = color
}

var nesPalette = [64]color.RGBA{
	{84, 84, 84, 255},    // 0x00
	{0, 30, 116, 255},    // 0x01
	{8, 16, 144, 255},    // 0x02
	{48, 0, 136, 255},    // 0x03
	{68, 0, 100, 255},    // 0x04
	{92, 0, 48, 255},     // 0x05
	{84, 4, 0, 255},      // 0x06
	{60, 24, 0, 255},     // 0x07
	{32, 42, 0, 255},     // 0x08
	{8, 58, 0, 255},      // 0x09
	{0, 64, 0, 255},      // 0x0A
	{0, 60, 0, 255},      // 0x0B
	{0, 50, 60, 255},     // 0x0C
	{0, 0, 0, 255},       // 0x0D
	{0, 0, 0, 255},       // 0x0E
	{0, 0, 0, 255},       // 0x0F
	{152, 150, 152, 255}, // 0x10
	{8, 76, 196, 255},    // 0x11
	{48, 50, 236, 255},   // 0x12
	{92, 30, 228, 255},   // 0x13
	{136, 20, 176, 255},  // 0x14
	{160, 20, 100, 255},  // 0x15
	{152, 34, 32, 255},   // 0x16
	{120, 60, 0, 255},    // 0x17
	{84, 90, 0, 255},     // 0x18
	{40, 114, 0, 255},    // 0x19
	{8, 124, 0, 255},     // 0x1A
	{0, 118, 40, 255},    // 0x1B
	{0, 102, 120, 255},   // 0x1C
	{0, 0, 0, 255},       // 0x1D
	{0, 0, 0, 255},       // 0x1E
	{0, 0, 0, 255},       // 0x1F
	{236, 238, 236, 255}, // 0x20
	{76, 154, 236, 255},  // 0x21
	{120, 124, 236, 255}, // 0x22
	{176, 98, 236, 255},  // 0x23
	{228, 84, 236, 255},  // 0x24
	{236, 88, 180, 255},  // 0x25
	{236, 106, 100, 255}, // 0x26
	{212, 136, 32, 255},  // 0x27
	{160, 170, 0, 255},   // 0x28
	{116, 196, 0, 255},   // 0x29
	{76, 208, 32, 255},   // 0x2A
	{56, 204, 108, 255},  // 0x2B
	{56, 180, 204, 255},  // 0x2C
	{60, 60, 60, 255},    // 0x2D
	{0, 0, 0, 255},       // 0x2E
	{0, 0, 0, 255},       // 0x2F
	{236, 238, 236, 255}, // 0x30
	{168, 204, 236, 255}, // 0x31
	{188, 188, 236, 255}, // 0x32
	{212, 178, 236, 255}, // 0x33
	{236, 174, 236, 255}, // 0x34
	{236, 174, 212, 255}, // 0x35
	{236, 180, 176, 255}, // 0x36
	{228, 196, 144, 255}, // 0x37
	{204, 210, 120, 255}, // 0x38
	{180, 222, 120, 255}, // 0x39
	{168, 226, 144, 255}, // 0x3A
	{152, 226, 180, 255}, // 0x3B
	{160, 214, 228, 255}, // 0x3C
	{160, 162, 160, 255}, // 0x3D
	{0, 0, 0, 255},       // 0x3E
	{0, 0, 0, 255},       // 0x3F
}

func (p *PPU) DrawToImage(dst *ebiten.Image) {
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			colorIndex := p.framebuffer[y][x] & 0x3F // 6 бит
			c := nesPalette[colorIndex]
			dst.Set(x, y, c)
		}
	}
}

func (ppu *PPU) incrementCoarseX() {
	if (ppu.v & 0x001F) == 31 {
		ppu.v &= ^uint16(0x001F) // coarse X = 0
		ppu.v ^= 0x0400          // switch horizontal nametable
	} else {
		ppu.v += 1 // increment coarse X
	}
}

func (ppu *PPU) incrementY() {
	if (ppu.v & 0x7000) != 0x7000 {
		ppu.v += 0x1000 // increment fine Y
	} else {
		ppu.v &= ^uint16(0x7000) // fine Y = 0
		y := (ppu.v & 0x03E0) >> 5
		if y == 29 {
			y = 0
			ppu.v ^= 0x0800 // switch vertical nametable
		} else if y == 31 {
			y = 0 // overflow — stays in same nametable
		} else {
			y += 1
		}
		ppu.v = (ppu.v & ^uint16(0x03E0)) | (y << 5)
	}
}

func (ppu *PPU) transferAddressX() {
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}

func (ppu *PPU) transferAddressY() {
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}
