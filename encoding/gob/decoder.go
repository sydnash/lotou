package gob

type Encoder struct {
	buffer []byte
	r, w   int
}
