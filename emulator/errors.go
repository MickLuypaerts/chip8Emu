package emulator

type DoubleKeyAssigmentError struct {
	Key string
}

func (e DoubleKeyAssigmentError) Error() string {
	return "double assignment of key: " + e.Key
}
