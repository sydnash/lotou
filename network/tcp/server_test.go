package tcp_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/binary"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
	"testing"
)

type M struct {
	*core.Base
	decoder *binary.Decoder
}

func TestServer(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	m := &M{Base: core.NewBase()}
	m.decoder = binary.NewDecoder()
	core.RegisterService(m)
	go func() {
		for s := range m.In() {
			if s.Type == core.MSG_TYPE_NORMAL {
				cmd := s.Data[0].(int)
				if cmd == tcp.AGENT_DATA {
					data := s.Data[1].([]byte)
					m.decoder.SetBuffer(data)
					var msg *[]byte = new([]byte)
					m.decoder.Decode(msg)
				}
			}
		}
	}()
	s := tcp.New("", "3333", m.Id())
	s.Listen()

	ch := make(chan int)
	<-ch
}
