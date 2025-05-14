package cpu

import (
	"testing"
)

type OpcodeTest struct {
	name   string
	init   func(*CPU)
	assert func(*CPU, *testing.T)
}

func runTests(t *testing.T, tests []OpcodeTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpuInstance := New()
			test.init(cpuInstance)

			// Выполняем инструкцию
			cpuInstance.Execute()

			// Проверяем результат
			test.assert(cpuInstance, t)
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpuInstance *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
			assert: func(cpu *CPU, t *testing.T) {
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
		{
			name: "Opcode: CA (DEX) regular decrement",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x02
				cpu.Memory[0x8000] = 0xCA // DEX
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.X != 0x01 {
					t.Errorf("Expected X = 0x01, got 0x%02X", cpu.X)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Zero flag should be clear")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Negative flag should be clear")
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: CA (DEX) result zero",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x01
				cpu.Memory[0x8000] = 0xCA // DEX
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.X != 0x00 {
					t.Errorf("Expected X = 0x00, got 0x%02X", cpu.X)
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Zero flag should be set")
				}
			},
		},
		{
			name: "Opcode: CA (DEX) result negative",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x00
				cpu.Memory[0x8000] = 0xCA // DEX
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.X != 0xFF {
					t.Errorf("Expected X = 0xFF, got 0x%02X", cpu.X)
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Negative flag should be set")
				}
			},
		},
		{
			name: "Opcode: E8 (INX) result zero",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0xFF
				cpu.Memory[0x8000] = 0xE8 // INX
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.X != 0x00 {
					t.Errorf("Expected X = 0x00, got 0x%02X", cpu.X)
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Zero flag should be set")
				}
			},
		},
		{
			name: "Opcode: 88 (DEY) regular decrement",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x05
				cpu.Memory[0x8000] = 0x88 // DEY
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Y != 0x04 {
					t.Errorf("Expected Y = 0x04, got %02X", cpu.Y)
				}
			},
		},
		{
			name: "Opcode: C8 (INY) result negative",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x7F
				cpu.Memory[0x8000] = 0xC8 // INY
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Y != 0x80 {
					t.Errorf("Expected Y = 0x80, got 0x%02X", cpu.Y)
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Negative flag should be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestFlagOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 18 (CLC)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Ставим Carry флаг заранее
				cpuInstance.SetFlag(FlagC, true)

				cpuInstance.Memory[0x8000] = 0x18 // CLC
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if cpuInstance.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: 38 (SEC)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				// Очищаем Carry флаг заранее
				cpuInstance.SetFlag(FlagC, false)

				cpuInstance.Memory[0x8000] = 0x38 // SEC
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if !cpuInstance.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
			},
		},
		{
			name: "Opcode: D8 (CLD)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				cpuInstance.SetFlag(FlagD, true)

				cpuInstance.Memory[0x8000] = 0xD8 // CLD
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if cpuInstance.GetFlag(FlagD) {
					t.Errorf("Expected Decimal Mode flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: F8 (SED)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				cpuInstance.SetFlag(FlagD, false)

				cpuInstance.Memory[0x8000] = 0xF8 // SED
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if !cpuInstance.GetFlag(FlagD) {
					t.Errorf("Expected Decimal Mode flag to be set")
				}
			},
		},
		{
			name: "Opcode: 58 (CLI)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				cpuInstance.SetFlag(FlagI, true)

				cpuInstance.Memory[0x8000] = 0x58 // CLI
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if cpuInstance.GetFlag(FlagI) {
					t.Errorf("Expected Interrupt Disable flag to be cleared")
				}
			},
		},
		{
			name: "Opcode: 78 (SEI)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				cpuInstance.SetFlag(FlagI, false)

				cpuInstance.Memory[0x8000] = 0x78 // SEI
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if !cpuInstance.GetFlag(FlagI) {
					t.Errorf("Expected Interrupt Disable flag to be set")
				}
			},
		},
		{
			name: "Opcode: B8 (CLV)",
			init: func(cpuInstance *CPU) {
				cpuInstance.Memory[ResetVector] = 0x00
				cpuInstance.Memory[ResetVector+1] = 0x80
				cpuInstance.Reset()

				cpuInstance.SetFlag(FlagV, true)

				cpuInstance.Memory[0x8000] = 0xB8 // CLV
			},
			assert: func(cpuInstance *CPU, t *testing.T) {
				if cpuInstance.GetFlag(FlagV) {
					t.Errorf("Expected Overflow flag to be cleared")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestADCOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 69 (ADC Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Memory[0x8000] = 0x69
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 65 (ADC ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Memory[0x20] = 0x20
				cpu.Memory[0x8000] = 0x65
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 75 (ADC ZeroPage,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.X = 0x01
				cpu.Memory[0x21] = 0x20
				cpu.Memory[0x8000] = 0x75
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 6D (ADC Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Memory[0x1234] = 0x20
				cpu.Memory[0x8000] = 0x6D
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 7D (ADC Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.X = 0x01
				cpu.Memory[0x1235] = 0x20
				cpu.Memory[0x8000] = 0x7D
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 79 (ADC Absolute,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Y = 0x01
				cpu.Memory[0x1235] = 0x20
				cpu.Memory[0x8000] = 0x79
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 61 (ADC Indirect,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.X = 0x01
				cpu.Memory[0x0021] = 0x34 // low byte of address
				cpu.Memory[0x0022] = 0x12 // high byte of address

				cpu.Memory[0x1234] = 0x20

				cpu.Memory[0x8000] = 0x61
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 71 (ADC Indirect,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Y = 0x01
				cpu.Memory[0x0020] = 0x34 // low byte of address
				cpu.Memory[0x0021] = 0x12 // high byte of address

				cpu.Memory[0x1235] = 0x20

				cpu.Memory[0x8000] = 0x71
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
	}

	runTests(t, tests)

	tests2 := []OpcodeTest{
		{
			name: "Opcode: 69 (ADC) carry flag set",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x8000] = 0x69 // ADC #$01
				cpu.Memory[0x8001] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x00 {
					t.Errorf("Expected A = 0x00, got 0x%02X", cpu.A)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
		{
			name: "Opcode: 69 (ADC) overflow flag set (positive + positive = negative)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x50
				cpu.Memory[0x8000] = 0x69 // ADC #$50
				cpu.Memory[0x8001] = 0x50
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xA0 {
					t.Errorf("Expected A = 0xA0, got 0x%02X", cpu.A)
				}
				if !cpu.GetFlag(FlagV) {
					t.Errorf("Expected Overflow flag to be set")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: 69 (ADC) with carry-in",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x01
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x8000] = 0x69 // ADC #$01
				cpu.Memory[0x8001] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x03 {
					t.Errorf("Expected A = 0x03, got 0x%02X", cpu.A)
				}
				if cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be clear")
				}
			},
		},
	}

	runTests(t, tests2)
}

func TestADCExecute(t *testing.T) {
	cpuInstance := New()
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Test ADC Immediate
	cpuInstance.A = 0x10
	cpuInstance.Memory[0x8000] = 0x69 // Opcode for ADC Immediate
	cpuInstance.Memory[0x8001] = 0x20 // Operand
	cpuInstance.Execute()
	if cpuInstance.A != 0x30 {
		t.Errorf("Expected A = 0x30, got 0x%02X", cpuInstance.A)
	}

	cpuInstance.Reset()
	// Test ADC Zero Page
	cpuInstance.A = 0x10
	cpuInstance.Memory[0x20] = 0x20
	cpuInstance.Memory[0x8000] = 0x65 // Opcode for ADC Zero Page
	cpuInstance.Memory[0x8001] = 0x20 // Operand
	cpuInstance.Execute()
	if cpuInstance.A != 0x30 {
		t.Errorf("Expected A = 0x30, got 0x%02X", cpuInstance.A)
	}
}

func TestSBCOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: E9 (SBC Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x50
				cpu.SetFlag(FlagC, true)  // Carry установлен
				cpu.Memory[0x8000] = 0xE9 // SBC #imm
				cpu.Memory[0x8001] = 0x10 // вычитаем 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x40 {
					t.Errorf("Expected A = 0x40, got 0x%02X", cpu.A)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: E5 (SBC Zero Page)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x30
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0020] = 0x10
				cpu.Memory[0x8000] = 0xE5
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x20 {
					t.Errorf("Expected A = 0x20, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: F5 (SBC Zero Page,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x20
				cpu.X = 0x05
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0025] = 0x10
				cpu.Memory[0x8000] = 0xF5
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x10 {
					t.Errorf("Expected A = 0x10, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: ED (SBC Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x30
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x1234] = 0x10
				cpu.Memory[0x8000] = 0xED
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x20 {
					t.Errorf("Expected A = 0x20, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: FD (SBC Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x20
				cpu.X = 0x01
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x1235] = 0x10
				cpu.Memory[0x8000] = 0xFD
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x10 {
					t.Errorf("Expected A = 0x10, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: F9 (SBC Absolute,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x20
				cpu.Y = 0x01
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x1235] = 0x10
				cpu.Memory[0x8000] = 0xF9
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x10 {
					t.Errorf("Expected A = 0x10, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: E1 (SBC (Indirect,X))",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x40
				cpu.X = 0x04
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0024] = 0x34
				cpu.Memory[0x0025] = 0x12
				cpu.Memory[0x1234] = 0x10
				cpu.Memory[0x8000] = 0xE1
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: F1 (SBC (Indirect),Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x40
				cpu.Y = 0x01
				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0020] = 0x34
				cpu.Memory[0x0021] = 0x12
				cpu.Memory[0x1235] = 0x10
				cpu.Memory[0x8000] = 0xF1
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x30 {
					t.Errorf("Expected A = 0x30, got 0x%02X", cpu.A)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestANDOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 29 (AND Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x8000] = 0x29 // AND #imm
				cpu.Memory[0x8001] = 0xCC // AND with 0xCC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 25 (AND Zero Page)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x0025] = 0xCC
				cpu.Memory[0x8000] = 0x25 // AND Zero Page
				cpu.Memory[0x8001] = 0x25 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 35 (AND Zero Page,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.X = 0x2
				cpu.Memory[0x0027] = 0xCC // Address is (25 + X)
				cpu.Memory[0x8000] = 0x35 // AND Zero Page,X
				cpu.Memory[0x8001] = 0x25 // Operand
				cpu.Memory[0x8002] = 0x00
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 2D (AND Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x1234] = 0xCC
				cpu.Memory[0x8000] = 0x2D // AND Absolute
				cpu.Memory[0x8001] = 0x34 // Operand
				cpu.Memory[0x8002] = 0x12 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 3D (AND Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.X = 1
				cpu.Memory[0x1235] = 0xCC // Address is (1234 + X)
				cpu.Memory[0x8000] = 0x3D // AND Absolute,X
				cpu.Memory[0x8001] = 0x34 // Operand
				cpu.Memory[0x8002] = 0x12 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},

		{
			name: "Opcode: 39 (AND Absolute,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Y = 1
				cpu.Memory[0x1235] = 0xCC // Address is (1234 + Y)
				cpu.Memory[0x8000] = 0x39 // AND Absolute,Y
				cpu.Memory[0x8001] = 0x34 // Operand
				cpu.Memory[0x8002] = 0x12 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 21 (AND (Indirect,X))",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.X = 1
				cpu.Memory[0x0021] = 0x34 // Low byte of address
				cpu.Memory[0x0022] = 0x12 // High byte of address
				cpu.Memory[0x1234] = 0xCC

				cpu.Memory[0x8000] = 0x21 // AND (Indirect,X)
				cpu.Memory[0x8001] = 0x20 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 31 (AND (Indirect),Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Y = 1
				cpu.Memory[0x0020] = 0x34 // Low byte of address
				cpu.Memory[0x0021] = 0x12 // High byte of address
				cpu.Memory[0x1235] = 0xCC

				cpu.Memory[0x8000] = 0x31 // AND (Indirect),Y
				cpu.Memory[0x8001] = 0x20 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xCC {
					t.Errorf("Expected A = 0xCC, got 0x%02X", cpu.A)
				}
			},
		},
	}
	runTests(t, tests)
}

func TestEOROpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 49 (EOR Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x8000] = 0x49 // EOR #imm
				cpu.Memory[0x8001] = 0x0F
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0xF0 {
					t.Errorf("Expected A = 0xF0, got 0x%02X", cpu.A)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: 45 (EOR Zero Page)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x0025] = 0xCC
				cpu.Memory[0x8000] = 0x45 // EOR Zero Page
				cpu.Memory[0x8001] = 0x25 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x33 {
					t.Errorf("Expected A = 0x33, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 55 (EOR Zero Page,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xAA
				cpu.X = 0x01
				cpu.Memory[0x0020+1] = 0xFF
				cpu.Memory[0x8000] = 0x55 // EOR ZeroPage,X
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x55 {
					t.Errorf("Expected A = 0x55, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 4D (EOR Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Memory[0x1234] = 0xCC
				cpu.Memory[0x8000] = 0x4D // EOR Absolute
				cpu.Memory[0x8001] = 0x34 // Operand
				cpu.Memory[0x8002] = 0x12 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x33 {
					t.Errorf("Expected A = 0x33, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 5D (EOR Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.X = 1
				cpu.Memory[0x1235] = 0xCC // Address is (1234 + X)
				cpu.Memory[0x8000] = 0x5D // EOR Absolute,X
				cpu.Memory[0x8001] = 0x34 // Operand
				cpu.Memory[0x8002] = 0x12 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x33 {
					t.Errorf("Expected A = 0x33, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 59 (EOR Absolute,Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Y = 1
				cpu.Memory[0x1235] = 0xCC // Address is (1234 + Y)
				cpu.Memory[0x8000] = 0x59 // EOR Absolute,Y
				cpu.Memory[0x8001] = 0x34 // Operand
				cpu.Memory[0x8002] = 0x12 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x33 {
					t.Errorf("Expected A = 0x33, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 41 (EOR (Indirect,X))",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.X = 1
				cpu.Memory[0x0021] = 0x34 // Low byte of address
				cpu.Memory[0x0022] = 0x12 // High byte of address
				cpu.Memory[0x1234] = 0xCC

				cpu.Memory[0x8000] = 0x41 // EOR (Indirect,X)
				cpu.Memory[0x8001] = 0x20 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x33 {
					t.Errorf("Expected A = 0x33, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 51 (EOR (Indirect),Y)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0xFF
				cpu.Y = 1
				cpu.Memory[0x0020] = 0x34 // Low byte of address
				cpu.Memory[0x0021] = 0x12 // High byte of address
				cpu.Memory[0x1235] = 0xCC

				cpu.Memory[0x8000] = 0x51 // EOR (Indirect),Y
				cpu.Memory[0x8001] = 0x20 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x33 {
					t.Errorf("Expected A = 0x33, got 0x%02X", cpu.A)
				}
			},
		},
	}
	runTests(t, tests)
}

func TestORAOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 09 (ORA Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Memory[0x8000] = 0x09 // ORA #imm
				cpu.Memory[0x8001] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x11 {
					t.Errorf("Expected A = 0x11, got 0x%02X", cpu.A)
				}
			},
		},
		{
			name: "Opcode: 15 (ORA ZeroPage,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x80
				cpu.X = 0x01
				cpu.Memory[0x0020+1] = 0x01
				cpu.Memory[0x8000] = 0x15 // ORA ZeroPage,X
				cpu.Memory[0x8001] = 0x20
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x81 {
					t.Errorf("Expected A = 0x81, got 0x%02X", cpu.A)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestCMPOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: C9 (CMP Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x50
				cpu.Memory[0x8000] = 0xC9
				cpu.Memory[0x8001] = 0x40
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
			},
		},
		{
			name: "Opcode: C5 (CMP ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x10
				cpu.Memory[0x0025] = 0x10
				cpu.Memory[0x8000] = 0xC5 // CMP Zero Page
				cpu.Memory[0x8001] = 0x25 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestCPXOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: E0 (CPX Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x10
				cpu.Memory[0x8000] = 0xE0 // CPX #imm
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
		{
			name: "Opcode: E4 (CPX ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x10
				cpu.Memory[0x0025] = 0x10
				cpu.Memory[0x8000] = 0xE4 // CPX Zero Page
				cpu.Memory[0x8001] = 0x25 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestCPYOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: C0 (CPY Immediate)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x10
				cpu.Memory[0x8000] = 0xC0 // CPY #imm
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
		{
			name: "Opcode: C4 (CPY ZeroPage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Y = 0x10
				cpu.Memory[0x0025] = 0x10
				cpu.Memory[0x8000] = 0xC4 // CPY Zero Page
				cpu.Memory[0x8001] = 0x25 // Operand
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestASLOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 0A (ASL A)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0b11000001        // старший бит установлен
				cpu.Memory[0x8000] = 0x0A // ASL A
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0b10000010 {
					t.Errorf("Expected A = 0b10000010, got 0x%08b", cpu.A)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: 06 (ASL Zeropage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x0010] = 0b01000001 // значение в zero page
				cpu.Memory[0x8000] = 0x06       // ASL $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x0010]
				if value != 0b10000010 {
					t.Errorf("Expected memory[0x10] = 0b10000010, got 0x%08b", value)
				}
				if cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: 0E (ASL Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x1234] = 0b01000001 // значение в absolute
				cpu.Memory[0x8000] = 0x0E       // ASL $1234
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x1234]
				if value != 0b10000010 {
					t.Errorf("Expected memory[0x1234] = 0b10000010, got 0x%08b", value)
				}
				if cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestLSROpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 4A (LSR A)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0b11000001        // старший бит установлен
				cpu.Memory[0x8000] = 0x4A // LSR A
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0b01100000 {
					t.Errorf("Expected A = 0b01100000, got 0x%08b", cpu.A)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: 46 (LSR Zeropage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x0010] = 0b01000001 // значение в zero page
				cpu.Memory[0x8000] = 0x46       // LSR $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x0010]
				if value != 0b00100000 {
					t.Errorf("Expected memory[0x10] = 0b00100000, got 0x%08b", value)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestROROpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 6A (ROR A)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0b11000001        // старший бит установлен
				cpu.Memory[0x8000] = 0x6A // ROR A
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0b01100000 {
					t.Errorf("Expected A = 0b01100000, got 0x%08b", cpu.A)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: 66 (ROR Zeropage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0010] = 0b01000001 // значение в zero page
				cpu.Memory[0x8000] = 0x66       // ROR $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x0010]
				if value != 0b10100000 {
					t.Errorf("Expected memory[0x10] = 0b10100000, got 0x%08b", value)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestROLOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 2A (ROL A)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0b11000001        // старший бит установлен
				cpu.Memory[0x8000] = 0x2A // ROL A
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0b10000010 {
					t.Errorf("Expected A = 0b10000010, got 0x%08b", cpu.A)
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
			},
		},
		{
			name: "Opcode: 26 (ROL Zeropage)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0010] = 0b01000001 // значение в zero page
				cpu.Memory[0x8000] = 0x26       // ROL $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x0010]
				if value != 0b10000011 {
					t.Errorf("Expected memory[0x10] = 0b10000011, got 0x%08b", value)
				}
				if cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: 36 (ROL Zeropage,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.X = 0x10
				cpu.Memory[0x0020] = 0b01000001 // значение в zero page + X
				cpu.Memory[0x8000] = 0x36       // ROL $10,X
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x0020] // после сдвига в адрес $30
				if value != 0b10000011 {
					t.Errorf("Expected memory[0x30] = 0b10000011, got 0x%08b", value)
				}
				if cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: 2E (ROL Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0030] = 0b01000001 // значение по абсолютному адресу
				cpu.Memory[0x8000] = 0x2E       // ROL $30
				cpu.Memory[0x8001] = 0x30
				cpu.Memory[0x8002] = 0x00
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x0030]
				if value != 0b10000011 {
					t.Errorf("Expected memory[0x30] = 0b10000011, got 0x%08b", value)
				}
				if cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestRTIOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 40 (RTI)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// эмуляция вызова прерывания: на стек положили статус и PC
				cpu.Push16(0x1234)      // PC
				cpu.Push(FlagZ | FlagC) // статус-регистр

				cpu.Memory[0x8000] = 0x40 // RTI
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x1234 {
					t.Errorf("Expected PC = 0x1234, got 0x%04X", cpu.PC)
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
				if !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Carry flag to be set")
				}
				if !cpu.GetFlag(FlagU) {
					t.Errorf("Expected Unused flag to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestRTSOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 60 (RTS)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// Подделываем возвратный адрес на стеке
				cpu.Push16(0x1233) // Возврат должен быть на 0x1234 (прибавляется 1)

				cpu.Memory[0x8000] = 0x60 // RTS
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x1234 {
					t.Errorf("Expected PC = 0x1234, got 0x%04X", cpu.PC)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestJSROpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 20 (JSR Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x8000] = 0x20 // JSR $1234
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x1234 {
					t.Errorf("Expected PC = 0x1234, got 0x%04X", cpu.PC)
				}
				retAddr := cpu.Pull16()
				if retAddr != 0x8002 {
					t.Errorf("Expected return address = 0x8002, got 0x%04X", retAddr)
				}
			},
		},
		{
			name: "Opcode: 20 + 60 (JSR + RTS)",
			init: func(cpu *CPU) {
				// Установим Reset Vector на $8000
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// Основной код
				cpu.Memory[0x8000] = 0x20 // JSR $9000
				cpu.Memory[0x8001] = 0x00
				cpu.Memory[0x8002] = 0x90

				cpu.Memory[0x8003] = 0xA9 // LDA #$42 (ожидаемая инструкция после возврата)
				cpu.Memory[0x8004] = 0x42

				// Подпрограмма по адресу $9000
				cpu.Memory[0x9000] = 0x60 // RTS
			},
			assert: func(cpu *CPU, t *testing.T) {
				cpu.Execute()
				cpu.Execute()

				if cpu.A != 0x42 {
					t.Errorf("Expected A = 0x42, got 0x%02X", cpu.A)
				}
				if cpu.PC != 0x8005 {
					t.Errorf("Expected PC = 0x8005 after RTS and LDA, got 0x%04X", cpu.PC)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBRKOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 00 (BRK)",
			init: func(cpu *CPU) {
				// Устанавливаем Reset-вектор на 0x8000 и делаем Reset
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// Кладём инструкцию BRK по адресу 0x8000
				cpu.Memory[0x8000] = 0x00 // BRK
				cpu.Memory[0x8001] = 0xFF // фиктивный байт (игнорируется)

				// Вектор прерываний BRK
				cpu.Memory[0xFFFE] = 0x34
				cpu.Memory[0xFFFF] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				// PC должен быть установлен из вектора
				if cpu.PC != 0x1234 {
					t.Errorf("Expected PC = 0x1234, got 0x%04X", cpu.PC)
				}

				// Снимаем статус-регистр с флагами
				status := cpu.Pull()
				expectedStatus := (cpu.P | FlagB | FlagU)
				if status != expectedStatus {
					t.Errorf("Expected status = 0x%02X, got 0x%02X", expectedStatus, status)
				}

				// Снимаем адрес возврата
				returnAddr := cpu.Pull16()
				if returnAddr != 0x8002 {
					t.Errorf("Expected return address = 0x8002, got 0x%04X", returnAddr)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestJMPOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 6C (JMP indirect normal)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x8000] = 0x6C // JMP ($1234)
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12

				// Задаём переходный адрес по $1234 = $40, $1235 = $80 -> PC = $8040
				cpu.Memory[0x1234] = 0x40
				cpu.Memory[0x1235] = 0x80
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8040 {
					t.Errorf("Expected PC = 0x8040, got 0x%04X", cpu.PC)
				}
			},
		},
		{
			name: "Opcode: 6C (JMP indirect with page boundary bug)",
			init: func(cpu *CPU) {
				// Установим начальный адрес выполнения
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// Указываем инструкцию JMP ($30FF)
				cpu.Memory[0x8000] = 0x6C // JMP (indirect)
				cpu.Memory[0x8001] = 0xFF // low byte of pointer = 0x30FF
				cpu.Memory[0x8002] = 0x30 // high byte of pointer

				// Устанавливаем значение по адресу $30FF (low) и $3000 (high)
				cpu.Memory[0x30FF] = 0x80 // low byte of target address
				cpu.Memory[0x3000] = 0x40 // high byte of target address (bug!)

				// А вот здесь то, что "ожидается логично", но не будет использовано
				cpu.Memory[0x3100] = 0x50
			},
			assert: func(cpu *CPU, t *testing.T) {
				expected := uint16(0x4080)
				if cpu.PC != expected {
					t.Errorf("Expected PC = 0x%04X due to page boundary bug, got 0x%04X", expected, cpu.PC)
				}
			},
		},
		{
			name: "Opcode: 4C (JMP absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x8000] = 0x4C // JMP $1234
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x1234 {
					t.Errorf("Expected PC = 0x1234, got 0x%04X", cpu.PC)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBEQOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: F0 (BEQ) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// Устанавливаем флаг Z
				cpu.SetFlag(FlagZ, true)

				cpu.Memory[0x80FE] = 0xF0 // BEQ
				cpu.Memory[0x80FF] = 0x00 // offset = +1 (0x8100)
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8100 {
					t.Errorf("Expected PC = 0x8100, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2+2 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: F0 (BEQ) same page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagZ, true)

				// PC начнётся с 0x8000, переход вперёд на 0x8002 + 2 = 0x8004
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xF0 // BEQ
				cpu.Memory[0x8001] = 0x02 // смещение +2
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8004 {
					t.Errorf("Expected PC = 0x8004, got 0x%04X", cpu.PC)
				}
			},
		},
		{
			name: "Opcode: F0 (BEQ) no branch when Z flag is clear",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				// Сбрасываем флаг Z (переход не должен произойти)
				cpu.SetFlag(FlagZ, false)

				// BEQ + offset
				cpu.Memory[0x8000] = 0xF0 // BEQ
				cpu.Memory[0x8001] = 0x10 // offset (направление не важно — он не сработает)
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles (branch not taken), got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: F0 (BEQ) branch forward same page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagZ, true)

				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xF0 // BEQ
				cpu.Memory[0x8001] = 0x02 // offset = +2 → 0x8000 + 2 + 2 = 0x8004
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8004 {
					t.Errorf("Expected PC = 0x8004, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2+1 { // переход, но та же страница
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: F0 (BEQ) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x81
				cpu.Reset()

				cpu.SetFlag(FlagZ, true)

				cpu.PC = 0x8100
				cpu.Memory[0x8100] = 0xF0 // BEQ
				cpu.Memory[0x8101] = 0xFE // offset = -2 → 0x8100 + 2 - 2 = 0x8100
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8100 {
					t.Errorf("Expected PC = 0x8100, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2+1 { // переход, та же страница
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: F0 (BEQ) no branch when Z flag is clear",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagZ, false)

				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xF0 // BEQ
				cpu.Memory[0x8001] = 0x10 // offset (ignored)
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles (branch not taken), got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBPLopcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 10 (BPL) branch taken same page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagN, false) // Positive
				cpu.Memory[0x8000] = 0x10 // BPL
				cpu.Memory[0x8001] = 0x10 // Offset +16 => 0x8012
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 10 (BPL) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagN, false)
				cpu.Memory[0x80FE] = 0x10
				cpu.Memory[0x80FF] = 0x01 // 0x8100 + 1 = 0x8101
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 10 (BPL) no branch when negative flag set",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagN, true)  // Negative
				cpu.Memory[0x8000] = 0x10 // BPL
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 10 (BPL) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagN, false)
				cpu.Memory[0x0100] = 0x10
				cpu.Memory[0x0101] = 0xFC // -4 → 0x00FE
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBMIOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 30 (BMI) branch taken",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagN, true)

				cpu.Memory[0x8000] = 0x30
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 30 (BMI) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagN, true)
				cpu.Memory[0x80FE] = 0x30
				cpu.Memory[0x80FF] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 30 (BMI) no branch",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagN, false)

				cpu.Memory[0x8000] = 0x30
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got 0x%04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 30 (BMI) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagN, true)
				cpu.Memory[0x0100] = 0x30
				cpu.Memory[0x0101] = 0xFC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBVCOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 50 (BVC) branch taken",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagV, false)
				cpu.Memory[0x8000] = 0x50
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got %04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 50 (BVC) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagV, false)
				cpu.Memory[0x80FE] = 0x50
				cpu.Memory[0x80FF] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 50 (BVC) no branch",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagV, true)
				cpu.Memory[0x8000] = 0x50
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got %04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 50 (BVC) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagV, false)
				cpu.Memory[0x0100] = 0x50
				cpu.Memory[0x0101] = 0xFC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBVSOpcodes(t *testing.T) {
	tests := []OpcodeTest{

		{
			name: "Opcode: 70 (BVS) branch taken",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagV, true)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0x70
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got %04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 70 (BVS) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagV, true)
				cpu.Memory[0x80FE] = 0x70
				cpu.Memory[0x80FF] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 70 (BVS) no branch",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagV, false)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0x70
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got %04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 70 (BVS) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagV, true)
				cpu.Memory[0x0100] = 0x70
				cpu.Memory[0x0101] = 0xFC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBCCOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 90 (BCC) branch taken",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, false)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0x90
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got %04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 90 (BCC) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, false)
				cpu.Memory[0x80FE] = 0x90
				cpu.Memory[0x80FF] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 90 (BCC) no branch",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0x90
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got %04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: 90 (BCC) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagC, false)
				cpu.Memory[0x0100] = 0x90
				cpu.Memory[0x0101] = 0xFC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBCSOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: B0 (BCS) branch taken",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xB0
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got %04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: B0 (BCS) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x80FE] = 0xB0
				cpu.Memory[0x80FF] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: B0 (BCS) no branch",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagC, false)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xB0
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got %04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: B0 (BCS) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagC, true)
				cpu.Memory[0x0100] = 0xB0
				cpu.Memory[0x0101] = 0xFC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBNEOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: D0 (BNE) branch taken",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagZ, false)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xD0
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8012 {
					t.Errorf("Expected PC = 0x8012, got %04X", cpu.PC)
				}
				if cpu.Cycles != 3 {
					t.Errorf("Expected 3 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: D0 (BNE) branch across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0xFE
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagZ, false)
				cpu.Memory[0x80FE] = 0xD0
				cpu.Memory[0x80FF] = 0x01
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8101 {
					t.Errorf("Expected PC = 0x8101, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: D0 (BNE) no branch",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.SetFlag(FlagZ, true)
				cpu.PC = 0x8000
				cpu.Memory[0x8000] = 0xD0
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8002 {
					t.Errorf("Expected PC = 0x8002, got %04X", cpu.PC)
				}
				if cpu.Cycles != 2 {
					t.Errorf("Expected 2 cycles, got %d", cpu.Cycles)
				}
			},
		},
		{
			name: "Opcode: D0 (BNE) branch backward across page",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x01
				cpu.Reset()

				cpu.SetFlag(FlagZ, false)
				cpu.Memory[0x0100] = 0xD0
				cpu.Memory[0x0101] = 0xFC
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x00FE {
					t.Errorf("Expected PC = 0x00FE, got %04X", cpu.PC)
				}
				if cpu.Cycles != 4 {
					t.Errorf("Expected 4 cycles, got %d", cpu.Cycles)
				}
			},
		},
	}

	runTests(t, tests)
}

func TestStackOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 48 (PHA)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0x42
				cpu.Memory[0x8000] = 0x48 // PHA
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Pull()
				if value != 0x42 {
					t.Errorf("Expected top of stack = 0x42, got 0x%02X", value)
				}
			},
		},
		{
			name: "Opcode: 68 (PLA)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Push(0x99)
				cpu.Memory[0x8000] = 0x68 // PLA
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.A != 0x99 {
					t.Errorf("Expected A = 0x99, got 0x%02X", cpu.A)
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: 08 (PHP)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.P = FlagZ | FlagC
				cpu.Memory[0x8000] = 0x08 // PHP
			},
			assert: func(cpu *CPU, t *testing.T) {
				status := cpu.Pull()
				expected := FlagZ | FlagC | FlagB | FlagU
				if status != byte(expected) {
					t.Errorf("Expected pushed status = 0x%02X, got 0x%02X", expected, status)
				}
			},
		},
		{
			name: "Opcode: 28 (PLP)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Push(FlagZ | FlagC)
				cpu.Memory[0x8000] = 0x28 // PLP
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) || !cpu.GetFlag(FlagC) {
					t.Errorf("Expected Z and C flags to be set")
				}
				if !cpu.GetFlag(FlagU) {
					t.Errorf("Expected unused flag (U) to be set")
				}
			},
		},
	}

	runTests(t, tests)
}

func TestINCOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: E6 (INC Zero Page)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x10] = 0x7F
				cpu.Memory[0x8000] = 0xE6
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				value := cpu.Memory[0x10]
				if value != 0x80 {
					t.Errorf("Expected memory[0x10] = 0x80, got 0x%02X", value)
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "Opcode: F6 (INC Zero Page,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x01
				cpu.Memory[0x11] = 0x00
				cpu.Memory[0x8000] = 0xF6
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x11] != 0x01 {
					t.Errorf("Expected memory[0x11] = 0x01, got 0x%02X", cpu.Memory[0x11])
				}
			},
		},
		{
			name: "Opcode: EE (INC Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x1234] = 0xFF
				cpu.Memory[0x8000] = 0xEE
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x1234] != 0x00 {
					t.Errorf("Expected memory[0x1234] = 0x00, got 0x%02X", cpu.Memory[0x1234])
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
		{
			name: "Opcode: FE (INC Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 0x01
				cpu.Memory[0x1235] = 0x40
				cpu.Memory[0x8000] = 0xFE
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x1235] != 0x41 {
					t.Errorf("Expected memory[0x1235] = 0x41, got 0x%02X", cpu.Memory[0x1235])
				}
			},
		},
	}

	runTests(t, tests)
}

func TestBITOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: 24 (BIT Zero Page)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0b00000000
				cpu.Memory[0x0010] = 0b11000000
				cpu.Memory[0x8000] = 0x24 // BIT $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
				if !cpu.GetFlag(FlagV) {
					t.Errorf("Expected Overflow flag to be set")
				}
			},
		},
		{
			name: "Opcode: 2C (BIT Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.A = 0b00000001
				cpu.Memory[0x1234] = 0b00000001
				cpu.Memory[0x8000] = 0x2C
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
				if cpu.GetFlag(FlagV) {
					t.Errorf("Expected Overflow flag to be clear")
				}
			},
		},
	}
	runTests(t, tests)
}

func TestDECOpcodes(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: C6 (DEC Zero Page)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x10] = 0x01
				cpu.Memory[0x8000] = 0xC6
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x10] != 0x00 {
					t.Errorf("Expected memory[0x10] = 0x00, got %02X", cpu.Memory[0x10])
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
			},
		},
		{
			name: "Opcode: D6 (DEC Zero Page,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 1
				cpu.Memory[0x11] = 0xFF
				cpu.Memory[0x8000] = 0xD6
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x11] != 0xFE {
					t.Errorf("Expected memory[0x11] = 0xFE, got %02X", cpu.Memory[0x11])
				}
			},
		},
		{
			name: "Opcode: CE (DEC Absolute)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x1234] = 0x01
				cpu.Memory[0x8000] = 0xCE
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x1234] != 0x00 {
					t.Errorf("Expected memory[0x1234] = 0x00, got %02X", cpu.Memory[0x1234])
				}
			},
		},
		{
			name: "Opcode: DE (DEC Absolute,X)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.X = 1
				cpu.Memory[0x1235] = 0x42
				cpu.Memory[0x8000] = 0xDE
				cpu.Memory[0x8001] = 0x34
				cpu.Memory[0x8002] = 0x12
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x1235] != 0x41 {
					t.Errorf("Expected memory[0x1235] = 0x41, got %02X", cpu.Memory[0x1235])
				}
			},
		},
	}
	runTests(t, tests)
}

