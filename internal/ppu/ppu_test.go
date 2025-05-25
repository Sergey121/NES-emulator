package ppu

// import "testing"

// func TestPPUVBlankFlagAndNMI(t *testing.T) {
// 	ppu := &PPU{
// 		PPUCTRL: 0x80, // Enable NMI
// 	}

// 	// Simulate the scanline and cycle for VBlank
// 	ppu.scanline = 241
// 	ppu.cycle = 1

// 	// Execute the PPU
// 	ppu.Step()

// 	// Check if the VBlank flag is set
// 	if ppu.PPUStatus&0x80 == 0 {
// 		t.Errorf("Expected VBlank flag to be set, but it was not")
// 	}

// 	// Check if NMI occurred
// 	if !ppu.nmiOccurred {
// 		t.Errorf("Expected NMI to occur, but it did not")
// 	}

// 	// Check if NMI output is enabled
// 	if !ppu.nmiOutput {
// 		t.Errorf("Expected NMI output to be enabled, but it was not")
// 	}
// }

// func TestPPUVBlankFlagAndNMI2(t *testing.T) {
// 	ppu := &PPU{}

// 	// Simulate the scanline and cycle for VBlank
// 	for ppu.Scanline() != 241 || ppu.Cycle() != 1 {
// 		ppu.Step()
// 	}

// 	ppu.Step()

// 	// Check if the VBlank flag is set
// 	if ppu.PPUStatus&0x80 == 0 {
// 		t.Errorf("Expected VBlank flag to be set, but it was not")
// 	}

// 	// Check if NMI occurred
// 	if !ppu.NMIOccurred() {
// 		t.Errorf("Expected NMI to occur, but it did not")
// 	}

// 	// Move scanline for a new frame
// 	for ppu.Scanline() != 261 || ppu.Cycle() != 1 {
// 		ppu.Step()
// 	}

// 	ppu.Step()

// 	if ppu.PPUStatus&0x80 != 0 {
// 		t.Errorf("Expected VBlank flag to be cleared, but it was not")
// 	}

// 	// Check if NMI occurred
// 	if ppu.NMIOccurred() {
// 		t.Errorf("Expected NMI to not occur, but it did")
// 	}
// }

// func TestPPU_PPUDATA_BufferedRead(t *testing.T) {
// 	ppu := &PPU{
// 		CHR:          make([]byte, 0x2000), // 8 KB CHR-ROM
// 		VRAM:         [0x800]byte{},
// 		PaletteTable: [32]byte{},
// 	}

// 	// Write test values by address 0x2000 (nametable)
// 	ppu.write(0x2000, 0xAB)
// 	ppu.write(0x2001, 0xCD)

// 	ppu.v = 0x2000 // Start from 0x2000

// 	// First read should return garbage (0x00), and buffer should contain 0xAB
// 	result1 := ppu.ReadRegister(0x2007)
// 	if result1 != 0x00 {
// 		t.Errorf("expected first read to return 0x00 (initial buffer), got %02X", result1)
// 	}

// 	// Second read should return 0xAB, and buffer should contain 0xCD
// 	result2 := ppu.ReadRegister(0x2007)
// 	if result2 != 0xAB {
// 		t.Errorf("expected second read to return 0xAB, got %02X", result2)
// 	}

// 	// Third read should return 0xCD, and buffer should be empty
// 	result3 := ppu.ReadRegister(0x2007)
// 	if result3 != 0xCD {
// 		t.Errorf("expected third read to return 0xCD, got %02X", result3)
// 	}
// }

// func TestPPU_PPUDATA_PaletteRead(t *testing.T) {
// 	ppu := &PPU{
// 		PaletteTable: [32]byte{},
// 	}

// 	ppu.write(0x3F00, 0x42)
// 	ppu.v = 0x3F00

// 	// Read should return 0x42 directly, without buffering
// 	result := ppu.ReadRegister(0x2007)
// 	if result != 0x42 {
// 		t.Errorf("expected palette read to return 0x42, got %02X", result)
// 	}
// }

// func TestPPU_PPUDATA_Increment(t *testing.T) {
// 	ppu := &PPU{
// 		CHR: make([]byte, 0x2000),
// 	}

// 	// Bit 2 = 0 → increment by 1
// 	ppu.PPUCTRL = 0x00
// 	ppu.v = 0x1000
// 	_ = ppu.ReadRegister(0x2007)
// 	if ppu.v != 0x1001 {
// 		t.Errorf("expected v to increment by 1, got %04X", ppu.v)
// 	}

// 	// Bit 2 = 1 → increment by 32
// 	ppu.PPUCTRL = 0x04
// 	ppu.v = 0x1000
// 	_ = ppu.ReadRegister(0x2007)
// 	if ppu.v != 0x1020 {
// 		t.Errorf("expected v to increment by 32, got %04X", ppu.v)
// 	}
// }

// func TestPPU_PPUDATA_WriteIncrement1(t *testing.T) {
// 	ppu := &PPU{
// 		VRAM: [0x800]byte{},
// 		CHR:  make([]byte, 0x2000),
// 	}

// 	ppu.PPUCTRL = 0x00 // Bit 2 = 0 → increment by 1
// 	ppu.v = 0x2000

// 	ppu.WriteRegister(0x2007, 0xAB)

// 	if ppu.read(0x2000) != 0xAB {
// 		t.Errorf("expected 0xAB written to 0x2000, got %02X", ppu.read(0x2000))
// 	}

// 	if ppu.v != 0x2001 {
// 		t.Errorf("expected v to increment to 0x2001, got %04X", ppu.v)
// 	}
// }

// func TestPPU_PPUDATA_WriteIncrement32(t *testing.T) {
// 	ppu := &PPU{
// 		VRAM: [0x800]byte{},
// 		CHR:  make([]byte, 0x2000),
// 	}

// 	ppu.PPUCTRL = 0x04 // Bit 2 = 1 → increment by 32
// 	ppu.v = 0x2000

// 	ppu.WriteRegister(0x2007, 0xCD)

// 	if ppu.read(0x2000) != 0xCD {
// 		t.Errorf("expected 0xCD written to 0x2000, got %02X", ppu.read(0x2000))
// 	}

// 	if ppu.v != 0x2020 {
// 		t.Errorf("expected v to increment to 0x2020, got %04X", ppu.v)
// 	}
// }

// func TestPPU_PPUDATA_WriteToPalette(t *testing.T) {
// 	ppu := &PPU{
// 		PaletteTable: [32]byte{},
// 	}

// 	ppu.v = 0x3F05
// 	ppu.WriteRegister(0x2007, 0x77)

// 	if ppu.PaletteTable[5] != 0x77 {
// 		t.Errorf("expected palette index 5 to be 0x77, got %02X", ppu.PaletteTable[5])
// 	}
// }
