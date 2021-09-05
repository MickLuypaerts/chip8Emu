package emulator

import "fmt"

type OpcodeInfo struct {
	opcode       uint16
	opcodeName   string
	opcodeType   string
	opcodeDesc   string
	programCount uint16
}

func (o *OpcodeInfo) String() string {
	progStats := fmt.Sprintf("OPCODE: 0x%04X\n", o.opcode) +
		fmt.Sprintf("Name:     %s\n", o.opcodeName) +
		fmt.Sprintf("Type: %s\n", o.opcodeType) +
		fmt.Sprintf("Desc: %s\n", o.opcodeDesc) +
		fmt.Sprintf("PC: %d\n", o.programCount)
	return progStats
}

func CreateOpcodeInfo(o uint16, n string, t string, d string, pc uint16) OpcodeInfo {
	var oInfo OpcodeInfo
	oInfo.opcode = o
	oInfo.opcodeName = n
	oInfo.opcodeType = t
	oInfo.opcodeDesc = d
	oInfo.programCount = pc
	return oInfo
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
