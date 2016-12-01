package gameserver

type DeskConrol struct {
	rc       *RoomControl
	posInfos [4]*DeskPosInfo
}

type DeskPosInfo struct {
	client *GameClient
	pos    int32
	isReay int32
}
