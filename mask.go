// Copyright 2015 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type memmory [0x1000]byte
type stack [16]uint16

type registers struct {
	v0 uint8
	v1 uint8
	v2 uint8
	v3 uint8
	v4 uint8
	v5 uint8
	v6 uint8
	v7 uint8
	v8 uint8
	v9 uint8
	va uint8
	vb uint8
	vc uint8
	vd uint8
	ve uint8
	vf uint8
	sp uint8
	i  uint16
	pc uint16

}

type cpu struct {
	r registers
}

type timer struct {
}

//CHIP8  ...
type chip8 struct {
	c cpu
	t timer
	m memmory
	s stack
}

func NewChip8() (c8 *chip8) {
	c8 = new(chip8)
	fonts := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0, //0
		0x20, 0x60, 0x20, 0x20, 0x70, //1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
		0x90, 0x90, 0xF0, 0x10, 0x10, //4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
		0xF0, 0x10, 0x20, 0x40, 0x40, //7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
		0xF0, 0x90, 0xF0, 0x90, 0x90, //A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
		0xF0, 0x80, 0x80, 0x80, 0xF0, //C
		0xE0, 0x90, 0x90, 0x90, 0xE0, //D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
		0xF0, 0x80, 0xF0, 0x80, 0x80} //F

	for i := 0; i < len(fonts); i++ {
		c8.m[i] = fonts[i]
	}
	return c8
}




func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	//go chip8.Run()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func layout(g *gocui.Gui) error {
	//	maxX, _ := g.Size()

	/*if v, err := g.SetView("Registers", 65, 0, 90, 32); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Code"
		fmt.Fprintln(v, "^a: Set mask")
		fmt.Fprintln(v, "^c: Exit")
	}
	*/
	r := Registers{123, 4, 25, 37, 67, 56, 2, 23, 45, 145, 6, 12, 13, 45, 64, 56, 94, 823, 842}
	if v, err := g.SetView("Registers", 65, 0, 90, 20); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Registers"
		fmt.Fprintf(v, "  v0 = 0x%02x\n", r.v0)
		fmt.Fprintf(v, "  v1 = 0x%02x\n", r.v1)
		fmt.Fprintf(v, "  v2 = 0x%02x\n", r.v2)
		fmt.Fprintf(v, "  v3 = 0x%02x\n", r.v3)
		fmt.Fprintf(v, "  v4 = 0x%02x\n", r.v4)
		fmt.Fprintf(v, "  v5 = 0x%02x\n", r.v5)
		fmt.Fprintf(v, "  v6 = 0x%02x\n", r.v6)
		fmt.Fprintf(v, "  v7 = 0x%02x\n", r.v7)
		fmt.Fprintf(v, "  v8 = 0x%02x\n", r.v8)
		fmt.Fprintf(v, "  v9 = 0x%02x\n", r.v9)
		fmt.Fprintf(v, "  va = 0x%02x\n", r.va)
		fmt.Fprintf(v, "  vb = 0x%02x\n", r.vb)
		fmt.Fprintf(v, "  vc = 0x%02x\n", r.vc)
		fmt.Fprintf(v, "  vd = 0x%02x\n", r.vd)
		fmt.Fprintf(v, "  ve = 0x%02x\n", r.ve)
		fmt.Fprintf(v, "  vf = 0x%02x\n", r.vf)
		fmt.Fprintf(v, "  sp = 0x%02x\n", r.sp)
		fmt.Fprintf(v, "  i  = 0x%04x\n", r.i)
		fmt.Fprintf(v, "  pc = 0x%04x\n", r.pc)

	}

	if v, err := g.SetView("Code", 65, 21, 90, 42); err !jira.inn.ru
		}
		v.Title = "Code"
		fmt.Fprintln(v, "^a: Set mask")
		fmt.Fprintln(v, "^c: Exit")
	}

	if v, err := g.SetView("Display", 0, 0, 64, 32); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Display"
		/*	if _, err := g.SetCurrentView("Display"); err != nil {
				return err
			}
			v.Editable = true
			v.Wrap = true
		*/
	}

	if v, err := g.SetView("Stack", 0, 33, 64, 42); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Stack"
		fmt.Fprintln(v, "^a: Set mask")
		fmt.Fprintln(v, "^c: Exit")
	}

	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyCtrlA, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.Mask ^= '*'
			return nil
		}); err != nil {
		return err
	}
	return nil
}
