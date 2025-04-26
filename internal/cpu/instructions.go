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
	initLDAInstruction()
	initSTAInstruction()
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

func initLDAInstruction() {
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

func initSTAInstruction() {
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
