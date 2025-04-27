package cpu

import "fmt"

type AddressingMode int

const (
	Implied AddressingMode = iota
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
	Relative
)

type Instruction struct {
	Name    string
	Opcode  byte
	Bytes   int
	Cycles  int
	Mode    AddressingMode
	Execute func(cpu *CPU, addr uint16)
}

var Instructions = map[byte]Instruction{}

func init() {
	initLDAInstructions()
	initLDXInstructions()
	initLDYInstructions()
	initSTAInstructions()
	initSTXInstructions()
	initSTYInstructions()
	initTransferInstructions()
}

func (inst *Instruction) GetAddress(c *CPU) uint16 {
	var addr uint16
	switch inst.Mode {
	case Immediate:
		addr = c.fetchImediate()
	case ZeroPage:
		addr = c.fetchZeroPage()
	case ZeroPageX:
		addr = c.fetchZeroPageX()
	case ZeroPageY:
		addr = c.fetchZeroPageY()
	case Absolute:
		addr = c.fetchAbsolute()
	case AbsoluteX:
		addr = c.fetchAbsoluteX()
	case AbsoluteY:
		addr = c.fetchAbsoluteY()
	case IndirectX:
		addr = c.fetchIndirectX()
	case IndirectY:
		addr = c.fetchIndirectY()
	case Indirect:
		addr = c.fetchIndirect()
	case Implied:
		addr = c.fetchImplied()

	default:
		str := fmt.Sprintf("Unknown addressing mode: %d", inst.Mode)
		panic(str)
	}
	return addr
}

func ldaExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.A = value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func initLDAInstructions() {
	Instructions[0xA9] = Instruction{
		Name:    "LDA Immediate",
		Opcode:  0xA9,
		Bytes:   2,
		Cycles:  2,
		Execute: ldaExecute,
		Mode:    Immediate,
	}

	Instructions[0xA5] = Instruction{
		Name:    "LDA Zero Page",
		Opcode:  0xA5,
		Bytes:   2,
		Cycles:  3,
		Execute: ldaExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xB5] = Instruction{
		Name:    "LDA Zero Page,X",
		Opcode:  0xB5,
		Bytes:   2,
		Cycles:  4,
		Execute: ldaExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0xAD] = Instruction{
		Name:    "LDA Absolute",
		Opcode:  0xAD,
		Bytes:   3,
		Cycles:  4,
		Execute: ldaExecute,
		Mode:    Absolute,
	}

	Instructions[0xBD] = Instruction{
		Name:    "LDA Absolute,X",
		Opcode:  0xBD,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: ldaExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0xB9] = Instruction{
		Name:    "LDA Absolute,Y",
		Opcode:  0xB9,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: ldaExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0xA1] = Instruction{
		Name:    "LDA (Indirect,X)",
		Opcode:  0xA1,
		Bytes:   2,
		Cycles:  6,
		Execute: ldaExecute,
		Mode:    IndirectX,
	}

	Instructions[0xB1] = Instruction{
		Name:    "LDA (Indirect),Y",
		Opcode:  0xB1,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: ldaExecute,
		Mode:    IndirectY,
	}
}

func staExecute(cpu *CPU, addr uint16) {
	cpu.Memory[addr] = cpu.A
	// STA does not affect any flags
}

func initSTAInstructions() {
	Instructions[0x85] = Instruction{
		Name:    "STA Zero Page",
		Opcode:  0x85,
		Bytes:   2,
		Cycles:  3,
		Mode:    ZeroPage,
		Execute: staExecute,
	}

	Instructions[0x95] = Instruction{
		Name:    "STA Zero Page,X",
		Opcode:  0x95,
		Bytes:   2,
		Cycles:  4,
		Mode:    ZeroPageX,
		Execute: staExecute,
	}

	Instructions[0x8D] = Instruction{
		Name:    "STA Absolute",
		Opcode:  0x8D,
		Bytes:   3,
		Cycles:  4,
		Mode:    Absolute,
		Execute: staExecute,
	}

	Instructions[0x9D] = Instruction{
		Name:    "STA Absolute,X",
		Opcode:  0x9D,
		Bytes:   3,
		Cycles:  5,
		Mode:    AbsoluteX,
		Execute: staExecute,
	}

	Instructions[0x99] = Instruction{
		Name:    "STA Absolute,Y",
		Opcode:  0x99,
		Bytes:   3,
		Cycles:  5,
		Mode:    AbsoluteY,
		Execute: staExecute,
	}

	Instructions[0x81] = Instruction{
		Name:    "STA (Indirect,X)",
		Opcode:  0x81,
		Bytes:   2,
		Cycles:  6,
		Mode:    IndirectX,
		Execute: staExecute,
	}

	Instructions[0x91] = Instruction{
		Name:    "STA (Indirect),Y",
		Opcode:  0x91,
		Bytes:   2,
		Cycles:  6,
		Mode:    IndirectY,
		Execute: staExecute,
	}

}

