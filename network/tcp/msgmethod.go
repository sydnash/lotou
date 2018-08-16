package tcp

import (
	"github.com/sydnash/lotou/core"
)

const (
	AGENT_CLOSED             core.CmdType = "CmdType.Tcp.AgentClosed"
	AGENT_DATA               core.CmdType = "CmdType.Tcp.AgentData"
	AGENT_ARRIVE             core.CmdType = "CmdType.Tcp.AgentArrive"
	AGENT_CMD_SEND           core.CmdType = "CmdType.Tcp.AgentSend"
	CLIENT_CONNECT_FAILED    core.CmdType = "CmdType.Tcp.ClientConnectFailed"
	CLIENT_CONNECTED         core.CmdType = "CmdType.Tcp.ClientConnected"
	CLIENT_SELF_CONNECTED    core.CmdType = "CmdType.Tcp.ClientSelfConnected"
	CLIENT_SELF_DISCONNECTED core.CmdType = "CmdType.Tcp.ClientSelfDisconnected"
	CLIENT_DISCONNECTED      core.CmdType = "CmdType.Tcp.ClientDisconnected"
	CLIENT_DATA              core.CmdType = "CmdType.Tcp.ClientData"
	CLIENT_CMD_CONNECT       core.CmdType = "CmdType.Tcp.ClientConnect"
	CLIENT_CMD_SEND          core.CmdType = "CmdType.Tcp.ClientSend"
)

const (
	MAX_PACKET_LEN        = 1024 * 1024 * 100 //100M
	DEFAULT_RECV_BUFF_LEN = 1 * 1024
)
