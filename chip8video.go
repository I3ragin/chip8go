package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle = "Text"
var imageName = "./favicon.bmp"
var winWidth, winHeight int = 960, 480

func main() {
	var window *sdl.Window
	var err error
	var renderer *sdl.Renderer
	var image *sdl.Surface
	var texture *sdl.Texture
	//var src, dst sdl.Rect
	var dst, dst2, dst3, dst4 sdl.Rect
	/*	var points []sdl.Point
		var rect sdl.Rect
		var rects []sdl.Rect
	*/
	if window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(1)
	}
	defer renderer.Destroy()
	image, err = sdl.LoadBMP(imageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s\n", err)
		os.Exit(3)
	}
	defer image.Free()

	texture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(4)
	}
	defer texture.Destroy()

	//src = sdl.Rect{0, 0, 32, 32}
	dst = sdl.Rect{100, 250, 15, 15}
	dst2 = sdl.Rect{135, 250, 15, 15}
	dst3 = sdl.Rect{170, 250, 15, 15}
	dst4 = sdl.Rect{205, 250, 15, 15}

	renderer.Clear()
	renderer.Copy(texture, nil, &dst)
	renderer.Copy(texture, nil, &dst2)
	renderer.Copy(texture, nil, &dst3)
	renderer.Copy(texture, nil, &dst4)
	//renderer.Copy(texture, &src, &dst)
	renderer.Present()

	sdl.Delay(2000)

}
