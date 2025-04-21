package cpu_test

import (
	"testing"

	"github.com/sergey121/nes-emulator/internal/cpu"
)

func TestCPUReset(t *testing.T) {
	cpuInstance := cpu.New()

	// Set up the Reset Vector in memory
	cpuInstance.Memory[cpu.ResetVector] = 0x34
	cpuInstance.Memory[cpu.ResetVector+1] = 0x12

	// Call the Reset method
	cpuInstance.Reset()

	// Verify the Program Counter (PC) is set correctly
	expectedPC := uint16(0x1234)
	if cpuInstance.PC != expectedPC {
		t.Errorf("PC expected to be 0x%04X, got 0x%04X", expectedPC, cpuInstance.PC)
	}

	// Verify the Stack Pointer (SP) is set correctly
	expectedSP := byte(0xFD)
	if cpuInstance.SP != expectedSP {
		t.Errorf("SP expected to be 0x%02X, got 0x%02X", expectedSP, cpuInstance.SP)
	}

	// Verify the Accumulator (A) is reset to 0
	if cpuInstance.A != 0 {
		t.Errorf("A expected to be 0, got %d", cpuInstance.A)
	}

	// Verify the Index Register X is reset to 0
	if cpuInstance.X != 0 {
		t.Errorf("X expected to be 0, got %d", cpuInstance.X)
	}

	// Verify the Index Register Y is reset to 0
	if cpuInstance.Y != 0 {
		t.Errorf("Y expected to be 0, got %d", cpuInstance.Y)
	}

	// Verify the Processor Status (P) is set correctly
	expectedP := cpu.FlagI | cpu.FlagU
	if cpuInstance.P != byte(expectedP) {
		t.Errorf("P expected to be 0x%02X, got 0x%02X", expectedP, cpuInstance.P)
	}
}
func TestSetFlag(t *testing.T) {
	cpuInstance := cpu.New()

	// Test setting a flag to true
	cpuInstance.SetFlag(cpu.FlagC, true)
	if !cpuInstance.GetFlag(cpu.FlagC) {
		t.Errorf("FlagC expected to be set, but it was not")
	}

	// Test setting a flag to false
	cpuInstance.SetFlag(cpu.FlagC, false)
	if cpuInstance.GetFlag(cpu.FlagC) {
		t.Errorf("FlagC expected to be cleared, but it was not")
	}

	// Test setting multiple flags
	cpuInstance.SetFlag(cpu.FlagZ, true)
	cpuInstance.SetFlag(cpu.FlagN, true)
	if !cpuInstance.GetFlag(cpu.FlagZ) {
		t.Errorf("FlagZ expected to be set, but it was not")
	}
	if !cpuInstance.GetFlag(cpu.FlagN) {
		t.Errorf("FlagN expected to be set, but it was not")
	}

	// Test clearing multiple flags
	cpuInstance.SetFlag(cpu.FlagZ, false)
	cpuInstance.SetFlag(cpu.FlagN, false)
	if cpuInstance.GetFlag(cpu.FlagZ) {
		t.Errorf("FlagZ expected to be cleared, but it was not")
	}
	if cpuInstance.GetFlag(cpu.FlagN) {
		t.Errorf("FlagN expected to be cleared, but it was not")
	}
}
