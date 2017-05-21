package tcp_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/binary"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
	"testing"
	"time"
)

type C struct {
	*core.Skeleton
	client  core.ServiceID
	encoder *binary.Encoder
	decoder *binary.Decoder
}

func (c *C) OnMainLoop(dt int) {
	var a []byte = []byte("alsdkjfladjflkasdjf")
	c.encoder.Reset()
	c.encoder.Encode(a)
	c.encoder.UpdateLen()
	t := c.encoder.Buffer()
	t1 := make([]byte, len(t))
	copy(t1[:], t[:])
	c.RawSend(c.client, core.MSG_TYPE_NORMAL, tcp.CLIENT_CMD_SEND, t1)
}

func (c *C) OnNormalMSG(msg *core.Message) {
	data := msg.Data
	if len(data) >= 2 {
		log.Info("recv data :%s", string(data[0].([]byte)))
	}
}

func (c *C) OnSocketMSG(msg *core.Message) {
	cmd := msg.Cmd
	data := msg.Data
	if cmd == tcp.CLIENT_DATA {
		data := data[0].([]byte)
		c.decoder.SetBuffer(data)
		var msg []byte = []byte{}
		c.decoder.Decode(&msg)
		log.Info(string(msg))
	}
}

func TestClient(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)

	for i := 0; i < 1; i++ {
		c := &C{Skeleton: core.NewSkeleton(10)}
		core.StartService(&core.ModuleParam{
			N: ".client",
			M: c,
			L: 0,
		})
		c.encoder = binary.NewEncoder()
		c.decoder = binary.NewDecoder()

		client := tcp.NewClient("127.0.0.1", "3333", c.Id)
		c.client = core.StartService(&core.ModuleParam{
			N: ".cc",
			M: client,
			L: 0,
		})
	}

	for {
		time.Sleep(time.Minute)
	}
}
