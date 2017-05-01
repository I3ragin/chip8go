package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

const PixOn, PixOff uint32 = 0x00EEEECC, 0

var vmem [64 * 32]uint32

type Display struct {
	winTitle  string
	winWidth  int
	winHeight int
	window    *sdl.Window
	renderer  *sdl.Renderer
	texture   *sdl.Texture
	//	vmem      [64 * 32]uint32
	src sdl.Rect
	dst sdl.Rect
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

	return d, nil
}

func (d *Display) Update(pixels unsafe.Pointer) {
	d.renderer.Clear()
	d.texture.Update(nil, pixels, 64*4)
	d.renderer.Copy(d.texture, &d.src, &d.dst)
	d.renderer.Present()
}

func (d *Display) Clear() {
	d.renderer.Clear()
	d.renderer.Present()
}

func (d *Display) Destroy() {
	defer d.window.Destroy()
	defer d.renderer.Destroy()
	defer d.texture.Destroy()
}

type memmory [0x1000]byte
type keyboard [0x10]byte
type stack struct {
	mem     [16]uint16
	pointer uint8
}

func (s *stack) push() {

}

func (s *stack) pop() {

}

//CHIP-8 was most commonly implemented on 4K systems, such as the Cosmac VIP and the Telmac 1800. These machines had 4096 (0x1000) memory locations, all of which are 8 bits (a byte) which is where the term CHIP-8 originated. However, the CHIP-8 interpreter itself occupies the first 512 bytes of the memory space on these machines. For this reason, most programs written for the original system begin at memory location 512 (0x200) and do not access any of the memory below the location 512 (0x200). The uppermost 256 bytes (0xF00-0xFFF) are reserved for display refresh, and the 96 bytes below that (0xEA0-0xEFF) were reserved for call stack, internal use, and other variables.

type registers struct {
	v  [16]uint8
	i  uint16 //The Index Register
	pc uint16 //Program Counter
	sp uint8  //The Stack Pointer
	st uint8  // the sound timer register
	dt uint8  // the delay timer register

}

type cpu struct {
	r registers
}

type timer struct {
}

/*
type display interface {
    bitmap [64*32]bool
}
func (c8 *chip8) drw(x, y, n uint8) {

}
*/

//CHIP8  ...
type chip8 struct {
	c cpu
	m memmory
	s stack
	k keyboard
	d *Display
	beep *mix.Music
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
	c8.c.r.pc = 0x200
	c8.c.r.i = 0x200
	c8.s.pointer = 0

	c8.d, _ = NewDisplay()

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		log.Println(err)
		return
	}
	c8.beep, _ = mix.LoadMUS("./beep.wav")

	return c8
}

func (c8 *chip8) cls() {
	//"\033[2J\033[H"
	//fmt.Print("\033[H\033[J")
	for i := 0; i < 64*32; i++ {
		vmem[i] = PixOff
	}
	c8.d.Clear()
}

