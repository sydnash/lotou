package core

type Module interface {
	OnInit()
	OnDestroy()
	OnMainLoop(dt int) //dt is the duration time(unit Millisecond)
	OnNormalMSG(src uint, data ...interface{})
	OnSocketMSG(src uint, data ...interface{})
	SetService(s *service)
}

type Skeleton struct {
	s    *service
	Id   uint
	Name string
}

func NewSkeleton() *Skeleton {
	return &Skeleton{}
}

func (s *Skeleton) SetService(ser *service) {
	s.s = ser
	s.Id = ser.getId()
	s.Name = ser.getName()
}

//use gob encode(not golang's standard library, see "github.com/sydnash/lotou/encoding/gob"
//only support basic types and Message
//user defined struct should encode and decode by user
func (s *Skeleton) Send(dst uint, msgType int, data ...interface{}) {
	send(s.s, dst, msgType, data...)
}

//not encode data, be careful use
//variables that passed by reference may be changed by others
func (s *Skeleton) RawSend(dst uint, msgType int, data ...interface{}) {
	sendNoEnc(s.s, dst, msgType, data...)
}

func (s *Skeleton) SendClose(dst uint) {
	sendNoEnc(s.s, dst, MSG_TYPE_CLOSE)
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
