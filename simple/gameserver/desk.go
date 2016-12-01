package gameserver

import (
	"github.com/sydnash/lotou/simple/btype"
	"github.com/sydnash/lotou/simple/global"
)

type DeskConrol struct {
	rc       *RoomControl
	posInfos [4]*DeskPosInfo
}

func (dc *DeskConrol) isCaneEnter() bool {
	for _, v := range dc.posInfos {
		if !v.hasPeople {
			return true
		}
	}
	return false
}

func (dc *DeskConrol) enter(client *GameClient) {
	dc.posInfos[0] = &DeskPosInfo{client: client, pos: 0, isReady: true, hasPeople: true}
	client.deskPos = 0
	client.dc = dc

	client.gs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_ENTER_DESK
	client.gs.encoder.Encode(head)

	ret := btype.SDeskPosInfo{}
	ret.Pos = client.deskPos
	ret.AcId = client.acId
	ret.QueBi = client.playerInfo.GetPropertyInt64(global.KPropertyTypeQueBiDC)
	ret.Name = client.playerInfo.GetPropertyString(global.KPropertyTypeNickname)
	ret.IsReady = dc.posInfos[client.deskPos].isReady
	ret.Sex = client.playerInfo.GetPropertyInt32(global.KPropertyTypeSex)
	ret.IsExit = false
	ret.TitleUrl = client.playerInfo.GetPropertyString(global.KPropertyTypeTitleUrl)
	ret.Coin = client.playerInfo.GetPropertyInt64(global.KPropertyTypeCoin)
	ret.JQ = client.playerInfo.GetPropertyInt64(global.KPropertyTypeJiangQuan)
	client.gs.encoder.Encode(true)
	client.gs.encoder.Encode(1)
	client.gs.encoder.Encode(ret)
	client.gs.sendToAgent(client.agentId)
}

type DeskPosInfo struct {
	client    *GameClient
	pos       int32
	isReady   bool
	hasPeople bool
}
