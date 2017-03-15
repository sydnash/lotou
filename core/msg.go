package core

import (
	"github.com/sydnash/lotou/encoding/gob"
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
	MSG_TYPE_DISTRIBUTE
	MSG_TYPE_MAX
)

const (
	MSG_ENC_TYPE_NO = iota
	MSG_ENC_TYPE_GO
)

//Message is the based struct of msg through all service
//by convention, the first value of Data is a string as the method name
type Message struct {
	Src     uint64
	Dst     uint64
	Type    int32
	EncType int32
	Data    []interface{}
}

func NewMessage(src, dst ServiceID, msgType, encType int32, data ...interface{}) *Message {
	msg := &Message{uint64(src), uint64(dst), msgType, encType, data}
	return msg
}

func init() {
	gob.RegisterStructType(Message{})
}

func sendNoEnc(src ServiceID, dst ServiceID, msgType int32, data ...interface{}) error {
	return rawSend(false, src, dst, msgType, data...)
}

func send(src ServiceID, dst ServiceID, msgType int32, data ...interface{}) error {
	return rawSend(true, src, dst, msgType, data...)
}

func rawSend(isNeedEnc bool, src, dst ServiceID, msgType int32, data ...interface{}) error {
	dsts, err := findServiceById(dst)
	isLocal := checkIsLocalId(dst)

	if err != nil && isLocal {
		return err
	}
	var msg *Message
	if isNeedEnc {
		msg = NewMessage(src, dst, msgType, MSG_ENC_TYPE_GO, gob.Pack(data))
	} else {
		msg = NewMessage(src, dst, msgType, MSG_ENC_TYPE_NO, data...)
	}
	if err != nil {
		//didn't find service and dst id is remote id, send a forward msg to master.
		sendToMaster("forward", msg)
		return nil
	}
	dsts.pushMSG(msg)
	return nil
}

//send msg to dst by dst's service name
func sendName(src ServiceID, dst string, msgType int32, data ...interface{}) error {
	dsts, err := findServiceByName(dst)
	if err != nil {
		return err
	}
	return rawSend(true, src, dsts.getId(), msgType, data...)
}

func ForwardLocal(m *Message) {
	dsts, err := findServiceById(ServiceID(m.Dst))
	if err != nil {
		return
	}
	switch m.Type {
	case MSG_TYPE_NORMAL, MSG_TYPE_REQUEST, MSG_TYPE_RESPOND, MSG_TYPE_CALL, MSG_TYPE_DISTRIBUTE:
		dsts.pushMSG(m)
	case MSG_TYPE_RET:
		if m.EncType == MSG_ENC_TYPE_GO {
			t := gob.Unpack(m.Data[0].([]byte))
			m.Data = t.([]interface{})
		}
		cid := m.Data[0].(int)
		data := m.Data[1].([]interface{})
		dsts.dispatchRet(cid, data...)
	}
}
func DistributeMSG(src ServiceID, data ...interface{}) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for dst, ser := range h.dic {
		if ServiceID(dst) != src {
			localSendWithNoMutex(src, ser, MSG_TYPE_DISTRIBUTE, MSG_ENC_TYPE_NO, data)
		}
	}
}

func localSendWithNoMutex(src ServiceID, dstService *service, msgType, encType int32, data ...interface{}) {
	msg := NewMessage(src, dstService.getId(), msgType, encType, data...)
	dstService.pushMSG(msg)
}
