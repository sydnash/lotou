package tcp_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/binary"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
	"testing"
)

type M struct {
	*core.Skeleton
	decoder *binary.Decoder
}

func (m *M) OnNormalMSG(msg *core.Message) {
	cmd := msg.Cmd
	if cmd == tcp.AGENT_CLOSED {
		log.Info("agent closed")
	}
}

func (m *M) OnSocketMSG(msg *core.Message) {
	src := msg.Src
	cmd := msg.Cmd
	data := msg.Data
	if cmd == tcp.AGENT_DATA {
		data := data[0].([]byte)
		m.decoder.SetBuffer(data)
		var msg []byte = []byte{}
		m.decoder.Decode(&msg)
		log.Info("%v, %v", src, string(msg))

		m.RawSend(src, core.MSG_TYPE_NORMAL, tcp.AGENT_CMD_SEND, data)
	}
}

func TestServer(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	m := &M{Skeleton: core.NewSkeleton(0)}
	m.decoder = binary.NewDecoder()
	core.StartService(&core.ModuleParam{
		N: ".m",
		M: m,
		L: 0,
	})

	s := tcp.NewServer("", "3333", m.Id)
	s.Listen()

	ch := make(chan int)
	<-ch

	s.Close()
}
