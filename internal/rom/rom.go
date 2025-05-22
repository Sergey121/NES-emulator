package rom

func (c *Cartridge) ReadPRG(addr uint16) byte {
	if c.PRG == nil {
		return 0
	}

	addr -= 0x8000 // PRG-ROM начинается с 0x8000

	if len(c.PRG) == 0x4000 {
		// Если только 16KB PRG, дублируем банк (mirror)
		addr %= 0x4000
	} else if len(c.PRG) == 0x8000 {
		// Если 32KB — прямой доступ
		addr %= 0x8000
	} else {
		addr %= uint16(len(c.PRG))
	}

	return c.PRG[addr]
}

func (c *Cartridge) WritePRG(addr uint16, value byte) {
	panic("WritePRG not implemented")
}
