package emulator

type Control struct {
	f     func()
	usage string
}

func NewControl(f func(), u string) Control {
	return Control{f: f, usage: u}
}