func initLDXInstructions() {
	Instructions[0xA2] = Instruction{
		Name:    "LDX Immediate",
		Opcode:  0xA2,
		Bytes:   2,
		Cycles:  2,
		Execute: ldxExecute,
		Mode:    Immediate,
	}

	Instructions[0xA6] = Instruction{
		Name:    "LDX Zero Page",
		Opcode:  0xA6,
		Bytes:   2,
		Cycles:  3,
		Execute: ldxExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xB6] = Instruction{
		Name:    "LDX Zero Page,Y",
		Opcode:  0xB6,
		Bytes:   2,
		Cycles:  4,
		Execute: ldxExecute,
		Mode:    ZeroPageY,
	}

	Instructions[0xAE] = Instruction{
		Name:    "LDX Absolute",
		Opcode:  0xAE,
		Bytes:   3,
		Cycles:  4,
		Execute: ldxExecute,
		Mode:    Absolute,
	}

	Instructions[0xBE] = Instruction{
		Name:    "LDX Absolute,Y",
		Opcode:  0xBE,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: ldxExecute,
		Mode:    AbsoluteY,
	}
}

func ldxExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.X = value
	cpu.SetFlag(FlagZ, cpu.X == 0)
	cpu.SetFlag(FlagN, (cpu.X&0x80) != 0)
}

func initLDYInstructions() {
	Instructions[0xA0] = Instruction{
		Name:    "LDY Immediate",
		Opcode:  0xA0,
		Bytes:   2,
		Cycles:  2,
		Execute: ldyExecute,
		Mode:    Immediate,
	}

	Instructions[0xA4] = Instruction{
		Name:    "LDY Zero Page",
		Opcode:  0xA4,
		Bytes:   2,
		Cycles:  3,
		Execute: ldyExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xB4] = Instruction{
		Name:    "LDY Zero Page,X",
		Opcode:  0xB4,
		Bytes:   2,
		Cycles:  4,
		Execute: ldyExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0xAC] = Instruction{
		Name:    "LDY Absolute",
		Opcode:  0xAC,
		Bytes:   3,
		Cycles:  4,
		Execute: ldyExecute,
		Mode:    Absolute,
	}

	Instructions[0xBC] = Instruction{
		Name:    "LDY Absolute,X",
		Opcode:  0xBC,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: ldyExecute,
		Mode:    AbsoluteX,
	}
}

func ldyExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.Y = value
	cpu.SetFlag(FlagZ, cpu.Y == 0)
	cpu.SetFlag(FlagN, (cpu.Y&0x80) != 0)
}

func initSTXInstructions() {
	Instructions[0x86] = Instruction{
		Name:    "STX Zero Page",
		Opcode:  0x86,
		Bytes:   2,
		Cycles:  3,
		Mode:    ZeroPage,
		Execute: stxExecute,
	}

	Instructions[0x96] = Instruction{
		Name:    "STX Zero Page,Y",
		Opcode:  0x96,
		Bytes:   2,
		Cycles:  4,
		Mode:    ZeroPageY,
		Execute: stxExecute,
	}

	Instructions[0x8E] = Instruction{
		Name:    "STX Absolute",
		Opcode:  0x8E,
		Bytes:   3,
		Cycles:  4,
		Mode:    Absolute,
		Execute: stxExecute,
	}
}

func stxExecute(cpu *CPU, addr uint16) {
	cpu.Memory[addr] = cpu.X
	// STX does not affect any flags
}

func initSTYInstructions() {
	Instructions[0x84] = Instruction{
		Name:    "STY Zero Page",
		Opcode:  0x84,
		Bytes:   2,
		Cycles:  3,
		Mode:    ZeroPage,
		Execute: styExecute,
	}

	Instructions[0x94] = Instruction{
		Name:    "STY Zero Page,X",
		Opcode:  0x94,
		Bytes:   2,
		Cycles:  4,
		Mode:    ZeroPageX,
		Execute: styExecute,
	}

	Instructions[0x8C] = Instruction{
		Name:    "STY Absolute",
		Opcode:  0x8C,
		Bytes:   3,
		Cycles:  4,
		Mode:    Absolute,
		Execute: styExecute,
	}
}

func styExecute(cpu *CPU, addr uint16) {
	cpu.Memory[addr] = cpu.Y
	// STY does not affect any flags
}

func initTransferInstructions() {
	Instructions[0xAA] = Instruction{
		Name:   "TAX",
		Opcode: 0xAA,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.X = cpu.A
			cpu.SetFlag(FlagZ, cpu.X == 0)
			cpu.SetFlag(FlagN, (cpu.X&0x80) != 0)
		},
	}

	Instructions[0xA8] = Instruction{
		Name:   "TAY",
		Opcode: 0xA8,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.Y = cpu.A
			cpu.SetFlag(FlagZ, cpu.Y == 0)
			cpu.SetFlag(FlagN, (cpu.Y&0x80) != 0)
		},
	}

	Instructions[0xBA] = Instruction{
		Name:   "TSX",
		Opcode: 0xBA,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.X = cpu.SP
			cpu.SetFlag(FlagZ, cpu.X == 0)
			cpu.SetFlag(FlagN, (cpu.X&0x80) != 0)
		},
	}

	Instructions[0x8A] = Instruction{
		Name:   "TXA",
		Opcode: 0x8A,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.A = cpu.X
			cpu.SetFlag(FlagZ, cpu.A == 0)
			cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
		},
	}

	Instructions[0x9A] = Instruction{
		Name:   "TXS",
		Opcode: 0x9A,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SP = cpu.X
			// TXS does not affect any flags
		},
	}

	Instructions[0x98] = Instruction{
		Name:   "TYA",
		Opcode: 0x98,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.A = cpu.Y
			cpu.SetFlag(FlagZ, cpu.A == 0)
			cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
		},
	}
}
