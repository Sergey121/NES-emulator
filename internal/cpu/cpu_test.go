package cpu

// import (
// 	"testing"
// )

// func TestCPUReset(t *testing.T) {
// 	cpuInstance := New()

// 	// Set up the Reset Vector in memory
// 	cpuInstance.Memory[ResetVector] = 0x34
// 	cpuInstance.Memory[ResetVector+1] = 0x12

// 	// Call the Reset method
// 	cpuInstance.Reset()

// 	// Verify the Program Counter (PC) is set correctly
// 	expectedPC := uint16(0x1234)
// 	if cpuInstance.PC != expectedPC {
// 		t.Errorf("PC expected to be 0x%04X, got 0x%04X", expectedPC, cpuInstance.PC)
// 	}

// 	// Verify the Stack Pointer (SP) is set correctly
// 	expectedSP := byte(0xFD)
// 	if cpuInstance.SP != expectedSP {
// 		t.Errorf("SP expected to be 0x%02X, got 0x%02X", expectedSP, cpuInstance.SP)
// 	}

// 	// Verify the Accumulator (A) is reset to 0
// 	if cpuInstance.A != 0 {
// 		t.Errorf("A expected to be 0, got %d", cpuInstance.A)
// 	}

// 	// Verify the Index Register X is reset to 0
// 	if cpuInstance.X != 0 {
// 		t.Errorf("X expected to be 0, got %d", cpuInstance.X)
// 	}

// 	// Verify the Index Register Y is reset to 0
// 	if cpuInstance.Y != 0 {
// 		t.Errorf("Y expected to be 0, got %d", cpuInstance.Y)
// 	}

// 	// Verify the Processor Status (P) is set correctly
// 	expectedP := FlagI | FlagU
// 	if cpuInstance.P != byte(expectedP) {
// 		t.Errorf("P expected to be 0x%02X, got 0x%02X", expectedP, cpuInstance.P)
// 	}
// }
// func TestSetFlag(t *testing.T) {
// 	cpuInstance := New()

// 	// Test setting a flag to true
// 	cpuInstance.SetFlag(FlagC, true)
// 	if !cpuInstance.GetFlag(FlagC) {
// 		t.Errorf("FlagC expected to be set, but it was not")
// 	}

// 	// Test setting a flag to false
// 	cpuInstance.SetFlag(FlagC, false)
// 	if cpuInstance.GetFlag(FlagC) {
// 		t.Errorf("FlagC expected to be cleared, but it was not")
// 	}

// 	// Test setting multiple flags
// 	cpuInstance.SetFlag(FlagZ, true)
// 	cpuInstance.SetFlag(FlagN, true)
// 	if !cpuInstance.GetFlag(FlagZ) {
// 		t.Errorf("FlagZ expected to be set, but it was not")
// 	}
// 	if !cpuInstance.GetFlag(FlagN) {
// 		t.Errorf("FlagN expected to be set, but it was not")
// 	}

// 	// Test clearing multiple flags
// 	cpuInstance.SetFlag(FlagZ, false)
// 	cpuInstance.SetFlag(FlagN, false)
// 	if cpuInstance.GetFlag(FlagZ) {
// 		t.Errorf("FlagZ expected to be cleared, but it was not")
// 	}
// 	if cpuInstance.GetFlag(FlagN) {
// 		t.Errorf("FlagN expected to be cleared, but it was not")
// 	}
// }

// func TestFetchImediate(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Call fetchImediate
// 	result := cpuInstance.fetchImediate()

// 	// Verify the result
// 	expected := uint16(0x1001)
// 	if result != expected {
// 		t.Errorf("fetchImediate expected to return 0x%04X, got 0x%04X", expected, result)
// 	}
// }

// func TestFetchZeroPage(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Set a value in memory at the zero page address
// 	cpuInstance.Memory[0x1001] = 0x42

// 	address := cpuInstance.fetchZeroPage()
// 	expectedAddress := uint16(0x0042)

// 	if address != expectedAddress {
// 		t.Errorf("fetchZeroPage expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}
// }

// func TestFetchZeroPageX(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Set a value in memory at the zero page address
// 	cpuInstance.Memory[0x1001] = 0x42
// 	cpuInstance.X = 0x10

// 	address := cpuInstance.fetchZeroPageX()
// 	expectedAddress := uint16(0x0052) // 0x42 + 0x10 = 0x52
// 	if address != expectedAddress {
// 		t.Errorf("fetchZeroPageX expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}

// 	cpuInstance.Memory[0x1001] = 0xFC
// 	cpuInstance.X = 0xFC // Test wraparound
// 	address = cpuInstance.fetchZeroPageX()
// 	expectedAddress = uint16(0x00F8) // 0xFC + 0xFC = 0x1F8 & 0xFF = 0xF8
// 	if address != expectedAddress {
// 		t.Errorf("fetchZeroPageX owerflow expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}
// }

// func TestFetchZeroPageY(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Set a value in memory at the zero page address
// 	cpuInstance.Memory[0x1001] = 0x42
// 	cpuInstance.Y = 0x10

