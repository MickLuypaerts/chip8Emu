build:
	go build
rIBM:
	chip8Emu ".\roms\IBM Logo.ch8"

rTest:
	chip8Emu ".\roms\test_opcode.ch8"

rTetris:
	chip8Emu ".\roms\Tetris [Fran Dachille, 1991].ch8"
rPong:
	chip8Emu ".\roms\Pong (1 player).ch8"
rSpaceInvaders:
	chip8Emu ".\roms\SpaceInvaders.ch8"

rTestKey:
	chip8Emu ".\roms\KeypadTest.ch8"
