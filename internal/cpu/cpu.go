package cpu

import (
	"fmt"
)

type CPUBus interface {
	CPURead(addr uint16) byte
	CPUWrite(addr uint16, value byte)
	ShouldTriggerNMI() bool
	AcknowledgeNMI()
	ClockPPU()
	StepPPU()
}

type CPU struct {
	A  byte   // Accumulator
	X  byte   // Index Register X
	Y  byte   // Index Register Y
	SP byte   // Stack Pointer
	PC uint16 // Program Counter
	P  byte   // Status flags

	// Memory [0x10000]byte
	Cycles int

	Bus        CPUBus
	CyclesLeft int
}

func New() *CPU {
	return &CPU{
		// Memory: [0x10000]byte{},
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

func (c *CPU) AttachBus(bus CPUBus) {
	c.Bus = bus
}

func (c *CPU) Reset() {
	c.PC = uint16(c.Bus.CPURead(ResetVector)) | (uint16(c.Bus.CPURead(ResetVector+1)) << 8)
	c.SP = 0xFD
	c.A = 0
	c.X = 0
	c.Y = 0
	c.P = FlagI | FlagU
	c.Cycles = 7
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
	// c.P = value | 0x20
	c.P = (value &^ FlagB) | FlagU
}

func (c *CPU) Push(value byte) {
	c.Bus.CPUWrite(0x0100+uint16(c.SP), value)
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
	return c.Bus.CPURead(0x0100 + uint16(c.SP))
}

func (c *CPU) Pull16() uint16 {
	lo := uint16(c.Pull())
	hi := uint16(c.Pull())
	return (hi << 8) | lo
}

func (c *CPU) Step() {
	if c.Bus.ShouldTriggerNMI() {
		c.TriggerNMI()
		c.Bus.AcknowledgeNMI()
	}

	c.Bus.ClockPPU()

	c.Execute()
}

func (cpu *CPU) Clock() {
	// PPU тикает 3 раза за каждый такт CPU
	cpu.Bus.StepPPU()
	cpu.Bus.StepPPU()
	cpu.Bus.StepPPU()

	if cpu.CyclesLeft == 0 {
		if cpu.Bus.ShouldTriggerNMI() {
			cpu.TriggerNMI()
			cpu.Bus.AcknowledgeNMI()
		}
		opcode := cpu.Bus.CPURead(cpu.PC)
		inst, ok := Instructions[opcode]
		if !ok {
			panic(fmt.Sprintf("Unknown opcode: %02X at %04X", opcode, cpu.PC))
		}
		addr, pageCrossed := inst.GetAddress(cpu)
		cpu.CyclesLeft = inst.Cycles
		inst.Execute(cpu, addr, pageCrossed)
		if !inst.ModifiesPC {
			cpu.PC += uint16(inst.Bytes)
		}
	}
	cpu.CyclesLeft--
	cpu.Cycles++
}

func (c *CPU) ClockOnce() {

	// for i := 0; i < 3; i++ {
	// c.Bus.ClockPPU()
	// c.Cycles++
	// }

	c.Clock()

}

func (c *CPU) Execute() {
	opcode := c.Bus.CPURead(c.PC)
	inst, ok := Instructions[opcode]
	if !ok {
		str := fmt.Sprintf("Unknown opcode: 0x%02X", opcode)
		panic(str)
	}

	addr, pageCrossed := inst.GetAddress(c)

	inst.Execute(c, addr, pageCrossed)
	c.Cycles += inst.Cycles

	if !inst.ModifiesPC {
		c.PC += uint16(inst.Bytes)
	}
}

func (cpuInstance *CPU) Trace(ppuScanline, ppuCycle int) string {
	opcode := cpuInstance.Bus.CPURead(cpuInstance.PC)
	inst, ok := Instructions[opcode]

	if !ok {
		return fmt.Sprintf("Unknown opcode: 0x%02X", opcode)
	}

	// Получить дизассемблированную строку инструкции (например, "JMP $C5F5")
	disasm := inst.Disassemble(cpuInstance, cpuInstance.PC)

	return fmt.Sprintf(
		"%04X  %-28sA:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d",
		cpuInstance.PC, disasm, cpuInstance.A, cpuInstance.X, cpuInstance.Y, cpuInstance.P, cpuInstance.SP, ppuScanline, ppuCycle, cpuInstance.Cycles,
	)
}

func (c *CPU) Read16(addr uint16) uint16 {
	low := c.Bus.CPURead(addr)
	high := c.Bus.CPURead(addr + 1)
	return uint16(low) | (uint16(high) << 8)
}

func (c *CPU) TriggerNMI() {
	fmt.Println("[CPU] >>> TriggerNMI выполнен!")
	c.Push16(c.PC)
	c.Push(c.P | 0x20)
	c.setInterruptDisable(true)
	c.PC = c.Read16(0xFFFA) // NMI vector
}

func (c *CPU) fetchImediate() uint16 {
	return uint16(c.PC + 1)
}

func (c *CPU) fetchZeroPage() uint16 {
	operand := c.Bus.CPURead(c.PC + 1)
	return uint16(operand)
}

func (cpu *CPU) fetchZeroPageX() uint16 {
	base := cpu.Bus.CPURead(cpu.PC + 1)
	addr := (uint16(base) + uint16(cpu.X)) & 0x00FF
	return addr
}

func (cpu *CPU) fetchZeroPageY() uint16 {
	base := cpu.Bus.CPURead(cpu.PC + 1)
	addr := (uint16(base) + uint16(cpu.Y)) & 0x00FF
	return addr
}

func (cpu *CPU) fetchAbsolute() uint16 {
	lo := cpu.Bus.CPURead(cpu.PC + 1)
	hi := cpu.Bus.CPURead(cpu.PC + 2)
	return uint16(lo) | (uint16(hi) << 8)
}

func (cpu *CPU) fetchAbsoluteX() (uint16, bool) {
	lo := cpu.Bus.CPURead(cpu.PC + 1)
	hi := cpu.Bus.CPURead(cpu.PC + 2)
	baseAddr := uint16(lo) | (uint16(hi) << 8)
	effectiveAddr := baseAddr + uint16(cpu.X)
	pageCrossed := (baseAddr & 0xFF00) != (effectiveAddr & 0xFF00)
	return effectiveAddr, pageCrossed
}

func (cpu *CPU) fetchAbsoluteY() (uint16, bool) {
	lo := cpu.Bus.CPURead(cpu.PC + 1)
	hi := cpu.Bus.CPURead(cpu.PC + 2)
	baseAddr := uint16(lo) | (uint16(hi) << 8)
	effectiveAddr := baseAddr + uint16(cpu.Y)
	pageCrossed := (baseAddr & 0xFF00) != (effectiveAddr & 0xFF00)
	return effectiveAddr, pageCrossed
}

func (cpu *CPU) fetchIndirectX() (uint16, bool) {
	zp := cpu.Bus.CPURead(cpu.PC + 1)
	ind := (zp + cpu.X) & 0xFF
	lo := cpu.Bus.CPURead(uint16(ind))
	hi := cpu.Bus.CPURead(uint16((ind + 1) & 0xFF))
	addr := uint16(lo) | (uint16(hi) << 8)
	return addr, false // pageCrossed никогда не нужен
}

func (cpu *CPU) fetchIndirectY() (uint16, bool) {
	base := cpu.Bus.CPURead(cpu.PC + 1)
	lo := cpu.Bus.CPURead(uint16(base))
	hi := cpu.Bus.CPURead(uint16(base+1) & 0x00FF)
	baseAddr := uint16(lo) | (uint16(hi) << 8)
	effectiveAddr := baseAddr + uint16(cpu.Y)
	pageCrossed := (baseAddr & 0xFF00) != (effectiveAddr & 0xFF00)
	return effectiveAddr, pageCrossed
}

func (cpu *CPU) fetchIndirect() uint16 {
	lo := cpu.Bus.CPURead(cpu.PC + 1)
	hi := cpu.Bus.CPURead(cpu.PC + 2)
	addr := uint16(lo) | (uint16(hi) << 8)
	// Специальная проверка на баг
	var indirectAddr uint16
	if lo == 0xFF {
		indirectAddr = uint16(cpu.Bus.CPURead(addr)) | (uint16(cpu.Bus.CPURead(addr&0xFF00)) << 8)
	} else {
		indirectAddr = uint16(cpu.Bus.CPURead(addr)) | (uint16(cpu.Bus.CPURead(addr+1)) << 8)
	}
	return indirectAddr
}

func (cpu *CPU) fetchImplied() uint16 {
	return 0
}

func (cpu *CPU) fetchAccumulator() uint16 {
	return 0
}

func (cpu *CPU) fetchRelative() (uint16, bool) {
	offset := int8(cpu.Bus.CPURead(cpu.PC + 1))
	target := uint16(int32(cpu.PC+2) + int32(offset))
	pageCrossed := ((cpu.PC + 2) & 0xFF00) != (target & 0xFF00)
	return target, pageCrossed
}

func (cpu *CPU) setZN(value byte) {
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, value&0x80 != 0)
}

func (cpu *CPU) setInterruptDisable(val bool) {
	cpu.SetFlag(FlagI, val)
}
