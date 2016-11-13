package tcp_test

import (
	"github.com/sydnash/majiang/core"
	"github.com/sydnash/majiang/log"
	"github.com/sydnash/majiang/network/tcp"
	"testing"
)

type M struct {
	*core.Base
}

func TestServer(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	m := &M{Base: core.NewBase()}
	core.RegisterService(m)
	go func() {
		for s := range m.In() {
			if s.Type == core.MSG_TYPE_NORMAL {
				log.Info("recv message: %v", string(s.Data[0].([]byte)))
			}
		}
	}()
	s := tcp.New("", "4000", m.Id())
	s.Listen()

	for {
	}
}
