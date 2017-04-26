package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type memmory [0x1000]byte
type stack struct {
	mem     [16]uint16
	pointer uint8
}

//CHIP-8 was most commonly implemented on 4K systems, such as the Cosmac VIP and the Telmac 1800. These machines had 4096 (0x1000) memory locations, all of which are 8 bits (a byte) which is where the term CHIP-8 originated. However, the CHIP-8 interpreter itself occupies the first 512 bytes of the memory space on these machines. For this reason, most programs written for the original system begin at memory location 512 (0x200) and do not access any of the memory below the location 512 (0x200). The uppermost 256 bytes (0xF00-0xFFF) are reserved for display refresh, and the 96 bytes below that (0xEA0-0xEFF) were reserved for call stack, internal use, and other variables.

type registers struct {
	v  [16]uint8
	i  uint16 //The Index Register
	pc uint16 //Program Counter
	sp uint8  //The Stack Pointer
	st timer  // the sound timer register
	dt timer  // the delay timer register

}

type cpu struct {
	r registers
}

type timer struct {
}

//CHIP8  ...
type chip8 struct {
	c cpu
	m memmory
	s stack
}

func NewChip8() (c8 *chip8) {
	c8 = new(chip8)

	fonts := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0, //0
		0x20, 0x60, 0x20, 0x20, 0x70, //1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
		0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
		0xF0, 0x80, 0x80, 0x80, 0xF0, //C
		0xE0, 0x90, 0x90, 0x90, 0xE0, //D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
		0xF0, 0x80, 0xF0, 0x80, 0x80} //F

	for i := 0; i < len(fonts); i++ {
		c8.m[i] = fonts[i]
	}
	c8.c.r.pc = 0x200
	c8.c.r.i = 0x200
	c8.s.pointer = 0
	return c8
}

func (c8 *chip8) cls() {
	//display.clear()
}

func (c8 *chip8) DRW(x, y, n uint8) {
	//display.draw_sprite()
}

func (c8 *chip8) Load(frimware string) error {

	dump, err := ioutil.ReadFile(frimware)
	if err != nil {
		log.Println(err)
		return err
	}

	for i := 0; i < len(dump); i++ {
		c8.m[0x200+i] = dump[i]
	}
	return nil
}

func (c8 *chip8) noop() {
}

