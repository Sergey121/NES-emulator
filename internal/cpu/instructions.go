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
	Accumulator
)

type Instruction struct {
	Name       string
	Opcode     byte
	Bytes      int
	Cycles     int
	Mode       AddressingMode
	Execute    func(cpu *CPU, addr uint16)
	ModifiesPC bool
}

var Instructions = map[byte]Instruction{}

func init() {
	initADCInstructions()
	initLDAInstructions()
	initLDXInstructions()
	initLDYInstructions()
	initSTAInstructions()
	initSTXInstructions()
	initSTYInstructions()
	initTransferInstructions()
	initFlagInstructions()
	initSBCInstructions()
	initANDInstructions()
	initEORInstructions()
	initORAInstructions()
	initCMPInstructions()
	initCPXInstructions()
	initCPYInstructions()
	initASLInstructions()
	initLSRInstructions()
	initRORInstructions()
	initROLInstructions()
	initRTIInstructions()
	initRTSInstructions()
	initJSRInstructions()
	initBRKInstructions()
	initJMPInstructions()
	initBEQInstructions()
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
	case Accumulator:
		addr = c.fetchAccumulator()
	case Relative:
		addr = c.fetchRelative()

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

func initFlagInstructions() {
	Instructions[0x18] = Instruction{
		Name:   "CLC",
		Opcode: 0x18,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagC, false)
		},
	}

	Instructions[0x38] = Instruction{
		Name:   "SEC",
		Opcode: 0x38,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagC, true)
		},
	}

	Instructions[0xD8] = Instruction{
		Name:   "CLD",
		Opcode: 0xD8,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagD, false)
		},
	}

	Instructions[0xF8] = Instruction{
		Name:   "SED",
		Opcode: 0xF8,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagD, true)
		},
	}

	Instructions[0x58] = Instruction{
		Name:   "CLI",
		Opcode: 0x58,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagI, false)
		},
	}

	Instructions[0x78] = Instruction{
		Name:   "SEI",
		Opcode: 0x78,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagI, true)
		},
	}

	Instructions[0xB8] = Instruction{
		Name:   "CLV",
		Opcode: 0xB8,
		Bytes:  1,
		Cycles: 2,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.SetFlag(FlagV, false)
		},
	}
}

func initADCInstructions() {
	Instructions[0x69] = Instruction{
		Name:    "ADC Immediate",
		Opcode:  0x69,
		Bytes:   2,
		Cycles:  2,
		Execute: adcExecute,
		Mode:    Immediate,
	}

	Instructions[0x65] = Instruction{
		Name:    "ADC Zero Page",
		Opcode:  0x65,
		Bytes:   2,
		Cycles:  3,
		Execute: adcExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x75] = Instruction{
		Name:    "ADC Zero Page,X",
		Opcode:  0x75,
		Bytes:   2,
		Cycles:  4,
		Execute: adcExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x6D] = Instruction{
		Name:    "ADC Absolute",
		Opcode:  0x6D,
		Bytes:   3,
		Cycles:  4,
		Execute: adcExecute,
		Mode:    Absolute,
	}

	Instructions[0x7D] = Instruction{
		Name:    "ADC Absolute,X",
		Opcode:  0x7D,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: adcExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0x79] = Instruction{
		Name:    "ADC Absolute,Y",
		Opcode:  0x79,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: adcExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0x61] = Instruction{
		Name:    "ADC (Indirect,X)",
		Opcode:  0x61,
		Bytes:   2,
		Cycles:  6,
		Execute: adcExecute,
		Mode:    IndirectX,
	}

	Instructions[0x71] = Instruction{
		Name:    "ADC (Indirect),Y",
		Opcode:  0x71,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: adcExecute,
		Mode:    IndirectY,
	}
}

func adcExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	carry := 0
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(value) + uint16(carry)

	cpu.SetFlag(FlagC, result > 0xFF)
	cpu.SetFlag(FlagZ, byte(result&0xFF) == 0)
	cpu.SetFlag(FlagN, (result&0x80) != 0)
	cpu.SetFlag(FlagV, (^(uint16(cpu.A)^uint16(value))&(uint16(cpu.A)^result)&0x80) != 0)

	cpu.A = byte(result & 0xFF)
}

