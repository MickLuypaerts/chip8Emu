build:
	go build
rIBM:
	chip8 ".\roms\IBM Logo.ch8"

rTest:
	chip8 ".\roms\test_opcode.ch8"

rTetris:
	chip8 ".\roms\Tetris [Fran Dachille, 1991].ch8"
rPong:
	chip8 ".\roms\Pong (1 player).ch8"

rTestKey:
	chip8 ".\roms\KeypadTest.ch8"