func (c8 *chip8) Run() {

	for ; c8.c.r.pc < 0x1000; c8.c.r.pc += 2 {

		cmd := uint16(c8.m[c8.c.r.pc+1]) | uint16(c8.m[c8.c.r.pc])<<8
		switch cmd & 0xF000 {
		case 0x0000:
			switch cmd & 0x00FF {
			case 0x10:
				fmt.Printf("%04x %04x\t MEGAOFF\n", c8.c.r.pc, cmd)
				c8.noop()
			case 0x11:
				fmt.Printf("%04x %04x\t MEGAON\n", c8.c.r.pc, cmd)
			case 0xE0:
				//00E0 - CLS
				//Clear the display
				fmt.Printf("%04x %04x\t CLS\n", c8.c.r.pc, cmd)
			case 0xEE:
				/*00EE - RET
				  Return from a subroutine.
				  The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
				*/
				fmt.Printf("%04x %04x\t RET\n", c8.c.r.pc, cmd)
			case 0xFB:
				fmt.Printf("%04x %04x\t SCR\n", c8.c.r.pc, cmd)
			case 0xFC:
				fmt.Printf("%04x %04x\t SCL\n", c8.c.r.pc, cmd)
			case 0xFD:
				fmt.Printf("%04x %04x\t EXIT\n", c8.c.r.pc, cmd)
			case 0xFE:
				fmt.Printf("%04x %04x\t LOW\n", c8.c.r.pc, cmd)
			case 0xFF:
				fmt.Printf("%04x %04x\t HIGH\n", c8.c.r.pc, cmd)
			default:
				switch cmd & 0x00F0 {
				case 0xC0:
					fmt.Printf("%04x %04x\t SCD  %1x\n", c8.c.r.pc, cmd, cmd&0x000F)
				default:
					/*
											  0nnn -  SYS addr
											  Jump to a machine code routine at nnn.
						            This instruction is only used on the old computers on which Chip-8 was originally implemented. It is ignored by modern interpreters.
					*/
					fmt.Printf("%04x %04x\t SYS  %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
				}

			}
			//HERE MEGA CHIP OPTCODE 02NN 03NN 04NN 05NN 06NN 07NN 08NN
		case 0x1000:
			/*
							  1nnn - JP addr
				        Jump to location nnn.
								The interpreter sets the program counter to nnn.
			*/
			c8.c.r.pc = cmd & 0x0FFF
			//fmt.Printf("%04x %04x\t JP   %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
		case 0x2000:
			/*
				2nnn - CALL addr
				Call subroutine at nnn.
				The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
			*/
			if c8.s.pointer < uint8(len(c8.s.mem)) {
				c8.s.mem[c8.s.pointer] = c8.c.r.pc
				c8.c.r.pc = cmd & 0x0FFF
				c8.s.pointer += 1
			}
			//fmt.Printf("%04x %04x\t CALL %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
		case 0x3000:
			/*
				3xkk - SE Vx, byte
				Skip next instruction if Vx = kk.
				The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
			*/
			i := uint8(cmd & 0x0F00 >> 4)
			b := uint8(cmd & 0x00FF)
			if c8.c.r.v[i] == b {
				c8.c.r.pc += 2
			}
			//fmt.Printf("%04x %04x\t SE   V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
		case 0x4000:
			/*
				4xkk - SNE Vx, byte
				Skip next instruction if Vx != kk.
				The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
			*/
			i := uint8(cmd & 0x0F00 >> 8)
			b := uint8(cmd & 0x00FF)
			if c8.c.r.v[i] != b {
				c8.c.r.pc += 2
			}
			//fmt.Printf("%04x %04x\t SNE  V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
		case 0x5000:
			/*
							5xy0 - SE Vx, Vy
							if c8.c.r.v[i] == c8.c.r.v[b] { c8.c.r.pc += 2 }
							//fmt.Printf("%04x %04x\t SE   V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
						case 0x6000:
							/*
							6xkk - LD Vx, byte
							Set Vx = kk.Cxkk - RND Vx, byte
				Set Vx = random byte AND kk.

				The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.

							The interpreter puts the value kk into register Vx.
			*/
			i := uint8(cmd & 0x0F00 >> 8)
			b := uint8(cmd & 0x00FF)
			c8.c.r.v[i] = b
			//fmt.Printf("%04x %04x\t LD   V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
		case 0x7000:
			/*
				7xkk - ADD Vx, byte
				Set Vx = Vx + kk.
				Adds the value kk to the value of register Vx, then stores the result in Vx.
			*/
			i := uint8(cmd & 0x0F00 >> 8)
			b := uint8(cmd & 0x00FF)
			c8.c.r.v[i] += b
			//fmt.Printf("%04x %04x\t ADD  V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
		case 0x8000:
			switch cmd & 0x000F {
			case 0:
				/*
					8xy0 - LD Vx, Vy
					Set Vx = Vy.
					Stores the value of register Vy in register Vx
				*/
				i := cmd & 0x0F00
				b := cmd & 0x00F0
				c8.c.r.v[i] = c8.c.r.v[b]
				//fmt.Printf("%04x %04x\t LD   V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 1:
				/*
					8xy1 - OR Vx, Vy
					Set Vx = Vx OR Vy.
					Performs a bitwise OR on the values of Vx and Vy, then stores the result
					in Vx.
				*/
				i := cmd & 0x0F00
				b := cmd & 0x00F0
				c8.c.r.v[i] |= c8.c.r.v[b]
				//fmt.Printf("%04x %04x\t OR   V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 2:
				/*
					8xy2 - AND Vx, Vy
					Set Vx = Vx AND Vy.
					Performs a bitwise AND on the values of Vx and Vy,
					then stores the result in Vx.
				*/
				i := cmd & 0x0F00
				b := cmd & 0x00F0
				c8.c.r.v[i] &= c8.c.r.v[b]
				//fmt.Printf("%04x %04x\t AND  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 3:
				/*
					Set Vx = Vx XOR Vy.
					Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
				*/
				i := cmd & 0x0F00
				b := cmd & 0x00F0
				c8.c.r.v[i] ^= c8.c.r.v[b]
				//fmt.Printf("%04x %04x\t XOR  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 4:
				/*
					8xy4 - ADD Vx, Vy
					Set Vx = Vx + Vy, set VF = carry.
					The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
				*/
				i := cmd & 0x0F00
				b := cmd & 0x00F0
				res := uint16(c8.c.r.v[i]) + uint16(c8.c.r.v[b])
				if res > 255 {
					c8.c.r.v[0xf] = 1
				} else {
					c8.c.r.v[0xf] = 0
				}
				c8.c.r.v[i] = uint8(res & 0x00FF)
				//fmt.Printf("%04x %04x\t ADD  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 5:
				/*
						Set Vx = Vx - Vy, set VF = NOT borrow.
						If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.

						i := cmd & 0x0F00
					  b := cmd & 0x00F0
						res := uint16(c8.c.r.v[i]) +  uint16(c8.c.r.v[b])
						if res > 255 {
							c8.c.r.vf = 1
						}
						else{
							c8.c.r.vf = 0
						}
						c8.c.r.v[i] = uint8(res & 0x00FF)*/
				//fmt.Printf("%04x %04x\t SUB  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 6:
				fmt.Printf("%04x %04x\t SHR  V%1X,{,V%1X}\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 7:
				fmt.Printf("%04x %04x\t SUBN V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			case 0xE:
				fmt.Printf("%04x %04x\t SHL  V%1X,{,V%1X}\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			default:
				fmt.Printf("%04x %04x\t ERROR\n", c8.c.r.pc, cmd)
			}
		case 0x9000:
			/*
				9xy0 - SNE Vx, Vy
				Skip next instruction if Vx != Vy.
				The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
			*/
			i := cmd & 0x0F00
			b := cmd & 0x00FF
			if c8.c.r.v[i] != c8.c.r.v[b] {
				c8.c.r.pc += 2
			}
			//fmt.Printf("%04x %04x\t SNE  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)

		case 0xA000:
			/*
				Annn - LD I, addr
				Set I = nnn.
				The value of register I is set to nnn.
			*/
			i := cmd & 0x0FFF
			c8.c.r.i = i
			//fmt.Printf("%04x %04x\t LD   I, %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
		case 0xB000:
			/*
				Bnnn - JP V0, addr
				Jump to location nnn + V0.
				The program counter is set to nnn plus the value of V0.
			*/
			i := uint16(cmd & 0x0FFF)
			c8.c.r.pc = uint16(c8.c.r.v[0]) + i
			//fmt.Printf("%04x %04x\t JP   V0, %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
		case 0xC000:
			/*
				Cxkk - RND Vx, byte
				Set Vx = random byte AND kk.
				The interpreter generates a random number from 0 to 255,
				which is then ANDed with the value kk. The results are stored in Vx.
				See instruction 8xy2 for more information on AND.
			*/
			//fmt.Printf("%04x %04x\t RND  V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
		case 0xD000:
			switch cmd & 0x000F {
			case 0:
				fmt.Printf("%04x %04x\t DRW  V%1X, V%1X, 0\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			default:
				fmt.Printf("%04x %04x\t DRW  V%1X, V%1X, %1x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4, cmd&0x000F)

			}

		case 0xE000:
			switch cmd & 0x00FF {
			case 0x9E:
				fmt.Printf("%04x %04x\t SKP  V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0xA1:
				fmt.Printf("%04x %04x\t SKNP V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			default:
				fmt.Printf("%04x %04x\t ERROR\n", c8.c.r.pc, cmd)
			}
		case 0xF000:
			switch cmd & 0x00FF {
			case 0x07:
				fmt.Printf("%04x %04x\t LD   V%1X,DT\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x0A:
				fmt.Printf("%04x %04x\t LD   V%1X,K\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x15:
				fmt.Printf("%04x %04x\t LD   DT,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x18:
				fmt.Printf("%04x %04x\t LD   ST,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x1E:
				fmt.Printf("%04x %04x\t ADD  I,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x29:
				fmt.Printf("%04x %04x\t LD   I,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x30:
				fmt.Printf("%04x %04x\t LD   HF,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x33:
				fmt.Printf("%04x %04x\t LD   B,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x55:
				fmt.Printf("%04x %04x\t LD   [I],V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x65:
				fmt.Printf("%04x %04x\t LD   V%1X,[I]\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x75:
				fmt.Printf("%04x %04x\t LD   R,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
			case 0x85:
				fmt.Printf("%04x %04x\t LD   V%1X,R\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)

			default:
				fmt.Printf("%04x %04x\t ERROR\n", c8.c.r.pc, cmd)
			}
		default:
			fmt.Printf("%04x %04x\n", c8.c.r.pc, cmd)
		}
	}
}

func main() {
	c8 := NewChip8()
	c8.Load("./dump")
	c8.Run()

}
