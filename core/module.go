package core

type Module interface {
	OnInit()
	OnDestroy()
	OnMainLoop(dt int) //dt is the duration time(unit Millisecond)
	OnNormalMSG(src uint, data ...interface{})
	OnSocketMSG(src uint, data ...interface{})
	OnRequestMSG(src uint, rid int, data ...interface{})
	OnCallMSG(src uint, rid int, data ...interface{})
	OnDistributeMSG(data ...interface{})
	SetService(s *service)
	GetDuration() int
}

type Skeleton struct {
	s    *service
	Id   uint
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
func (s *Skeleton) Send(dst uint, msgType int, data ...interface{}) {
	send(s.s.getId(), dst, msgType, data...)
}

//not encode data, be careful use
//variables that passed by reference may be changed by others
func (s *Skeleton) RawSend(dst uint, msgType int, data ...interface{}) {
	sendNoEnc(s.s.getId(), dst, msgType, data...)
}

func (s *Skeleton) SendClose(dst uint) {
	sendNoEnc(s.s.getId(), dst, MSG_TYPE_CLOSE)
}

func (s *Skeleton) Request(dst uint, timeout int, responseCb interface{}, timeoutCb interface{}, data ...interface{}) {
	s.s.request(dst, timeout, responseCb, timeoutCb, data...)
}

func (s *Skeleton) Respond(dst uint, rid int, data ...interface{}) {
	s.s.respond(dst, rid, data...)
}

func (s *Skeleton) Call(dst uint, data ...interface{}) ([]interface{}, error) {
	return s.s.call(dst, data...)
}

func (s *Skeleton) Ret(dst uint, cid int, data ...interface{}) {
	s.s.ret(dst, cid, data...)
}

func (s *Skeleton) OnDestroy() {
}
func (s *Skeleton) OnMainLoop(dt int) {
}
func (s *Skeleton) OnNormalMSG(src uint, data ...interface{}) {
}
func (s *Skeleton) OnInit() {
}
func (s *Skeleton) OnSocketMSG(src uint, data ...interface{}) {
}
func (s *Skeleton) OnRequestMSG(src uint, rid int, data ...interface{}) {
}
func (s *Skeleton) OnCallMSG(src uint, rid int, data ...interface{}) {
}
func (s *Skeleton) OnDistributeMSG(data ...interface{}) {
}
