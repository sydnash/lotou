package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type master struct {
	*core.Skeleton
	nodesMap      map[uint64]core.ServiceID //nodeid : agentSID
	globalNameMap map[string]core.ServiceID
	tcpServer     *tcp.Server
}

func StartMaster(ip, port string) {
	m := &master{Skeleton: core.NewSkeleton(0)}
	m.nodesMap = make(map[uint64]core.ServiceID)
	m.globalNameMap = make(map[string]core.ServiceID)
	core.StartService(".router", m)

	m.tcpServer = tcp.NewServer(ip, port, m.Id)
	m.tcpServer.Listen()
}

func (m *master) OnNormalMSG(msg *core.Message) {
	//cmd such as (registerName, getIdByName, syncName, forward ...)
	cmd := msg.MethodId.(string)
	data := msg.Data
	if cmd == "forward" {
		msg := data[0].(*core.Message)
		m.forwardM(msg, nil)
	} else if cmd == "registerName" {
		id := data[0].(uint64)
		name := data[1].(string)
		m.onRegisterName(core.ServiceID(id), name)
	} else if cmd == "getIdByName" {
		name := data[0].(string)
		rid := data[1].(uint)
		id, ok := m.globalNameMap[name]
		core.DispatchGetIdByNameRet(id, ok, name, rid)
	}
}

func (m *master) onRegisterNode(src core.ServiceID) {
	//generate node id
	nodeId := core.GenerateNodeId()
	m.nodesMap[nodeId] = src
	msg := core.NewMessage(0, 0, 0, 0, 0, "registerNodeRet", nodeId)
	sendData := gob.Pack(msg)
	m.RawSend(src, core.MSG_TYPE_NORMAL, tcp.AGENT_CMD_SEND, sendData)
}

func (m *master) onRegisterName(serviceId core.ServiceID, serviceName string) {
	m.globalNameMap[serviceName] = serviceId
	m.distributeM("nameAdd", serviceName, serviceId)
}

func (m *master) onGetIdByName(src core.ServiceID, name string, rId uint) {
	id, ok := m.globalNameMap[name]
	msg := core.NewMessage(0, 0, 0, 0, 0, "getIdByNameRet", id, ok, name, rId)
	sendData := gob.Pack(msg)
	m.RawSend(src, core.MSG_TYPE_NORMAL, tcp.AGENT_CMD_SEND, sendData)
}

func (m *master) OnSocketMSG(msg *core.Message) {
	//src is slave's agent's serviceid
	src := msg.Src
	//cmd is socket status
	cmd := msg.MethodId.(int)
	//data[0] is a gob encode with message
	data := msg.Data
	//it's first encode value is cmd such as (registerNode, regeisterName, getIdByName, forword...)
	if cmd == tcp.AGENT_DATA {
		sdata := gob.Unpack(data[0].([]byte))
		slaveMSG := sdata.([]interface{})[0].(*core.Message)
		scmd := slaveMSG.MethodId.(string)
		array := slaveMSG.Data
		if scmd == "registerNode" {
			m.onRegisterNode(src)
		} else if scmd == "registerName" {
			serviceId := array[0].(uint64)
			serviceName := array[1].(string)
			m.onRegisterName(core.ServiceID(serviceId), serviceName)
		} else if scmd == "getIdByName" {
			name := array[0].(string)
			rId := array[1].(uint)
			m.onGetIdByName(src, name, rId)
		} else if scmd == "forward" { //find correct agent and send msg to that node.
			forwardMsg := array[0].(*core.Message)
			m.forwardM(forwardMsg, data[0].([]byte))
		}
	} else if cmd == tcp.AGENT_CLOSED {
		//on agent disconnected
		//delet node from nodesMap
		var nodeId uint64 = 0
		for id, v := range m.nodesMap {
			if v == src {
				nodeId = id
			}
		}
		delete(m.nodesMap, nodeId)

		//notify other services delete name's id on agent which is disconnected.
		deletedNames := []string{}
		for name, id := range m.globalNameMap {
			nid := core.ParseNodeId(id)
			if nid == nodeId {
				deletedNames = append(deletedNames, name)
			}
		}
		m.distributeM("nameDeleted", deletedNames)
	}
}

func (m *master) distributeM(methodId string, data ...interface{}) {
	for _, agent := range m.nodesMap {
		msg := &core.Message{}
		msg.MethodId = "distribute"
		msg.Data = append(msg.Data, methodId)
		msg.Data = append(msg.Data, data...)
		sendData := gob.Pack(msg)
		m.RawSend(agent, core.MSG_TYPE_NORMAL, tcp.AGENT_CMD_SEND, sendData)
	}
	core.DistributeMSG(m.Id, methodId, data...)
}

func (m *master) forwardM(msg *core.Message, data []byte) {
	nodeId := core.ParseNodeId(core.ServiceID(msg.Dst))
	isLcoal := core.CheckIsLocalServiceId(core.ServiceID(msg.Dst))
	//log.Debug("master forwardM is send to master: %v, nodeid: %d", isLcoal, nodeId)
	if isLcoal {
		core.ForwardLocal(msg)
		return
	}
	agent, ok := m.nodesMap[nodeId]
	if !ok {
		log.Debug("node:%v is disconnected.", nodeId)
		return
	}
	//if has no encode data, encode it first.
	if data == nil {
		ret := &core.Message{}
		ret.MethodId = "forward"
		ret.Data = append(ret.Data, msg)
		data = gob.Pack(ret)
	}
	m.RawSend(agent, core.MSG_TYPE_NORMAL, tcp.AGENT_CMD_SEND, data)
}

func (m *master) OnDestroy() {
	if m.tcpServer != nil {
		m.tcpServer.Close()
	}
	for _, v := range m.nodesMap {
		m.SendClose(v, false)
	}
}
