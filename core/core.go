package core

import (
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"sync"
)

const (
	NATIVE_PRE     = 0xFF
	NATIVE_CORE_ID = NATIVE_PRE<<24 | 0
	REMOTE_CORE_ID = 0X00
)

type Service interface {
	Send(m *Message)
	SetId(id uint)
	Request(cb interface{}) int
	Id() uint
	Call() []interface{}
	Ret([]interface{})
}

type manager struct {
	id        uint
	nodeId    uint
	mutex     sync.RWMutex
	nameMutex sync.RWMutex
	dictory   map[uint]Service
	nameDic   map[string]uint
}

var (
	c *manager
)

func init() {
	c = new(manager)
	c.dictory = make(map[uint]Service)
	c.nameDic = make(map[string]uint)
	gob.RegisterStructType(Message{})
}

func GetService(id uint) Service {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ser, ok := c.dictory[id]
	if !ok {
		log.Warn("GetService: service %d is not exist.\n", id)
	}
	return ser
}
func RegisterService(s Service) uint {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id++
	id := selfNodeId<<24 | c.id&0xFFFFFF
	c.dictory[id] = s
	s.SetId(id)
	return id
}

func Send(dest uint, src uint, data ...interface{}) bool {
	return send(dest, src, MSG_TYPE_NORMAL, "go", data...)
}
func SendSocket(dest, src uint, data ...interface{}) bool {
	return send(dest, src, MSG_TYPE_NORMAL, "socket", data...)
}

func GetIdByName(name string) (uint, bool) {
	c.nameMutex.RLock()
	id, ok := c.nameDic[name]
	c.nameMutex.RUnlock()
	if !ok {
		return getGlobalIdByName(name)
	}
	return id, true
}

func SendName(name string, src uint, data ...interface{}) bool {
	id, ok := GetIdByName(name)
	if !ok {
		return false
	}
	return send(id, src, MSG_TYPE_NORMAL, "go", data...)
}

func Name(id uint, name string) bool {
	c.nameMutex.Lock()
	c.nameDic[name] = id
	c.nameMutex.Unlock()
	if name[0] != '.' {
		globalName(id, name)
	}
	return true
}

func Close(dest uint, src uint) bool {
	ret := send(dest, src, MSG_TYPE_CLOSE, "go")
	remove(dest)
	return ret
}

func remove(id uint) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.dictory, id)
}

func send(dest, src uint, msgType int, msgEncodeType string, data ...interface{}) bool {
	m := &Message{dest, src, msgType, msgEncodeType, data}

	isLocal := CheckIsLocalServiceId(dest)
	var ser Service
	if isLocal {
		ser = GetService(dest)
		if ser == nil {
			return false
		}
		ser.Send(m)
		return true
	}
	sendToMaster(m)
	return true
}

const (
	MSG_TYPE_NORMAL = iota
	MSG_TYPE_REQUEST
	MSG_TYPE_RESPOND
	MSG_TYPE_CALL
	MSG_TYPE_RET
	MSG_TYPE_CLOSE
)

type Message struct {
	Dest          uint
	Src           uint
	Type          int
	MsgEncodeType string
	Data          []interface{}
}

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("%s", err)
			}
		}()
		f()
	}()
}
func SafeCall(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("%s", err)
			}
		}()
		f()
	}()
}

func Request(dest uint, self Service, cb interface{}, data ...interface{}) {
	rid := self.Request(cb)
	sid := self.Id()
	param := make([]interface{}, 2)
	param[0] = rid
	param[1] = data
	send(dest, sid, MSG_TYPE_REQUEST, "go", param...)
}

func Respond(dest, src uint, rid int, data ...interface{}) {
	rid = -rid
	param := make([]interface{}, 2)
	param[0] = rid
	param[1] = data
	send(dest, src, MSG_TYPE_RESPOND, "go", param...)
}

func Call(dest uint, self Service, data ...interface{}) []interface{} {
	sid := self.Id()
	send(dest, sid, MSG_TYPE_CALL, "go", data...)
	ret := self.Call()
	return ret
}
func Ret(dest, src uint, data ...interface{}) {
	isLocal := CheckIsLocalServiceId(dest)
	var ser Service
	if isLocal {
		ser = GetService(dest)
		ser.Ret(data)
		return
	}
	m := &Message{dest, src, MSG_TYPE_RET, "go", data}
	sendToMaster(m)
	return
}

func ForwardLocal(m *Message) {
	switch m.Type {
	case MSG_TYPE_NORMAL:
		fallthrough
	case MSG_TYPE_REQUEST:
		fallthrough
	case MSG_TYPE_CLOSE:
		fallthrough
	case MSG_TYPE_CALL:
		fallthrough
	case MSG_TYPE_RESPOND:
		send(m.Dest, m.Src, m.Type, m.MsgEncodeType, m.Data...)
	case MSG_TYPE_RET:
		Ret(m.Dest, m.Src, m.Data...)
	}
}
