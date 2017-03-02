package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type slave struct {
	*core.Skeleton
	client uint
}

func StartSlave(ip, port string) {
	m := &slave{Skeleton: core.NewSkeleton(0)}
	core.StartService(".slave", m)
	c := tcp.NewClient(ip, port, m.Id)
	m.client = core.StartService("", c)
}

func (s *slave) OnNormalMSG(dst uint, data ...interface{}) {
	//dest is master's id, src is core's id
	//data[0] is cmd such as (registerNode, regeisterName, getIdByName...)
	//data[1] is dest nodeService's id
	//parse node's id, and choose correct agent and send msg to that node.
	t1 := gob.Pack(data)
	s.RawSend(s.client, core.MSG_TYPE_NORMAL, tcp.CLIENT_CMD_SEND, t1)
}
func (s *slave) OnSocketMSG(dst uint, data ...interface{}) {
	//dest is master's id, src is agent's id
	//data[0] is socket status
	//data[1] is a gob encode data
	//it's first encode value is cmd such as (registerNodeRet, regeisterNameRet, getIdByNameRet, forword...)
	//it's second encode value is dest service's id.
	//find correct agent and send msg to that node.
	cmd := data[0].(int)
	if cmd == tcp.CLIENT_DATA {
		sdata := gob.Unpack(data[1].([]byte))
		array := sdata.([]interface{})
		scmd := array[0].(string)
		if scmd == "registerNodeRet" {
			nodeId := array[1].(uint)
			core.RegisterNodeRet(nodeId)
		} else if scmd == "distibute" {
			data := array[1].([]interface{})
			core.DistributeMSG(s.Id, data...)
		} else if scmd == "getIdByNameRet" {
			id := array[1].(uint)
			ok := array[2].(bool)
			name := array[3].(string)
			rid := array[4].(uint)
			core.GetIdByNameRet(id, ok, name, rid)
		} else if scmd == "forward" {
			msg := array[1].(*core.Message)
			s.forwardM(msg)
		}
	}
}

func (s *slave) forwardM(msg *core.Message) {
	isLcoal := core.CheckIsLocalServiceId(msg.Dst)
	if isLcoal {
		core.ForwardLocal(msg)
		return
	}
	log.Warn("recv msg not forward to this node.")
}

func (s *slave) OnDestroy() {
	s.SendClose(s.client)
}
