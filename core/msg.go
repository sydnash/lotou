package core

const (
	MSG_TYPE_NORMAL = iota
	MSG_TYPE_REQUEST
	MSG_TYPE_RESPOND
	MSG_TYPE_TIMEOUT
	MSG_TYPE_CALL
	MSG_TYPE_RET
	MSG_TYPE_CLOSE
	MSG_TYPE_ERR
	MSG_TYPE_MAX
)

type Message struct {
	Dest          uint
	Src           uint
	Type          int
	MsgEncodeType string
	Data          []interface{}
}
