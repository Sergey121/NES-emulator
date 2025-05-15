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
	Cycles int
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

func (c *CPU) SetStatus(value byte) {
	c.P = value | 0x20
}

func (c *CPU) Push(value byte) {
	c.Memory[0x0100+uint16(c.SP)] = value
	c.SP--
}

func (c *CPU) Push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0x00FF)
	c.Push(hi)
	c.Push(lo)
}

func (c *CPU) Pull() byte {
	c.SP++
	return c.Memory[0x0100+uint16(c.SP)]
}

func (c *CPU) Pull16() uint16 {
	lo := uint16(c.Pull())
	hi := uint16(c.Pull())
	return (hi << 8) | lo
}

func (c *CPU) Execute() {
	opcode := c.Memory[c.PC]
	inst, ok := Instructions[opcode]
	if !ok {
		str := fmt.Sprintf("Unknown opcode: 0x%02X", opcode)
		panic(str)
	}

	addr := inst.GetAddress(c)

	inst.Execute(c, addr)
	c.Cycles += inst.Cycles

	if !inst.ModifiesPC {
		c.PC += uint16(inst.Bytes)
	}
}

func (cpu *CPU) Trace() string {
	opcode := cpu.Memory[cpu.PC]
	return fmt.Sprintf(
		"PC: %04X  OPCODE: %02X  A:%02X X:%02X Y:%02X P:%02X SP:%02X",
		cpu.PC, opcode,
		cpu.A, cpu.X, cpu.Y, cpu.P, cpu.SP,
	)
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

func (cpu *CPU) fetchZeroPageY() uint16 {
	base := cpu.Memory[cpu.PC+1]
	addr := (uint16(base) + uint16(cpu.Y)) & 0x00FF
	return addr
}

func (cpu *CPU) fetchAbsolute() uint16 {
	lo := cpu.Memory[cpu.PC+1]
	hi := cpu.Memory[cpu.PC+2]
	return uint16(lo) | (uint16(hi) << 8)
}

func (cpu *CPU) fetchAbsoluteX() uint16 {
	lo := cpu.Memory[cpu.PC+1]
	hi := cpu.Memory[cpu.PC+2]
	addr := uint16(lo) | (uint16(hi) << 8)
	addr += uint16(cpu.X)
	return addr
}

func (cpu *CPU) fetchAbsoluteY() uint16 {
	lo := cpu.Memory[cpu.PC+1]
	hi := cpu.Memory[cpu.PC+2]
	addr := uint16(lo) | (uint16(hi) << 8)
	addr += uint16(cpu.Y)
	return addr
}

func (cpu *CPU) fetchIndirectX() uint16 {
	base := cpu.Memory[cpu.PC+1]
	addr := (uint16(base) + uint16(cpu.X)) & 0x00FF
	lo := cpu.Memory[addr]
	hi := cpu.Memory[(addr+1)&0x00FF]
	return uint16(lo) | (uint16(hi) << 8)
}

func (cpu *CPU) fetchIndirectY() uint16 {
	base := cpu.Memory[cpu.PC+1]
	lo := cpu.Memory[base]
	hi := cpu.Memory[(base+1)&0x00FF]
	addr := uint16(lo) | (uint16(hi) << 8)
	addr += uint16(cpu.Y)
	return addr
}

func (cpu *CPU) fetchIndirect() uint16 {
	lo := cpu.Memory[cpu.PC+1]
	hi := cpu.Memory[cpu.PC+2]
	addr := uint16(lo) | (uint16(hi) << 8)
	// Специальная проверка на баг
	var indirectAddr uint16
	if lo == 0xFF {
		indirectAddr = uint16(cpu.Memory[addr]) | (uint16(cpu.Memory[addr&0xFF00]) << 8)
	} else {
		indirectAddr = uint16(cpu.Memory[addr]) | (uint16(cpu.Memory[addr+1]) << 8)
	}
	return indirectAddr
}

func (cpu *CPU) fetchImplied() uint16 {
	return 0
}

func (cpu *CPU) fetchAccumulator() uint16 {
	return 0
}

func (cpu *CPU) fetchRelative() uint16 {
	offset := int8(cpu.Memory[cpu.PC+1])
	return uint16(int32(cpu.PC+2) + int32(offset))
}

func (cpu *CPU) setZN(value byte) {
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, value&0x80 != 0)
}
