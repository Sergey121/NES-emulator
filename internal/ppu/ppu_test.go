package ppu

import "testing"

func TestPPUVBlankFlagAndNMI(t *testing.T) {
	ppu := &PPU{
		PPUCTRL: 0x80, // Enable NMI
	}

	// Simulate the scanline and cycle for VBlank
	ppu.scanline = 241
	ppu.cycle = 1

	// Execute the PPU
	ppu.Step()

	// Check if the VBlank flag is set
	if ppu.PPUStatus&0x80 == 0 {
		t.Errorf("Expected VBlank flag to be set, but it was not")
	}

	// Check if NMI occurred
	if !ppu.nmiOccurred {
		t.Errorf("Expected NMI to occur, but it did not")
	}

	// Check if NMI output is enabled
	if !ppu.nmiOutput {
		t.Errorf("Expected NMI output to be enabled, but it was not")
	}
}

func TestPPUVBlankFlagAndNMI2(t *testing.T) {
	ppu := &PPU{}

	// Simulate the scanline and cycle for VBlank
	for ppu.Scanline() != 241 || ppu.Cycle() != 1 {
		ppu.Step()
	}

	ppu.Step()

	// Check if the VBlank flag is set
	if ppu.PPUStatus&0x80 == 0 {
		t.Errorf("Expected VBlank flag to be set, but it was not")
	}

	// Check if NMI occurred
	if !ppu.NMIOccurred() {
		t.Errorf("Expected NMI to occur, but it did not")
	}

	// Move scanline for a new frame
	for ppu.Scanline() != 261 || ppu.Cycle() != 1 {
		ppu.Step()
	}

	ppu.Step()

	if ppu.PPUStatus&0x80 != 0 {
		t.Errorf("Expected VBlank flag to be cleared, but it was not")
	}

	// Check if NMI occurred
	if ppu.NMIOccurred() {
		t.Errorf("Expected NMI to not occur, but it did")
	}
}
