package cpu

import (
	"fmt"
)

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

const ResetVector = 0xFFFC

const (
	FlagC = 1 << 0 // Carry Flag
	FlagZ = 1 << 1 // Zero Flag
	FlagI = 1 << 2 // Interrupt Disable Flag
	FlagD = 1 << 3 // Decimal Mode Flag
	FlagB = 1 << 4 // Break Command Flag
	FlagU = 1 << 5 // Unused Flag
	FlagV = 1 << 6 // Overflow Flag
	FlagN = 1 << 7 // Negative Flag
)

func (c *CPU) Reset() {
	c.PC = uint16(c.Memory[ResetVector]) | (uint16(c.Memory[ResetVector+1]) << 8)
	c.SP = 0xFD
	c.A = 0
	c.X = 0
	c.Y = 0
	c.P = FlagI | FlagU
}

func (c *CPU) SetFlag(flag byte, value bool) {
	if value {
		c.P |= flag
	} else {
		c.P &^= flag
	}
}

func (c *CPU) GetFlag(flag byte) bool {
	return c.P&flag != 0
}

func (c *CPU) Execute() {
	opcode := c.Memory[c.PC]
	inst, ok := Instructions[opcode]
	if !ok {
		str := fmt.Sprintf("Unknown opcode: 0x%02X", opcode)
		panic(str)
	}

	var addr uint16
	switch inst.Mode {
	case Immediate:
		addr = c.fetchImediate()
	case ZeroPage:
		addr = c.fetchZeroPage()
	case ZeroPageX:
		addr = c.fetchZeroPageX()
	case Absolute:
		addr = c.fetchAbsolute()
	default:
		str := fmt.Sprintf("Unknown addressing mode: %d", inst.Mode)
		panic(str)
	}
	inst.Execute(c, addr)
	c.PC += uint16(inst.Bytes)
}

func (c *CPU) fetchImediate() uint16 {
	return uint16(c.PC + 1)
}

func (c *CPU) fetchZeroPage() uint16 {
	operand := c.Memory[c.PC+1]
	return uint16(operand)
}

func (cpu *CPU) fetchZeroPageX() uint16 {
	base := cpu.Memory[cpu.PC+1]
	addr := (uint16(base) + uint16(cpu.X)) & 0x00FF
	return addr
}

func (cpu *CPU) fetchAbsolute() uint16 {
	lo := cpu.Memory[cpu.PC+1]
	hi := cpu.Memory[cpu.PC+2]
	return uint16(lo) | (uint16(hi) << 8)
}
