package chip8

import (
	"encoding/binary"
	"io"
)

type Chip8 struct {
	rom io.ReadSeekCloser

	iRegister uint16
	memory    Ram
	registers [16]byte
}

func (c *Chip8) readNextInstruction() (ins uint16, err error) {
	err = binary.Read(c.rom, binary.BigEndian, &ins)
	return
}
