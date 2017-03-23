package core

const (
	Cmd_Forward         = "forward"
	Cmd_Distribute      = "distribute"
	Cmd_RegisterNode    = "registerNode"
	Cmd_RegisterNodeRet = "registerNodeRet"
	Cmd_RegisterName    = "registerName"
	Cmd_GetIdByName     = "getIdByName"
	Cmd_GetIdByNameRet  = "getIdByNameRet"
	Cmd_NameAdd         = "nameAdd"
	Cmd_NameDeleted     = "nameDeleted"
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
