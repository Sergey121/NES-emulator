package cpu

import (
	"testing"
)

type OpcodeTest struct {
	name   string
	init   func(*CPU)
	assert func(*CPU)
}

func runTests(t *testing.T, tests []OpcodeTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpuInstance := New()
			test.init(cpuInstance)

			// Выполняем инструкцию
			cpuInstance.Execute()

			// Проверяем результат
			test.assert(cpuInstance)
		})
	}
}

func TestLDAOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: A9",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Пишем LDA #$42 по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xA9 // LDA Immediate
				cpuInstance.Memory[0x8001] = 0xC0
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0xC0 {
					t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: A5",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем значение в Zero Page
				cpuInstance.Memory[0x0042] = 0xC0

				// Пишем LDA $42 по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xA5 // LDA Zero Page
				cpuInstance.Memory[0x8001] = 0x42
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0xC0 {
					t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: B5",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем значение в Zero Page
				cpuInstance.Memory[0x10] = 0xC0

				// Устанавливаем X-регистр
				cpuInstance.X = 0x20

				// Пишем LDA $F0,X по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xB5 // LDA ZeroPage,X
				cpuInstance.Memory[0x8001] = 0xF0
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0xC0 {
					t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: B5 Wraparound",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем X так, чтобы было переполнение
				cpuInstance.X = 0x20

				// Пишем LDA $F0,X по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xB5 // LDA ZeroPage,X
				cpuInstance.Memory[0x8001] = 0xF0

				// Результат адресации: (0xF0 + 0x20) & 0xFF = 0x10
				cpuInstance.Memory[0x10] = 0xAB // Значение, которое должно загрузиться в A
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0xAB {
					t.Errorf("Expected A = 0xAB, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if (cpuInstance.A&0x80) == 0 && cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: AD",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Пишем LDA $1234 по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xAD // LDA Absolute
				cpuInstance.Memory[0x8001] = 0x34
				cpuInstance.Memory[0x8002] = 0x12

				// В памяти по адресу 0x1234 кладем значение, которое LDA должна загрузить
				cpuInstance.Memory[0x1234] = 0x99
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0x99 {
					t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: BD",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем X-регистр
				cpuInstance.X = 0x01

				// Пишем LDA $1234,X по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xBD // LDA Absolute,X
				cpuInstance.Memory[0x8001] = 0x34
				cpuInstance.Memory[0x8002] = 0x12

				// В памяти по адресу 0x1235 (0x1234 + X) кладем значение, которое LDA должна загрузить
				cpuInstance.Memory[0x1235] = 0x99
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0x99 {
					t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: B9",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Y-регистр
				cpuInstance.Y = 0x01

				// Пишем LDA $1234,Y по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xB9 // LDA Absolute,Y
				cpuInstance.Memory[0x8001] = 0x34
				cpuInstance.Memory[0x8002] = 0x12

				// В памяти по адресу 0x1235 (0x1234 + Y) кладем значение, которое LDA должна загрузить
				cpuInstance.Memory[0x1235] = 0x99
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0x99 {
					t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: A1",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем X-регистр
				cpuInstance.X = 0x04

				// Пишем LDA (0x10,X) — фактически читаем адрес по адресу (0x14)
				cpuInstance.Memory[0x8000] = 0xA1 // LDA (Indirect,X)
				cpuInstance.Memory[0x8001] = 0x10 // Operand: 0x10

				// В Zero Page по адресу 0x14 и 0x15 кладем адрес 0x1234
				cpuInstance.Memory[0x14] = 0x34 // low byte
				cpuInstance.Memory[0x15] = 0x12 // high byte

				// По адресу 0x1234 кладем значение, которое LDA должна загрузить
				cpuInstance.Memory[0x1234] = 0x99
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0x99 {
					t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
				}
				if cpuInstance.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be cleared")
				}
				if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: B1",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Вставляем инструкцию LDA ($10),Y по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0xB1 // opcode
				cpuInstance.Memory[0x8001] = 0x10 // операнд
				cpuInstance.Y = 0x01              // индекс

				// В Zero Page по адресу 0x10 и 0x11 лежит указатель на 0x1234
				cpuInstance.Memory[0x0010] = 0x34
				cpuInstance.Memory[0x0011] = 0x12

				// В память по адресу 0x1235 (0x1234 + Y) кладем значение
				cpuInstance.Memory[0x1235] = 0x99
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.A != 0x99 {
					t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestLDAExecute(t *testing.T) {
	tests := []struct {
		name          string
		memoryValue   byte
		expectedA     byte
		expectedZFlag bool
		expectedNFlag bool
	}{
		{
			name:          "Load positive value",
			memoryValue:   0x42,
			expectedA:     0x42,
			expectedZFlag: false,
			expectedNFlag: false,
		},
		{
			name:          "Load zero value",
			memoryValue:   0x00,
			expectedA:     0x00,
			expectedZFlag: true,
			expectedNFlag: false,
		},
		{
			name:          "Load negative value",
			memoryValue:   0x80,
			expectedA:     0x80,
			expectedZFlag: false,
			expectedNFlag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpuInstance := New()

			// Set the memory value at a test address
			testAddr := uint16(0x1234)
			cpuInstance.Memory[testAddr] = tt.memoryValue

			// Execute ldaExecute
			ldaExecute(cpuInstance, testAddr)

			// Validate the accumulator value
			if cpuInstance.A != tt.expectedA {
				t.Errorf("Expected A = 0x%02X, got 0x%02X", tt.expectedA, cpuInstance.A)
			}

			// Validate the Zero flag
			if cpuInstance.GetFlag(FlagZ) != tt.expectedZFlag {
				t.Errorf("Expected Zero flag = %v, got %v", tt.expectedZFlag, cpuInstance.GetFlag(FlagZ))
			}

			// Validate the Negative flag
			if cpuInstance.GetFlag(FlagN) != tt.expectedNFlag {
				t.Errorf("Expected Negative flag = %v, got %v", tt.expectedNFlag, cpuInstance.GetFlag(FlagN))
			}
		})
	}
}

func TestSTAOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 85",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x99

				// Пишем STA $20 по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x85
				cpuInstance.Memory[0x8001] = 0x20

			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x0020] != 0x99 {
					t.Errorf("Expected memory at 0x0020 to be 0x99, got 0x%02X", cpuInstance.Memory[0x0020])
				}
			},
		},
		{
			name: "Opcode: 95",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x99

				// Устанавливаем X-регистр
				cpuInstance.X = 0x10

				// Пишем STA $20,X по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x95
				cpuInstance.Memory[0x8001] = 0x20
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x0030] != 0x99 {
					t.Errorf("Expected memory at 0x0030 to be 0x99, got 0x%02X", cpuInstance.Memory[0x0030])
				}
			},
		},
		{
			name: "Opcode: 8D",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x42

				// Пишем STA $1234 по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x8D
				cpuInstance.Memory[0x8001] = 0x34
				cpuInstance.Memory[0x8002] = 0x12
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x1234] != 0x42 {
					t.Errorf("Expected memory at 0x1234 to be 0x42, got 0x%02X", cpuInstance.Memory[0x1234])
				}
			},
		},
		{
			name: "Opcode: 9D",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x42

				// Устанавливаем X-регистр
				cpuInstance.X = 0x01

				// Пишем STA $1234,X по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x9D
				cpuInstance.Memory[0x8001] = 0x34
				cpuInstance.Memory[0x8002] = 0x12
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x1235] != 0x42 {
					t.Errorf("Expected memory at 0x1235 to be 0x42, got 0x%02X", cpuInstance.Memory[0x1235])
				}
			},
		},
		{
			name: "Opcode: 99",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x42

				// Устанавливаем Y-регистр
				cpuInstance.Y = 0x01

				// Пишем STA $1234,Y по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x99
				cpuInstance.Memory[0x8001] = 0x34
				cpuInstance.Memory[0x8002] = 0x12
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x1235] != 0x42 {
					t.Errorf("Expected memory at 0x1235 to be 0x42, got 0x%02X", cpuInstance.Memory[0x1235])
				}
			},
		},
		{
			name: "Opcode: 81",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x42

				// Устанавливаем X-регистр
				cpuInstance.X = 0x04

				// Пишем STA ($10,X) по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x81
				cpuInstance.Memory[0x8001] = 0x10

				// В Zero Page по адресу 0x14 и 0x15 кладем адрес 0x1234
				cpuInstance.Memory[0x14] = 0x34 // low byte
				cpuInstance.Memory[0x15] = 0x12 // high byte
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x1234] != 0x42 {
					t.Errorf("Expected memory at 0x1234 to be 0x42, got 0x%02X", cpuInstance.Memory[0x1234])
				}
			},
		},
		{
			name: "Opcode: 91",
			init: func(cpuInstance *CPU) {
				// Устанавливаем Reset Vector на 0x8000
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Устанавливаем Accumulator
				cpuInstance.A = 0x42

				// Устанавливаем Y-регистр
				cpuInstance.Y = 0x01

				// Пишем STA ($10),Y по адресу 0x8000
				cpuInstance.Memory[0x8000] = 0x91
				cpuInstance.Memory[0x8001] = 0x10

				// В Zero Page по адресу 0x10 и 0x11 кладем адрес 0x1234
				cpuInstance.Memory[0x0010] = 0x34 // low byte
				cpuInstance.Memory[0x0011] = 0x12 // high byte
			},
			assert: func(cpuInstance *CPU) {
				if cpuInstance.Memory[0x1235] != 0x42 {
					t.Errorf("Expected memory at 0x1235 to be 0x42, got 0x%02X", cpuInstance.Memory[0x1235])
				}
			},
		},
	}

	runTests(t, tests)
}

