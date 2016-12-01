package core

import (
	"errors"
	"github.com/sydnash/lotou/log"
	"reflect"
	"runtime/debug"
)

type Base struct {
	id                  uint
	in                  chan *Message
	baseMessageDispatch map[int](*reflect.Value)
	requestId           int
	requestMap          map[int]reflect.Value
	timeoutMap          map[int]func()
	callCh              chan []interface{}
	dispatcher          MSGDispatcher
}

func (self *Base) Id() uint {
	return self.id
}

//if in is full, how can i process this.(return a error or panic right hear?)
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
func (self *Base) SetDispatcher(dispatcher MSGDispatcher) {
	self.dispatcher = dispatcher
}

/*
func findBasicV(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return findBasicV(v.Elem())
	}
	return v
}
func findNestMethodByName(v reflect.Value, name string) reflect.Value {
	f := v.MethodByName(name)
	if !f.IsValid() {
		t := findBasicV(v)
		n := t.NumField()
		for i := 0; i < n; i++ {
			field := t.Field(i)
			if field.Kind() == reflect.Struct {
				return findNestMethodByName(field, name)
			}
		}
	} else {
		return f
	}
	return reflect.Value{}
}
*/

//request function
//func (dest, src, encodetype, rid, data...)
func (self *Base) dispatchRequest(m *Message) {
	rid := m.Data[0].(int)
	data := m.Data[1].([]interface{})
	self.dispatcher.RequestMSG(m.Dest, m.Src, rid, data...)
}

var RespondNotExist = errors.New("respond cb is not exist.")

func (self *Base) getAndDeleteRespond(rid int) (reflect.Value, func(), error) {
	cb, ok1 := self.requestMap[rid]
	to, ok2 := self.timeoutMap[rid]
	if !ok1 && !ok2 {
		return reflect.Value{}, nil, RespondNotExist
	}
	delete(self.requestMap, rid)
	delete(self.timeoutMap, rid)
	return cb, to, nil
}

//func (istimeout, data ...)
func (self *Base) dispatchRespond(m *Message) {
	rid := m.Data[0].(int)
	if rid < 0 {
		rid = -rid
	}
	cb, timeout, err := self.getAndDeleteRespond(rid)
	if err != nil {
		return
	}

	log.Debug("type:%v", m.Type)
	if m.Type == MSG_TYPE_TIMEOUT {
		if timeout != nil {
			timeout()
		}
		return
	}

	data := m.Data[1].([]interface{})
	n := len(data)
	param := make([]reflect.Value, n)
	for i := 0; i < n; i++ {
		param[i] = reflect.ValueOf(data[i])
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
	if (m.Type == MSG_TYPE_RESPOND) || (m.Type == MSG_TYPE_TIMEOUT) {
		self.dispatchRespond(m)
		return true
	}

	switch m.Type {
	case MSG_TYPE_CLOSE:
		self.dispatcher.CloseMSG(m.Dest, m.Src)
	case MSG_TYPE_NORMAL:
		self.dispatcher.NormalMSG(m.Dest, m.Src, m.MsgEncodeType, m.Data...)
	case MSG_TYPE_CALL:
		self.dispatcher.CallMSG(m.Dest, m.Src, m.Data...)
	case MSG_TYPE_REQUEST:
		self.dispatchRequest(m)
	default:
		panic("use on supported msg type.")
	}
	return true
}

func NewBase() *Base {
	a := &Base{}
	a.in = make(chan *Message, 1024)
	a.baseMessageDispatch = make(map[int](*reflect.Value))
	a.requestMap = make(map[int]reflect.Value)
	a.timeoutMap = make(map[int]func())
	a.callCh = make(chan []interface{})
	return a
}

func NewBaseLen(blen int) *Base {
	a := &Base{}
	a.in = make(chan *Message, blen)
	a.baseMessageDispatch = make(map[int](*reflect.Value))
	a.requestMap = make(map[int]reflect.Value)
	a.timeoutMap = make(map[int]func())
	a.callCh = make(chan []interface{})
	return a
}

//func Request: generate request id and save callback
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
		dv := reflect.ValueOf(self.dispatcher)
		if !dv.IsValid() {
			panic("base:request>> self.dispatcher must not be nil")
		}
		rcb = dv.MethodByName(cb.(string))
	default:
		panic("base:request>> cb must be func or string.")
	}
	self.requestMap[id] = rcb
	return id
}
func (self *Base) Timeout(f func(), id int) {
	self.timeoutMap[id] = f
}

//Call : block until dest service return
func (self *Base) Call() []interface{} {
	ret := <-self.callCh
	return ret
}

//Ret : tell Call something has been returned
func (self *Base) Ret(data []interface{}) {
	select {
	case self.callCh <- data:
	default:
		log.Warn("there is no current call")
	}
}
