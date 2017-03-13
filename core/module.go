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
	OnRequestMSG(src ServiceID, rid int, data ...interface{})
	//OnCallMSG is called when received msg from Call()
	OnCallMSG(src ServiceID, rid int, data ...interface{})
	//OnDistributeMSG is called when received msg from Send() or RawSend() with MSG_TYPE_DISTRIBUTE
	OnDistributeMSG(data ...interface{})
	//OnCloseNotify is called when received msg from SendClose() with false param.
	OnCloseNotify()
	SetService(s *service)
	GetDuration() int
}

type Skeleton struct {
	s    *service
	Id   ServiceID
	Name string
	D    int
}

func NewSkeleton(d int) *Skeleton {
	return &Skeleton{D: d}
}

func (s *Skeleton) SetService(ser *service) {
	s.s = ser
	s.Id = ser.getId()
	s.Name = ser.getName()
}

func (s *Skeleton) GetDuration() int {
	return s.D
}

//use gob encode(not golang's standard library, see "github.com/sydnash/lotou/encoding/gob"
//only support basic types and Message
//user defined struct should encode and decode by user
func (s *Skeleton) Send(dst ServiceID, msgType int, data ...interface{}) {
	send(s.s.getId(), dst, msgType, data...)
}

//RawSend not encode variables, be careful use
//variables that passed by reference may be changed by others
func (s *Skeleton) RawSend(dst ServiceID, msgType int, data ...interface{}) {
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
func (s *Skeleton) Respond(dst ServiceID, rid int, data ...interface{}) {
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
func (s *Skeleton) Ret(dst ServiceID, cid int, data ...interface{}) {
	s.s.ret(dst, cid, data...)
}

func (s *Skeleton) OnDestroy() {
}
func (s *Skeleton) OnMainLoop(dt int) {
}
func (s *Skeleton) OnNormalMSG(src ServiceID, data ...interface{}) {
}
func (s *Skeleton) OnInit() {
}
func (s *Skeleton) OnSocketMSG(src ServiceID, data ...interface{}) {
}
func (s *Skeleton) OnRequestMSG(src ServiceID, rid int, data ...interface{}) {
}
func (s *Skeleton) OnCallMSG(src ServiceID, rid int, data ...interface{}) {
}
func (s *Skeleton) OnDistributeMSG(data ...interface{}) {
}
func (s *Skeleton) OnCloseNotify() {
	s.SendClose(s.s.getId(), true)
}
