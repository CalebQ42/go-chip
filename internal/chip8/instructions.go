package chip8

import (
	"log"
	"math/rand/v2"
	"time"
)

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
		sprite := make([]byte, ins&0xF)
		copy(sprite, c.memory[c.iRegister:])
		if c.screen.AddSprite(sprite, c.registers[(ins>>8)&0xF], c.registers[(ins>>4)&0xF]) {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
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
		time.Sleep(20 * time.Second)
	case 0x15:
		c.setDelayTime(c.registers[x])
	case 0x18:
		c.setSoundTime(c.registers[x])
	case 0x1E:
		c.iRegister += uint16(c.registers[x])
	case 0x29:
		c.iRegister = 5 * uint16(c.registers[x])
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
