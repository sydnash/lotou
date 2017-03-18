package core

import (
	"github.com/sydnash/lotou/timer"
)

type Module interface {
	//OnInit is called within StartService
	OnInit()
	//OnDestory is called when service is closed
	OnDestroy()
	//OnMainLoop is called ever main loop, the delta time is specific by GetDuration()
	OnMainLoop(dt int) //dt is the duration time(unit Millisecond)
	//OnNormalMSG is called when received msg from Send() or RawSend() with MSG_TYPE_NORMAL
	OnNormalMSG(src ServiceID, data ...interface{})
	//OnSocketMSG is called when received msg from Send() or RawSend() with MSG_TYPE_SOCKET
	OnSocketMSG(src ServiceID, data ...interface{})
	//OnRequestMSG is called when received msg from Request()
	OnRequestMSG(src ServiceID, rid uint64, data ...interface{})
	//OnCallMSG is called when received msg from Call()
	OnCallMSG(src ServiceID, rid uint64, data ...interface{})
	//OnDistributeMSG is called when received msg from Send() or RawSend() with MSG_TYPE_DISTRIBUTE
	OnDistributeMSG(data ...interface{})
	//OnCloseNotify is called when received msg from SendClose() with false param.
	OnCloseNotify()

	setService(s *service)
	getDuration() int
}

type Skeleton struct {
	s                 *service
	Id                ServiceID
	Name              string
	D                 int
	normalDispatcher  *CallHelper
	requestDispatcher *CallHelper
	callDispatcher    *CallHelper
}

func NewSkeleton(d int) *Skeleton {
	ret := &Skeleton{D: d}
	ret.normalDispatcher = NewCallHelper()
	ret.requestDispatcher = NewCallHelper()
	ret.callDispatcher = NewCallHelper()
	return ret
}

func (s *Skeleton) setService(ser *service) {
	s.s = ser
	s.Id = ser.getId()
	s.Name = ser.getName()
}

func (s *Skeleton) getDuration() int {
	return s.D
}

//use gob encode(not golang's standard library, see "github.com/sydnash/lotou/encoding/gob"
//only support basic types and Message
//user defined struct should encode and decode by user
func (s *Skeleton) Send(dst ServiceID, msgType int32, data ...interface{}) {
	send(s.s.getId(), dst, msgType, data...)
}

//RawSend not encode variables, be careful use
//variables that passed by reference may be changed by others
func (s *Skeleton) RawSend(dst ServiceID, msgType int32, data ...interface{}) {
	sendNoEnc(s.s.getId(), dst, msgType, data...)
}

//if isForce is false, then it will just notify the sevice it need to close
//then service can do choose close immediate or close after self clean.
//if isForce is true, then it close immediate
func (s *Skeleton) SendClose(dst ServiceID, isForce bool) {
	sendNoEnc(s.s.getId(), dst, MSG_TYPE_CLOSE, isForce)
}

//Request send a request msg to dst, and start timeout function if timeout > 0
//after receiver call Respond, the responseCb will be called
func (s *Skeleton) Request(dst ServiceID, timeout int, responseCb interface{}, timeoutCb interface{}, data ...interface{}) {
	s.s.request(dst, timeout, responseCb, timeoutCb, data...)
}

//Respond used to respond request msg
func (s *Skeleton) Respond(dst ServiceID, rid uint64, data ...interface{}) {
	s.s.respond(dst, rid, data...)
}

//Call send a call msg to dst, and start a timeout function with the conf.CallTimeOut
//after receiver call Ret, it will return
func (s *Skeleton) Call(dst ServiceID, data ...interface{}) ([]interface{}, error) {
	return s.s.call(dst, data...)
}

func (s *Skeleton) Schedule(interval, repeat int, cb timer.TimerCallback) *timer.Timer {
	if s.s == nil {
		panic("Schedule must call after OnInit is called(contain OnInit)")
	}
	return s.s.schedule(interval, repeat, cb)
}

//Ret used to ret call msg
func (s *Skeleton) Ret(dst ServiceID, cid uint64, data ...interface{}) {
	s.s.ret(dst, cid, data...)
}

func (s *Skeleton) OnDestroy() {
}
func (s *Skeleton) OnMainLoop(dt int) {
}
func (s *Skeleton) OnNormalMSG(src ServiceID, data ...interface{}) {
	id := data[0]
	data[0] = src
	s.normalDispatcher.Call(id, data...)
}
func (s *Skeleton) OnInit() {
}
func (s *Skeleton) OnSocketMSG(src ServiceID, data ...interface{}) {
}
func (s *Skeleton) OnRequestMSG(src ServiceID, rid uint64, data ...interface{}) {
	id := data[0]
	data[0] = src
	ret := s.requestDispatcher.Call(id, data...)
	s.Respond(src, rid, ret...)
}
func (s *Skeleton) OnCallMSG(src ServiceID, rid uint64, data ...interface{}) {
	id := data[0]
	data[0] = src
	ret := s.callDispatcher.Call(id, data...)
	s.Ret(src, rid, ret...)
}

func (s *Skeleton) findCallerByType(infoType int32) *CallHelper {
	var caller *CallHelper
	switch infoType {
	case MSG_TYPE_NORMAL:
		caller = s.normalDispatcher
	case MSG_TYPE_REQUEST:
		caller = s.requestDispatcher
	case MSG_TYPE_CALL:
		caller = s.callDispatcher
	default:
		panic("not support infoType")
	}
	return caller
}

//function's first parameter must ServiceID
func (s *Skeleton) SubscribeFunc(infoType int32, id interface{}, fun interface{}) {
	caller := s.findCallerByType(infoType)
	switch key := id.(type) {
	case int:
		caller.AddFuncInt(key, fun)
	case string:
		caller.AddFunc(key, fun)
	}
}

//method's first parameter must ServiceID
func (s *Skeleton) SubscribeMethod(infoType int32, id interface{}, v interface{}, methodName string) {
	caller := s.findCallerByType(infoType)
	switch key := id.(type) {
	case int:
		caller.AddMethodInt(key, v, methodName)
	case string:
		caller.AddMethod(key, v, methodName)
	}
}

func (s *Skeleton) OnDistributeMSG(data ...interface{}) {
}
func (s *Skeleton) OnCloseNotify() {
	s.SendClose(s.s.getId(), true)
}
