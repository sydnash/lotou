package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type master struct {
	*core.Base
	decoder  *gob.Decoder
	encoder  *gob.Encoder
	nodesMap map[uint]uint
}

func StartMaster(ip, port string) {
	m := &master{Base: core.NewBase()}
	m.decoder = gob.NewDecoder()
	m.encoder = gob.NewEncoder()
	m.nodesMap = make(map[uint]uint)
	core.RegisterService(m)
	core.Name(m.Id(), ".master")
	s := tcp.New(ip, port, m.Id())
	s.Listen()
	m.run()
}

func (m *master) run() {
	m.SetSelf(m)
	m.RegisterBaseCB(core.MSG_TYPE_CLOSE, (*master).close, true)
	m.RegisterBaseCB(core.MSG_TYPE_NORMAL, (*master).normalMSG, true)
	go func() {
		for msg := range m.In() {
			m.DispatchM(msg)
		}
	}()
}

func (m *master) normalMSG(dest, src uint, msgEncode string, data ...interface{}) {
	if msgEncode == "go" {
		//dest is master's id, src is core's id
		//data[0] is cmd such as (registerNodeRet, regeisterNameRet, getIdByNameRet...)
		//data[1] is dest nodeService's id
		//parse node's id, and choose correct agent and send msg to that node.
		t1 := m.encode(data)
		for _, v := range m.nodesMap {
			core.Send(v, m.Id(), tcp.AGENT_CMD_SEND, t1)
		}
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
			sdata, _ := m.decoder.Decode()
			array := sdata.([]interface{})
			scmd := array[0].(string)
			log.Debug("recv cmd:%s", scmd)
			if scmd == "registerNode" {
				nodeId := core.GenerateNodeId()
				m.nodesMap[nodeId] = src
				ret := make([]interface{}, 2, 2)
				ret[0] = "registerNodeRet"
				ret[1] = nodeId
				sendData := m.encode(ret)
				core.Send(src, m.Id(), tcp.AGENT_CMD_SEND, sendData)
			} else if scmd == "registerName" {
				serviceId := array[1].(uint)
				serviceName := array[2].(string)
				core.Name(serviceId, serviceName)
			}
		} else if cmd == tcp.AGENT_CLOSED {
		}
	}
}

func (m *master) encode(d []interface{}) []byte {
	m.encoder.Reset()
	m.encoder.Encode(d)
	m.encoder.UpdateLen()
	t := m.encoder.Buffer()
	//make a copy to be send.
	t1 := make([]byte, len(t))
	copy(t1, t)
	return t1
}

func (m *master) close(dest, src uint) {
	_, _ = dest, src
	m.Close()
}