func TestSTAExecute(t *testing.T) {
	cpuInstance := New()
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	cpuInstance.A = 0x42
	cpuInstance.Memory[0x8000] = 0x85 // Opcode for STA Zero Page
	cpuInstance.Memory[0x8001] = 0x20 // Operand
	cpuInstance.Execute()
	if cpuInstance.Memory[0x20] != 0x42 {
		t.Errorf("Expected memory at 0x0020 to be 0x42, got 0x%02X", cpuInstance.Memory[0x20])
	}

	cpuInstance.Reset()

	cpuInstance.A = 0x99
	cpuInstance.Memory[0x8000] = 0x8D // Opcode for STA Absolute
	cpuInstance.Memory[0x8001] = 0x34
	cpuInstance.Memory[0x8002] = 0x12 // Operand
	cpuInstance.Execute()
	if cpuInstance.Memory[0x1234] != 0x99 {
		t.Errorf("Expected memory at 0x1234 to be 0x99, got 0x%02X", cpuInstance.Memory[0x1234])
	}
}

func TestLDXOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: A2 (LDX Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x8000] = 0xA2
				cpu.Memory[0x8001] = 0x55
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0x55 {
					t.Errorf("Expected X = 0x55, got 0x%02X", cpu.X)
				}
			},
		},
		{
			name: "Opcode: A6 (LDX ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x10] = 0x66
				cpu.Memory[0x8000] = 0xA6
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0x66 {
					t.Errorf("Expected X = 0x66, got %02X", cpu.X)
				}
			},
		},
		{
			name: "Opcode: B6 (LDX ZeroPage,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x01
				cpu.Memory[0x11] = 0x77
				cpu.Memory[0x8000] = 0xB6
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0x77 {
					t.Errorf("Expected X = 0x77, got %02X", cpu.X)
				}
			},
		},
		{
			name: "Opcode: AE (LDX Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x1234] = 0x88
				cpu.Memory[0x8000] = 0xAE
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0x88 {
					t.Errorf("Expected X = 0x88, got %02X", cpu.X)
				}
			},
		},
		{
			name: "Opcode: BE (LDX Absolute,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x01
				cpu.Memory[0x1235] = 0x99
				cpu.Memory[0x8000] = 0xBE
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0x99 {
					t.Errorf("Expected X = 0x99, got %02X", cpu.X)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestLDXExecute(t *testing.T) {
	cpuInstance := New()
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Test LDX Immediate
	cpuInstance.Memory[0x8000] = 0xA2 // Opcode for LDX Immediate
	cpuInstance.Memory[0x8001] = 0x55 // Operand
	cpuInstance.Execute()
	if cpuInstance.X != 0x55 {
		t.Errorf("Expected X = 0x55, got 0x%02X", cpuInstance.X)
	}

	cpuInstance.Reset()
	// Test LDX Zero Page
	cpuInstance.Memory[0x10] = 0x66
	cpuInstance.Memory[0x8000] = 0xA6 // Opcode for LDX Zero Page
	cpuInstance.Memory[0x8001] = 0x10 // Operand
	cpuInstance.Execute()
	if cpuInstance.X != 0x66 {
		t.Errorf("Expected X = 0x66, got 0x%02X", cpuInstance.X)
	}
}

func TestLDYOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: A0 (LDY Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x8000] = 0xA0
				cpu.Memory[0x8001] = 0x77
			},
			assert: func(cpu *CPU) {
				if cpu.Y != 0x77 {
					t.Errorf("Expected Y = 0x77, got 0x%02X", cpu.Y)
				}
			},
		},
		{
			name: "Opcode: A4 (LDY ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x10] = 0x88
				cpu.Memory[0x8000] = 0xA4
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.Y != 0x88 {
					t.Errorf("Expected Y = 0x88, got 0x%02X", cpu.Y)
				}
			},
		},
		{
			name: "Opcode: B4 (LDY ZeroPage,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x01
				cpu.Memory[0x11] = 0x99
				cpu.Memory[0x8000] = 0xB4
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.Y != 0x99 {
					t.Errorf("Expected Y = 0x99, got 0x%02X", cpu.Y)
				}
			},
		},
		{
			name: "Opcode: AC (LDY Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x1234] = 0xAB
				cpu.Memory[0x8000] = 0xAC
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU) {
				if cpu.Y != 0xAB {
					t.Errorf("Expected Y = 0xAB, got 0x%02X", cpu.Y)
				}
			},
		},
		{
			name: "Opcode: BC (LDY Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x01
				cpu.Memory[0x1235] = 0xCD
				cpu.Memory[0x8000] = 0xBC
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU) {
				if cpu.Y != 0xCD {
					t.Errorf("Expected Y = 0xCD, got 0x%02X", cpu.Y)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestLDYExecute(t *testing.T) {
	cpuInstance := New()
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Test LDY Immediate
	cpuInstance.Memory[0x8000] = 0xA0 // Opcode for LDY Immediate
	cpuInstance.Memory[0x8001] = 0x77 // Operand
	cpuInstance.Execute()
	if cpuInstance.Y != 0x77 {
		t.Errorf("Expected Y = 0x77, got 0x%02X", cpuInstance.Y)
	}

	cpuInstance.Reset()
	// Test LDY Zero Page
	cpuInstance.Memory[0x10] = 0x88
	cpuInstance.Memory[0x8000] = 0xA4 // Opcode for LDY Zero Page
	cpuInstance.Memory[0x8001] = 0x10 // Operand
	cpuInstance.Execute()
	if cpuInstance.Y != 0x88 {
		t.Errorf("Expected Y = 0x88, got 0x%02X", cpuInstance.Y)
	}
}

func TestSTXOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 86 (STX ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x55
				cpu.Memory[0x8000] = 0x86
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.Memory[0x0010] != 0x55 {
					t.Errorf("Expected memory at 0x0010 to be 0x55, got 0x%02X", cpu.Memory[0x0010])
				}
			},
		},
		{
			name: "Opcode: 96 (STX ZeroPage,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x66
				cpu.Y = 0x04
				cpu.Memory[0x8000] = 0x96
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.Memory[0x0014] != 0x66 {
					t.Errorf("Expected memory at 0x0014 to be 0x66, got 0x%02X", cpu.Memory[0x0014])
				}
			},
		},
		{
			name: "Opcode: 8E (STX Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x77
				cpu.Memory[0x8000] = 0x8E
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU) {
				if cpu.Memory[0x1234] != 0x77 {
					t.Errorf("Expected memory at 0x1234 to be 0x77, got 0x%02X", cpu.Memory[0x1234])
				}
			},
		},
	}

	runTests(t, tests)
}

