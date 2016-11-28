package core

type CloseMSGDispatcher interface {
	CloseMSG(dest, src uint)
}
type NormalMSGDispatcher interface {
	NormalMSG(dest, src uint, enType string, data ...interface{})
}
type CallMSGDispatcher interface {
	CallMSG(dest, src uint, data ...interface{})
}
type RequestMSGDispatcher interface {
	RequestMSG(dest, src uint, rid int, data ...interface{})
}
type MSGDispatcher interface {
	CloseMSGDispatcher
	NormalMSGDispatcher
	CallMSGDispatcher
	RequestMSGDispatcher
}
type EmptyClose struct{}

func (p *EmptyClose) CloseMSG(dest, src uint) {}

type EmptyNormal struct{}

func (p *EmptyNormal) NormalMSG(dest, src uint, enType string, data ...interface{}) {
}

type EmptyCall struct{}

func (p *EmptyCall) CallMSG(dest, src uint, data ...interface{}) {}

type EmptyRequest struct{}

func (p *EmptyRequest) RequestMSG(dest, src uint, rid int, data ...interface{}) {}
