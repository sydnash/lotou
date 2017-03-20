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
	Src      ServiceID
	Dst      ServiceID
	Type     int32
	EncType  int32
	Id       uint64 //request id or call id
	MethodId interface{}
	Data     []interface{}
}

func NewMessage(src, dst ServiceID, msgType, encType int32, id uint64, methodId interface{}, data ...interface{}) *Message {
	switch encType {
	case MSG_ENC_TYPE_NO:
	case MSG_ENC_TYPE_GO:
		data = append([]interface{}(nil), gob.Pack(data...))
	}
	msg := &Message{src, dst, msgType, encType, id, methodId, data}
	return msg
}

func init() {
	gob.RegisterStructType(Message{})
}

func sendNoEnc(src ServiceID, dst ServiceID, msgType int32, id uint64, methodId interface{}, data ...interface{}) error {
	return lowLevelSend(src, dst, msgType, MSG_ENC_TYPE_NO, id, methodId, data...)
}

func send(src ServiceID, dst ServiceID, msgType, encType int32, id uint64, methodId interface{}, data ...interface{}) error {
	return lowLevelSend(src, dst, msgType, encType, id, methodId, data...)
}

func lowLevelSend(src, dst ServiceID, msgType, encType int32, id uint64, methodId interface{}, data ...interface{}) error {
	dsts, err := findServiceById(dst)
	isLocal := checkIsLocalId(dst)

	if err != nil && isLocal {
		return err
	}
	var msg *Message
	msg = NewMessage(src, dst, msgType, encType, id, methodId, data...)
	if err != nil {
		//doesn't find service and dstid is remote id, send a forward msg to master.
		route(Cmd_Forward, msg)
		return nil
	}
	dsts.pushMSG(msg)
	return nil
}

//send msg to dst by dst's service name
func sendName(src ServiceID, dst string, msgType int32, methodId interface{}, data ...interface{}) error {
	dsts, err := findServiceByName(dst)
	if err != nil {
		return err
	}
	return lowLevelSend(src, dsts.getId(), msgType, MSG_ENC_TYPE_GO, 0, methodId, data...)
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
		cid := m.Id
		dsts.dispatchRet(cid, m.Data...)
	}
}
func DistributeMSG(src ServiceID, methodId interface{}, data ...interface{}) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for dst, ser := range h.dic {
		if ServiceID(dst) != src {
			localSendWithoutMutex(src, ser, MSG_TYPE_DISTRIBUTE, MSG_ENC_TYPE_NO, 0, methodId, data...)
		}
	}
}

func localSendWithoutMutex(src ServiceID, dstService *service, msgType, encType int32, id uint64, methodId interface{}, data ...interface{}) {
	msg := NewMessage(src, dstService.getId(), msgType, encType, id, methodId, data...)
	dstService.pushMSG(msg)
}
