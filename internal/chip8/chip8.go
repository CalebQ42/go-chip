package chip8

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Chip8 struct {
	delayTicker   *time.Ticker
	soundTicker   *time.Ticker
	pc            *uint16
	screen        *Screen
	keyboardBlock chan byte
	keyState      map[byte]bool
	stack         []*uint16
	iRegister     uint16
	memory        Ram
	registers     [16]byte
	running       bool
	delayReg      byte
	soundReg      byte
}

func NewChip8(rom io.Reader) (*Chip8, error) {
	romDat, err := io.ReadAll(rom)
	if err != nil {
		return nil, err
	}
	initPC := uint16(512)
	return &Chip8{
		memory: InitRam(romDat),
		screen: &Screen{},
		pc:     &initPC,
	}, nil
}

func (c *Chip8) Start() {
	c.running = true

	if c.delayTicker == nil {
		c.delayTicker = time.NewTicker(time.Second / 60)
		go func() {
			c.delayTicker.Stop()
			for c.running {
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
			for c.running {
				//TODO: Make a noise
				<-c.delayTicker.C
				c.soundReg--
				if c.delayReg == 0 {
					c.delayTicker.Stop()
					//TODO: stop sound
				}
			}
		}()
	}
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("Chip-8")
	ebiten.SetTPS(15)
	err := ebiten.RunGame(c)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Chip8) Update() error {
	if c.keyboardBlock == nil {
		for n := 0; n < 100; n++ {
			stop := c.nextInstruction()
			if c.keyboardBlock != nil {
				c.keyState = make(map[byte]bool)
				for i := byte(0); i <= 0xF; i++ {
					c.keyState[i] = c.checkKeyDown(i)
				}
				break
			}
			if stop {
				break
			}
		}
	} else {
		for i := byte(0); i <= 0xF; i++ {
			if c.keyState[i] != c.checkKeyDown(i) {
				if c.keyState[i] {
					c.keyState[i] = false
				} else {
					c.keyboardBlock <- i
					break
				}
			}
		}
	}
	return nil
}

func (c *Chip8) Draw(screen *ebiten.Image) {
	c.screen.DrawTo(screen)
}

func (c *Chip8) Layout(_, _ int) (int, int) {
	return 64, 32
}

func (c *Chip8) checkKeyDown(key byte) bool {
	switch key {
	case 0:
		return ebiten.IsKeyPressed(ebiten.KeyX)
	case 1:
		return ebiten.IsKeyPressed(ebiten.Key1)
	case 2:
		return ebiten.IsKeyPressed(ebiten.Key2)
	case 3:
		return ebiten.IsKeyPressed(ebiten.Key3)
	case 4:
		return ebiten.IsKeyPressed(ebiten.KeyQ)
	case 5:
		return ebiten.IsKeyPressed(ebiten.KeyW)
	case 6:
		return ebiten.IsKeyPressed(ebiten.KeyE)
	case 7:
		return ebiten.IsKeyPressed(ebiten.KeyA)
	case 8:
		return ebiten.IsKeyPressed(ebiten.KeyS)
	case 9:
		return ebiten.IsKeyPressed(ebiten.KeyD)
	case 0xA:
		return ebiten.IsKeyPressed(ebiten.KeyZ)
	case 0xB:
		return ebiten.IsKeyPressed(ebiten.KeyC)
	case 0xC:
		return ebiten.IsKeyPressed(ebiten.Key4)
	case 0xD:
		return ebiten.IsKeyPressed(ebiten.KeyR)
	case 0xE:
		return ebiten.IsKeyPressed(ebiten.KeyF)
	case 0xF:
		return ebiten.IsKeyPressed(ebiten.KeyV)
	}
	return false
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

func (c *Chip8) nextInstruction() bool {
	ins := binary.BigEndian.Uint16(c.memory[*c.pc:])
	*c.pc += 2
	// fmt.Printf("%x\n", ins)
	return c.handleInstruction(ins)
}
