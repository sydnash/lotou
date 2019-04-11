package core

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/helper"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/timer"
)

type ServiceID uint64

func (id ServiceID) parseNodeId() uint64 {
	return (uint64(id) & NODE_ID_MASK) >> NODE_ID_OFF
}
func (id ServiceID) parseBaseId() uint64 {
	return uint64(id) & (^uint64(NODE_ID_MASK))
}
func (id ServiceID) IsValid() bool {
	return !(id == INVALID_SERVICE_ID || id == 0)
}
func (id ServiceID) InValid() bool {
	return id == INVALID_SERVICE_ID || id == 0
}

type requestCB struct {
	respond reflect.Value
	//timeout reflect.Value
}
type service struct {
	id           ServiceID
	name         string
	msgChan      chan *Message
	loopTicker   *time.Ticker
	loopDuration int //unit is Millisecond
	m            Module
	requestId    uint64
	requestMap   map[uint64]requestCB
	requestMutex sync.Mutex
	callId       uint64
	callChanMap  map[uint64]chan []interface{}
	callMutex    sync.Mutex
	ts           *timer.TimerSchedule
}

var (
	ServiceCallTimeout = errors.New("call time out")
)

func newService(name string, len int) *service {
	s := &service{name: name}
	if len <= 1024 {
		len = 1024
	}
	s.msgChan = make(chan *Message, len)
	s.requestId = 0
	s.requestMap = make(map[uint64]requestCB)
	s.callChanMap = make(map[uint64]chan []interface{})
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
	select {
	case s.msgChan <- m:
	default:
		if s.msgChan == nil {
			log.Warn("msg chan is closed.<%s>", s.getName())
		} else {
			panic(fmt.Sprintf("service is full.<%s>", s.getName()))
		}
	}
}

func (s *service) destroy() {
	unregisterService(s)
	msgChan := s.msgChan
	s.msgChan = nil
	close(msgChan)
	if s.loopTicker != nil {
		s.loopTicker.Stop()
	}
}

func (s *service) dispatchMSG(msg *Message) bool {
	if msg.EncType == MSG_ENC_TYPE_GO {
		t, err := gob.Unpack(msg.Data[0].([]byte))
		if err != nil {
			panic(err)
		}
		msg.Data = t.([]interface{})
	}
	switch msg.Type {
	case MSG_TYPE_NORMAL:
		s.m.OnNormalMSG(msg)
	case MSG_TYPE_CLOSE:
		if msg.Data[0].(bool) {
			return true
		}
		s.m.OnCloseNotify()
	case MSG_TYPE_SOCKET:
		s.m.OnSocketMSG(msg)
	case MSG_TYPE_REQUEST:
		s.dispatchRequest(msg)
	case MSG_TYPE_RESPOND:
		s.dispatchRespond(msg)
	case MSG_TYPE_CALL:
		s.dispatchCall(msg)
	case MSG_TYPE_DISTRIBUTE:
		s.m.OnDistributeMSG(msg)
	case MSG_TYPE_TIMEOUT:
		s.dispatchTimeout(msg)
	}
	return false
}

//select on msgChan only
func (s *service) loopSelect() (ret bool) {
	ret = true
	defer func() {
		if err := recover(); err != nil {
			log.Error("error in service<%v>", s.getName())
			log.Error("recover: stack: %v\n, %v", helper.GetStack(), err)
		}
	}()
	select {
	case msg, ok := <-s.msgChan:
		if !ok {
			return false
		}
		isClose := s.dispatchMSG(msg)
		if isClose {
			return false
		}
	}
	return true
}

func (s *service) loop() {
	s.m.OnInit()
	for {
		if !s.loopSelect() {
			break
		}
	}
	s.m.OnDestroy()
	s.destroy()
}

//select on msgChan and a loop ticker
func (s *service) loopWithLoopSelect() (ret bool) {
	ret = true
	defer func() {
		if err := recover(); err != nil {
			log.Error("error in service<%v>", s.getName())
			log.Error("recover: stack: %v\n, %v", helper.GetStack(), err)
		}
	}()
	select {
	case msg, ok := <-s.msgChan:
		if !ok {
			return false
		}
		isClose := s.dispatchMSG(msg)
		if isClose {
			return false
		}
	case <-s.loopTicker.C:
		s.ts.Update(s.loopDuration)
		s.m.OnMainLoop(s.loopDuration)
	}
	return true
}

func (s *service) loopWithLoop() {
	s.m.OnInit()
	for {
		if !s.loopWithLoopSelect() {
			break
		}
	}
	s.loopTicker.Stop()
	s.m.OnDestroy()
	s.destroy()
}

//start a goroutinue with no ticker for main loop
func (s *service) run() {
	SafeGo(s.loop)
}

//start a goroutinue with ticker for main loop
func (s *service) runWithLoop(d int) {
	s.loopDuration = d
	s.loopTicker = time.NewTicker(time.Duration(d) * time.Millisecond)
	s.ts = timer.NewTS()
	SafeGo(s.loopWithLoop)
}

