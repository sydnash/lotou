package main

import (
	"encoding/json"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/simple/btype"
	"github.com/sydnash/lotou/simple/global"
	"github.com/sydnash/lotou/simple/utils"
	"math/rand"
	"strconv"
	"time"
)

const (
	kHallClientNormal = iota
	kHallClientDisConnected
	kHallClientClosed
)

type HallClient struct {
	agentId      uint
	playerInfo   *global.PropertySet
	session      int32
	acId         int32
	lastSyncTime int64
	server       *HallService
	state        int
}

func (hc *HallClient) saveInfoToString() ([]byte, error) {
	sendMap := make(map[string]string)
	for k, v := range hc.playerInfo.Property {
		base, err := global.TypeToKey(k)
		if err != nil {
			continue
		}
		kstr := base.PropertyName
		var vstr string
		switch base.ValueType {
		case global.KValueTypeInt32:
			a := []byte{}
			a = strconv.AppendInt(a, int64(v.(int32)), 10)
			vstr = string(a)
		case global.KValueTypeInt64:
			a := []byte{}
			a = strconv.AppendInt(a, v.(int64), 10)
			vstr = string(a)
		case global.KValueTypeString:
			vstr = v.(string)
		}
		sendMap[kstr] = vstr
	}
	jsonStr, err := json.Marshal(sendMap)
	return jsonStr, err
}

func (hc *HallClient) syncPlayerInfoToDb() {
	jsonStr, err := hc.saveInfoToString()
	if err != nil {
		log.Error("HallClient:update:%s", err)
	} else {
		core.Send(hc.server.dbId, hc.server.Id(), "UpdatePlayerData", hc.acId, hc.playerInfo.GetPropertyString(global.KPropertyTypeNickname), jsonStr)
	}
}

func (hc *HallClient) update() {
	now := utils.Now()
	diff := now - hc.lastSyncTime
	if diff > int64(time.Minute) {
		hc.lastSyncTime = now
		//hc.syncPlayerInfoToDb()
	}
}

func (hc *HallClient) reLogin(agentId uint) {
	hc.agentId = agentId
	hc.session = rand.Int31()
	hc.state = kHallClientNormal
}
func (hc *HallClient) init() {
	hc.session = rand.Int31()
	hc.state = kHallClientNormal
}

func (hc *HallClient) replyLogin(basic *btype.PHead) {
	hc.server.encoder.Reset()
	basic.Type = btype.S_MSG_CHECK_SESSION
	hc.server.encoder.Encode(*basic)
	var ret btype.SCheckSession
	ret.Session = hc.session
	ret.AcId = hc.acId
	ret.Coin = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeCoin)
	ret.YuanBao = 0
	ret.JQ = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeJiangQuan)
	ret.NickName = hc.playerInfo.GetPropertyString(global.KPropertyTypeNickname)
	ret.Hasdc = hc.playerInfo.HasFlag(global.KPropertyTypeFlags1, global.KFlagsIsHaveMjDC)
	ret.Tili = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeTiliDC)
	ret.MaxTili = 6
	ret.Quebi = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeQueBiDC)
	ret.Jifen = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeScoreDC)
	ret.Roomid = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeRoomIdDC)
	ret.BeiShu = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeJBSRewardBeiShu)
	ret.DayBuzu = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeDayBuZuNum)
	ret.UserType = hc.playerInfo.GetPropertyInt32(global.KPropertyTypePlayerType)
	ret.Sex = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeSex)
	ret.HeadURL = hc.playerInfo.GetPropertyString(global.KPropertyTypeTitleUrl)
	hc.server.encoder.Encode(ret)
	hc.server.sendToAgent(hc.agentId)
}

func (hc *HallClient) queryDCInfo(basic *btype.PHead) {
	hc.server.encoder.Reset()
	basic.Type = btype.S_MSG_RET_DC_INFO
	hc.server.encoder.Encode(*basic)
	var ret btype.SQueryDCInfo
	ret.Type = 0
	ret.Ok = true
	ret.Reason = 0
	ret.HasDC = hc.playerInfo.HasFlag(global.KPropertyTypeFlags1, global.KFlagsIsHaveMjDC)
	ret.HasFanpai = false
	ret.Tili = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeTiliDC)
	ret.MaxTili = 6
	ret.Quebi = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeQueBiDC)
	ret.Jifen = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeScoreDC)
	ret.Roomid = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeRoomIdDC)
	ret.BeiShu = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeJBSRewardBeiShu)
	ret.Coin = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeCoin)
	hc.server.encoder.Encode(ret)
	hc.server.sendToAgent(hc.agentId)
}

func (hc *HallClient) applyDC(basic *btype.PHead, param *btype.CQueryDCInfo) {
	hc.playerInfo.SetFlag(global.KPropertyTypeFlags1, global.KFlagsIsHaveMjDC)
	hc.playerInfo.SetPropertyInt32(global.KPropertyTypeTiliDC, 0)
	hc.playerInfo.SetPropertyInt64(global.KPropertyTypeQueBiDC, 10000)
	hc.playerInfo.SetPropertyInt32(global.KPropertyTypeScoreDC, 0)
	hc.playerInfo.SetPropertyInt32(global.KPropertyTypeRoomIdDC, param.RoomId)
	hc.playerInfo.SetPropertyInt32(global.KPropertyTypeJBSRewardBeiShu, param.BeiShu)

	hc.server.encoder.Reset()
	basic.Type = btype.S_MSG_RET_DC_INFO
	hc.server.encoder.Encode(*basic)
	var ret btype.SQueryDCInfo
	ret.Type = 1
	ret.Ok = true
	ret.Reason = 0
	ret.HasDC = hc.playerInfo.HasFlag(global.KPropertyTypeFlags1, global.KFlagsIsHaveMjDC)
	ret.HasFanpai = false
	ret.Tili = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeTiliDC)
	ret.MaxTili = 6
	ret.Quebi = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeQueBiDC)
	ret.Jifen = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeScoreDC)
	ret.Roomid = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeRoomIdDC)
	ret.BeiShu = hc.playerInfo.GetPropertyInt32(global.KPropertyTypeJBSRewardBeiShu)
	ret.Coin = hc.playerInfo.GetPropertyInt64(global.KPropertyTypeCoin)
	hc.server.encoder.Encode(ret)
	hc.server.sendToAgent(hc.agentId)
}
