package chip8

import (
	"encoding/binary"
	"fmt"
	"io"
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
				fmt.Print('\a')
				<-c.delayTicker.C
				c.soundReg--
				if c.delayReg == 0 {
					c.delayTicker.Stop()
					//TODO: stop sound
				}
			}
		}()
	}

	c.running = true
	go c.screen.StartPrinting()
	for c.running {
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
