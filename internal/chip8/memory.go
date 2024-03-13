package chip8

var (
	hexSprites = [][]byte{
		{0xF0, 0x90, 0x90, 0x90, 0xF0}, //0
		{0x20, 0x60, 0x20, 0x20, 0x70}, //1
		{0xF0, 0x10, 0xF0, 0x80, 0xF0}, //2
		{0xF0, 0x10, 0xf0, 0x10, 0xF0}, //3
		{0x90, 0x90, 0xF0, 0x10, 0x10}, //4
		{0xF0, 0x80, 0xF0, 0x10, 0xF0}, //5
		{0xF0, 0x80, 0xF0, 0x90, 0xF0}, //6
		{0xF0, 0x10, 0x20, 0x40, 0x40}, //7
		{0xF0, 0x90, 0xF0, 0x90, 0x90}, //8
		{0xF0, 0x90, 0xF0, 0x10, 0xF0}, //9
		{0xF0, 0x90, 0xF0, 0x90, 0x90}, //a
		{0xE0, 0x90, 0xE0, 0x90, 0xE0}, //b
		{0xF0, 0x90, 0x80, 0x80, 0xF0}, //c
		{0xE0, 0x90, 0x90, 0x90, 0xE0}, //d
		{0xF0, 0x80, 0xF0, 0x80, 0xF0}, //e
		{0xF0, 0x80, 0xF0, 0x80, 0x80}, //f
	}
)

type Ram [4096]byte

func InitRam() *Ram {
	r := Ram{}
	for i := 0; i < len(hexSprites); i++ {
		copy(r[5*i:], hexSprites[i])
	}
	return &r
}
