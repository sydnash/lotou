package core

import (
	"fmt"
	"sync"
)

type Service interface {
	Send(m *Message)
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
		fmt.Printf("service %d is not exist.\n", id)
	}
	return ser
}
func RegisterService(s Service) uint {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id++
	c.dictory[c.id] = s
	return c.id
}

func Send(dest uint, src uint, data ...interface{}) bool {
	return send(dest, src, MSG_TYPE_NORMAL, data...)
}

func Close(dest uint, src uint) bool {
	ret := send(dest, src, MSG_TYPE_CLOSE)
	close(dest)
	return ret
}

func close(id uint) {
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
