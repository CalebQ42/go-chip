package main

import (
	"encoding/binary"
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
	var ins uint16
	binary.Read(rom, binary.BigEndian, &ins)
	log.Printf("%x", ins)
	ram := chip8.InitRam()
	_ = ram
	// TODO
}