func TestSTXExecute(t *testing.T) {
	cpuInstance := New()
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Test STX Zero Page
	cpuInstance.X = 0x55
	cpuInstance.Memory[0x8000] = 0x86 // Opcode for STX Zero Page
	cpuInstance.Memory[0x8001] = 0x10 // Operand
	cpuInstance.Execute()
	if cpuInstance.Memory[0x0010] != 0x55 {
		t.Errorf("Expected memory at 0x0010 to be 0x55, got 0x%02X", cpuInstance.Memory[0x0010])
	}

	cpuInstance.Reset()
	// Test STX Absolute
	cpuInstance.X = 0x99
	cpuInstance.Memory[0x8000] = 0x8E // Opcode for STX Absolute
	cpuInstance.Memory[0x8001] = 0x34 // Operand
	cpuInstance.Memory[0x8002] = 0x12 // Operand
	cpuInstance.Execute()
	if cpuInstance.Memory[0x1234] != 0x99 {
		t.Errorf("Expected memory at 0x1234 to be 0x99, got 0x%02X", cpuInstance.Memory[0x1234])
	}
}

func TestSTYOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 84 (STY ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x11
				cpu.Memory[0x8000] = 0x84
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU) {
				if cpu.Memory[0x0020] != 0x11 {
					t.Errorf("Expected memory at 0x0020 to be 0x11, got 0x%02X", cpu.Memory[0x0020])
				}
			},
		},
		{
			name: "Opcode: 94 (STY ZeroPage,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x22
				cpu.X = 0x05
				cpu.Memory[0x8000] = 0x94
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU) {
				if cpu.Memory[0x0015] != 0x22 {
					t.Errorf("Expected memory at 0x0015 to be 0x22, got 0x%02X", cpu.Memory[0x0015])
				}
			},
		},
		{
			name: "Opcode: 8C (STY Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x33
				cpu.Memory[0x8000] = 0x8C
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU) {
				if cpu.Memory[0x1234] != 0x33 {
					t.Errorf("Expected memory at 0x1234 to be 0x33, got 0x%02X", cpu.Memory[0x1234])
				}
			},
		},
	}

	runTests(t, tests)
}

