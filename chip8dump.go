package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type opcode [2]byte

func decode(point, cmd uint16) {
	//fmt.Printf("%04x %04x\n", point, cmd & 0xF000)
	switch cmd & 0xF000 {
	case 0x0000:
		switch cmd & 0x00FF {
		case 0x10:
			fmt.Printf("%04x %04x\t MEGAOFF\n", point, cmd)
		case 0x11:
			fmt.Printf("%04x %04x\t MEGAON\n", point, cmd)
		case 0xE0:
			fmt.Printf("%04x %04x\t CLS\n", point, cmd)
		case 0xEE:
			fmt.Printf("%04x %04x\t RET\n", point, cmd)
		case 0xFB:
			fmt.Printf("%04x %04x\t SCR\n", point, cmd)
		case 0xFC:
			fmt.Printf("%04x %04x\t SCL\n", point, cmd)
		case 0xFD:
			fmt.Printf("%04x %04x\t EXIT\n", point, cmd)
		case 0xFE:
			fmt.Printf("%04x %04x\t LOW\n", point, cmd)
		case 0xFF:
			fmt.Printf("%04x %04x\t HIGH\n", point, cmd)
		default:
			switch cmd & 0x00F0 {
			case 0xC0:
				fmt.Printf("%04x %04x\t SCD  %1x\n", point, cmd, cmd&0x000F)
			default:
				fmt.Printf("%04x %04x\t SYS  %03x\n", point, cmd, cmd&0x0FFF)
			}

		}
		//HERE MEGA CHIP OPTCODE 02NN 03NN 04NN 05NN 06NN 07NN 08NN
	case 0x1000:
		fmt.Printf("%04x %04x\t JP   %03x\n", point, cmd, cmd&0x0FFF)
	case 0x2000:
		fmt.Printf("%04x %04x\t CALL %03x\n", point, cmd, cmd&0x0FFF)
	case 0x3000:
		fmt.Printf("%04x %04x\t SE   V%1X,%02x\n", point, cmd, cmd&0x0F00>>8, cmd&0x00FF)
	case 0x4000:
		fmt.Printf("%04x %04x\t SNE  V%1X,%02x\n", point, cmd, cmd&0x0F00>>8, cmd&0x00FF)
	case 0x5000:
		fmt.Printf("%04x %04x\t SE   V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
	case 0x6000:
		fmt.Printf("%04x %04x\t LD   V%1X,%02x\n", point, cmd, cmd&0x0F00>>8, cmd&0x00FF)
	case 0x7000:
		fmt.Printf("%04x %04x\t ADD  V%1X,%02x\n", point, cmd, cmd&0x0F00>>8, cmd&0x00FF)
	case 0x8000:
		switch cmd & 0x000F {
		case 0:
			fmt.Printf("%04x %04x\t LD   V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 1:
			fmt.Printf("%04x %04x\t OR   V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 2:
			fmt.Printf("%04x %04x\t AND  V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 3:
			fmt.Printf("%04x %04x\t XOR  V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 4:
			fmt.Printf("%04x %04x\t ADD  V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 5:
			fmt.Printf("%04x %04x\t SUB  V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 6:
			fmt.Printf("%04x %04x\t SHR  V%1X,{,V%1X}\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 7:
			fmt.Printf("%04x %04x\t SUBN V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		case 0xE:
			fmt.Printf("%04x %04x\t SHL  V%1X,{,V%1X}\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		default:
			fmt.Printf("%04x %04x\t ERROR\n", point, cmd)
		}
	case 0x9000:
		fmt.Printf("%04x %04x\t SNE  V%1X,V%1X\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)

	case 0xA000:
		fmt.Printf("%04x %04x\t LD   I, %03x\n", point, cmd, cmd&0x0FFF)
	case 0xB000:
		fmt.Printf("%04x %04x\t JP   V0, %03x\n", point, cmd, cmd&0x0FFF)
	case 0xC000:
		fmt.Printf("%04x %04x\t RND  V%1X,%02x\n", point, cmd, cmd&0x0F00>>8, cmd&0x00FF)
	case 0xD000:
		switch cmd & 0x000F {
		case 0:
			fmt.Printf("%04x %04x\t DRW  V%1X, V%1X, 0\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
		default:
			fmt.Printf("%04x %04x\t DRW  V%1X, V%1X, %1x\n", point, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4, cmd&0x000F)

		}

	case 0xE000:
		switch cmd & 0x00FF {
		case 0x9E:
			fmt.Printf("%04x %04x\t SKP  V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0xA1:
			fmt.Printf("%04x %04x\t SKNP V%1X\n", point, cmd, cmd&0x0F00>>8)
		default:
			fmt.Printf("%04x %04x\t ERROR\n", point, cmd)
		}
	case 0xF000:
		switch cmd & 0x00FF {
		case 0x07:
			fmt.Printf("%04x %04x\t LD   V%1X,DT\n", point, cmd, cmd&0x0F00>>8)
		case 0x0A:
			fmt.Printf("%04x %04x\t LD   V%1X,K\n", point, cmd, cmd&0x0F00>>8)
		case 0x15:
			fmt.Printf("%04x %04x\t LD   DT,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x18:
			fmt.Printf("%04x %04x\t LD   ST,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x1E:
			fmt.Printf("%04x %04x\t ADD  I,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x29:
			fmt.Printf("%04x %04x\t LD   I,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x30:
			fmt.Printf("%04x %04x\t LD   HF,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x33:
			fmt.Printf("%04x %04x\t LD   B,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x55:
			fmt.Printf("%04x %04x\t LD   [I],V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x65:
			fmt.Printf("%04x %04x\t LD   V%1X,[I]\n", point, cmd, cmd&0x0F00>>8)
		case 0x75:
			fmt.Printf("%04x %04x\t LD   R,V%1X\n", point, cmd, cmd&0x0F00>>8)
		case 0x85:
			fmt.Printf("%04x %04x\t LD   V%1X,R\n", point, cmd, cmd&0x0F00>>8)

		default:
			fmt.Printf("%04x %04x\t ERROR\n", point, cmd)
		}
	default:
		fmt.Printf("%04x %04x\n", point, cmd)
	}

}

func objdump(b []byte) {
	for i := 0; i < len(b); i = i + 2 {
		decode(uint16(0x200+i), uint16(b[i+1])|uint16(b[i])<<8)
	}
}

func main() {
	dump, err := ioutil.ReadFile("./dump")
	if err != nil {
		log.Fatal(err)
	}
	objdump(dump)
	fmt.Println("Ok")
}
