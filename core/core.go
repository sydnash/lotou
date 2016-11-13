package core

import (
	"github.com/sydnash/majiang/log"
	"sync"
)

type Service interface {
	Send(m *Message)
	SetId(id uint)
}

type manager struct {
	id      uint
	mutex   sync.Mutex
	dictory map[uint]Service
}

var c *manager

func init() {
	c = new(manager)
	c.dictory = make(map[uint]Service)
}

func GetService(id uint) Service {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ser, ok := c.dictory[id]
	if !ok {
		log.Info("service %d is not exist.", id)
	}
	return ser
}
func RegisterService(s Service) uint {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id++
	c.dictory[c.id] = s
	s.SetId(c.id)
	return c.id
}

func Send(dest uint, src uint, data ...interface{}) bool {
	return send(dest, src, MSG_TYPE_NORMAL, data...)
}

func Close(dest uint, src uint) bool {
	ret := send(dest, src, MSG_TYPE_CLOSE)
	remove(dest)
	return ret
}

func remove(id uint) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.dictory, id)
}

func send(dest, src uint, msgType int, data ...interface{}) bool {
	ser := GetService(dest)
	if ser == nil {
		return false
	}
	m := &Message{dest, src, msgType, data}
	ser.Send(m)
	return true
}

const (
	MSG_TYPE_NORMAL = iota
	MSG_TYPE_CLOSE
)

type Message struct {
	Dest uint
	Src  uint
	Type int
	Data []interface{}
}
