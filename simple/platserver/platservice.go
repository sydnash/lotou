package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"reflect"
	"time"
)

type PlatService struct {
	*core.Base
	sessionMap map[int]uint64
}

func (ps *PlatService) CloseMSG(dest, src uint) {
	ps.Base.Close()
}
func (ps *PlatService) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	log.Info("%x, %x, %v", src, dest, data)
	if msgType == "go" {
		cmd := data[0].(string)
		if cmd == "SessionTimeout" {
			ps.SessionTimeout(data[1].(int))
		}
	}
}
func (ps *PlatService) CallMSG(dest, src uint, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)

	cmd := data[0].(string)
	psv := reflect.ValueOf(ps)
	fv := psv.MethodByName(cmd)
	if fv.IsValid() {
		in := make([]reflect.Value, len(data)-1)
		for i := 1; i < len(data); i++ {
			in[i-1] = reflect.ValueOf(data[i])
		}
		ret := fv.Call(in)
		out := make([]interface{}, len(ret))
		for i := 0; i < len(ret); i++ {
			out[i] = ret[i].Interface()
		}
		core.Ret(src, dest, out...)
	} else {
		core.Ret(src, dest, 0)
	}
}
func (ps *PlatService) Login(ac string, acType, qdType int, mac string, loginType int) (iret int, msg, ip, port, pwd_sec string, session uint64, acid int) {
	row := db.QueryRow("call Logon(?,?,?,?,?,?,?,?)", ac, "", mac, acType, loginType, "", 1, qdType)
	err := row.Scan(&acid)
	if err != nil {
		iret = 0
		msg = err.Error()
		return
	}
	iret = 1
	ip = "192.168.23.7"
	port = "20001"
	session = core.UUID()
	ps.sessionMap[acid] = session

	time.AfterFunc(time.Second*10, func() {
		core.Send(ps.Id(), ps.Id(), "SessionTimeout", acid)
	})
	return
}
func (ps *PlatService) CheckSeesion(acid int, session uint64) (ok bool) {
	var cur uint64
	cur, ok = ps.sessionMap[acid]
	if !ok {
		return false
	}
	if cur == session {
		return true
	}
	return false
}
func (ps *PlatService) SessionTimeout(acid int) {
	delete(ps.sessionMap, acid)
	log.Debug("delete session of acid:%v", acid)
}

func (ps *PlatService) RequestMSG(dest, src uint, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	core.Respond(src, dest, rid, data...)
}

func NewPS() *PlatService {
	ps := &PlatService{Base: core.NewBase()}
	ps.sessionMap = make(map[int]uint64)
	ps.SetDispatcher(ps)
	return ps
}

func (ps *PlatService) Run() {
	core.RegisterService(ps)
	core.Name(ps.Id(), "platservice")
	go func() {
		for msg := range ps.In() {
			ps.DispatchM(msg)
		}
	}()
}
