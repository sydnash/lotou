package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/binary"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
	"github.com/sydnash/lotou/simple/btype"
	"strconv"
)

type HallService struct {
	*core.Base
	platId  uint
	decoder *binary.Decoder
	encoder *binary.Encoder
}

func (hs *HallService) CloseMSG(dest, src uint) {
	hs.Base.Close()
}
func (hs *HallService) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	log.Info("%x, %x, %v", src, dest, data)
	if msgType == "socket" {
		cmd := data[0].(int)
		var d []byte
		if len(data) >= 2 {
			d = data[1].([]byte)
		}
		hs.socketMSG(src, cmd, d)
	} else if msgType == "go" {
	}
}
func (hs *HallService) CallMSG(dest, src uint, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)
	core.Ret(src, dest, data...)
}
func (hs *HallService) RequestMSG(dest, src uint, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	core.Respond(src, dest, rid, data...)
}

func (hs *HallService) socketMSG(src uint, cmd int, data []byte) {
	switch cmd {
	case tcp.AGENT_DATA:
		hs.socketData(src, data)
	}
}

func (hs *HallService) socketData(src uint, data []byte) {
	var basic btype.PHead
	hs.decoder.SetBuffer(data)
	hs.decoder.Decode(&basic)
	log.Debug("recv package: %v", basic)
	ctype := basic.Type
	switch ctype {
	case btype.C_MSG_CHECK_SESSION:
		hs.ccheckSession(src, &basic)
	}
}
func (hs *HallService) ccheckSession(src uint, basic *btype.PHead) {
	var param btype.CCheckSession
	hs.decoder.Decode(&param)
	log.Debug("check session param:%v", param)

	cb := func(ok bool) {
		log.Info("check session isok:%v", ok)
		if ok {
			hs.encoder.Reset()
			basic.Type = btype.S_MSG_CHECK_SESSION
			hs.encoder.Encode(*basic)
			var ret btype.SCheckSession
			hs.encoder.Encode(ret)
			hs.sendToAgent(src)
		}
	}
	session, _ := strconv.ParseUint(param.Session, 10, 64)
	log.Debug("platid:%x, src id :%x", hs.platId, hs.Id())
	core.Request(hs.platId, hs, cb, "CheckSeesion", int(param.AcId), session)
}

func (hs *HallService) sendToAgent(dest uint) {
	hs.encoder.UpdateLen()
	b := hs.encoder.Buffer()
	nb := make([]byte, len(b))
	copy(nb, b)
	core.Send(dest, hs.Id(), tcp.AGENT_CMD_SEND, nb)
}

func NewHS(platid uint) *HallService {
	hs := &HallService{Base: core.NewBase()}
	hs.platId = platid
	hs.decoder = binary.NewDecoder()
	hs.encoder = binary.NewEncoder()
	hs.SetDispatcher(hs)
	return hs
}

func (hs *HallService) Run() {
	core.RegisterService(hs)
	core.Name(hs.Id(), "hallService")
	go func() {
		for msg := range hs.In() {
			hs.DispatchM(msg)
		}
	}()

	s := tcp.New("", "20001", hs.Id())
	s.Listen()
}
