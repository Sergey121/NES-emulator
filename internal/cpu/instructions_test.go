package cpu_test

import (
	"testing"

	"github.com/sergey121/nes-emulator/internal/cpu"
)

func TestLDAImmediate(t *testing.T) {
	cpuInstance := cpu.New()

	// Устанавливаем Reset Vector на 0x8000
	cpuInstance.Memory[cpu.ResetVector] = 0x00
	cpuInstance.Memory[cpu.ResetVector+1] = 0x80
	cpuInstance.Reset()

	// Пишем LDA #$42 по адресу 0x8000
	cpuInstance.Memory[0x8000] = 0xA9 // LDA Immediate
	cpuInstance.Memory[0x8001] = 0xC0

	// Выполняем
	cpuInstance.Execute()

	if cpuInstance.A != 0xC0 {
		t.Errorf("Expected A = 0xC0, got 0x%02X", cpuInstance.A)
	}
	if cpuInstance.GetFlag(cpu.FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if !cpuInstance.GetFlag(cpu.FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}
func TestLDAZeroPage(t *testing.T) {
	cpuInstance := cpu.New()

	// Устанавливаем Reset Vector на 0x8000
	cpuInstance.Memory[cpu.ResetVector] = 0x00
	cpuInstance.Memory[cpu.ResetVector+1] = 0x80
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
	if cpuInstance.GetFlag(cpu.FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if !cpuInstance.GetFlag(cpu.FlagN) {
		t.Errorf("Expected Negative flag to be set")
	}
}

func TestLDAZeroPageX_Wraparound(t *testing.T) {
	cpuInstance := cpu.New()

	// Устанавливаем Reset Vector на 0x8000
	cpuInstance.Memory[cpu.ResetVector] = 0x00
	cpuInstance.Memory[cpu.ResetVector+1] = 0x80
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
	if cpuInstance.GetFlag(cpu.FlagZ) {
		t.Errorf("Expected Zero flag to be cleared")
	}
	if (cpuInstance.A&0x80) == 0 && cpuInstance.GetFlag(cpu.FlagN) {
		t.Errorf("Expected Negative flag to be cleared")
	}
}