func TestSTYExecute(t *testing.T) {
	cpuInstance := New()
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Test STY Zero Page
	cpuInstance.Y = 0x11
	cpuInstance.Memory[0x8000] = 0x84 // Opcode for STY Zero Page
	cpuInstance.Memory[0x8001] = 0x20 // Operand
	cpuInstance.Execute()
	if cpuInstance.Memory[0x0020] != 0x11 {
		t.Errorf("Expected memory at 0x0020 to be 0x11, got 0x%02X", cpuInstance.Memory[0x0020])
	}

	cpuInstance.Reset()
	// Test STY Absolute
	cpuInstance.Y = 0x99
	cpuInstance.Memory[0x8000] = 0x8C // Opcode for STY Absolute
	cpuInstance.Memory[0x8001] = 0x34 // Operand
	cpuInstance.Memory[0x8002] = 0x12 // Operand
	cpuInstance.Execute()
	if cpuInstance.Memory[0x1234] != 0x99 {
		t.Errorf("Expected memory at 0x1234 to be 0x99, got 0x%02X", cpuInstance.Memory[0x1234])
	}
}

func TestTransferOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: AA - TAX",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x00              // чтобы тестировать установку флага Zero
				cpu.Memory[0x8000] = 0xAA // TAX
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0x00 {
					t.Errorf("TAX failed: expected X = 0x00, got 0x%02X", cpu.X)
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("TAX failed: expected Zero flag to be set")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("TAX failed: expected Negative flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: A8 - TAY",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x80              // чтобы тестировать установку флага Negative
				cpu.Memory[0x8000] = 0xA8 // TAY
			},
			assert: func(cpu *CPU) {
				if cpu.Y != 0x80 {
					t.Errorf("TAY failed: expected Y = 0x80, got 0x%02X", cpu.Y)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("TAY failed: expected Zero flag to be cleared")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("TAY failed: expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: BA - TSX",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SP = 0xFF             // чтобы тестировать установку флага Negative
				cpu.Memory[0x8000] = 0xBA // TSX
			},
			assert: func(cpu *CPU) {
				if cpu.X != 0xFF {
					t.Errorf("TSX failed: expected X = 0xFF, got 0x%02X", cpu.X)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("TSX failed: expected Zero flag to be cleared")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("TSX failed: expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: 8A - TXA",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x00              // проверим Zero флаг
				cpu.Memory[0x8000] = 0x8A // TXA
			},
			assert: func(cpu *CPU) {
				if cpu.A != 0x00 {
					t.Errorf("TXA failed: expected A = 0x00, got 0x%02X", cpu.A)
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("TXA failed: expected Zero flag to be set")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("TXA failed: expected Negative flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: 9A - TXS",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0xFF              // чтобы тестировать установку флага Negative
				cpu.Memory[0x8000] = 0x9A // TXS
			},
			assert: func(cpu *CPU) {
				if cpu.SP != 0xFF {
					t.Errorf("TXS failed: expected SP = 0xFF, got 0x%02X", cpu.SP)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("TXS failed: expected Zero flag to be cleared")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("TXS failed: expected Negative flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: 98 - TYA",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0xFF              // тестируем Negative флаг
				cpu.Memory[0x8000] = 0x98 // TYA
			},
			assert: func(cpu *CPU) {
				if cpu.A != 0xFF {
					t.Errorf("TYA failed: expected A = 0xFF, got 0x%02X", cpu.A)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("TYA failed: expected Zero flag to be cleared")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("TYA failed: expected Negative flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}
