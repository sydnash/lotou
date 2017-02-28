package core

import (
	"github.com/sydnash/lotou/encoding/gob"
	"sync"
)

const (
	MSG_TYPE_NORMAL = iota
	MSG_TYPE_REQUEST
	MSG_TYPE_RESPOND
	MSG_TYPE_TIMEOUT
	MSG_TYPE_CALL
	MSG_TYPE_RET
	MSG_TYPE_CLOSE
	MSG_TYPE_SOCKET
	MSG_TYPE_ERR
	MSG_TYPE_MAX
)

const (
	MSG_ENC_TYPE_NO = iota
	MSG_ENC_TYPE_GO
)

type Message struct {
	Src     uint
	Dst     uint
	Type    int
	EncType int
	Data    []interface{}
}

func NewMessage(src, dst uint, msgType, encType int, data ...interface{}) *Message {
	msg := &Message{src, dst, msgType, encType, data}
	return msg
}

var (
	gobDecoder   *gob.Decoder
	gobEncoder   *gob.Encoder
	encoderMutex sync.Mutex
	decoderMutex sync.Mutex
)

func init() {
	gobDecoder = gob.NewDecoder()
	gobEncoder = gob.NewEncoder()
	gob.RegisterStructType(Message{})
}

func pack(data []interface{}) []byte {
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
func unpack(data []byte) interface{} {
	decoderMutex.Lock()
	defer decoderMutex.Unlock()
	gobDecoder.SetBuffer(data)
	sdata, ok := gobDecoder.Decode()
	PanicWhen(!ok)
	return sdata
}

func sendNoEnc(src *service, dst uint, msgType int, data ...interface{}) error {
	dsts, err := findServiceById(dst)
	if err != nil {
		return err
	}
	msg := NewMessage(src.getId(), dst, msgType, MSG_ENC_TYPE_NO, data...)
	dsts.pushMSG(msg)
	return nil
}

func send(src *service, dst uint, msgType int, data ...interface{}) error {
	dsts, err := findServiceById(dst)
	if err != nil {
		return err
	}
	msg := NewMessage(src.getId(), dst, msgType, MSG_ENC_TYPE_GO, pack(data))
	dsts.pushMSG(msg)
	return nil
}
