package rom

type MirroringType int

const (
	Horizontal MirroringType = iota
	Vertical
)

type Cartridge struct {
	PRG       []byte        // Program ROM
	CHR       []byte        // Character ROM
	Mapper    byte          // Mapper type
	Mirroring MirroringType // Mirroring type
	HasCHRROM bool          // Indicates if the cartridge has CHR ROM
}
