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
	*core.Base
	client  uint
	encoder *binary.Encoder
}

func TestClient(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)

	for i := 0; i < 10000; i++ {
		c := &C{Base: core.NewBase()}
		core.RegisterService(c)
		c.encoder = binary.NewEncoder()

		client := tcp.NewClient("127.0.0.1", "3333", c.Id())
		c.client = client.Run()

		go func() {
			for m := range c.In() {
				if m.Type == core.MSG_TYPE_NORMAL {
					//cmd := m.Data[0].(int)
					if len(m.Data) >= 2 {
						log.Info("recv data : %s", string(m.Data[1].([]byte)))
					}
				}
			}
		}()

		go func() {
			for {
				var a []byte = []byte("alsdkjfladjflkasdjf")
				c.encoder.Reset()
				c.encoder.Encode(a)
				c.encoder.UpdateLen()
				t := c.encoder.Buffer()
				t1 := make([]byte, len(t))
				copy(t1[:], t[:])
				core.Send(c.client, 0, tcp.CLIENT_CMD_SEND, t1)
				time.Sleep(time.Second)
			}
		}()
	}

	for {
		time.Sleep(time.Minute)
	}

}