func (c8 *chip8) drw(x, y, n uint8) {
	/*x = 30
	y = 2
	n = 30
	c8.c.r.i = 0*/
	var pixel int
	c8.c.r.v[0xf] = 0
	fmt.Printf("!!!y=%d,x=%d\n",y,x)
	for dy := uint16(0); dy < uint16(n); dy++ {

		bt := c8.m[c8.c.r.i+dy]

		for dx := uint(0); dx < 8; dx++ {

			if ((bt >> (8 - dx)) & 0x01) > 0 {
				//fmt.Printf("\033[%d;%dH#", y+uint8(dy), x+uint8(dx))
				index := (uint32(x+uint8(dx))+(uint32(y+uint8(dy))*64) )
				fmt.Printf("y=%d,x=%d,z=%d\n",y+uint8(dy),x+uint8(dx),index)
				if vmem[index]> 0 {
					c8.c.r.v[0xf] = 1
					pixel = 1
				}else {
					pixel = 0
				}
				pixel ^=1

				if pixel > 0 {
					vmem[index] = PixOn
				}else{
					vmem[index] = PixOff
				}
				//fmt.Printf("\033[%d;%dH ", y+uint8(i), x+uint8(n)
			}
		}
	}
	//fmt.Printf("x=%d,y=%d\n",x,y)
	c8.d.Update(unsafe.Pointer(&vmem))

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

func (c8 *chip8) keyboard() {
	var event sdl.Event
	var running bool
	kmap := map[sdl.Keycode]uint8{
		'1': 0x1,
		'2': 0x2,
		'3': 0x3,
		'c': 0xc,

		'4': 0x4,
		'5': 0x5,
		'6': 0x6,
		'd': 0xd,

		'7': 0x7,
		'8': 0x8,
		'9': 0x9,
		'e': 0xe,

		'a': 0xa,
		'0': 0x0,
		'b': 0xb,
		'f': 0xf,
	}

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyDownEvent:
				if key, ok := kmap[t.Keysym.Sym]; ok {
					//fmt.Printf("sym:%c\tstate:%d\n", t.Keysym.Sym, t.State)
					c8.k[key] = t.State
				}
			case *sdl.KeyUpEvent:
				if key, ok := kmap[t.Keysym.Sym]; ok {
					//fmt.Printf("sym:%c\tstate:%d\n", t.Keysym.Sym, t.State)
					c8.k[key] = t.State
				}
			}
		}
		sdl.Delay(16)
	}

}

func (c8 *chip8) noop() {
}

func (c8 *chip8) timer() {
	for {
		if c8.c.r.st > 0 {
			c8.c.r.st--
			//beep()
			c8.beep.Play(1);
			fmt.Println("Beep")
		}
		if c8.c.r.dt > 0 {
			c8.c.r.dt--
		}
		//60Hz
		time.Sleep(16666 * time.Microsecond)
	}
}

