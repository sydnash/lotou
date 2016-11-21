package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type master struct {
	*core.Base
	decoder *gob.Decoder
	encoder *gob.Encoder
}

func StartMaster(ip, port) {
	m := &master{Base: core.NewBase()}
	m.decoder = gob.NewDecoder()
	m.encoder = gob.NewEncoder()
	core.RegisterService(m)
	core.Name(".master", m.Id())
	s := tcp.New(ip, port, m.Id())
	s.Listen()
	m.run()
}

func (m *master) run() {
	m.SetSelf(m)
	m.RegisterBaseCB(MSG_TYPE_CLOSE, (*master).close, true)
	m.RegisterBaseCB(MSG_TYPE_NORMAL, (*master).normalMSG, false)
	for msg := range m.In() {
		m.DispatchM(msg)
	}
}

func (m *master) normalMSG(dest, src uint, msgEncode string, data ...interface{}) {
	if msgEncode == "go" {
		//dest is master's id, src is core's id
		//data[0] is cmd such as (registerNodeRet, regeisterNameRet, getIdByNameRet...)
		//data[1] is dest nodeService's id
		//parse node's id, and choose correct agent and send msg to that node.
	} else if msgEncode == "socket" {
		//dest is master's id, src is agent's id
		//data[0] is socket status
		//data[1] is a gob encode data
		//it's first encode value is cmd such as (registerNode, regeisterName, getIdByName, forword...)
		//it's second encode value is dest service's id.
		//find correct agent and send msg to that node.
		cmd := data[0].(int)
		if cmd == tcp.AGENT_DATA {
			m.decoder.SetBuffer(data[1].([]byte))
			sdata := m.decoder.Deocde()
			scmd := sdata[0].(string)
			if smcd == "registerNode" {
				nodeId := core.GenerateNodeId()
			}
		} else if cmd == tcp.AGENT_CLOSED {
		}
	}
}

func (m *master) close(dest, src uint) {
	_, _ = dest, src
	m.Close()
}
