package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"reflect"
	"time"
)

type HallService struct {
	*core.Base
}

func (hs *HallService) CloseMSG(dest, src uint, msgType string) {
	hs.Base.Close()
}
func (hs *HallService) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	if msgType == "socket" {
		cmd := data[0].(int)
		d := data[1].([]byte)
		hs.socketMSG(cmd, d)
	} else if msgType == "go" {
	}
	log.Info("%x, %x, %v", src, dest, data)
}
func (hs *HallService) CallMSG(dest, src uint, msgType string, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)
	core.Ret(src, dest, data...)
}
func (hs *HallService) RequestMSG(dest, src uint, msgType string, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	core.Respond(src, dest, rid, data...)
}

func (hs *HallService) socketMSG(cmd int, data []byte) {
}

func NewHS() *HallService {
	hs := &HallService{Base: core.NewBase()}
	hs.sessionMap = make(map[int]uint64)
	hs.SetSelf(hs)
	hs.RegisterBaseCB(core.MSG_TYPE_CLOSE, "CloseMSG")
	hs.RegisterBaseCB(core.MSG_TYPE_NORMAL, "NormalMSG")
	hs.RegisterBaseCB(core.MSG_TYPE_CALL, "CallMSG")
	hs.RegisterBaseCB(core.MSG_TYPE_REQUEST, "RequestMSG")
	return hs
}

func (hs *HallService) Run() {
	core.RegisterService(hs)
	core.Name(hs.Id(), "platservice")
	go func() {
		for msg := range hs.In() {
			hs.DispatchM(msg)
		}
	}()

	s := tcp.New("", "20001", hs.Id())
	s.Listen()
}
