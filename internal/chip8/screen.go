package chip8

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Screen struct {
	screen [32][8]byte
}

func (s *Screen) DrawTo(screen *ebiten.Image) {
	for y := range s.screen {
		for x := range s.screen[y] {
			for i := 0; i < 8; i++ {
				if s.screen[y][x]>>(7-i)&1 == 1 {
					screen.Set((x*8)+i, y, color.White)
				} else {
					screen.Set((x*8)+i, y, color.Black)
				}
			}
		}
	}
}

func (s *Screen) String() string {
	out := "--------------------------------------------------------------------\n"
	for y := range s.screen {
		out += "| "
		for x := range s.screen[y] {
			out += fmt.Sprintf("%08b", s.screen[y][x])
		}
		out += " |\n"
	}
	out += "--------------------------------------------------------------------"
	out = strings.ReplaceAll(out, "0", " ")
	out = strings.ReplaceAll(out, "1", "█")
	return out
}

func (s *Screen) PrintScreen() {
	fmt.Print(s)
}

func (s *Screen) Clear() {
	for y := range s.screen {
		for x := range s.screen[y] {
			s.screen[y][x] = 0
		}
	}
}

func (s *Screen) AddSprite(sprite []byte, x, y byte) bool {
	xByte := x / 8
	xOffset := x % 8
	erased := false
	for i := byte(0); i < byte(len(sprite)); i++ {
		yCoord := y + i
		if yCoord >= 32 {
			yCoord -= 32
		}
		orig1 := s.screen[yCoord][xByte]
		s.screen[yCoord][xByte] ^= sprite[i] >> xOffset
		if !erased && s.screen[yCoord][xByte]&orig1 != orig1 {
			erased = true
		}
		if xOffset != 0 {
			x2 := xByte + 1
			if x2 >= 8 {
				x2 = 0
			}
			orig2 := s.screen[yCoord][x2]
			s.screen[yCoord][x2] ^= sprite[i] << (8 - xOffset)
			if !erased && s.screen[yCoord][x2]&orig2 != orig2 {
				erased = true
			}
		}
	}
	return erased
}
