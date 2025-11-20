package ppu

import "testing"

func TestSprite0Hit(t *testing.T) {
	ppu := &PPU{
		// Enable background (bit 3) and sprites (bit 4)
		// Enable leftmost 8 pixels for both (bit 1 and 2)
		PPUMASK: 0x1E,
	}
	ppu.Reset()
	ppu.PPUMASK = 0x1E // Reset clears it

	// Setup Sprite 0
	// Y = 10 (so visible on scanline 11)
	// Tile = 1
	// Attr = 0
	// X = 5
	ppu.OAM[0] = 10
	ppu.OAM[1] = 1
	ppu.OAM[2] = 0
	ppu.OAM[3] = 5

	// We need to mock Read to return opaque pixels for both BG and Sprite
	// For Sprite 0 (Tile 1):
	// Fetch happens at... well, we can just manually set the internal buffers if we want to test renderPixel directly.
	// But renderPixel uses `spritePatternsLow/High` which are populated by `fetchSpritePatterns`.
	// And `fetchSpritePatterns` uses `secondaryOAM` populated by `evaluateSprites`.

	// So we need to run the pipeline or manually populate the intermediate state.
	// Let's manually populate for simplicity of unit testing `renderPixel`.

	ppu.scanline = 11
	ppu.cycle = 6 // X=5 corresponds to cycle 6 (cycle 1 is X=0)

	// Populate Sprite buffers
	ppu.spriteCount = 1
	ppu.spritePositions[0] = 5
	ppu.spriteAttributes[0] = 0
	ppu.spriteIndexes[0] = 0 // This is Sprite 0
	// Pattern: Solid block (all 1s)
	ppu.spritePatternsLow[0] = 0xFF
	ppu.spritePatternsHigh[0] = 0xFF

	// Populate BG buffers
	// We need `bgPatternLow` and `bgPatternHigh` to have a bit at the right position.
	// X=5. Shift = 15 - 5 = 10.
	// We want bit 10 to be 1.
	ppu.x = 0
	ppu.bgPatternLow = (1 << 15)
	ppu.bgPatternHigh = (1 << 15)

	// Clear Status
	ppu.PPUStatus = 0

	// Run renderPixel
	ppu.renderPixel()

	// Check Sprite 0 Hit
	if ppu.PPUStatus&0x40 == 0 {
		t.Errorf("Sprite 0 Hit failed. Status: 0x%02X", ppu.PPUStatus)
	}
}
