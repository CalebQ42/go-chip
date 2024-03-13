package main

import (
	"flag"
	"log"
	"os"

	"github.com/CalebQ42/go-chip/internal/chip8"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalln("Please provide the path to the rom")
	}
	rom, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}
	comp, err := chip8.NewChip8(rom)
	if err != nil {
		log.Fatalln(err)
	}
	comp.Start()
	// TODO
}
