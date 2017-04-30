package main

import (
	"chip8"
	"time"
	"unsafe"
)

//type vmem [64 * 32]uint32;

const pixelOn, pixelOff uint32 = 0x00EEEECC, 0

func main() {
	var vmem [64 * 32]uint32

	for i := 0; i < 64*32; i++ {
		vmem[i] = chip8.PixOn
		i ++
	}

	display, _ := chip8.NewDisplay()
	display.Update(unsafe.Pointer(&vmem))
	time.Sleep(2 * time.Second)
	display.Clear()
	time.Sleep(2 * time.Second)
	display.Destroy()

}