func initSBCInstructions() {
	Instructions[0xE9] = Instruction{
		Name:    "SBC Immediate",
		Opcode:  0xE9,
		Bytes:   2,
		Cycles:  2,
		Execute: sbcExecute,
		Mode:    Immediate,
	}

	Instructions[0xE5] = Instruction{
		Name:    "SBC Zero Page",
		Opcode:  0xE5,
		Bytes:   2,
		Cycles:  3,
		Execute: sbcExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xF5] = Instruction{
		Name:    "SBC Zero Page,X",
		Opcode:  0xF5,
		Bytes:   2,
		Cycles:  4,
		Execute: sbcExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0xED] = Instruction{
		Name:    "SBC Absolute",
		Opcode:  0xED,
		Bytes:   3,
		Cycles:  4,
		Execute: sbcExecute,
		Mode:    Absolute,
	}

	Instructions[0xFD] = Instruction{
		Name:    "SBC Absolute,X",
		Opcode:  0xFD,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: sbcExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0xF9] = Instruction{
		Name:    "SBC Absolute,Y",
		Opcode:  0xF9,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: sbcExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0xE1] = Instruction{
		Name:    "SBC (Indirect,X)",
		Opcode:  0xE1,
		Bytes:   2,
		Cycles:  6,
		Execute: sbcExecute,
		Mode:    IndirectX,
	}

	Instructions[0xF1] = Instruction{
		Name:    "SBC (Indirect),Y",
		Opcode:  0xF1,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: sbcExecute,
		Mode:    IndirectY,
	}
}

func sbcExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	carryIn := 0
	if cpu.GetFlag(FlagC) {
		carryIn = 1
	}
	// Сначала делаем инверсию carry: (1 - carryIn)
	result := uint16(cpu.A) - uint16(value) - (1 - uint16(carryIn))

	cpu.SetFlag(FlagC, result < 0x100)
	cpu.SetFlag(FlagZ, byte(result&0xFF) == 0)
	cpu.SetFlag(FlagN, (result&0x80) != 0)
	cpu.SetFlag(FlagV, ((uint16(cpu.A)^result)&(uint16(cpu.A)^uint16(value))&0x80) != 0)

	cpu.A = byte(result & 0xFF)
}

func initANDInstructions() {
	Instructions[0x29] = Instruction{
		Name:    "AND Immediate",
		Opcode:  0x29,
		Bytes:   2,
		Cycles:  2,
		Execute: andExecute,
		Mode:    Immediate,
	}

	Instructions[0x25] = Instruction{
		Name:    "AND Zero Page",
		Opcode:  0x25,
		Bytes:   2,
		Cycles:  3,
		Execute: andExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x35] = Instruction{
		Name:    "AND Zero Page,X",
		Opcode:  0x35,
		Bytes:   2,
		Cycles:  4,
		Execute: andExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x2D] = Instruction{
		Name:    "AND Absolute",
		Opcode:  0x2D,
		Bytes:   3,
		Cycles:  4,
		Execute: andExecute,
		Mode:    Absolute,
	}

	Instructions[0x3D] = Instruction{
		Name:    "AND Absolute,X",
		Opcode:  0x3D,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: andExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0x39] = Instruction{
		Name:    "AND Absolute,Y",
		Opcode:  0x39,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: andExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0x21] = Instruction{
		Name:    "AND (Indirect,X)",
		Opcode:  0x21,
		Bytes:   2,
		Cycles:  6,
		Execute: andExecute,
		Mode:    IndirectX,
	}

	Instructions[0x31] = Instruction{
		Name:    "AND (Indirect),Y",
		Opcode:  0x31,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: andExecute,
		Mode:    IndirectY,
	}
}

func andExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.A &= value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func initEORInstructions() {
	Instructions[0x49] = Instruction{
		Name:    "EOR Immediate",
		Opcode:  0x49,
		Bytes:   2,
		Cycles:  2,
		Execute: eorExecute,
		Mode:    Immediate,
	}

	Instructions[0x45] = Instruction{
		Name:    "EOR Zero Page",
		Opcode:  0x45,
		Bytes:   2,
		Cycles:  3,
		Execute: eorExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x55] = Instruction{
		Name:    "EOR Zero Page,X",
		Opcode:  0x55,
		Bytes:   2,
		Cycles:  4,
		Execute: eorExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x4D] = Instruction{
		Name:    "EOR Absolute",
		Opcode:  0x4D,
		Bytes:   3,
		Cycles:  4,
		Execute: eorExecute,
		Mode:    Absolute,
	}

	Instructions[0x5D] = Instruction{
		Name:    "EOR Absolute,X",
		Opcode:  0x5D,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: eorExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0x59] = Instruction{
		Name:    "EOR Absolute,Y",
		Opcode:  0x59,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: eorExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0x41] = Instruction{
		Name:    "EOR (Indirect,X)",
		Opcode:  0x41,
		Bytes:   2,
		Cycles:  6,
		Execute: eorExecute,
		Mode:    IndirectX,
	}

	Instructions[0x51] = Instruction{
		Name:    "EOR (Indirect),Y",
		Opcode:  0x51,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: eorExecute,
		Mode:    IndirectY,
	}
}

func eorExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.A ^= value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func initORAInstructions() {
	Instructions[0x09] = Instruction{
		Name:    "ORA Immediate",
		Opcode:  0x09,
		Bytes:   2,
		Cycles:  2,
		Execute: oraExecute,
		Mode:    Immediate,
	}

	Instructions[0x05] = Instruction{
		Name:    "ORA Zero Page",
		Opcode:  0x05,
		Bytes:   2,
		Cycles:  3,
		Execute: oraExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x15] = Instruction{
		Name:    "ORA Zero Page,X",
		Opcode:  0x15,
		Bytes:   2,
		Cycles:  4,
		Execute: oraExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x0D] = Instruction{
		Name:    "ORA Absolute",
		Opcode:  0x0D,
		Bytes:   3,
		Cycles:  4,
		Execute: oraExecute,
		Mode:    Absolute,
	}

	Instructions[0x1D] = Instruction{
		Name:    "ORA Absolute,X",
		Opcode:  0x1D,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: oraExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0x19] = Instruction{
		Name:    "ORA Absolute,Y",
		Opcode:  0x19,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: oraExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0x01] = Instruction{
		Name:    "ORA (Indirect,X)",
		Opcode:  0x01,
		Bytes:   2,
		Cycles:  6,
		Execute: oraExecute,
		Mode:    IndirectX,
	}

	Instructions[0x11] = Instruction{
		Name:    "ORA (Indirect),Y",
		Opcode:  0x11,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: oraExecute,
		Mode:    IndirectY,
	}
}

func oraExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.A |= value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func initCMPInstructions() {
	Instructions[0xC9] = Instruction{
		Name:    "CMP Immediate",
		Opcode:  0xC9,
		Bytes:   2,
		Cycles:  2,
		Execute: cmpExecute,
		Mode:    Immediate,
	}

	Instructions[0xC5] = Instruction{
		Name:    "CMP Zero Page",
		Opcode:  0xC5,
		Bytes:   2,
		Cycles:  3,
		Execute: cmpExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xD5] = Instruction{
		Name:    "CMP Zero Page,X",
		Opcode:  0xD5,
		Bytes:   2,
		Cycles:  4,
		Execute: cmpExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0xCD] = Instruction{
		Name:    "CMP Absolute",
		Opcode:  0xCD,
		Bytes:   3,
		Cycles:  4,
		Execute: cmpExecute,
		Mode:    Absolute,
	}

	Instructions[0xDD] = Instruction{
		Name:    "CMP Absolute,X",
		Opcode:  0xDD,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: cmpExecute,
		Mode:    AbsoluteX,
	}

	Instructions[0xD9] = Instruction{
		Name:    "CMP Absolute,Y",
		Opcode:  0xD9,
		Bytes:   3,
		Cycles:  4, // add 1 to cycles if page boundary is crossed
		Execute: cmpExecute,
		Mode:    AbsoluteY,
	}

	Instructions[0xC1] = Instruction{
		Name:    "CMP (Indirect,X)",
		Opcode:  0xC1,
		Bytes:   2,
		Cycles:  6,
		Execute: cmpExecute,
		Mode:    IndirectX,
	}

	Instructions[0xD1] = Instruction{
		Name:    "CMP (Indirect),Y",
		Opcode:  0xD1,
		Bytes:   2,
		Cycles:  5, // add 1 to cycles if page boundary is crossed
		Execute: cmpExecute,
		Mode:    IndirectY,
	}
}

func cmpExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	result := uint16(cpu.A) - uint16(value)

	cpu.SetFlag(FlagC, cpu.A >= value)
	cpu.SetFlag(FlagZ, byte(result&0xFF) == 0)
	cpu.SetFlag(FlagN, (result&0x80) != 0)
}

func initCPXInstructions() {
	Instructions[0xE0] = Instruction{
		Name:    "CPX Immediate",
		Opcode:  0xE0,
		Bytes:   2,
		Cycles:  2,
		Execute: cpxExecute,
		Mode:    Immediate,
	}

	Instructions[0xE4] = Instruction{
		Name:    "CPX Zero Page",
		Opcode:  0xE4,
		Bytes:   2,
		Cycles:  3,
		Execute: cpxExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xEC] = Instruction{
		Name:    "CPX Absolute",
		Opcode:  0xEC,
		Bytes:   3,
		Cycles:  4,
		Execute: cpxExecute,
		Mode:    Absolute,
	}
}

func cpxExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	result := uint16(cpu.X) - uint16(value)

	cpu.SetFlag(FlagC, cpu.X >= value)
	cpu.SetFlag(FlagZ, byte(result&0xFF) == 0)
	cpu.SetFlag(FlagN, (result&0x80) != 0)
}

func initCPYInstructions() {
	Instructions[0xC0] = Instruction{
		Name:    "CPY Immediate",
		Opcode:  0xC0,
		Bytes:   2,
		Cycles:  2,
		Execute: cpyExecute,
		Mode:    Immediate,
	}

	Instructions[0xC4] = Instruction{
		Name:    "CPY Zero Page",
		Opcode:  0xC4,
		Bytes:   2,
		Cycles:  3,
		Execute: cpyExecute,
		Mode:    ZeroPage,
	}

	Instructions[0xCC] = Instruction{
		Name:    "CPY Absolute",
		Opcode:  0xCC,
		Bytes:   3,
		Cycles:  4,
		Execute: cpyExecute,
		Mode:    Absolute,
	}
}

func cpyExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	result := uint16(cpu.Y) - uint16(value)

	cpu.SetFlag(FlagC, cpu.Y >= value)
	cpu.SetFlag(FlagZ, byte(result&0xFF) == 0)
	cpu.SetFlag(FlagN, (result&0x80) != 0)
}

func initASLInstructions() {
	Instructions[0x0A] = Instruction{
		Name:    "ASL Accumulator",
		Opcode:  0x0A,
		Bytes:   1,
		Cycles:  2,
		Execute: aslAExecute,
		Mode:    Accumulator,
	}

	Instructions[0x06] = Instruction{
		Name:    "ASL Zero Page",
		Opcode:  0x06,
		Bytes:   2,
		Cycles:  5,
		Execute: aslExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x16] = Instruction{
		Name:    "ASL Zero Page,X",
		Opcode:  0x16,
		Bytes:   2,
		Cycles:  6,
		Execute: aslExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x0E] = Instruction{
		Name:    "ASL Absolute",
		Opcode:  0x0E,
		Bytes:   3,
		Cycles:  6,
		Execute: aslExecute,
		Mode:    Absolute,
	}

	Instructions[0x1E] = Instruction{
		Name:    "ASL Absolute,X",
		Opcode:  0x1E,
		Bytes:   3,
		Cycles:  7,
		Execute: aslExecute,
		Mode:    AbsoluteX,
	}
}

func aslAExecute(cpu *CPU, _ uint16) {
	value := cpu.A
	cpu.SetFlag(FlagC, (value&0x80) != 0)
	value <<= 1
	cpu.A = value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func aslExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.SetFlag(FlagC, (value&0x80) != 0)
	value <<= 1
	cpu.Memory[addr] = value
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, (value&0x80) != 0)
}

