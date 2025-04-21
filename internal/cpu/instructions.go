package cpu

type Instruction struct {
	Name    string
	Opcode  byte
	Bytes   int
	Cycles  int
	Execute func(cpu *CPU, operand byte)
}

var Instructions = map[byte]Instruction{}

func init() {

	Instructions[0xA9] = Instruction{
		Name:    "LDA",
		Opcode:  0xA9,
		Bytes:   2,
		Cycles:  2,
		Execute: ldaImmediate,
	}
}

func ldaImmediate(cpu *CPU, operand byte) {
	cpu.A = operand
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, (cpu.A&0x80) != 0)
}