//respndCb is a function like: func(isok bool, ...interface{})  the first param must be a bool
func (s *service) request(dst ServiceID, encType EncType, timeout int, respondCb interface{}, cmd CmdType, data ...interface{}) {
	s.requestMutex.Lock()
	id := s.requestId
	s.requestId++
	cbp := requestCB{reflect.ValueOf(respondCb)}
	s.requestMap[id] = cbp
	s.requestMutex.Unlock()
	helper.PanicWhen(cbp.respond.Kind() != reflect.Func, "respond cb must function.")

	lowLevelSend(s.getId(), dst, MSG_TYPE_REQUEST, encType, id, cmd, data...)

	if timeout > 0 {
		time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
			s.requestMutex.Lock()
			_, ok := s.requestMap[id]
			s.requestMutex.Unlock()
			if ok {
				lowLevelSend(INVALID_SERVICE_ID, s.getId(), MSG_TYPE_TIMEOUT, MSG_ENC_TYPE_NO, id, Cmd_None)
			}
		})
	}
}

func (s *service) dispatchTimeout(m *Message) {
	rid := m.Id
	cbp, ok := s.getDeleteRequestCb(rid)
	if !ok {
		return
	}
	cb := cbp.respond
	var param []reflect.Value
	param = append(param, reflect.ValueOf(true))
	plen := cb.Type().NumIn()
	for i := 1; i < plen; i++ {
		param = append(param, reflect.New(cb.Type().In(i)).Elem())
	}
	cb.Call(param)
}

func (s *service) dispatchRequest(msg *Message) {
	s.m.OnRequestMSG(msg)
}

func (s *service) respond(dst ServiceID, encType EncType, rid uint64, data ...interface{}) {
	lowLevelSend(s.getId(), dst, MSG_TYPE_RESPOND, encType, rid, Cmd_None, data...)
}

//return request callback by request id
func (s *service) getDeleteRequestCb(id uint64) (requestCB, bool) {
	s.requestMutex.Lock()
	cb, ok := s.requestMap[id]
	delete(s.requestMap, id)
	s.requestMutex.Unlock()
	return cb, ok
}

func (s *service) dispatchRespond(m *Message) {
	var rid uint64
	var data []interface{}
	rid = m.Id
	data = m.Data

	cbp, ok := s.getDeleteRequestCb(rid)
	if !ok {
		return
	}
	cb := cbp.respond
	n := len(data)
	param := make([]reflect.Value, n+1)
	param[0] = reflect.ValueOf(false)
	HelperFunctionToUseReflectCall(cb, param, 1, data)
	cb.Call(param)
}

func (s *service) call(dst ServiceID, encType EncType, cmd CmdType, data ...interface{}) ([]interface{}, error) {
	helper.PanicWhen(dst == s.getId(), "dst must equal to s's id")
	s.callMutex.Lock()
	id := s.callId
	s.callId++
	s.callMutex.Unlock()

	//ch has one buffer, make ret service not block on it.
	ch := make(chan []interface{}, 1)
	s.callMutex.Lock()
	s.callChanMap[id] = ch
	s.callMutex.Unlock()
	if err := lowLevelSend(s.getId(), dst, MSG_TYPE_CALL, encType, id, cmd, data...); err != nil {
		return nil, err
	}
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

func (s *service) callWithTimeout(dst ServiceID, encType EncType, timeout int, cmd CmdType, data ...interface{}) ([]interface{}, error) {
	helper.PanicWhen(dst == s.getId(), "dst must equal to s's id")
	s.callMutex.Lock()
	id := s.callId
	s.callId++
	s.callMutex.Unlock()

	//ch has one buffer, make ret service not block on it.
	ch := make(chan []interface{}, 1)
	s.callMutex.Lock()
	s.callChanMap[id] = ch
	s.callMutex.Unlock()
	if err := lowLevelSend(s.getId(), dst, MSG_TYPE_CALL, encType, id, cmd, data...); err != nil {
		return nil, err
	}
	if timeout > 0 {
		time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
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

func (s *service) dispatchCall(msg *Message) {
	s.m.OnCallMSG(msg)
}

func (s *service) ret(dst ServiceID, encType EncType, cid uint64, data ...interface{}) {
	var dstService *service
	dstService, err := findServiceById(dst)
	if err != nil {
		lowLevelSend(s.getId(), dst, MSG_TYPE_RET, encType, cid, Cmd_None, data...)
		return
	}
	dstService.dispatchRet(cid, data...)
}

func (s *service) dispatchRet(cid uint64, data ...interface{}) {
	s.callMutex.Lock()
	ch, ok := s.callChanMap[cid]
	s.callMutex.Unlock()

	if ok {
		select {
		case ch <- data:
		default:
			helper.PanicWhen(true, "dispatchRet failed on ch.")
		}
	}
}

func (s *service) schedule(interval, repeat int, cb timer.TimerCallback) *timer.Timer {
	helper.PanicWhen(s.loopDuration <= 0, "loopDuraton must greater than zero.")
	return s.ts.Schedule(interval, repeat, cb)
}
