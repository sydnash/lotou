package core

import (
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/timer"
	"runtime/debug"
)

//StartService start a given module with specific name
//call module's OnInit interface after register
//and register name to master if name is a global name
//start msg loop
func StartService(name string, m Module) ServiceID {
	s := newService(name)
	s.m = m
	id := registerService(s)
	m.SetService(s)
	d := m.GetDuration()
	s.loopDuration = d
	if d > 0 {
		s.ts = timer.NewTS()
	}
	m.OnInit()
	if !checkIsLocalName(name) {
		globalName(id, name)
	}
	if d > 0 {
		s.runWithLoop(d)
	} else {
		s.run()
	}
	return id
}

//Parse Node Id parse node id from service id
func ParseNodeId(id ServiceID) uint {
	return parseNodeIdFromId(id)
}

func Send(dst ServiceID, msgType int, data ...interface{}) error {
	return rawSend(true, INVALID_SERVICE_ID, dst, msgType, data...)
}

//SendCloseToAll simple send a close msg to all service
func SendCloseToAll() {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for _, ser := range h.dic {
		localSendWithNoMutex(INVALID_SERVICE_ID, ser, MSG_TYPE_CLOSE, MSG_ENC_TYPE_NO, false)
	}
}

func Wait() {
	exitGroup.Wait()
}

func CheckIsLocalServiceId(id ServiceID) bool {
	return checkIsLocalId(id)
}

//SafeGo start a groutine, and handle all panic within it.
func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("recover: stack: %v\n, %v", string(debug.Stack()), err)
			}
			return
		}()
		f()
	}()
}
