package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/binary"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type C struct {
	*core.Base
	client  uint
	encoder *binary.Encoder
}

func main() {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	c := &C{Base: core.NewBase()}
	core.RegisterService(c)
	c.encoder = binary.NewEncoder()

	client := tcp.NewClient("127.0.0.1", "4000", c.Id())
	c.client = client.Run()

	go func() {
		for m := range c.In() {
			if m.Type == core.MSG_TYPE_NORMAL {
				cmd := m.Data[0].(int)
				log.Info("recv message: %v", cmd)
				if len(m.Data) >= 3 {
					log.Info("recv data : %s", string(m.Data[1].([]byte)))
				}
			}
		}
	}()

	for {
		var a []byte = []byte("alsdkjfladjflkasdjf")
		c.encoder.Reset()
		c.encoder.Encode(a)
		c.encoder.UpdateLen()
		t := c.encoder.Buffer()
		t1 := make([]byte, len(t))
		copy(t1[:], t[:])
		core.Send(c.client, 0, tcp.CLIENT_CMD_SEND, t1)
	}
	for {
	}
}
