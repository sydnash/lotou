package core

type Module interface {
	OnInit()
	OnDestroy()
	OnMainLoop(dt int) //dt is the duration time(unit Millisecond)
	OnNormalMSG(src uint, data ...interface{})
	SetService(s *service)
}

type Skeleton struct {
	s    *service
	Id   uint
	Name string
}

func (s *Skeleton) SetService(ser *service) {
	s.s = ser
	s.Id = ser.getId()
	s.Name = ser.getName()
}

func (s *Skeleton) Send(dst uint, msgType int, data ...interface{}) {
	send(s.s, dst, msgType, data...)
}

func (s *Skeleton) OnDestroy() {
}
func (s *Skeleton) OnMainLoop(dt int) {
}
func (s *Skeleton) OnNormalMSG(src uint, data ...interface{}) {
}
func (s *Skeleton) OnInit() {
}
