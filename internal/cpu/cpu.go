package cpu

type CPU struct {
	A  byte   // Accumulator
	X  byte   // Index Register X
	Y  byte   // Index Register Y
	SP byte   // Stack Pointer
	PC uint16 // Program Counter
	P  byte   // Status flags

	Memory [0x10000]byte
}

func New() *CPU {
	return &CPU{
		Memory: [0x10000]byte{},
	}
}
