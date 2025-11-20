package ppu

import "testing"

func TestRenderPixel_FineX(t *testing.T) {
	ppu := &PPU{
		// Enable background rendering (Bit 3 = 1)
		// Enable leftmost 8 pixels (Bit 1 = 1)
		PPUMASK: 0x0A,
	}

	// Setup shift registers with a known pattern
	// We want to see which bit is picked.
	// Let's put 1s at specific positions to identify them.
	// Bit 15 is the leftmost pixel (Fine X = 0)
	// Bit 14 is the next pixel (Fine X = 1)
	// ...
	// Bit 8 is the 8th pixel (Fine X = 7)

	// Case 1: Fine X = 0. Should pick bit 15.
	ppu.x = 0
	ppu.bgPatternLow = 0x8000 // Bit 15 set
	ppu.bgPatternHigh = 0x0000
	ppu.bgAttributeLow = 0x0000
	ppu.bgAttributeHigh = 0x0000

	// We need to mock the palette lookup or check the intermediate pixel value.
	// Since renderPixel writes to framebuffer, let's check the framebuffer.
	// But renderPixel logic is:
	// bit0 := (bgPatternLow >> (15-x)) & 1
	// bit1 := (bgPatternHigh >> (15-x)) & 1
	// pixel = bit0 | (bit1 << 1)

	// If logic is correct:
	// x=0 -> shift 15 -> bit 15
	// x=1 -> shift 14 -> bit 14

	// Let's run renderPixel for cycle 1, scanline 0
	ppu.scanline = 0
	ppu.cycle = 1

	// Clear framebuffer
	ppu.framebuffer[0][0] = 0

	// We need to ensure Read(0x3F00 + ...) returns something predictable.
	// But wait, renderPixel calls ppu.Read which might be complex.
	// However, if we just want to verify the bit selection, we can rely on the fact that
	// if pixel value is 0, it reads 0x3F00 (universal bg).
	// If pixel value is != 0, it reads palette.

	// Let's make sure Read returns different values for different addresses if possible,
	// or just check if it calculates the correct index.
	// Since we can't easily mock Read without changing the struct or using an interface,
	// let's just assume PaletteTable is 0s and we check if it tries to read from the right place?
	// Actually, we can just set the PaletteTable!

	// Set Universal BG color (index 0) to 0x01
	ppu.PaletteTable[0] = 0x01
	// Set Palette 0 Color 1 (index 1) to 0x02
	ppu.PaletteTable[1] = 0x02

	// Run renderPixel
	ppu.renderPixel()

	// With bgPatternLow=0x8000 and x=0, pixel value should be 1.
	// So it should use color from palette index 1.
	// We expect framebuffer[0][0] to be 0x02 (if Read works as expected on PaletteTable)

	// Wait, ppu.Read maps 0x3F00-0x3FFF to PaletteTable.
	// We need to ensure ppu.Read works for this test.
	// ppu.Read is not fully shown in the snippets, but usually it handles palette mirroring.
	// Let's assume standard behavior or check ppu.go again if needed.
	// But for now, let's try to run it.

	// Actually, let's check the failing case first (current implementation).
	// Current impl: always shifts by 15.
	// So if we set x=1, and put bit at 14, current impl will miss it (read bit 15 which is 0).

	// Test Case 2: Fine X = 1. Should pick bit 14.
	ppu.x = 1
	ppu.bgPatternLow = 0x4000 // Bit 14 set (0100 0000 0000 0000)
	ppu.bgPatternHigh = 0x0000

	ppu.cycle = 2 // Move to next pixel to avoid overwriting if we were tracking that
	// But renderPixel uses ppu.cycle-1 for index.

	ppu.renderPixel()

	// If fixed: pixel=1 -> color 0x02.
	// If broken: pixel=0 (reads bit 15 which is 0) -> color 0x01.

	if ppu.framebuffer[0][1] != 0x02 {
		t.Errorf("Fine X=1 failed. Expected color 0x02 (pixel=1), got 0x%02X. The bit selection might be wrong.", ppu.framebuffer[0][1])
	}
}
