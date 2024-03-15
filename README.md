# go-chip

A chip-8 interpreter. Special thanks to the technical documents found [here](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM). A collection of roms can be found [here](https://johnearnest.github.io/chip8Archive/). Make sure the rom is for chip8 not schip or xochip.

## Running

Simply pass the rom as an argument

> go-chip ~/Downloads/rom.ch8

## Keyboard

Chip-8 uses a 16 key keypad for input which is mapped to the keyboard:

```text
---------       ---------
|1|2|3|C|   \   |1|2|3|4|
|4|5|6|D|  ==\  |Q|W|E|R|
|7|8|9|E|  ==/  |A|S|D|F|
|A|0|B|F|   /   |Z|X|C|V|
---------       ---------
```

## Possible future improvements

* Implement super chip-8 instructions set.
