package emulator

import "fmt"

// TODO: why is programCount here move programCount somewhere else or rename struc
type OpcodeInfo struct {
	programCount uint16
	opcode       uint16
	opcodeName   string
	opcodeType   string
	opcodeDesc   string
}

func (o OpcodeInfo) String() string {
	return fmt.Sprintf("PC: %d\n", o.programCount) +
		fmt.Sprintf("OPCODE: 0x%04X\n", o.opcode) +
		fmt.Sprintf("Name:     %s\n", o.opcodeName) +
		fmt.Sprintf("Type: %s\n", o.opcodeType) +
		fmt.Sprintf("Desc: %s\n", o.opcodeDesc)
}

func CreateOpcodeInfo(o uint16, n string, t string, d string, pc uint16) OpcodeInfo {
	return OpcodeInfo{programCount: pc, opcode: o, opcodeName: n, opcodeType: t, opcodeDesc: d}
}

func (o OpcodeInfo) Opcode() uint16 {
	return o.opcode
}
func (o OpcodeInfo) OpcodeName() string {
	return o.opcodeName
}
func (o OpcodeInfo) OpcodeType() string {
	return o.opcodeType
}
func (o OpcodeInfo) OpcodeDesc() string {
	return o.opcodeDesc
}

func (o OpcodeInfo) ProgramCount() uint16 {
	return o.programCount
}
