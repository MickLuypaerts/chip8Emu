A Chip8 emulator in GO

# Chip8
## Memory
### Font
in modern CHIP-8 implementations, where the interpreter is running natively outside the 4K memory space, there is no need to avoid the lower 512 bytes of memory (0x000-0x200), and it is common to store font data there.

### stack
The stack is an array of 16 16-bit values, used to store the address that the interpreter shoud return to when finished with a subroutine. Chip-8 allows for up to 16 levels of nested subroutines.

### Timers 
Timers
CHIP-8 has two timers. They both count down at 60 hertz, until they reach 0.

Delay timer: This timer is intended to be used for timing the events of games. Its value can be set and read.
Sound timer: This timer is used for sound effects. When its value is nonzero, a beeping sound is made.



# TODO
[X] Screen  
[X] Input  
[ ] Fix buggy input (reseting keyboard)  
[ ] Add rest of the controls  
[ ] Sound  
[ ] Opcodes  
[X] Stepping Through the program  
[X] Running the program  
[ ] let user choose Hertz 
[ ] Fix weird passing functions between Chip8 and TUI  
[ ] Fix Bug in decode.go where we assign the PC for the TUI with setOpcodeInfo() but don't update it when we Skip  
[X] Print controls in usage  
[ ] Clean up Print usage

## Opcodes
[X] 0NNN  
[X] 00E0  
[X] 00EE  
[X] 1NNN  
[X] 2NNN  
[X] 3XNN  
[X] 4XNN  
[X] 5XY0  
[X] 6XNN  
[X] 7XNN  
[X] 8XY0  
[X] 8XY1  
[X] 8XY2  
[X] 8XY3  
[X] 8XY4  
[X] 8XY5  
[X] 8XY6  
[X] 8XY7  
[X] 8XYE  
[X] 9XY0  
[X] ANNN  
[X] BNNN  
[X] CXNN  
[X] DXYN  
[X] EX9E  
[X] EXA1  
[X] FX07  
[ ] FX0A  
[X] FX15  
[X] FX18  
[X] FX1E  
[X] FX29  
[X] FX33  
[X] FX55  
[X] FX65  


1 Ticker for the clock cycles => every tick chip8.EmulateCycle()
non blocking channel for keyboard input => select + default
blocking channel for keyboard opcode: FX0A



# packages
https://github.com/gizak/termui/

# Links
[Roms](https://github.com/kripod/chip8-roms)
[audio](https://dev.to/ik5/quick-and-dirty-audio-playing-in-golang-3n7c)
[azul3 keyboard package](https://github.com/azul3d/engine/tree/master/keyboard)
[azul3 keyboard package docs](https://pkg.go.dev/azul3d.org/engine/keyboard?utm_source=godoc)