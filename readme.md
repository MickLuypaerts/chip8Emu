A Chip8 emulator in GO


# Memory
## Font
in modern CHIP-8 implementations, where the interpreter is running natively outside the 4K memory space, there is no need to avoid the lower 512 bytes of memory (0x000-0x200), and it is common to store font data there.