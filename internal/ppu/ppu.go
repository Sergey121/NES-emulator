package ppu

import (
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

	// Sprite rendering
	secondaryOAM       [32]byte // 8 sprites * 4 bytes
	spriteCount        int      // Number of sprites found on next scanline
	spritePatternsLow  [8]byte  // Low byte of sprite pattern
	spritePatternsHigh [8]byte  // High byte of sprite pattern
	spritePositions    [8]byte  // X position of sprite
	spriteAttributes   [8]byte  // Attributes of sprite
	spriteIndexes      [8]byte  // Index in OAM (for Sprite 0 hit)
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
	// Reset sprite data
	ppu.spriteCount = 0
	for i := 0; i < 8; i++ {
		ppu.spritePatternsLow[i] = 0
		ppu.spritePatternsHigh[i] = 0
		ppu.spritePositions[i] = 0
		ppu.spriteAttributes[i] = 0
		ppu.spriteIndexes[i] = 0
	}
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

func (ppu *PPU) Step() {
	// 1. Инкремент циклов/сканлайнов/кадров
	ppu.cycle++
	if ppu.cycle >= 341 { // 0-340 циклов
		ppu.cycle = 0
		ppu.scanline++

		if ppu.scanline >= 262 { // 0-261 сканлайнов
			ppu.scanline = 0
			ppu.frame++
			// Здесь можно сигнализировать об окончании кадра для рендеринга на главном потоке
		}
	}

	renderingEnabled := (ppu.PPUMASK&0x08 != 0) || (ppu.PPUMASK&0x10 != 0) // Background or Sprite enable

	// 2. Обработка PPUSTATUS (флаги NMI, Sprite Zero Hit, Sprite Overflow)
	// Эти флаги сбрасываются в начале пред-рендеринг сканлайна (261)
	if ppu.scanline == 261 && ppu.cycle == 1 {
		ppu.PPUStatus &= (^(byte(1 << 7))) // Clear VBlank flag
		ppu.PPUStatus &= (^(byte(1 << 6))) // Clear Sprite 0 Hit flag
		ppu.PPUStatus &= (^(byte(1 << 5))) // Clear Sprite Overflow flag
		ppu.nmiOccurred = false
	}

	// 3. Логика NMI
	// NMI генерируется на scanline 241, cycle 1, если включен в PPUCTRL
	if ppu.scanline == 241 && ppu.cycle == 1 {
		ppu.PPUStatus |= (1 << 7)    // Set VBlank flag
		if ppu.PPUCTRL&(1<<7) != 0 { // Если NMI включен
			ppu.nmiOccurred = true
		}
	}

	// 4. Логика рендеринга (Visible Scanlines 0-239 and Pre-render Scanline 261)
	if (ppu.scanline >= 0 && ppu.scanline <= 239) || ppu.scanline == 261 {
		// --- Фаза предвыборки данных (Background Fetch) ---
		// PPU Fetch sequence: NT byte -> AT byte -> Low Tile byte -> High Tile byte (every 8 cycles)
		// These happen on cycles: 1, 9, 17, ... 257, 321, 329
		// And also on cycles 257, 337, 339 (for reloading shifters for next row/first two tiles)

		if (ppu.cycle >= 1 && ppu.cycle <= 256) || (ppu.cycle >= 321 && ppu.cycle <= 340) { // Cycles where background fetch/render occurs
			// Каждый 8-й цикл: Загрузка следующего тайла (NameTable, Attribute, Pattern Low, Pattern High)
			// и загрузка шифтеров
			switch ppu.cycle % 8 {
			case 1: // Fetch Nametable Byte
				ppu.nameTableByte = ppu.Read(0x2000 | (ppu.v & 0x0FFF))
				// Load shift registers with previous fetched data
				ppu.loadShiftRegisters()
			case 3: // Fetch Attribute Table Byte
				// Address: 0x23C0 | (NT_Y << 6) | (NT_X >> 2) (coarse X/Y for 4x4 group)
				attrAddr := 0x23C0 | (ppu.v & 0x0C00) | ((ppu.v >> 4) & 0x38) | ((ppu.v >> 2) & 0x07)
				attrByte := ppu.Read(attrAddr)
				// Determine 2-bit palette based on coarse X/Y within the 4x4 group
				if (ppu.v & 0x40) != 0 { // Coarse Y bit 1
					attrByte >>= 4
				}
				if (ppu.v & 0x02) != 0 { // Coarse X bit 1
					attrByte >>= 2
				}
				ppu.attributeTableByte = (attrByte & 0x03) // 2-bit palette index
			case 5: // Fetch Low Tile Byte
				// Base pattern table address: PPUCTRL bit 4 (0 for $0000, 1 for $1000)
				// Tile index: ppu.nameTableByte
				// Fine Y scroll: (ppu.v >> 12) & 0x07
				patternTableAddrLow := (uint16(ppu.PPUCTRL&0x10) << 8) | (uint16(ppu.nameTableByte) << 4) | ((ppu.v >> 12) & 0x07)
				ppu.lowTileByte = ppu.Read(patternTableAddrLow)
			case 7: // Fetch High Tile Byte
				patternTableAddrHigh := ((uint16(ppu.PPUCTRL&0x10) << 8) | (uint16(ppu.nameTableByte) << 4) | ((ppu.v >> 12) & 0x07)) + 8
				ppu.highTileByte = ppu.Read(patternTableAddrHigh)
			case 0: // (or 8, if using 1-indexed cycles) - Increments horizontal VRAM address
				// Increment coarse X
				ppu.incrementScrollX()
			}
		}

		// --- Фаза инкремента адреса VRAM ---
		// Горизонтальный инкремент coarse X происходит на каждом 8-м цикле
		// (уже обрабатывается в switch ppu.cycle % 8)

		// Вертикальный инкремент coarse Y (происходит на cycle 256)
		if ppu.cycle == 256 && renderingEnabled {
			ppu.incrementScrollY()
		}

		// Копирование горизонтальных битов VRAM из T в V (после 256 цикла)
		if ppu.cycle == 257 && renderingEnabled {
			ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F) // V_horz = T_horz
		}

		// Копирование вертикальных битов VRAM из T в V (на пред-рендеринг сканлайне 261)
		if ppu.scanline == 261 && ppu.cycle >= 280 && ppu.cycle <= 304 && renderingEnabled {
			ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0) // V_vert = T_vert
		}

		// --- Сдвиг шифтеров (на каждом пикселе) ---
		// И рендеринг пикселя
		// Fetch loop: 1-256 (visible) and 321-340 (pre-fetch)
		if (ppu.cycle >= 1 && ppu.cycle <= 256) || (ppu.cycle >= 321 && ppu.cycle <= 340) {
			// Render pixel (only visible)
			if ppu.cycle <= 256 {
				ppu.renderPixel()
			}

			// Background shift (visible 1-256 and pre-fetch 321-336)
			// STOP shifting at 336 to keep the loaded tiles in place for next line
			if renderingEnabled {
				if (ppu.cycle >= 1 && ppu.cycle <= 256) || (ppu.cycle >= 321 && ppu.cycle <= 336) {
					ppu.shiftBackgroundRegisters()
				}
			}
		}

		// Sprite evaluation (Cycles 65-256)
		if ppu.cycle == 257 && renderingEnabled {
			ppu.evaluateSprites()
		}

		// Sprite fetches (Cycles 257-320)
		// Simplified: do all fetches at 257 for now or spread them if needed.
		// For accurate timing we should spread them, but for functionality 257 is okay to prepare for next line.
		if ppu.cycle == 320 && renderingEnabled {
			ppu.fetchSpritePatterns()
		}
	}
}