func initLSRInstructions() {
	Instructions[0x4A] = Instruction{
		Name:    "LSR Accumulator",
		Opcode:  0x4A,
		Bytes:   1,
		Cycles:  2,
		Execute: lsrAExecute,
		Mode:    Accumulator,
	}

	Instructions[0x46] = Instruction{
		Name:    "LSR Zero Page",
		Opcode:  0x46,
		Bytes:   2,
		Cycles:  5,
		Execute: lsrExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x56] = Instruction{
		Name:    "LSR Zero Page,X",
		Opcode:  0x56,
		Bytes:   2,
		Cycles:  6,
		Execute: lsrExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x4E] = Instruction{
		Name:    "LSR Absolute",
		Opcode:  0x4E,
		Bytes:   3,
		Cycles:  6,
		Execute: lsrExecute,
		Mode:    Absolute,
	}

	Instructions[0x5E] = Instruction{
		Name:    "LSR Absolute,X",
		Opcode:  0x5E,
		Bytes:   3,
		Cycles:  7,
		Execute: lsrExecute,
	}
}

func lsrAExecute(cpu *CPU, _ uint16) {
	value := cpu.A
	cpu.SetFlag(FlagC, (value&0x01) != 0)
	value >>= 1
	cpu.A = value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
}

func lsrExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	cpu.SetFlag(FlagC, (value&0x01) != 0)
	value >>= 1
	cpu.Memory[addr] = value
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, false)
}

func initRORInstructions() {
	Instructions[0x6A] = Instruction{
		Name:    "ROR Accumulator",
		Opcode:  0x6A,
		Bytes:   1,
		Cycles:  2,
		Execute: rorAExecute,
		Mode:    Accumulator,
	}

	Instructions[0x66] = Instruction{
		Name:    "ROR Zero Page",
		Opcode:  0x66,
		Bytes:   2,
		Cycles:  5,
		Execute: rorExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x76] = Instruction{
		Name:    "ROR Zero Page,X",
		Opcode:  0x76,
		Bytes:   2,
		Cycles:  6,
		Execute: rorExecute,
		Mode:    ZeroPageX,
	}
	Instructions[0x6E] = Instruction{
		Name:    "ROR Absolute",
		Opcode:  0x6E,
		Bytes:   3,
		Cycles:  6,
		Execute: rorExecute,
		Mode:    Absolute,
	}
	Instructions[0x7E] = Instruction{
		Name:    "ROR Absolute,X",
		Opcode:  0x7E,
		Bytes:   3,
		Cycles:  7,
		Execute: rorExecute,
		Mode:    AbsoluteX,
	}
}

func rorAExecute(cpu *CPU, _ uint16) {
	value := cpu.A
	carry := cpu.GetFlag(FlagC)
	cpu.SetFlag(FlagC, (value&0x01) != 0)
	value >>= 1
	if carry {
		value |= 0x80
	}
	cpu.A = value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func rorExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	carry := cpu.GetFlag(FlagC)
	cpu.SetFlag(FlagC, (value&0x01) != 0)
	value >>= 1
	if carry {
		value |= 0x80
	}
	cpu.Memory[addr] = value
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, (value&0x80) != 0)
}

func initROLInstructions() {
	Instructions[0x2A] = Instruction{
		Name:    "ROL Accumulator",
		Opcode:  0x2A,
		Bytes:   1,
		Cycles:  2,
		Execute: rolAExecute,
		Mode:    Accumulator,
	}

	Instructions[0x26] = Instruction{
		Name:    "ROL Zero Page",
		Opcode:  0x26,
		Bytes:   2,
		Cycles:  5,
		Execute: rolExecute,
		Mode:    ZeroPage,
	}

	Instructions[0x36] = Instruction{
		Name:    "ROL Zero Page,X",
		Opcode:  0x36,
		Bytes:   2,
		Cycles:  6,
		Execute: rolExecute,
		Mode:    ZeroPageX,
	}

	Instructions[0x2E] = Instruction{
		Name:    "ROL Absolute",
		Opcode:  0x2E,
		Bytes:   3,
		Cycles:  6,
		Execute: rolExecute,
		Mode:    Absolute,
	}

	Instructions[0x3E] = Instruction{
		Name:    "ROL Absolute,X",
		Opcode:  0x3E,
		Bytes:   3,
		Cycles:  7,
		Execute: rolExecute,
		Mode:    AbsoluteX,
	}
}

