package core

import (
	"errors"
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/timer"
	"reflect"
	"sync"
	"time"
)

type ServiceID uint
type requestCB struct {
	respond reflect.Value
	timeout reflect.Value
}
type service struct {
	id           ServiceID
	name         string
	msgChan      chan *Message
	loopTicker   *time.Ticker
	loopDuration int //unit is Millisecond
	m            Module
	requestId    int
	requestMap   map[int]requestCB
	requestMutex sync.Mutex
	callId       int
	callChanMap  map[int]chan []interface{}
	callMutex    sync.Mutex
	ts           *timer.TimerSchedule
}

var (
	ServiceCallTimeout = errors.New("call time out")
)

func newService(name string) *service {
	s := &service{name: name}
	s.msgChan = make(chan *Message, 1024)
	s.requestId = 0
	s.requestMap = make(map[int]requestCB)
	s.callChanMap = make(map[int]chan []interface{})
	return s
}

func (s *service) setModule(m Module) {
	s.m = m
}

func (s *service) getName() string {
	return s.name
}

func (s *service) setId(id ServiceID) {
	s.id = id
}

func (s *service) getId() ServiceID {
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
		s.m.OnNormalMSG(ServiceID(msg.Src), msg.Data...)
	case MSG_TYPE_CLOSE:
		if msg.Data[0].(bool) {
			return true
		}
		s.m.OnCloseNotify()
	case MSG_TYPE_SOCKET:
		s.m.OnSocketMSG(ServiceID(msg.Src), msg.Data...)
	case MSG_TYPE_REQUEST:
		s.dispatchRequest(msg)
	case MSG_TYPE_RESPOND:
		s.dispatchRespond(msg)
	case MSG_TYPE_CALL:
		s.dispatchCall(msg)
	case MSG_TYPE_DISTRIBUTE:
		s.m.OnDistributeMSG(msg.Data...)
	case MSG_TYPE_TIMEOUT:
		s.dispatchTimeout(msg)
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
			s.ts.Update(s.loopDuration)
			s.m.OnMainLoop(s.loopDuration)
		}
	}
	s.loopTicker.Stop()
	s.m.OnDestroy()
	s.destroy()
}

func (s *service) run() {
	SafeGo(s.loop)
}

func (s *service) runWithLoop(d int) {
	s.loopDuration = d
	s.loopTicker = time.NewTicker(time.Duration(d) * time.Millisecond)
	SafeGo(s.loopWithLoop)
}

//respndCb is a function like: func(isok bool, ...interface{})  the first param must be a bool
//timeoutCb is a function with no param : func()
func (s *service) request(dst ServiceID, timeout int, respondCb interface{}, timeoutCb interface{}, data ...interface{}) {
	s.requestMutex.Lock()
	id := s.requestId
	s.requestId++
	cbp := requestCB{reflect.ValueOf(respondCb), reflect.ValueOf(timeoutCb)}
	s.requestMap[id] = cbp
	s.requestMutex.Unlock()
	PanicWhen(cbp.respond.Kind() != reflect.Func, "respond cb must function.")
	PanicWhen(cbp.timeout.Kind() != reflect.Func, "timeout cb must function.")

	param := make([]interface{}, 2)
	param[0] = id
	param[1] = data
	rawSend(true, s.getId(), dst, MSG_TYPE_REQUEST, param...)

	if timeout > 0 {
		time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
			rawSend(false, INVALID_SERVICE_ID, s.getId(), MSG_TYPE_TIMEOUT, id)
		})
	}
}

func (s *service) dispatchTimeout(m *Message) {
	rid := m.Data[0].(int)
	cbp, ok := s.getDeleteRequestCb(rid)
	if !ok {
		return
	}
	cb := cbp.timeout
	cb.Call([]reflect.Value{})
}

func (s *service) dispatchRequest(m *Message) {
	rid := m.Data[0].(int)
	data := m.Data[1].([]interface{})
	s.m.OnRequestMSG(ServiceID(m.Src), rid, data...)
}

func (s *service) respond(dst ServiceID, rid int, data ...interface{}) {
	param := make([]interface{}, 2)
	param[0] = rid
	param[1] = data
	rawSend(true, s.getId(), dst, MSG_TYPE_RESPOND, param...)
}

func (s *service) getDeleteRequestCb(id int) (requestCB, bool) {
	s.requestMutex.Lock()
	cb, ok := s.requestMap[id]
	delete(s.requestMap, id)
	s.requestMutex.Unlock()
	return cb, ok
}

func (s *service) dispatchRespond(m *Message) {
	rid := m.Data[0].(int)
	data := m.Data[1].([]interface{})

	cbp, ok := s.getDeleteRequestCb(rid)
	if !ok {
		return
	}
	cb := cbp.respond
	n := len(data)
	param := make([]reflect.Value, n+1)
	param[0] = reflect.ValueOf(false)
	for i := 0; i < n; i++ {
		param[i+1] = reflect.ValueOf(data[i])
	}
	cb.Call(param)
}

func (s *service) call(dst ServiceID, data ...interface{}) ([]interface{}, error) {
	PanicWhen(dst == s.getId(), "dst must equal to s's id")
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
	if conf.CallTimeOut > 0 {
		time.AfterFunc(time.Duration(conf.CallTimeOut)*time.Millisecond, func() {
			s.dispatchRet(id, ServiceCallTimeout)
		})
	}
	ret := <-ch
	s.callMutex.Lock()
	delete(s.callChanMap, id)
	s.callMutex.Unlock()

	close(ch)
	if err, ok := ret[0].(error); ok {
		return ret[1:], err
	}
	return ret, nil
}

func (s *service) dispatchCall(m *Message) {
	cid := m.Data[0].(int)
	data := m.Data[1].([]interface{})
	s.m.OnCallMSG(ServiceID(m.Src), cid, data...)
}

func (s *service) ret(dst ServiceID, cid int, data ...interface{}) {
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

func (s *service) schedule(interval, repeat int, cb timer.TimerCallback) *timer.Timer {
	PanicWhen(s.loopDuration <= 0, "loopDuraton must greater than zero.")
	return s.ts.Schedule(interval, repeat, cb)
}
