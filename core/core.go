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
	isLocal := CheckIsLocalServiceId(dest)
	var ser Service
	if isLocal {
		ser = GetService(dest)
		if ser == nil {
			return false
		}
	}
	m := &Message{dest, src, msgType, msgEncodeType, data}

	if isLocal {
		ser.Send(m)
		return true
	}
	sendToMaster(m)
	return true
}

const (
	MSG_TYPE_NORMAL = iota
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
