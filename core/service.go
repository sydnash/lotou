package core

import (
	"github.com/sydnash/lotou/log"
	"reflect"
	"runtime/debug"
)

type callback struct {
	hasSelf bool
	cb      reflect.Value
}
type Base struct {
	id                  uint
	in                  chan *Message
	self                reflect.Value //a pointer a concreat value
	baseMessageDispatch map[int](*callback)
}

func (self *Base) Id() uint {
	return self.id
}
func (self *Base) Send(m *Message) {
	self.in <- m
}
func (self *Base) SetId(id uint) {
	self.id = id
}
func (self *Base) In() chan *Message {
	return self.in
}
func (self *Base) Close() {
	close(self.in)
}
func (self *Base) SetSelf(i interface{}) {
	self.self = reflect.ValueOf(i)
}
func (self *Base) RegisterBaseCB(typ int, i interface{}, isMethod bool) {
	self.baseMessageDispatch[typ] = &callback{isMethod, reflect.ValueOf(i)}
}
func (self *Base) DispatchM(m *Message) (ret bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("recover base dispatchm: stack: %v\n, %v", string(debug.Stack()), err)
		}
		return
	}()
	cb, ok := self.baseMessageDispatch[m.Type]
	if !ok {
		log.Warn("message type %d is has no cb", m.Type)
		return false
	}
	cbv := cb.cb
	n := 3 + len(m.Data)
	if cb.hasSelf {
		n++
	}
	param := make([]reflect.Value, n)
	start := 0
	if cb.hasSelf {
		param[start] = self.self
		start++
	}
	param[start] = reflect.ValueOf(m.Dest)
	start++
	param[start] = reflect.ValueOf(m.Src)
	start++
	param[start] = reflect.ValueOf(m.MsgEncodeType)
	start++
	for i := start; i < n; i++ {
		param[i] = reflect.ValueOf(m.Data[i-start])
	}
	cbv.Call(param)
	return true
}

func NewBase() *Base {
	a := &Base{}
	a.in = make(chan *Message, 1024)
	a.baseMessageDispatch = make(map[int](*callback))
	return a
}
