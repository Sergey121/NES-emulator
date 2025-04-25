package cpu

import (
	"testing"
)

func TestOpcode_A9(t *testing.T) {
	cpuInstance := New()

	// Устанавливаем Reset Vector на 0x8000
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Пишем LDA #$42 по адресу 0x8000
	cpuInstance.Memory[0x8000] = 0xA9 // LDA Immediate
	cpuInstance.Memory[0x8001] = 0xC0

	// Выполняем
	cpuInstance.Execute()

	if cpuInstance.A != 0xC0 {
		t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}
func TestOpcode_A5(t *testing.T) {
	cpuInstance := New()

	// Устанавливаем Reset Vector на 0x8000
	cpuInstance.Memory[ResetVector] = 0x00
	cpuInstance.Memory[ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Устанавливаем значение в Zero Page
	cpuInstance.Memory[0x0042] = 0xC0

	// Пишем LDA $42 по адресу 0x8000
	cpuInstance.Memory[0x8000] = 0xA5 // LDA Zero Page
	cpuInstance.Memory[0x8001] = 0x42

	// Выполняем
	cpuInstance.Execute()

	if cpuInstance.A != 0xC0 {
		t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestOpcode_B5(t *testing.T) {
	cpuInstance := New()

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

	// Выполняем инструкцию
	cpuInstance.Execute()

	if cpuInstance.A != 0xC0 {
		t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestOpcode_B5_Wraparound(t *testing.T) {
	cpuInstance := New()

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

	// Выполняем инструкцию
	cpuInstance.Execute()

	if cpuInstance.A != 0xAB {
		t.Errorf("Expected A = 0xAB, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if (cpuInstance.A&0x80) == 0 && cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be cleared")
	}
}

func TestOpcode_AD(t *testing.T) {
	cpuInstance := New()

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

	// Выполняем инструкцию
	cpuInstance.Execute()

	if cpuInstance.A != 0x99 {
		t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestOpcode_BD(t *testing.T) {
	cpuInstance := New()

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

	// Выполняем инструкцию
	cpuInstance.Execute()

	if cpuInstance.A != 0x99 {
		t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestOpcode_B9(t *testing.T) {
	cpuInstance := New()

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

	// Выполняем инструкцию
	cpuInstance.Execute()

	if cpuInstance.A != 0x99 {
		t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestOpcode_A1(t *testing.T) {
	cpuInstance := New()

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

	// Выполняем инструкцию
	cpuInstance.Execute()

	if cpuInstance.A != 0x99 {
		t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if (cpuInstance.A&0x80) != 0 && !cpuInstance.GetFlag(FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestOpcode_B1(t *testing.T) {
	cpuInstance := New()

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

	// Выполнение
	cpuInstance.Execute()

	// Проверка
	if cpuInstance.A != 0x99 {
		t.Errorf("Expected A = 0x99, got 0x%02X", cpuInstance.A)
	}
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
