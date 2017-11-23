package topology

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/gob"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
)

type slave struct {
	*core.Skeleton
	client core.ServiceID
}

func StartSlave(ip, port string) {
	m := &slave{Skeleton: core.NewSkeleton(0)}
	core.StartService(&core.ModuleParam{
		N: ".router",
		M: m,
		L: 0,
	})
	c := tcp.NewClient(ip, port, m.Id)
	m.client = core.StartService(&core.ModuleParam{
		N: "",
		M: c,
		L: 0,
	})
}

func (s *slave) OnNormalMSG(msg *core.Message) {
	//dest is master's id, src is core's id
	//data[0] is cmd such as (registerNode, regeisterName, getIdByName...)
	t1 := gob.Pack(msg)
	s.RawSend(s.client, core.MSG_TYPE_NORMAL, tcp.CLIENT_CMD_SEND, t1)
}
func (s *slave) OnSocketMSG(msg *core.Message) {
	//cmd is socket status
	cmd := msg.Cmd
	//data[0] is a gob encode data of Message
	data := msg.Data
	if cmd == tcp.CLIENT_DATA {
		sdata, err := gob.Unpack(data[0].([]byte))
		if err != nil {
			return
		}
		masterMSG := sdata.([]interface{})[0].(*core.Message)
		scmd := masterMSG.Cmd
		array := masterMSG.Data
		switch scmd {
		case core.Cmd_RegisterNodeRet:
			nodeId := array[0].(uint64)
			core.DispatchRegisterNodeRet(nodeId)
		case core.Cmd_Distribute:
			core.DistributeMSG(s.Id, core.CmdType(array[0].(string)), array[1:]...)
		case core.Cmd_GetIdByNameRet:
			id := array[0].(uint64)
			ok := array[1].(bool)
			name := array[2].(string)
			rid := array[3].(uint)
			core.DispatchGetIdByNameRet(core.ServiceID(id), ok, name, rid)
		case core.Cmd_Forward:
			msg := array[0].(*core.Message)
			s.forwardM(msg)
		case core.Cmd_Exit:
			log.Info("receive exit command, node will exit now.")
			core.SendCloseToAll()
		}
	}
}

func (s *slave) forwardM(msg *core.Message) {
	isLcoal := core.CheckIsLocalServiceId(core.ServiceID(msg.Dst))
	if isLcoal {
		core.ForwardLocal(msg)
		return
	}
	log.Warn("recv msg not forward to this node.")
}

func (s *slave) OnDestroy() {
	s.SendClose(s.client, false)
}