func TestNOPOpcode(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "Opcode: EA (NOP)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x8000] = 0xEA // NOP
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.PC != 0x8001 {
					t.Errorf("Expected PC = 0x8001 after NOP, got %04X", cpu.PC)
				}
			},
		},
	}
	runTests(t, tests)
}

func TestDECExecute(t *testing.T) {
	tests := []OpcodeTest{
		{
			name: "DEC to zero (sets Z flag)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x0010] = 0x01
				cpu.Memory[0x8000] = 0xC6 // DEC $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x0010] != 0x00 {
					t.Errorf("Expected memory[0x10] = 0x00, got 0x%02X", cpu.Memory[0x0010])
				}
				if !cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be set")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
			},
		},
		{
			name: "DEC wraps to 0xFF (sets N flag)",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x0010] = 0x00
				cpu.Memory[0x8000] = 0xC6 // DEC $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x0010] != 0xFF {
					t.Errorf("Expected memory[0x10] = 0xFF, got 0x%02X", cpu.Memory[0x0010])
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
				if !cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be set")
				}
			},
		},
		{
			name: "DEC clears Z and N",
			init: func(cpu *CPU) {
				cpu.Memory[ResetVector] = 0x00
				cpu.Memory[ResetVector+1] = 0x80
				cpu.Reset()

				cpu.Memory[0x0010] = 0x10
				cpu.Memory[0x8000] = 0xC6 // DEC $10
				cpu.Memory[0x8001] = 0x10
			},
			assert: func(cpu *CPU, t *testing.T) {
				if cpu.Memory[0x0010] != 0x0F {
					t.Errorf("Expected memory[0x10] = 0x0F, got 0x%02X", cpu.Memory[0x0010])
				}
				if cpu.GetFlag(FlagZ) {
					t.Errorf("Expected Zero flag to be clear")
				}
				if cpu.GetFlag(FlagN) {
					t.Errorf("Expected Negative flag to be clear")
				}
			},
		},
	}

	runTests(t, tests)
}
