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
		cmd := data[0].(string)
		if cmd == "syncName" {
			for _, v := range m.nodesMap {
				core.Send(v, m.Id(), tcp.AGENT_CMD_SEND, t1)
			}
		} else if cmd == "forward" {
			msg := data[1].(*core.Message)
			m.forwardM(msg, nil)
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

				nameDic := core.NameDic()
				for k, v := range nameDic {
					if k[0] != '.' {
						sendArray := make([]interface{}, 3)
						sendArray[0] = "syncName"
						sendArray[1] = v
						sendArray[2] = k
						t1 := m.encode(sendArray)
						core.Send(src, m.Id(), tcp.AGENT_CMD_SEND, t1)
					}
				}
			} else if scmd == "registerName" {
				serviceId := array[1].(uint)
				serviceName := array[2].(string)
				log.Debug("%v", serviceName)
				core.Name(serviceId, serviceName)
			} else if scmd == "getIdByName" {
				name := array[1].(string)
				id, ok := core.GetIdByName(name)
				ret := make([]interface{}, 5, 5)
				ret[0] = "getIdByNameRet"
				ret[1] = id
				ret[2] = ok
				ret[3] = name
				ret[4] = array[2].(uint)
				sendData := m.encode(ret)
				core.Send(src, m.Id(), tcp.AGENT_CMD_SEND, sendData)
			} else if scmd == "forward" {
				msg := array[1].(*core.Message)
				m.forwardM(msg, data[1].([]byte))
			}
		} else if cmd == tcp.AGENT_CLOSED {
			var nodeId uint = 0
			for id, v := range m.nodesMap {
				if v == src {
					nodeId = id
				}
			}
			delete(m.nodesMap, nodeId)
		}
	}
}

func (m *master) forwardM(msg *core.Message, data []byte) {
	nodeId := core.ParseNodeId(msg.Dest)
	isLcoal := core.CheckIsLocalServiceId(msg.Dest)
	log.Debug("master forwardM is send to master: %v", isLcoal)
	if isLcoal {
		core.ForwardLocal(msg)
		return
	}
	agent, ok := m.nodesMap[nodeId]
	if !ok {
		log.Debug("node:%v is disconnected.", nodeId)
		return
	}
	if data == nil {
		ret := make([]interface{}, 2, 2)
		ret[0] = "forward"
		ret[1] = msg
		data = m.encode(ret)
	}
	core.Send(agent, m.Id(), tcp.AGENT_CMD_SEND, data)
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
