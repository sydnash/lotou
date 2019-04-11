package core

const (
	Cmd_None                    CmdType = "CmdType.Core.None"
	Cmd_Forward                 CmdType = "CmdType.Core.Forward"
	Cmd_Distribute              CmdType = "CmdType.Core.Distribute"
	Cmd_RegisterNode            CmdType = "CmdType.Core.RegisterNode"
	Cmd_RegisterNodeRet         CmdType = "CmdType.Core.RegisterNodeRet"
	Cmd_RegisterName            CmdType = "CmdType.Core.RegisterName"
	Cmd_GetIdByName             CmdType = "CmdType.Core.GetIdByName"
	Cmd_GetIdByNameRet          CmdType = "CmdType.Core.GetIdByNameRet"
	Cmd_NameAdd                 CmdType = "CmdType.Core.NameAdd"
	Cmd_NameDeleted             CmdType = "CmdType.Core.NameDeleted"
	Cmd_Exit                    CmdType = "CmdType.Core.Exit"
	Cmd_Exit_Node               CmdType = "CmdType.Core.ExitNode"
	Cmd_Default                 CmdType = "CmdType.Core.Default"
	Cmd_RefreshSlaveWhiteIPList CmdType = "CmdType.Core.RefreshSlaveWhiteIPList"
)

/*
const (
	MSG_TYPE_NORMAL = iota
	MSG_TYPE_REQUEST
	MSG_TYPE_RESPOND
	MSG_TYPE_TIMEOUT
	MSG_TYPE_CALL
	MSG_TYPE_RET
	MSG_TYPE_CLOSE
	MSG_TYPE_SOCKET
	MSG_TYPE_ERR
	MSG_TYPE_DISTRIBUTE
	MSG_TYPE_MAX
)*/

const (
	MSG_TYPE_NORMAL     MsgType = "MsgType.Normal"
	MSG_TYPE_REQUEST    MsgType = "MsgType.Request"
	MSG_TYPE_RESPOND    MsgType = "MsgType.Respond"
	MSG_TYPE_TIMEOUT    MsgType = "MsgType.TimeOut"
	MSG_TYPE_CALL       MsgType = "MsgType.Call"
	MSG_TYPE_RET        MsgType = "MsgType.Ret"
	MSG_TYPE_CLOSE      MsgType = "MsgType.Close"
	MSG_TYPE_SOCKET     MsgType = "MsgType.Socket"
	MSG_TYPE_ERR        MsgType = "MsgType.Error"
	MSG_TYPE_DISTRIBUTE MsgType = "MsgType.Distribute"
	//MSG_TYPE_MAX="MsgType.Max"
)

/*
const (
	MSG_ENC_TYPE_NO = iota
	MSG_ENC_TYPE_GO
)
*/
const (
	MSG_ENC_TYPE_NO EncType = "EncType.No"
	MSG_ENC_TYPE_GO EncType = "EncType.LotouGob"
)
