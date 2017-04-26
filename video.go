package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	winTitle  string
	winWidth  int
	winHeight int
	window    *sdl.Window
	renderer  *sdl.Renderer
	texture   *sdl.Texture
	vmem      [64 * 32]uint32
	src       sdl.Rect
	dst       sdl.Rect
}

func NewDisplay() (d *Display, err error) {
	d = &Display{
		winTitle:  "chip8",
		winWidth:  960,
		winHeight: 480,
		src:       sdl.Rect{0, 0, 640, 320},
		dst:       sdl.Rect{160, 80, 640, 320},
	}

	d.window, err = sdl.CreateWindow(d.winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, d.winWidth, d.winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return nil, err
	}

	d.renderer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return nil, err
	}

	d.texture, err = d.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STATIC, 64, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return nil, err
	}

	for i := 0; i < 64*32; i++ {
		d.vmem[i] = 0x00EEEECC
	}

	return d, nil
}

func (d *Display) Update() {
	d.renderer.Clear()
	d.texture.Update(nil, unsafe.Pointer(&d.vmem), 64*4)
	d.renderer.Copy(d.texture, &d.src, &d.dst)
	d.renderer.Present()
}

func (d *Display) Destroy() {
	defer d.window.Destroy()
	defer d.renderer.Destroy()
	defer d.texture.Destroy()
}
