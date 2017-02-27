package core

func PanicWhen(b bool) {
	if b {
		panic("")
	}
}
