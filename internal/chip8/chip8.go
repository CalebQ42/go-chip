package chip8

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"time"
)

type Chip8 struct {
	delayTicker *time.Ticker
	soundTicker *time.Ticker
	pc          *uint16
	screen      *Screen
	stack       []*uint16
	running     bool
	iRegister   uint16
	memory      Ram
	registers   [16]byte
	delayReg    byte
	soundReg    byte
}

func NewChip8(rom io.Reader) (*Chip8, error) {
	romDat, err := io.ReadAll(rom)
	if err != nil {
		return nil, err
	}
	initPC := uint16(512)
	return &Chip8{
		memory: InitRam(romDat),
		screen: NewScreen(),
		pc:     &initPC,
	}, nil
}

func (c *Chip8) Start() {
	if c.delayTicker == nil {
		c.delayTicker = time.NewTicker(time.Second / 60)
		go func() {
			c.delayTicker.Stop()
			for {
				<-c.delayTicker.C
				c.delayReg--
				if c.delayReg == 0 {
					c.delayTicker.Stop()
				}
			}
		}()
	}
	if c.soundTicker == nil {
		c.soundTicker = time.NewTicker(time.Second / 60)
		go func() {
			c.soundTicker.Stop()
			for {
				<-c.delayTicker.C
				c.soundReg--
				if c.delayReg == 0 {
					c.delayTicker.Stop()
					//TODO: stop sound
				}
			}
		}()
	}
	// c.screen.StartPrinting()
	c.running = true
	for c.running {
		fmt.Println(c.screen)
		time.Sleep(time.Second / 2)
		c.nextInstruction()
	}

}
func (c *Chip8) Stop() {
	c.running = false
}

func (c *Chip8) PrintRegisters() {
	fmt.Println(c.registers)
}

func (c *Chip8) setDelayTime(t byte) {
	if c.delayReg > 0 {
		c.delayReg = t
		return
	}
	c.delayReg = t
	c.delayTicker.Reset(time.Second / 60)
}

func (c *Chip8) setSoundTime(t byte) {
	if c.soundReg > 0 {
		c.soundReg = t
		return
	}
	c.soundReg = t
	c.soundTicker.Reset(time.Second / 60)
}

func (c *Chip8) nextInstruction() {
	ins := binary.BigEndian.Uint16(c.memory[*c.pc:])
	*c.pc += 2
	// fmt.Printf("%x\n", ins)
	c.handleInstruction(ins)
}

func (c *Chip8) handleInstruction(ins uint16) {
	switch ins >> 12 {
	case 0:
		if ins == 0x00EE {
			c.pc = c.stack[len(c.stack)-1]
			c.stack = c.stack[:len(c.stack)-1]
		} else if ins == 0x00E0 {
			c.screen.Clear()
		} else {
			log.Fatalf("invalid instruction: %x", ins)
		}
	case 1:
		*c.pc = ins & 0xFFF
	case 2:
		newPC := uint16(ins & 0xFFF)
		c.stack = append(c.stack, c.pc)
		c.pc = &newPC
	case 3:
		if c.registers[(ins>>8)&0xF] == byte(ins&0xFF) {
			*c.pc += 2
		}
	case 4:
		if c.registers[(ins>>8)&0xF] != byte(ins&0xFF) {
			*c.pc += 2
		}
	case 5:
		if c.registers[(ins>>8)&0xF] == c.registers[(ins>>4)&0xF] {
			*c.pc += 2
		}
	case 6:
		c.registers[(ins>>8)&0xF] = byte(ins & 0xFF)
	case 7:
		c.registers[(ins>>8)&0xF] += byte(ins & 0xFF)
	case 8:
		c.multiRegisterMath(ins)
	case 9:
		if c.registers[(ins>>8)&0xF] != c.registers[(ins>>4)&0xF] {
			*c.pc += 2
		}
	case 0xA:
		c.iRegister = ins & 0xFFF
	case 0xB:
		*c.pc = ins&0xFFF + uint16(c.registers[0])
	case 0xC:
		c.registers[(ins>>8)&0xF] = byte(rand.UintN(256)) & byte(ins&0xFF)
	case 0xD:
		//TODO: Add sprite to screen
	case 0xE:
		//TODO: Handle keyboard. skip instructions
	case 0xF:
		//0xF instructions are kind of random. we just group them together
		c.leftovers(ins)
	default:
		log.Fatalln("invalid instruction:", ins)
	}
}

func (c *Chip8) multiRegisterMath(ins uint16) {
	x := (ins >> 8) & 0xF
	y := (ins >> 4) & 0xF
	switch ins & 0xF {
	case 0:
		c.registers[x] = c.registers[y]
	case 1:
		c.registers[x] |= c.registers[y]
	case 2:
		c.registers[x] &= c.registers[y]
	case 3:
		c.registers[x] ^= c.registers[y]
	case 4:
		if uint16(c.registers[x])+uint16(c.registers[y]) > 255 {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
		c.registers[x] += c.registers[y]
	case 5:
		if c.registers[x] > c.registers[y] {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
		c.registers[x] -= c.registers[y]
	case 6:
		if c.registers[x]&0x1 == 0x1 {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
		c.registers[x] /= 2
	case 7:
		if c.registers[y] > c.registers[x] {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
		c.registers[x] = c.registers[y] - c.registers[x]
	case 0xE:
		if c.registers[x]&0x80 == 0x80 {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
		c.registers[x] *= 2
	default:
		log.Fatalln("invalid instruction:", ins)
	}
}

func (c *Chip8) leftovers(ins uint16) {
	x := (ins >> 8) & 0xF
	switch ins & 0xFF {
	case 7:
		c.registers[x] = c.delayReg
	case 0xA:
		//TODO: wait until keyboard button is pressed
	case 0x15:
		c.setDelayTime(c.registers[x])
	case 0x18:
		c.setSoundTime(c.registers[x])
	case 0x1E:
		c.iRegister += uint16(c.registers[x])
	case 0x29:
		//TODO: set iRegister to sprite represented in register[x]
	case 0x33:
		c.memory[c.iRegister] = (c.registers[x] / 100) % 10
		c.memory[c.iRegister+1] = (c.registers[x] / 10) % 10
		c.memory[c.iRegister+2] = c.registers[x] % 10
	case 0x55:
		for i := uint16(0); i <= x; i++ {
			c.memory[c.iRegister+i] = c.registers[i]
		}
	case 0x65:
		for i := uint16(0); i <= x; i++ {
			c.registers[i] = c.memory[c.iRegister+i]
		}
	default:
		log.Fatalln("invalid instruction:", ins)
	}
}
