package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type slave struct {
	*core.Base
	decoder *gob.Decoder
	encoder *gob.Encoder
	client  *uint
}

func StartSlave(ip, port) {
	m := &slave{Base: core.NewBase()}
	m.decoder = gob.NewDecoder()
	m.encoder = gob.NewEncoder()
	core.RegisterService(m)
	core.Name(".slave", m.Id())
	c := tcp.NewClient(ip, port, m.Id())
	m.client = c.Run()
	m.run()
}

func (s *slave) run() {
	s.SetSelf(s)
	s.RegisterBaseCB(MSG_TYPE_CLOSE, (*master).close, true)
	s.RegisterBaseCB(MSG_TYPE_NORMAL, (*master).normalMSG, false)
	for msg := range s.In() {
		s.DispatchM(msg)
	}
}

func (s *slave) normalMSG(dest, src uint, msgEncode string, data ...interface{}) {
	if msgEncode == "go" {
		//dest is master's id, src is core's id
		//data[0] is cmd such as (registerNode, regeisterName, getIdByName...)
		//data[1] is dest nodeService's id
		//parse node's id, and choose correct agent and send msg to that node.
		s.encoder.Reset()
		s.encoder.Encode(data)
		s.encoder.UpdateLen()
		t := s.encoder.Buffer()
		t1 := make([]byte, len(t))
		copy(t1, t)
		core.Send(c.client, s.Id(), tcp.CLIENT_CMD_SNED, t1)
	} else if msgEncode == "socket" {
		//dest is master's id, src is agent's id
		//data[0] is socket status
		//data[1] is a gob encode data
		//it's first encode value is cmd such as (registerNodeRet, regeisterNameRet, getIdByNameRet, forword...)
		//it's second encode value is dest service's id.
		//find correct agent and send msg to that node.
	}
}

func (s *slave) close(dest, src uint) {
	_, _ = dest, src
	core.Close(s.client, s.Id())
	s.Close()
}