func (ppu *PPU) loadShiftRegisters() {
	// Загрузка байтов паттерна в младшие 8 бит 16-битных шифтеров
	// (Старые данные сдвинутся в старшие 8 бит)
	ppu.bgPatternLow = (ppu.bgPatternLow & 0xFF00) | uint16(ppu.lowTileByte)
	ppu.bgPatternHigh = (ppu.bgPatternHigh & 0xFF00) | uint16(ppu.highTileByte)

	// Загрузка атрибутов: каждый бит палитры должен быть расширен на 8 бит.
	// 0x00 или 0xFF.
	// Если attributeTableByte & 0x01 = 1, то low bits = 0xFF.
	// Если attributeTableByte & 0x02 = 1, то high bits = 0xFF.
	// Мы сохраняем старшие 8 бит (предыдущий тайл) и загружаем новые 8 бит.
	attrValLow := byte(0)
	attrValHigh := byte(0)
	if (ppu.attributeTableByte & 0x01) != 0 {
		attrValLow = 0xFF // 0b11111111
	}
	if (ppu.attributeTableByte & 0x02) != 0 {
		attrValHigh = 0xFF // 0b11111111
	}

	ppu.bgAttributeLow = (ppu.bgAttributeLow & 0xFF00) | uint16(attrValLow)
	ppu.bgAttributeHigh = (ppu.bgAttributeHigh & 0xFF00) | uint16(attrValHigh)
}

