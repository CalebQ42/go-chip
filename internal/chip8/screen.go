package chip8

import (
	"fmt"
	"strings"
	"time"
)

type Screen struct {
	t        *time.Ticker
	waitChan chan struct{}
	screen   [32][8]byte
}

func NewScreen() *Screen {
	return &Screen{
		t:        time.NewTicker(time.Second / 60),
		waitChan: make(chan struct{}),
	}
}

func (s *Screen) String() string {
	out := "--------------------------------------------------------------------\n"
	for x := range s.screen {
		out += "| "
		for y := range s.screen[x] {
			out += fmt.Sprintf("%08b", s.screen[x][y])
		}
		out += " |\n"
	}
	out += "--------------------------------------------------------------------"
	out = strings.ReplaceAll(out, "0", " ")
	out = strings.ReplaceAll(out, "1", "█")
	return out
}

func (s *Screen) Clear() {
	for x := range s.screen {
		for y := range s.screen[x] {
			s.screen[x][y] = 0
		}
	}
	s.waitChan <- struct{}{}
}

func (s *Screen) StartPrinting() {
	for {
		<-s.t.C
		for i := 0; i < len(s.waitChan); i++ {
			<-s.waitChan
		}
		fmt.Println(s)
	}
}