// 	address := cpuInstance.fetchZeroPageY()
// 	expectedAddress := uint16(0x0052) // 0x42 + 0x10 = 0x52
// 	if address != expectedAddress {
// 		t.Errorf("fetchZeroPageY expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}

// 	cpuInstance.Memory[0x1001] = 0xFC
// 	cpuInstance.Y = 0xFC // Test wraparound
// 	address = cpuInstance.fetchZeroPageY()
// 	expectedAddress = uint16(0x00F8) // 0xFC + 0xFC = 0x1F8 & 0xFF = 0xF8
// 	if address != expectedAddress {
// 		t.Errorf("fetchZeroPageY owerflow expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}
// }

// func TestFetchAbsolute(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Set a value in memory at the absolute address
// 	cpuInstance.Memory[0x1001] = 0x42
// 	cpuInstance.Memory[0x1002] = 0x84 // 0x8442

// 	address := cpuInstance.fetchAbsolute()
// 	expectedAddress := uint16(0x8442) // 0x42 + (0x84 << 8)
// 	if address != expectedAddress {
// 		t.Errorf("fetchAbsolute expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}
// }
// func TestFetchAbsoluteX(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Set a value in memory at the absolute address
// 	cpuInstance.Memory[0x1001] = 0x34
// 	cpuInstance.Memory[0x1002] = 0x12 // 0x1234
// 	cpuInstance.X = 0x10

// 	address := cpuInstance.fetchAbsoluteX()
// 	expectedAddress := uint16(0x1244) // 0x1234 + 0x10
// 	if address != expectedAddress {
// 		t.Errorf("fetchAbsoluteX expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}

// 	// Test with wraparound
// 	cpuInstance.Memory[0x1001] = 0xFF
// 	cpuInstance.Memory[0x1002] = 0xFF // 0xFFFF
// 	cpuInstance.X = 0x01

// 	address = cpuInstance.fetchAbsoluteX()
// 	expectedAddress = uint16(0x0000) // 0xFFFF + 0x01 = 0x10000 & 0xFFFF = 0x0000
// 	if address != expectedAddress {
// 		t.Errorf("fetchAbsoluteX wraparound expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}
// }
// func TestFetchAbsoluteY(t *testing.T) {
// 	cpuInstance := New()

// 	// Set the Program Counter (PC) to a known value
// 	cpuInstance.PC = 0x1000

// 	// Set a value in memory at the absolute address
// 	cpuInstance.Memory[0x1001] = 0x34
// 	cpuInstance.Memory[0x1002] = 0x12 // 0x1234
// 	cpuInstance.Y = 0x10

// 	address := cpuInstance.fetchAbsoluteY()
// 	expectedAddress := uint16(0x1244) // 0x1234 + 0x10
// 	if address != expectedAddress {
// 		t.Errorf("fetchAbsoluteY expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}

// 	// Test with wraparound
// 	cpuInstance.Memory[0x1001] = 0xFF
// 	cpuInstance.Memory[0x1002] = 0xFF // 0xFFFF
// 	cpuInstance.Y = 0x01

// 	address = cpuInstance.fetchAbsoluteY()
// 	expectedAddress = uint16(0x0000) // 0xFFFF + 0x01 = 0x10000 & 0xFFFF = 0x0000
// 	if address != expectedAddress {
// 		t.Errorf("fetchAbsoluteY wraparound expected to return 0x%04X, got 0x%04X", expectedAddress, address)
// 	}
// }

// func TestFetchIndirectY(t *testing.T) {
// 	cpuInstance := New()

// 	// Настройка: PC указывает на инструкцию
// 	cpuInstance.PC = 0x0200
// 	cpuInstance.Y = 0x05

// 	// В память по адресу PC + 1 кладем указатель на Zero Page: 0x10
// 	cpuInstance.Memory[0x0201] = 0x10

// 	// В Zero Page по адресу 0x10 и 0x11 лежит базовый адрес: 0x34 + 0x12<<8 = 0x1234
// 	cpuInstance.Memory[0x0010] = 0x34
// 	cpuInstance.Memory[0x0011] = 0x12

// 	// Ожидаемый адрес: 0x1234 + 0x05 = 0x1239
// 	addr := cpuInstance.fetchIndirectY()
// 	expected := uint16(0x1239)

// 	if addr != expected {
// 		t.Errorf("Expected fetchIndirectY to return 0x%04X, got 0x%04X", expected, addr)
// 	}

// 	// Тест wraparound: указатель в Zero Page — 0xFF, следующий байт — 0x00
// 	cpuInstance.Memory[0x0201] = 0xFF
// 	cpuInstance.Memory[0x00FF] = 0x78
// 	cpuInstance.Memory[0x0000] = 0x56 // wraparound!

// 	cpuInstance.Y = 0x01
// 	addr = cpuInstance.fetchIndirectY()
// 	expected = (uint16(0x56)<<8 | 0x78) + 0x01

// 	if addr != expected {
// 		t.Errorf("Expected fetchIndirectY with wraparound to return 0x%04X, got 0x%04X", expected, addr)
// 	}
// }
