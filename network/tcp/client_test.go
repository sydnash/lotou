package tcp_test

import (
	"github.com/sydnash/majiang/core"
	"github.com/sydnash/majiang/log"
	"github.com/sydnash/majiang/network/tcp"
	"testing"
	"time"
)

type C struct {
	*core.Base
	client uint
}

func TestClient(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	log.Info("start test")
	c := &C{Base: core.NewBase()}
	core.RegisterService(c)

	client := tcp.NewClient("", "4000", c.Id())
	c.client = client.Run()

	go func() {
		for m := range c.In() {
			if m.Type == core.MSG_TYPE_NORMAL {
				cmd := m.Data[0].(int)
				log.Info("recv message: %v", cmd)
				if len(m.Data) >= 2 {
					log.Info("recv data : %s", string(m.Data[1].([]byte)))
				}
			}
		}
	}()

	for {
		var a []byte = []byte("alsdkjfladjflkasdjf")
		core.Send(c.client, 0, 1, a)
		time.Sleep(time.Second)
	}
	for {
	}
}