func (c8 *chip8) Run() {
	go c8.timer()
	go c8.keyboard()

	for c8.c.r.pc < 0x1000 {
		time.Sleep(1000 * time.Microsecond)
		cmd := uint16(c8.m[c8.c.r.pc+1]) | uint16(c8.m[c8.c.r.pc])<<8
		switch cmd & 0xF000 {
		case 0x0000:
			switch cmd & 0x00FF {
			case 0x10:
				fmt.Printf("%04x %04x\t MEGAOFF\n", c8.c.r.pc, cmd)
				c8.noop()
				c8.c.r.pc += 2
			case 0x11:
				fmt.Printf("%04x %04x\t MEGAON\n", c8.c.r.pc, cmd)
				c8.noop()
				c8.c.r.pc += 2
			case 0xE0:
				fmt.Printf("%04x %04x\t CLS\n", c8.c.r.pc, cmd)
				//00E0 - CLS
				//Clear the display
				c8.cls()
				c8.c.r.pc += 2
			case 0xEE:
				fmt.Printf("%04x %04x\t RET\n", c8.c.r.pc, cmd)
				/*00EE - RET
				  Return from a subroutine.
				  The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
				*/

				if c8.s.pointer >= 0 {
					c8.s.pointer -= 1
					c8.c.r.pc = c8.s.mem[c8.s.pointer]
				}

				c8.c.r.pc += 2
			case 0xFB:
				fmt.Printf("%04x %04x\t SCR\n", c8.c.r.pc, cmd)
				//Super Chip-48 Instructions
				c8.c.r.pc += 2
			case 0xFC:
				fmt.Printf("%04x %04x\t SCL\n", c8.c.r.pc, cmd)
				//Super Chip-48 Instructions
				c8.c.r.pc += 2
			case 0xFD:
				fmt.Printf("%04x %04x\t EXIT\n", c8.c.r.pc, cmd)
				//Super Chip-48 Instructions
				c8.c.r.pc += 2
			case 0xFE:
				fmt.Printf("%04x %04x\t LOW\n", c8.c.r.pc, cmd)
				c8.c.r.pc += 2
			case 0xFF:
				fmt.Printf("%04x %04x\t HIGH\n", c8.c.r.pc, cmd)
				//Super Chip-48 Instructions
				c8.c.r.pc += 2
			default:
				switch cmd & 0x00F0 {
				case 0xC0:
					fmt.Printf("%04x %04x\t SCD  %1x\n", c8.c.r.pc, cmd, cmd&0x000F)
					//Super Chip-48 Instructions
					c8.c.r.pc += 2
				default:
					fmt.Printf("%04x %04x\t SYS  %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
					/*
											  0nnn -  SYS addr
											  Jump to a machine code routine at nnn.
						            This instruction is only used on the old computers on which Chip-8 was originally implemented. It is ignored by modern interpreters.
					*/
					c8.noop()
					c8.c.r.pc += 2
				}

			}
			//HERE MEGA CHIP OPTCODE 02NN 03NN 04NN 05NN 06NN 07NN 08NN
		case 0x1000:
			fmt.Printf("%04x %04x\t JP   %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
			/*
							  1nnn - JP addr
				        Jump to location nnn.
								The interpreter sets the program counter to nnn.
			*/
			c8.c.r.pc = cmd & 0x0FFF
		case 0x2000:
			fmt.Printf("%04x %04x\t CALL %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
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

		case 0x3000:
			fmt.Printf("%04x %04x\t SE   V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
			/*
				3xkk - SE Vx, byte
				Skip next instruction if Vx = kk.
				The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
			*/
			x := uint8((cmd & 0x0F00) >> 8)
			kk := uint8(cmd & 0x00FF)
			if c8.c.r.v[x] == kk {
				c8.c.r.pc += 2
			}
			c8.c.r.pc += 2
		case 0x4000:
			fmt.Printf("%04x %04x\t SNE  V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
			/*
				4xkk - SNE Vx, byte
				Skip next instruction if Vx != kk.
				The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
			*/
			x := uint8((cmd & 0x0F00) >> 8)
			kk := uint8(cmd & 0x00FF)
			if c8.c.r.v[x] != kk {
				c8.c.r.pc += 2
			}
			c8.c.r.pc += 2
		case 0x5000:
			fmt.Printf("%04x %04x\t SE   V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			/*
					5xy0 - SE Vx, Vy
				  Skip next instruction if Vx = Vy.
					The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
			*/
			x := uint8((cmd & 0x0F00) >> 8)
			y := uint8((cmd & 0x00F0) >> 4)
			if c8.c.r.v[x] == c8.c.r.v[y] {
				c8.c.r.pc += 2
			}
			c8.c.r.pc += 2
		case 0x6000:
			fmt.Printf("%04x %04x\t LD   V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
			/*
				6xkk - LD Vx, byte
				Set Vx = kk.
				The interpreter puts the value kk into register Vx.
			*/
			x := uint8((cmd & 0x0F00) >> 8)
			kk := uint8(cmd & 0x00FF)
			c8.c.r.v[x] = kk
			c8.c.r.pc += 2
		case 0x7000:
			fmt.Printf("%04x %04x\t ADD  V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
			/*
				7xkk - ADD Vx, byte
				Set Vx = Vx + kk.
				Adds the value kk to the value of register Vx, then stores the result in Vx.
			*/

			x := uint8(cmd & 0x0F00 >> 8)
			kk := uint8(cmd & 0x00FF)
			fmt.Printf("V%1X = %X\n", x, c8.c.r.v[x])
			c8.c.r.v[x] += kk
			fmt.Printf("V%1X = %X\n", x, c8.c.r.v[x])
			c8.c.r.pc += 2
		case 0x8000:
			switch cmd & 0x000F {
			case 0:
				fmt.Printf("%04x %04x\t LD   V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					8xy0 - LD Vx, Vy
					Set Vx = Vy.
					Stores the value of register Vy in register Vx
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				c8.c.r.v[x] = c8.c.r.v[y]
				c8.c.r.pc += 2
			case 1:
				fmt.Printf("%04x %04x\t OR   V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					8xy1 - OR Vx, Vy
					Set Vx = Vx OR Vy.
					Performs a bitwise OR on the values of Vx and Vy, then stores the result
					in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				c8.c.r.v[x] |= c8.c.r.v[y]
				c8.c.r.pc += 2
			case 2:
				fmt.Printf("%04x %04x\t AND  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					8xy2 - AND Vx, Vy
					Set Vx = Vx AND Vy.
					Performs a bitwise AND on the values of Vx and Vy,
					then stores the result in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				c8.c.r.v[x] &= c8.c.r.v[y]
				c8.c.r.pc += 2
			case 3:
				fmt.Printf("%04x %04x\t XOR  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					  8xy3 - XOR Vx, Vy
						Set Vx = Vx XOR Vy.
						Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				c8.c.r.v[x] ^= c8.c.r.v[y]
				c8.c.r.pc += 2
			case 4:
				fmt.Printf("%04x %04x\t ADD  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					8xy4 - ADD Vx, Vy
					Set Vx = Vx + Vy, set VF = carry.
					The values of Vx and Vy are added together.
					If the result is greater than 8 bits (i.e., > 255,) VF is set to 1,
					otherwise 0. Only the lowest 8 bits of the result are kept,
					and stored in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				res := uint16(c8.c.r.v[x]) + uint16(c8.c.r.v[y])
				if res > 255 {
					c8.c.r.v[0xf] = 1
				} else {
					c8.c.r.v[0xf] = 0
				}
				c8.c.r.v[x] = uint8(res & 0x00FF)
				c8.c.r.pc += 2
			case 5:
				fmt.Printf("%04x %04x\t SUB  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					    8xy5 - SUB Vx, Vy
							Set Vx = Vx - Vy, set VF = NOT borrow.
							If Vx > Vy, then VF is set to 1, otherwise 0.
							Then Vy is subtracted from Vx, and the results stored in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				if c8.c.r.v[x] > c8.c.r.v[y] {
					c8.c.r.v[0xf] = 1
				} else {
					c8.c.r.v[0xf] = 0
				}
				c8.c.r.v[x] = c8.c.r.v[x] - c8.c.r.v[y]
				c8.c.r.pc += 2
			case 6:
				fmt.Printf("%04x %04x\t SHR  V%1X,{,V%1X}\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					8xy6 - SHR Vx {, Vy}
					Set Vx = Vx SHR 1.
					If the least-significant bit of Vx is 1, then VF is set to 1,
					otherwise 0. Then Vx is divided by 2.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.v[0xf] &= 0x01
				c8.c.r.v[x] = c8.c.r.v[x] >> 1
				c8.c.r.pc += 2
			case 7:
				fmt.Printf("%04x %04x\t SUBN V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
					Set Vx = Vy - Vx, set VF = NOT borrow.
					If Vy > Vx, then VF is set to 1, otherwise 0.
					Then Vx is subtracted from Vy, and the results stored in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				if c8.c.r.v[y] > c8.c.r.v[x] {
					c8.c.r.v[0xf] = 1
				} else {
					c8.c.r.v[0xf] = 0
				}
				c8.c.r.v[x] = c8.c.r.v[x] - c8.c.r.v[y]
				c8.c.r.pc += 2
			case 0xE:
				fmt.Printf("%04x %04x\t SHL  V%1X,{,V%1X}\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				/*
									8xyE - SHL Vx {, Vy}
									Set Vx = Vx SHL 1.
					        If the most-significant bit of Vx is 1, then VF is set to 1,
									otherwise to 0. Then Vx is multiplied by 2.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.v[0xf] &= 0x01
				c8.c.r.v[x] = c8.c.r.v[x] << 1
				c8.c.r.pc += 2
			default:
				fmt.Printf("%04x %04x\t ERROR\n", c8.c.r.pc, cmd)
				c8.c.r.pc += 2
			}
		case 0x9000:
			fmt.Printf("%04x %04x\t SNE  V%1X,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
			/*
				9xy0 - SNE Vx, Vy
				Skip next instruction if Vx != Vy.
				The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
			*/
			x := uint8((cmd & 0x0F00) >> 8)
			y := uint8((cmd & 0x00F0) >> 4)
			if c8.c.r.v[x] != c8.c.r.v[y] {
				c8.c.r.pc += 2
			}
			c8.c.r.pc += 2

		case 0xA000:
			fmt.Printf("%04x %04x\t LD   I, %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
			/*
				Annn - LD I, addr
				Set I = nnn.
				The value of register I is set to nnn.
			*/
			c8.c.r.i = cmd & 0x0FFF
			c8.c.r.pc += 2
		case 0xB000:
			fmt.Printf("%04x %04x\t JP   V0, %03x\n", c8.c.r.pc, cmd, cmd&0x0FFF)
			/*
				Bnnn - JP V0, addr
				Jump to location nnn + V0.
				The program counter is set to nnn plus the value of V0.
			*/
			c8.c.r.pc = uint16(c8.c.r.v[0]) + uint16(cmd&0x0FFF)
		case 0xC000:
			fmt.Printf("%04x %04x\t RND  V%1X,%02x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00FF)
			/*
				Cxkk - RND Vx, byte
				Set Vx = random byte AND kk.
				The interpreter generates a random number from 0 to 255,
				which is then ANDed with the value kk. The results are stored in Vx.
				See instruction 8xy2 for more information on AND.
			*/
			x := uint8((cmd & 0x0F00) >> 8)
			kk := uint8(cmd & 0x00FF)
			c8.c.r.v[x] = kk & uint8(rand.Intn(255))
			c8.c.r.pc += 2
		case 0xD000:
			switch cmd & 0x000F {
			case 0:
				fmt.Printf("%04x %04x\t DRW  V%1X, V%1X, 0\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4)
				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				c8.drw(x, y, 0)
				c8.c.r.pc += 2
			default:
				/*
					Dxyn - DRW Vx, Vy, nibble
					Display n-byte sprite starting at memory location I at (Vx, Vy),
					set VF = collision. The interpreter reads n bytes from memory, starting
					at the address stored in I. These bytes are then displayed as sprites
					on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing
					screen. If this causes any pixels to be erased, VF is set to 1,
					otherwise it is set to 0. If the sprite is positioned so part of it is
					outside the coordinates of the display, it wraps around to the opposite
					side of the screen. See instruction 8xy3 for more information on XOR,
					and section 2.4, Display, for more information on the Chip-8 screen and
					sprites.
				*/
				fmt.Printf("%04x %04x\t DRW  V%1X, V%1X, %1x\n", c8.c.r.pc, cmd, cmd&0x0F00>>8, cmd&0x00F0>>4, cmd&0x000F)

				x := uint8((cmd & 0x0F00) >> 8)
				y := uint8((cmd & 0x00F0) >> 4)
				n := uint8(cmd & 0x000F)

				c8.drw(c8.c.r.v[x], c8.c.r.v[y], n)
				c8.c.r.pc += 2
			}

		case 0xE000:
			switch cmd & 0x00FF {
			case 0x9E:
				fmt.Printf("%04x %04x\t SKP  V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Ex9E - SKP Vx
					Skip next instruction if key with the value of Vx is pressed.
					Checks the keyboard, and if the key corresponding to the value of Vx
					is currently in the down position, PC is increased by 2.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				k := c8.c.r.v[x]
				if c8.k[k] == 1 {
					c8.c.r.pc += 2
				}

				c8.c.r.pc += 2
			case 0xA1:
				fmt.Printf("%04x %04x\t SKNP V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Skip next instruction if key with the value of Vx is not pressed.
					Checks the keyboard, and if the key corresponding to the value of Vx
					is currently in the up position, PC is increased by 2.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				k := c8.c.r.v[x]
				if c8.k[k] != 1 {
					c8.c.r.pc += 2
				}

				c8.c.r.pc += 2
			default:
				fmt.Printf("%04x %04x\t ERROR\n", c8.c.r.pc, cmd)
				c8.c.r.pc += 2
			}
		case 0xF000:
			switch cmd & 0x00FF {
			case 0x07:
				fmt.Printf("%04x %04x\t LD   V%1X,DT\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx07 - LD Vx, DT
					Set Vx = delay timer value.
					The value of DT is placed into Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.v[x] = c8.c.r.dt
				c8.c.r.pc += 2
			case 0x0A:
				fmt.Printf("%04x %04x\t LD   V%1X,K\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx0A - LD Vx, K
					Wait for a key press, store the value of the key in Vx.
					All execution stops until a key is pressed,
					then the value of that key is stored in Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				//k := c8.c.r.v[x]
				for i := uint8(0); i < 15; i++ {
					if c8.k[i] == 1 {
						c8.c.r.v[x] = i
						c8.c.r.pc += 2
					}
				}
			case 0x15:
				fmt.Printf("%04x %04x\t LD   DT,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx15 - LD DT, Vx
					Set delay timer = Vx.

					DT is set equal to the value of Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.dt = c8.c.r.v[x]
				c8.c.r.pc += 2
			case 0x18:
				fmt.Printf("%04x %04x\t LD   ST,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx18 - LD ST, Vx
					Set sound timer = Vx.
					ST is set equal to the value of Vx.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.st = c8.c.r.v[x]
				c8.c.r.pc += 2
			case 0x1E:
				fmt.Printf("%04x %04x\t ADD  I,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx1E - ADD I, Vx
					Set I = I + Vx.
					The values of I and Vx are added, and the results are stored in I.
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.i += uint16(c8.c.r.v[x])
				c8.c.r.pc += 2
			case 0x29:
				fmt.Printf("%04x %04x\t LD   I,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx29 - LD F, Vx
					Set I = location of sprite for digit Vx.
					The value of I is set to the location for the hexadecimal sprite corresponding
					to the value of Vx. See section 2.4, Display, for more information on the Chip-8 hexadecimal font
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.c.r.i = uint16(c8.c.r.v[x]) * 5
				c8.c.r.pc += 2
			case 0x30:
				fmt.Printf("%04x %04x\t LD   HF,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
				 */
				c8.c.r.pc += 2
			case 0x33:
				fmt.Printf("%04x %04x\t LD   B,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx33 - LD B, Vx
					Store BCD representation of Vx in memory locations I, I+1,
					and I+2. The interpreter takes the decimal value of Vx,
					and places the hundreds digit in memory at location in I,
					the tens digit at location I+1, and the ones digit at
					location I+2.
					http://bit.ly/2oWXbjX (BCD = Binary-coded decimal)
				*/
				x := uint8((cmd & 0x0F00) >> 8)
				c8.m[c8.c.r.i] = c8.c.r.v[x] / 100
				c8.m[c8.c.r.i+1] = (c8.c.r.v[x] / 10) % 10
				c8.m[c8.c.r.i+2] = (c8.c.r.v[x] % 100) % 10
				c8.c.r.pc += 2
			case 0x55:
				fmt.Printf("%04x %04x\t LD   [I],V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					Fx55 - LD [I], Vx
					Store registers V0 through Vx in memory starting at location I.
					The interpreter copies the values of registers V0 through Vx
					into memory, starting at the address in I.
				*/
				x := uint16((cmd & 0x0F00) >> 8)
				for i := uint16(0); i <= x; i++ {
					c8.m[c8.c.r.i+i] = c8.c.r.v[i]
				}
				c8.c.r.pc += 2
			case 0x65:
				fmt.Printf("%04x %04x\t LD   V%1X,[I]\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				/*
					        Fx65 - LD Vx, [I]
									Read registers V0 through Vx from memory starting at location I.
									The interpreter reads values from memory starting at location I
									into registers V0 through Vx.
				*/
				x := uint16((cmd & 0x0F00) >> 8)
				for i := uint16(0); i <= x; i++ {
					c8.c.r.v[i] = c8.m[c8.c.r.i+i]
				}
				c8.c.r.pc += 2
			case 0x75:
				fmt.Printf("%04x %04x\t LD   R,V%1X\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				c8.c.r.pc += 2
			case 0x85:
				fmt.Printf("%04x %04x\t LD   V%1X,R\n", c8.c.r.pc, cmd, cmd&0x0F00>>8)
				c8.c.r.pc += 2

			default:
				fmt.Printf("%04x %04x\t ERROR\n", c8.c.r.pc, cmd)
				c8.c.r.pc += 2
			}
		default:
			fmt.Printf("%04x %04x\n", c8.c.r.pc, cmd)
			c8.c.r.pc += 2
		}
	}
}

func main() {
	c8 := NewChip8()
	c8.Load("./dump")
	c8.Run()

}