func rolAExecute(cpu *CPU, _ uint16) {
	value := cpu.A
	carry := cpu.GetFlag(FlagC)
	cpu.SetFlag(FlagC, (value&0x80) != 0)
	value <<= 1
	if carry {
		value |= 0x01
	}
	cpu.A = value
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}

func rolExecute(cpu *CPU, addr uint16) {
	value := cpu.Memory[addr]
	carry := cpu.GetFlag(FlagC)
	cpu.SetFlag(FlagC, (value&0x80) != 0)
	value <<= 1
	if carry {
		value |= 0x01
	}
	cpu.Memory[addr] = value
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, (value&0x80) != 0)
}

func initRTIInstructions() {
	Instructions[0x40] = Instruction{
		Name:   "RTI",
		Opcode: 0x40,
		Bytes:  1,
		Cycles: 6,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			status := cpu.Pull()
			cpu.SetStatus(status)
			cpu.PC = cpu.Pull16()
		},
		ModifiesPC: true,
	}
}

func initRTSInstructions() {
	Instructions[0x60] = Instruction{
		Name:   "RTS",
		Opcode: 0x60,
		Bytes:  1,
		Cycles: 6,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			cpu.PC = cpu.Pull16() + 1
		},
		ModifiesPC: true,
	}
}

func initJSRInstructions() {
	Instructions[0x20] = Instruction{
		Name:   "JSR",
		Opcode: 0x20,
		Bytes:  3,
		Cycles: 6,
		Mode:   Absolute,
		Execute: func(cpu *CPU, addr uint16) {
			// Push address of last byte of JSR instruction (PC+2)
			cpu.Push16(cpu.PC + 2)
			cpu.PC = addr
		},
		ModifiesPC: true,
	}
}

func initBRKInstructions() {
	Instructions[0x00] = Instruction{
		Name:   "BRK",
		Opcode: 0x00,
		Bytes:  2, // формально 2 байта
		Cycles: 7,
		Mode:   Implied,
		Execute: func(cpu *CPU, _ uint16) {
			// Push PC + 2 на стек (указатель на следующую инструкцию)
			cpu.Push16(cpu.PC + 2)

			// Сохраняем статус-регистр с установленным флагом B и U
			cpu.Push(cpu.P | FlagB | FlagU)

			// Переход по адресу из вектора прерываний (0xFFFE/F)
			lo := cpu.Memory[0xFFFE]
			hi := cpu.Memory[0xFFFF]
			cpu.PC = uint16(lo) | uint16(hi)<<8
		},
		ModifiesPC: true,
	}
}

func initJMPInstructions() {
	Instructions[0x6C] = Instruction{
		Name:       "JMP (indirect)",
		Opcode:     0x6C,
		Bytes:      3,
		Cycles:     5,
		Mode:       Indirect,
		Execute:    jmpExecute,
		ModifiesPC: true,
	}

	Instructions[0x4C] = Instruction{
		Name:       "JMP Absolute",
		Opcode:     0x4C,
		Bytes:      3,
		Cycles:     3,
		Mode:       Absolute,
		Execute:    jmpExecute,
		ModifiesPC: true,
	}
}

func jmpExecute(cpu *CPU, addr uint16) {
	cpu.PC = addr
}

func initBEQInstructions() {
	Instructions[0xF0] = Instruction{
		Name:   "BEQ",
		Opcode: 0xF0,
		Bytes:  2,
		Cycles: 2,
		Mode:   Relative,
		Execute: func(cpu *CPU, addr uint16) {
			if cpu.GetFlag(FlagZ) {
				cpu.Cycles += 1

				oldPC := cpu.PC + 1

				if (oldPC & 0xFF00) != (addr & 0xFF00) {
					cpu.Cycles += 1
				}

				cpu.PC = addr
			} else {
				cpu.PC += 2
			}
		},
		ModifiesPC: true,
	}
}
