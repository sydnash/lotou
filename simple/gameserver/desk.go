package gameserver

import (
	"github.com/sydnash/lotou/simple/btype"
	"github.com/sydnash/lotou/simple/global"
)

const (
	kDeskStateNone = iota
	KDeskStateWaitReady
	KDeskStateDingQue
	KDeskStatePlaying
	KDeskStateEnded
)

type DeskConrol struct {
	rc           *RoomControl
	posInfos     [4]*DeskPosInfo
	deskId       int32
	deskState    int32
	curHasPeople int16
	mjLogicInfo  MjLogicInfo
}

func (dc *DeskConrol) isCaneEnter() bool {
	if dc.deskState > KDeskStateWaitReady && dc.deskState < KDeskStateEnded {
		return false
	}
	if dc.curHasPeople >= dc.rc.roomInfo.StartMinMax {
		return false
	}
	for _, v := range dc.posInfos {
		if !v.hasPeople {
			return true
		}
	}
	return false
}

func (dc *DeskConrol) getEnterPos() int32 {
	for pos, v := range dc.posInfos {
		if !v.hasPeople {
			return int32(pos)
		}
	}
	panic("there is no pos to enter.")
}

func (dc *DeskConrol) exit(client *GameClient) {
	pos := client.deskPos
	dc.posInfos[pos].hasPeople = false
	client.gs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_OTHER_EXIT_DESK
	client.gs.encoder.Encode(head)
	client.gs.encoder.Encode(pos)
	for _, v := range dc.posInfos {
		if v.hasPeople {
			client.gs.sendToAgent(v.client.agentId)
		}
	}

	client.gs.encoder.Reset()
	head.Type = btype.S_MSG_EXIT_DESK
	client.gs.encoder.Encode(head)
	client.gs.encoder.Encode(true)
	client.gs.encoder.Encode(0)
	client.gs.sendToAgent(client.agentId)

	dc.curHasPeople--
	if dc.curHasPeople < 0 {
		dc.curHasPeople = 0
	}
	if dc.curHasPeople == 0 {
		dc.deskState = kDeskStateNone
	}
}

func (dc *DeskConrol) sendDeskBaseInfo(client *GameClient) {
	client.gs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_DESK_BASIC_INFO
	client.gs.encoder.Encode(head)
	client.gs.encoder.Encode(*(dc.rc.roomInfo))
	client.gs.sendToAgent(client.agentId)
}

func (dc *DeskConrol) noticePeopleEnter(client *GameClient) {
	client.gs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_ENTER_DESK
	client.gs.encoder.Encode(head)

	client.gs.encoder.Encode(true)

	cli := []btype.SDeskPosInfo{}
	var nenter *btype.SDeskPosInfo
	for _, v := range dc.posInfos {
		if v.hasPeople {
			c := v.client
			ret := btype.SDeskPosInfo{}
			ret.Pos = c.deskPos
			ret.AcId = c.acId
			ret.QueBi = c.playerInfo.GetPropertyInt64(global.KPropertyTypeQueBiDC)
			ret.Name = c.playerInfo.GetPropertyString(global.KPropertyTypeNickname)
			ret.IsReady = dc.posInfos[c.deskPos].isReady
			ret.Sex = c.playerInfo.GetPropertyInt32(global.KPropertyTypeSex)
			ret.IsExit = false
			ret.TitleUrl = c.playerInfo.GetPropertyString(global.KPropertyTypeTitleUrl)
			ret.Coin = c.playerInfo.GetPropertyInt64(global.KPropertyTypeCoin)
			ret.JQ = c.playerInfo.GetPropertyInt64(global.KPropertyTypeJiangQuan)
			if c.acId == client.acId {
				nenter = &ret
			}
			cli = append(cli, ret)
		}
	}

	client.gs.encoder.Encode(cli)
	client.gs.sendToAgent(client.agentId)

	client.gs.encoder.Reset()
	head.Type = btype.S_MSG_OTHER_ENTER_DESK
	client.gs.encoder.Encode(head)
	client.gs.encoder.Encode(1)
	client.gs.encoder.Encode(*nenter)
	for _, v := range dc.posInfos {
		if v.hasPeople && v.client.acId != client.acId {
			client.gs.sendToAgent(v.client.agentId)
		}
	}

	dc.curHasPeople++
	dc.sendDeskBaseInfo(client)
}
func (dc *DeskConrol) enter(client *GameClient) {
	pos := dc.getEnterPos()
	dc.posInfos[pos] = &DeskPosInfo{client: client, pos: pos, isReady: true, hasPeople: true}
	client.deskPos = pos
	client.dc = dc
	dc.noticePeopleEnter(client)

	dc.checkIsCanStart()
}

type DeskPosInfo struct {
	client         *GameClient
	pos            int32
	isReady        bool
	hasPeople      bool
	mjLogicPosInfo *MJLogicPosInfo
}

func NewDC(rc *RoomControl) *DeskConrol {
	dc := &DeskConrol{}
	dc.rc = rc
	for i, p := range dc.posInfos {
		p = &DeskPosInfo{}
		p.pos = int32(i)
		p.hasPeople = false
		dc.posInfos[i] = p
	}
	dc.mjLogicInfo.init()
	dc.deskState = kDeskStateNone
	dc.mjLogicInfo.dc = dc
	dc.curHasPeople = 0
	return dc
}
