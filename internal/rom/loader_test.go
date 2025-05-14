package rom

import "testing"

func TestCreateCartridge(t *testing.T) {
	prgSize := make([]byte, 16*1024)
	// Test with a valid ROM file
	data := []byte{
		'N', 'E', 'S', 0x1A, // NES header
		1, 0, // PRG and CHR banks
		0x00, 0x00, // Mapper and mirroring flags
		0, 0, 0, 0, 0, 0, 0, 0, // Padding
	}
	data = append(data, prgSize...)
	cartridge, err := createCartridge(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cartridge.PRG == nil || cartridge.CHR == nil {
		t.Fatal("expected PRG and CHR to be initialized")
	}

	// Test with an invalid ROM file
	data = []byte{
		'F', 'A', 'K', 0x1A, // Invalid header
	}
	_, err = createCartridge(data)
	if err == nil {
		t.Fatal("expected an error for invalid ROM file")
	}
}
