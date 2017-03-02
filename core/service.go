package core

import (
	"github.com/sydnash/lotou/encoding/gob"
	"reflect"
	"sync"
	"time"
)

type service struct {
	id           uint
	name         string
	msgChan      chan *Message
	loopTicker   *time.Ticker
	loopDuration int //unit is Millisecond
	m            Module
	requestId    int
	requestMap   map[int]reflect.Value
	requestMutex sync.Mutex
	callId       int
	callChanMap  map[int]chan []interface{}
	callMutex    sync.Mutex
}

func newService(name string) *service {
	s := &service{name: name}
	s.msgChan = make(chan *Message, 1024)
	s.requestId = 0
	s.requestMap = make(map[int]reflect.Value)
	s.callChanMap = make(map[int]chan []interface{})
	return s
}

func (s *service) setModule(m Module) {
	s.m = m
}

func (s *service) getName() string {
	return s.name
}

func (s *service) setId(id uint) {
	s.id = id
}

func (s *service) getId() uint {
	return s.id
}

func (s *service) pushMSG(m *Message) {
	s.msgChan <- m
}

func (s *service) destroy() {
	unregisterService(s)
	close(s.msgChan)
	if s.loopTicker != nil {
		s.loopTicker.Stop()
	}
}

func (s *service) dispatchMSG(msg *Message) bool {
	if msg.EncType == MSG_ENC_TYPE_GO {
		t := gob.Unpack(msg.Data[0].([]byte))
		msg.Data = t.([]interface{})
	}
	switch msg.Type {
	case MSG_TYPE_NORMAL:
		s.m.OnNormalMSG(msg.Src, msg.Data...)
	case MSG_TYPE_CLOSE:
		return true
	case MSG_TYPE_SOCKET:
		s.m.OnSocketMSG(msg.Src, msg.Data...)
	case MSG_TYPE_REQUEST:
		s.dispatchRequest(msg)
	case MSG_TYPE_RESPOND:
		s.dispatchRespond(msg)
	case MSG_TYPE_CALL:
		s.dispatchCall(msg)
	case MSG_TYPE_DISTRIBUTE:
		s.m.OnDistributeMSG(msg.Data...)
	}
	return false
}

func (s *service) loop() {
EXIT:
	for {
		select {
		case msg, ok := <-s.msgChan:
			if !ok {
				break EXIT
			}
			isClose := s.dispatchMSG(msg)
			if isClose {
				break EXIT
			}
		}
	}
	s.m.OnDestroy()
	s.destroy()
}

func (s *service) loopWithLoop() {
EXIT:
	for {
		select {
		case msg, ok := <-s.msgChan:
			if !ok {
				break EXIT
			}
			isClose := s.dispatchMSG(msg)
			if isClose {
				break EXIT
			}
		case <-s.loopTicker.C:
			s.m.OnMainLoop(s.loopDuration)
		}
	}
	s.loopTicker.Stop()
	s.m.OnDestroy()
	s.destroy()
}

func (s *service) run() {
	go s.loop()
}

func (s *service) runWithLoop(d int) {
	s.loopTicker = time.NewTicker(time.Duration(d) * time.Millisecond)
	go s.loopWithLoop()
}

func (s *service) request(dst uint, timeout int, cb interface{}, data ...interface{}) {
	s.requestMutex.Lock()
	id := s.requestId
	s.requestId++
	v := reflect.ValueOf(cb)
	s.requestMap[id] = v
	s.requestMutex.Unlock()
	PanicWhen(v.Kind() != reflect.Func)

	param := make([]interface{}, 2)
	param[0] = id
	param[1] = data
	rawSend(true, s.getId(), dst, MSG_TYPE_REQUEST, param...)
}

func (s *service) dispatchRequest(m *Message) {
	rid := m.Data[0].(int)
	data := m.Data[1].([]interface{})
	s.m.OnRequestMSG(m.Src, rid, data...)
}

func (s *service) respond(dst uint, rid int, data ...interface{}) {
	param := make([]interface{}, 2)
	param[0] = rid
	param[1] = data
	rawSend(true, s.getId(), dst, MSG_TYPE_RESPOND, param...)
}

func (s *service) dispatchRespond(m *Message) {
	rid := m.Data[0].(int)
	data := m.Data[1].([]interface{})

	s.requestMutex.Lock()
	cb, ok := s.requestMap[rid]
	delete(s.requestMap, rid)
	s.requestMutex.Unlock()

	if !ok {
		return
	}
	n := len(data)
	param := make([]reflect.Value, n+1)
	param[0] = reflect.ValueOf(false)
	for i := 0; i < n; i++ {
		param[i+1] = reflect.ValueOf(data[i])
	}
	cb.Call(param)
}

func (s *service) call(dst uint, data ...interface{}) ([]interface{}, error) {
	PanicWhen(dst == s.getId())
	s.callMutex.Lock()
	id := s.callId
	s.callId++
	s.callMutex.Unlock()
	param := make([]interface{}, 2)
	param[0] = id
	param[1] = data
	if err := rawSend(true, s.getId(), dst, MSG_TYPE_CALL, param...); err != nil {
		return nil, err
	}
	ch := make(chan []interface{})

	s.callMutex.Lock()
	s.callChanMap[id] = ch
	s.callMutex.Unlock()

	ret := <-ch
	s.callMutex.Lock()
	delete(s.callChanMap, id)
	s.callMutex.Unlock()

	close(ch)
	return ret, nil
}

func (s *service) dispatchCall(m *Message) {
	cid := m.Data[0].(int)
	data := m.Data[1].([]interface{})
	s.m.OnCallMSG(m.Src, cid, data...)
}

func (s *service) ret(dst uint, cid int, data ...interface{}) {
	var dstService *service
	dstService, err := findServiceById(dst)
	if err != nil {
		param := make([]interface{}, 2)
		param[0] = cid
		param[1] = data
		rawSend(true, s.getId(), dst, MSG_TYPE_RET, param...)
		return
	}
	dstService.dispatchRet(cid, data...)
}

func (s *service) dispatchRet(cid int, data ...interface{}) {
	s.callMutex.Lock()
	ch, ok := s.callChanMap[cid]
	s.callMutex.Unlock()

	if ok {
		ch <- data
	}
}