func (ppu *PPU) shiftBackgroundRegisters() {
	// Сдвигаем все шифтеры на 1 бит влево.
	// Это перемещает биты к старшему концу, где мы будем их извлекать.
	ppu.bgPatternLow <<= 1
	ppu.bgPatternHigh <<= 1
	ppu.bgAttributeLow <<= 1
	ppu.bgAttributeHigh <<= 1
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

func (ppu *PPU) Read(addr uint16) byte {
	addr %= 0x4000 // PPU Memory Map is 0x0000 - 0x3FFF

	var val byte

	if addr >= 0x0000 && addr <= 0x1FFF {
		// CHR ROM/RAM
		val = ppu.CHR[addr] // Или от маппера, если CHR RAM
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		// Nametables (VRAM)
		// Обработка зеркалирования nametables (горизонтальное/вертикальное/1-screen/4-screen)
		// Это зависит от маппера. По умолчанию для iNES Mapper 0:
		// 0x2000-0x23FF -> Nametable 0
		// 0x2400-0x27FF -> Nametable 1
		// 0x2800-0x2BFF -> Nametable 0 (mirror of 0)
		// 0x2C00-0x2FFF -> Nametable 1 (mirror of 1)
		// 0x3000-0x3EFF -> Mirrors of 0x2000-0x2EFF
		// Пример для вертикального зеркалирования:
		// addr = 0x2000 + (addr-0x2000)%0x800
		// Пример для горизонтального зеркалирования:
		// addr = 0x2000 + (addr-0x2000)%0x400 (для первой пары) или 0x2400 + (addr-0x2400)%0x400 (для второй пары)

		// Предположим, что VRAM (nametables) находится в ppu.VRAM
		// Вам нужно будет передать информацию о зеркалировании из маппера
		// Для начала можно просто использовать:
		val = ppu.VRAM[addr%0x800] // Простейшее 2KB зеркалирование
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		// Palette RAM
		addr %= 0x20 // 32 байта палитры
		// Зеркалирование: 0x3F00, 0x3F04, 0x3F08, 0x3F0C зеркалят 0x3F10, 0x3F14, 0x3F18, 0x3F1C
		if addr == 0x10 || addr == 0x14 || addr == 0x18 || addr == 0x1C {
			addr -= 0x10
		}
		val = ppu.PaletteTable[addr]
	}
	return val // Should not happen
}

func (ppu *PPU) Write(addr uint16, data byte) {
	addr %= 0x4000 // PPU Memory Map is 0x0000 - 0x3FFF

	if addr >= 0x0000 && addr <= 0x1FFF {
		// CHR ROM/RAM (если CHR RAM, то сюда пишем)
		// if ppu.IsCHRRAM { ppu.CHR[addr] = data }
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		// Nametables (VRAM)
		// Смотри логику зеркалирования из Read
		ppu.VRAM[addr%0x800] = data
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		// Palette RAM
		addr %= 0x20
		if addr == 0x10 || addr == 0x14 || addr == 0x18 || addr == 0x1C {
			addr -= 0x10
		}
		ppu.PaletteTable[addr] = data
	}
}

// Read registers
func (ppu *PPU) ReadRegister(addr uint16) byte {
	switch addr {
	case 0x2002: // PPUSTATUS
		status := ppu.PPUStatus
		ppu.PPUStatus &= (^(byte(1 << 7))) // Clear VBlank flag
		ppu.w = false                      // Clear write toggle
		return status
	case 0x2004: // OAMDATA
		return ppu.OAM[ppu.OAMADDR]
	case 0x2007: // PPUDATA
		data := ppu.bufferedRead
		ppu.bufferedRead = ppu.Read(ppu.v) // Загружаем следующее значение
		// Для чтения из палитры - нет буфера, возвращаем сразу
		if ppu.v >= 0x3F00 {
			data = ppu.Read(ppu.v)
		}
		// Инкремент VRAM адреса
		if ppu.PPUCTRL&(1<<2) != 0 { // Bit 2 of PPUCTRL (VRAM address increment)
			ppu.v += 32 // Vertical increment
		} else {
			ppu.v += 1 // Horizontal increment
		}
		return data
	default:
		// Open bus behavior: last value written or last read
		// For simplicity, return 0 for now or whatever was last on the bus
		return ppu.bufferedRead // Or some other default
	}
}

// Write registers
func (ppu *PPU) WriteRegister(addr uint16, data byte) {
	switch addr {
	case 0x2000: // PPUCTRL
		ppu.PPUCTRL = data
		ppu.t = (ppu.t & 0xF3FF) | ((uint16(data) & 0x03) << 10) // Update nametable select in T
	case 0x2001: // PPUMASK
		ppu.PPUMASK = data
	case 0x2003: // OAMADDR
		ppu.OAMADDR = data
	case 0x2004: // OAMDATA
		ppu.OAM[ppu.OAMADDR] = data
		ppu.OAMADDR++ // OAMADDR auto-increments
	case 0x2005: // PPUSCROLL
		if !ppu.w { // First write: X scroll
			ppu.x = data & 0x07                            // Fine X scroll
			ppu.t = (ppu.t & 0xFFE0) | (uint16(data) >> 3) // Coarse X scroll
		} else { // Second write: Y scroll
			ppu.t = (ppu.t & 0x8FFF) | ((uint16(data) & 0x07) << 12) // Fine Y scroll
			ppu.t = (ppu.t & 0xFC1F) | ((uint16(data) & 0xF8) << 2)  // Coarse Y scroll
		}
		ppu.w = !ppu.w
	case 0x2006: // PPUADDR
		if !ppu.w { // First write: High byte
			ppu.t = (ppu.t & 0x80FF) | ((uint16(data) & 0x3F) << 8)
		} else { // Second write: Low byte
			ppu.t = (ppu.t & 0xFF00) | uint16(data)
			ppu.v = ppu.t // Transfer temporary address to current VRAM address
		}
		ppu.w = !ppu.w
	case 0x2007: // PPUDATA
		ppu.Write(ppu.v, data)
		// Инкремент VRAM адреса
		if ppu.PPUCTRL&(1<<2) != 0 { // Bit 2 of PPUCTRL (VRAM address increment)
			ppu.v += 32
		} else {
			ppu.v += 1
		}
	case 0x4014: // OAMDMA (handled by CPU, not PPU itself)
		ppu.OAMDMA = data // CPU будет читать это и выполнять DMA
	}
}

func (ppu *PPU) incrementScrollX() {
	// Если рендеринг выключен, не инкрементируем
	if ppu.PPUMASK&0x08 == 0 && ppu.PPUMASK&0x10 == 0 {
		return
	}

	// Инкремент coarse X
	if (ppu.v & 0x001F) == 31 { // Если coarse X == 31 (достигли конца nametable)
		ppu.v &= ^uint16(0x001F) // Сбросить coarse X
		ppu.v ^= uint16(0x0400)  // Переключить horizontal nametable (bit 10)
	} else {
		ppu.v++ // Инкремент coarse X
	}
}

func (ppu *PPU) incrementScrollY() {
	// Если рендеринг выключен, не инкрементируем
	if ppu.PPUMASK&0x08 == 0 && ppu.PPUMASK&0x10 == 0 {
		return
	}

	// Инкремент fine Y
	if (ppu.v & 0x7000) != 0x7000 { // Если fine Y < 7
		ppu.v += 0x1000 // Инкремент fine Y
	} else {
		ppu.v &= ^uint16(0x7000) // Сбросить fine Y
		// Инкремент coarse Y
		coarseY := (ppu.v & 0x03E0) >> 5
		if coarseY == 29 { // Если coarse Y == 29 (достигли конца видимой области nametable)
			coarseY = 0             // Сбросить coarse Y
			ppu.v ^= uint16(0x0800) // Переключить vertical nametable (bit 11)
		} else if coarseY == 31 { // Если coarse Y == 31 (обычно не используется)
			coarseY = 0 // Сбросить coarse Y, но не переключать nametable (для специальных эффектов)
		} else {
			coarseY++ // Инкремент coarse Y
		}
		ppu.v = (ppu.v & ^uint16(0x03E0)) | (coarseY << 5)
	}
}

func (ppu *PPU) evaluateSprites() {
	// Clear secondary OAM
	ppu.spriteCount = 0
	for i := 0; i < 32; i++ {
		ppu.secondaryOAM[i] = 0xFF
	}

	// Sprite height (8x8 or 8x16)
	spriteHeight := 8
	if ppu.PPUCTRL&0x20 != 0 {
		spriteHeight = 16
	}

	count := 0
	for i := 0; i < 64; i++ {
		y := int(ppu.OAM[i*4])
		// Check if sprite is on next scanline
		// Note: Sprite Y is delayed by one scanline, so we check against scanline (which is currently being rendered,
		// but we are preparing for scanline + 1? No, evaluation happens for the *next* scanline).
		// Actually, PPU renders current scanline, and evaluates sprites for the NEXT scanline.
		// So we check if sprite is in range [scanline, scanline + height).
		// Wait, standard behavior: Y byte in OAM is Y-1.
		// So sprite is visible on lines Y+1 to Y+height.
		// We are currently at `ppu.scanline`. We are preparing for `ppu.scanline`.
		// Wait, evaluation happens on scanline N for scanline N.
		// The sprites evaluated on scanline N are rendered on scanline N+1.

		// Let's assume we are preparing for the NEXT scanline (ppu.scanline).
		// The Y coordinate in OAM is the top of the sprite minus 1.
		// So if OAM_Y = 10, sprite starts at 11.
		// If we are on scanline 11, we want to render it.
		// Evaluation happens on scanline 10.

		// So, diff = scanline - y.
		// If 0 <= diff < height, then it is visible on the NEXT scanline?
		// Let's stick to: we are evaluating for the current `ppu.scanline` to be rendered?
		// No, evaluation runs on scanline N to prepare for N+1.
		// So we check if sprite will be on scanline N+1?
		// Or does the `scanline` variable track the one being rendered?
		// Usually `scanline` is the one being output.

		// Correct logic:
		// On scanline N:
		// We evaluate sprites that will appear on scanline N.
		// Wait, no. We evaluate on N for N.
		// Actually, let's look at the Y coordinate.
		// If OAM Y = 10. Sprite is on lines 11-18 (for 8x8).
		// We are currently rendering line 10. We evaluate for line 10? No.
		// We evaluate for line 11.

		// Let's use the standard logic:
		// Target line = ppu.scanline.
		// diff = targetLine - y
		// If 0 <= diff < height -> Visible.

		diff := ppu.scanline - y
		if diff >= 0 && diff < spriteHeight {
			if count < 8 {
				// Copy sprite data to secondary OAM
				ppu.secondaryOAM[count*4+0] = ppu.OAM[i*4+0]
				ppu.secondaryOAM[count*4+1] = ppu.OAM[i*4+1]
				ppu.secondaryOAM[count*4+2] = ppu.OAM[i*4+2]
				ppu.secondaryOAM[count*4+3] = ppu.OAM[i*4+3]

				// Store index for Sprite 0 detection
				ppu.spriteIndexes[count] = byte(i)
				count++
			} else {
				// Sprite overflow flag would be set here
				ppu.PPUStatus |= 0x20
				break
			}
		}
	}
	ppu.spriteCount = count
}

func (ppu *PPU) fetchSpritePatterns() {
	spriteHeight := 8
	if ppu.PPUCTRL&0x20 != 0 {
		spriteHeight = 16
	}

	for i := 0; i < ppu.spriteCount; i++ {
		y := int(ppu.secondaryOAM[i*4+0])
		tileIndex := ppu.secondaryOAM[i*4+1]
		attributes := ppu.secondaryOAM[i*4+2]
		x := ppu.secondaryOAM[i*4+3]

		ppu.spritePositions[i] = x
		ppu.spriteAttributes[i] = attributes

		// Calculate row within the sprite
		// We are preparing for the current scanline (which we just evaluated? No, we evaluated for THIS scanline).
		// Wait, if we evaluated for `ppu.scanline`, then we are fetching for `ppu.scanline`.
		row := ppu.scanline - y

		// Vertical flip
		if attributes&0x80 != 0 {
			row = spriteHeight - 1 - row
		}

		var addr uint16
		if spriteHeight == 8 {
			// 8x8 Sprites
			// Table from PPUCTRL bit 3
			table := uint16(0)
			if ppu.PPUCTRL&0x08 != 0 {
				table = 0x1000
			}
			addr = table | (uint16(tileIndex) << 4) | uint16(row)
		} else {
			// 8x16 Sprites
			// Table from bit 0 of tile index
			table := uint16(0)
			if tileIndex&1 != 0 {
				table = 0x1000
			}
			tileIndex &= 0xFE // Ignore last bit
			if row >= 8 {
				tileIndex++
				row -= 8
			}
			addr = table | (uint16(tileIndex) << 4) | uint16(row)
		}

		ppu.spritePatternsLow[i] = ppu.Read(addr)
		ppu.spritePatternsHigh[i] = ppu.Read(addr + 8)
	}
}

func (ppu *PPU) renderPixel() {
	if ppu.scanline >= 0 && ppu.scanline <= 239 && ppu.cycle >= 1 && ppu.cycle <= 256 {
		renderBackground := (ppu.PPUMASK & 0x08) != 0
		renderSprites := (ppu.PPUMASK & 0x10) != 0

		bgPixel := byte(0)
		bgPalette := byte(0)

		if renderBackground {
			if ppu.cycle <= 8 && (ppu.PPUMASK&0x02) == 0 {
				bgPixel = 0
			} else {
				shift := 15 - ppu.x
				bit0 := byte((ppu.bgPatternLow >> shift) & 0x01)
				bit1 := byte(((ppu.bgPatternHigh >> shift) & 0x01) << 1)
				bgPixel = bit0 | bit1

				attrBit0 := byte((ppu.bgAttributeLow >> shift) & 0x01)
				attrBit1 := byte(((ppu.bgAttributeHigh >> shift) & 0x01) << 1)
				bgPalette = attrBit0 | attrBit1
			}
		}

		spritePixel := byte(0)
		spritePalette := byte(0)
		spritePriority := false // false = front, true = back
		isSprite0 := false

		if renderSprites {
			if ppu.cycle <= 8 && (ppu.PPUMASK&0x04) == 0 {
				spritePixel = 0
			} else {
				for i := 0; i < ppu.spriteCount; i++ {
					x := int(ppu.spritePositions[i])
					// Check if pixel is within sprite X range
					// cycle-1 is the current X coordinate
					diff := (ppu.cycle - 1) - x
					if diff >= 0 && diff < 8 {
						attributes := ppu.spriteAttributes[i]

						// Horizontal flip
						col := diff
						if attributes&0x40 != 0 {
							col = 7 - col
						}

						bit0 := (ppu.spritePatternsLow[i] >> (7 - col)) & 1
						bit1 := (ppu.spritePatternsHigh[i] >> (7 - col)) & 1
						pixel := bit0 | (bit1 << 1)

						if pixel != 0 {
							spritePixel = pixel
							spritePalette = (attributes & 0x03) + 4 // +4 for sprite palettes
							spritePriority = (attributes & 0x20) != 0
							if ppu.spriteIndexes[i] == 0 {
								isSprite0 = true
							}
							break // Priority to first sprite found
						}
					}
				}
			}
		}

		// Sprite 0 Hit Detection
		if isSprite0 && bgPixel != 0 && renderBackground && renderSprites {
			// Check if we are at x=255 (right edge), usually ignored?
			// Check if cycle != 255?
			// Standard behavior: if both opaque, set flag.
			if ppu.cycle != 255 { // Some docs say 255 is excluded
				ppu.PPUStatus |= 0x40
			}
		}

		finalColorIndex := byte(0)
		if bgPixel == 0 && spritePixel == 0 {
			finalColorIndex = ppu.Read(0x3F00)
		} else if bgPixel != 0 && spritePixel == 0 {
			finalColorIndex = ppu.Read(0x3F00 + uint16((bgPalette<<2)+bgPixel))
		} else if bgPixel == 0 && spritePixel != 0 {
			finalColorIndex = ppu.Read(0x3F00 + uint16((spritePalette<<2)+spritePixel))
		} else {
			// Both opaque
			if spritePriority {
				// Sprite behind background
				finalColorIndex = ppu.Read(0x3F00 + uint16((bgPalette<<2)+bgPixel))
			} else {
				// Sprite in front
				finalColorIndex = ppu.Read(0x3F00 + uint16((spritePalette<<2)+spritePixel))
			}
		}

		ppu.framebuffer[ppu.scanline][ppu.cycle-1] = finalColorIndex
	}
}
