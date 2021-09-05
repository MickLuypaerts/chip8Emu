package emulator

type OpcodeInfo struct {
	opcodeName string
	opcodeType string
	opcodeDesc string
}

func CreateOpcodeInfo(n string, t string, d string) OpcodeInfo {
	var o OpcodeInfo
	o.opcodeName = n
	o.opcodeType = t
	o.opcodeDesc = d
	return o
}
