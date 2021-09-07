package emulator

import "fmt"

type EmulatorInfo struct {
	programCount uint16
	opcode       uint16
	opcodeName   string
	opcodeType   string
	opcodeDesc   string
}

func (o EmulatorInfo) String() string {
	return fmt.Sprintf("PC: %d\n", o.programCount) +
		fmt.Sprintf("OPCODE: 0x%04X\n", o.opcode) +
		fmt.Sprintf("Name:     %s\n", o.opcodeName) +
		fmt.Sprintf("Type: %s\n", o.opcodeType) +
		fmt.Sprintf("Desc: %s\n", o.opcodeDesc)
}

func CreateEmulatorInfo(o uint16, n string, t string, d string, pc uint16) EmulatorInfo {
	return EmulatorInfo{programCount: pc, opcode: o, opcodeName: n, opcodeType: t, opcodeDesc: d}
}

func (o EmulatorInfo) ProgramCount() uint16 {
	return o.programCount
}
