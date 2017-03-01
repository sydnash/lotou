package gob

import (
	"sync"
)

var (
	gobDecoder   *Decoder
	gobEncoder   *Encoder
	encoderMutex sync.Mutex
	decoderMutex sync.Mutex
)

func init() {
	gobDecoder = NewDecoder()
	gobEncoder = NewEncoder()
}

func Pack(data []interface{}) []byte {
	encoderMutex.Lock()
	defer encoderMutex.Unlock()
	gobEncoder.Reset()
	gobEncoder.Encode(data)
	gobEncoder.UpdateLen()
	buf := gobEncoder.Buffer()
	ret := make([]byte, len(buf))
	copy(ret, buf)
	return ret
}
func Unpack(data []byte) interface{} {
	decoderMutex.Lock()
	defer decoderMutex.Unlock()
	gobDecoder.SetBuffer(data)
	sdata, ok := gobDecoder.Decode()
	if !ok {
		panic("gob unpack failed")
	}
	return sdata
}
