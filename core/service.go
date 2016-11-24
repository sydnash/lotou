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
	requestId           int
	requestMap          map[int]reflect.Value
	callCh              chan []interface{}
	//callId              int
	//callMapMutex        sync.Mutex
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

//request function
//func (dest, src, encodetype, rid, data...)
func (self *Base) dispatchRequest(cb *callback, m *Message) {
	rid := m.Data[0]
	data := m.Data[1].([]interface{})
	n := 4 + len(data)
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
	param[start] = reflect.ValueOf(rid)
	start++
	for i := start; i < n; i++ {
		param[i] = reflect.ValueOf(data[i-start])
	}
	cb.cb.Call(param)
}

//func (encodetype, data ...)
func (self *Base) dispatchRespond(m *Message) {
	rid := m.Data[0].(int)
	if rid < 0 {
		rid = -rid
	}
	data := m.Data[1].([]interface{})
	cb, ok := self.requestMap[rid]
	if !ok {
		log.Warn("dispatchRespond: id %v is not find.", rid)
		return
	}
	delete(self.requestMap, rid)
	n := len(data)
	param := make([]reflect.Value, n+1)
	param[0] = reflect.ValueOf(m.MsgEncodeType)
	for i := 0; i < n; i++ {
		param[i+1] = reflect.ValueOf(data[i])
	}
	cb.Call(param)
}

//func (dest, src, encodetype, data...)
func (self *Base) DispatchM(m *Message) (ret bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("recover base dispatchm: stack: %v\n, %v", string(debug.Stack()), err)
		}
		return
	}()
	if m.Type == MSG_TYPE_RESPOND {
		self.dispatchRespond(m)
		return true
	}
	cb, ok := self.baseMessageDispatch[m.Type]
	if !ok {
		log.Warn("message type %d is has no cb", m.Type)
		return false
	}
	if m.Type == MSG_TYPE_REQUEST {
		self.dispatchRequest(cb, m)
		return true
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
	a.requestMap = make(map[int]reflect.Value)
	a.callCh = make(chan []interface{})
	return a
}

func (self *Base) Request(cb interface{}) (id int) {
	self.requestId++
	id = self.requestId
	v := reflect.ValueOf(cb)
	kind := v.Kind()

	var rcb reflect.Value
	switch kind {
	case reflect.Func:
		rcb = v
	case reflect.String:
		if !self.self.IsValid() {
			panic("base:request>> self.self must not be nil")
		}
		rcb = self.self.MethodByName(cb.(string))
	default:
		panic("base:request>> cb must be func or string.")
	}
	self.requestMap[id] = rcb
	return id
}
func (self *Base) Call() []interface{} {
	ret := <-self.callCh
	return ret
}
func (self *Base) Ret(data []interface{}) {
	select {
	case self.callCh <- data:
	default:
		log.Warn("there is no current call")
	}
}
