A Chip8 emulator in GO


# Memory
## Font
in modern CHIP-8 implementations, where the interpreter is running natively outside the 4K memory space, there is no need to avoid the lower 512 bytes of memory (0x000-0x200), and it is common to store font data there.

# stack
The stack is an array of 16 16-bit values, used to store the address that the interpreter shoud return to when finished with a subroutine. Chip-8 allows for up to 16 levels of nested subroutines.

# Timers 
Timers
CHIP-8 has two timers. They both count down at 60 hertz, until they reach 0.

Delay timer: This timer is intended to be used for timing the events of games. Its value can be set and read.
Sound timer: This timer is used for sound effects. When its value is nonzero, a beeping sound is made.


# Opcodes
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
[ ] BNNN  
[ ] CXNN  
[X] DXYN  
[ ] EX9E  
[ ] EXA1  
[ ] FX07  
[ ] FX0A  
[ ] FX15  
[ ] FX18  
[ ] FX1E  
[ ] FX29  
[ ] FX33  
[ ] FX55  
[ ] FX65  

# packages
https://github.com/gizak/termui/