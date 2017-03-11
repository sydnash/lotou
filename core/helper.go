package core

func PanicWhen(b bool, s string) {
	if b {
		panic(s)
	}
}
