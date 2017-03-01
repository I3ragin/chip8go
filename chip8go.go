package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type memmory [4096]byte

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
	i  uint16
}

type cpu struct {
	r registers
	m memmory
}

type timer struct {
}

//CHIP8  ...
type CHIP8 struct {
	c cpu
	t timer
}

/*
func (chip *CHIP8) Read(b []byte) (n int, err error){

}

func (chip *CHIP8) Load(frimware string) err error{

}
*/

func main() {
	dump, err := ioutil.ReadFile("./dump")
	if err != nil {
		log.Fatal(err)
	}
	objdump(dump)
	fmt.Println("Ok")

}
