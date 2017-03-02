package core

import (
	"github.com/sydnash/lotou/log"
	"runtime/debug"
)

func StartService(name string, m Module) uint {
	s := newService(name)
	s.m = m
	id := registerService(s)
	m.SetService(s)
	m.OnInit()
	if !checkIsLocalName(name) {
		globalName(id, name)
	}
	d := m.GetDuration()
	if d > 0 {
		s.runWithLoop(d)
	} else {
		s.run()
	}
	return id
}
func ParseNodeId(id uint) uint {
	return parseNodeIdFromId(id)
}

func CheckIsLocalServiceId(id uint) bool {
	return checkIsLocalId(id)
}

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
