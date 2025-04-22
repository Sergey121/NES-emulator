package cpu

import (
	"testing"
)

func TestCPUReset(t *testing.T) {
	cpuInstance := New()

	// Set up the Reset Vector in memory
	cpuInstance.Memory[ResetVector] = 0x34
	cpuInstance.Memory[ResetVector+1] = 0x12

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
	expectedP := FlagI | FlagU
	if cpuInstance.P != byte(expectedP) {
		t.Errorf("P expected to be 0x%02X, got 0x%02X", expectedP, cpuInstance.P)
	}
}
func TestSetFlag(t *testing.T) {
	cpuInstance := New()

	// Test setting a flag to true
	cpuInstance.SetFlag(FlagC, true)
	if !cpuInstance.GetFlag(FlagC) {
		t.Errorf("FlagC expected to be set, but it was not")
	}

	// Test setting a flag to false
	cpuInstance.SetFlag(FlagC, false)
	if cpuInstance.GetFlag(FlagC) {
		t.Errorf("FlagC expected to be cleared, but it was not")
	}

	// Test setting multiple flags
	cpuInstance.SetFlag(FlagZ, true)
	cpuInstance.SetFlag(FlagN, true)
	if !cpuInstance.GetFlag(FlagZ) {
		t.Errorf("FlagZ expected to be set, but it was not")
	}
	if !cpuInstance.GetFlag(FlagN) {
		t.Errorf("FlagN expected to be set, but it was not")
	}

	// Test clearing multiple flags
	cpuInstance.SetFlag(FlagZ, false)
	cpuInstance.SetFlag(FlagN, false)
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("FlagZ expected to be cleared, but it was not")
	}
	if cpuInstance.GetFlag(FlagN) {
		t.Errorf("FlagN expected to be cleared, but it was not")
	}
}

func TestFetchImediate(t *testing.T) {
	cpuInstance := New()

	// Set the Program Counter (PC) to a known value
	cpuInstance.PC = 0x1000

	// Call fetchImediate
	result := cpuInstance.fetchImediate()

	// Verify the result
	expected := uint16(0x1001)
	if result != expected {
		t.Errorf("fetchImediate expected to return 0x%04X, got 0x%04X", expected, result)
	}
}

func TestFetchZeroPage(t *testing.T) {
	cpuInstance := New()

	// Set the Program Counter (PC) to a known value
	cpuInstance.PC = 0x1000

	// Set a value in memory at the zero page address
	cpuInstance.Memory[0x1001] = 0x42

	address := cpuInstance.fetchZeroPage()
	expectedAddress := uint16(0x0042)

	if address != expectedAddress {
		t.Errorf("fetchZeroPage expected to return 0x%04X, got 0x%04X", expectedAddress, address)
	}
}

func TestFetchZeroPageX(t *testing.T) {
	cpuInstance := New()

	// Set the Program Counter (PC) to a known value
	cpuInstance.PC = 0x1000

	// Set a value in memory at the zero page address
	cpuInstance.Memory[0x1001] = 0x42
	cpuInstance.X = 0x10

	address := cpuInstance.fetchZeroPageX()
	expectedAddress := uint16(0x0052) // 0x42 + 0x10 = 0x52
	if address != expectedAddress {
		t.Errorf("fetchZeroPageX expected to return 0x%04X, got 0x%04X", expectedAddress, address)
	}

	cpuInstance.Memory[0x1001] = 0xFC
	cpuInstance.X = 0xFC // Test wraparound
	address = cpuInstance.fetchZeroPageX()
	expectedAddress = uint16(0x00F8) // 0xFC + 0xFC = 0x1F8 & 0xFF = 0xF8
	if address != expectedAddress {
		t.Errorf("fetchZeroPageX owerflow expected to return 0x%04X, got 0x%04X", expectedAddress, address)
	}
}

func TestFetchZeroPageY(t *testing.T) {
	cpuInstance := New()

	// Set the Program Counter (PC) to a known value
	cpuInstance.PC = 0x1000

	// Set a value in memory at the zero page address
	cpuInstance.Memory[0x1001] = 0x42
	cpuInstance.Y = 0x10

	address := cpuInstance.fetchZeroPageY()
	expectedAddress := uint16(0x0052) // 0x42 + 0x10 = 0x52
	if address != expectedAddress {
		t.Errorf("fetchZeroPageY expected to return 0x%04X, got 0x%04X", expectedAddress, address)
	}

	cpuInstance.Memory[0x1001] = 0xFC
	cpuInstance.Y = 0xFC // Test wraparound
	address = cpuInstance.fetchZeroPageY()
	expectedAddress = uint16(0x00F8) // 0xFC + 0xFC = 0x1F8 & 0xFF = 0xF8
	if address != expectedAddress {
		t.Errorf("fetchZeroPageY owerflow expected to return 0x%04X, got 0x%04X", expectedAddress, address)
	}
}

func TestFetchAbsolute(t *testing.T) {
	cpuInstance := New()

	// Set the Program Counter (PC) to a known value
	cpuInstance.PC = 0x1000

	// Set a value in memory at the absolute address
	cpuInstance.Memory[0x1001] = 0x42
	cpuInstance.Memory[0x1002] = 0x84 // 0x8442

	address := cpuInstance.fetchAbsolute()
	expectedAddress := uint16(0x8442) // 0x42 + (0x84 << 8)
	if address != expectedAddress {
		t.Errorf("fetchAbsolute expected to return 0x%04X, got 0x%04X", expectedAddress, address)
	}
}
