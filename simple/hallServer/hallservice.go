package main

import (
	"encoding/json"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/encoding/binary"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/network/tcp"
	"github.com/sydnash/lotou/simple/btype"
	"github.com/sydnash/lotou/simple/global"
	"reflect"
	"strconv"
	"time"
)

type HallService struct {
	*core.Base
	platId    uint
	dbId      uint
	decoder   *binary.Decoder
	encoder   *binary.Encoder
	ticker    *time.Ticker
	clientMap map[int32]*HallClient
}

func (hs *HallService) CloseMSG(dest, src uint) {
	hs.Base.Close()
	hs.ticker.Stop()
}
func (hs *HallService) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	log.Info("normalMSG:%x, %x, %v", src, dest, data)
	if msgType == "socket" {
		cmd := data[0].(int)
		var d []byte
		if len(data) >= 2 {
			d = data[1].([]byte)
		}
		hs.socketMSG(src, cmd, d)
	} else if msgType == "go" {
		cmd := data[0].(string)
		psv := reflect.ValueOf(hs)
		fv := psv.MethodByName(cmd)
		if fv.IsValid() {
			in := make([]reflect.Value, len(data)-1)
			for i := 1; i < len(data); i++ {
				in[i-1] = reflect.ValueOf(data[i])
			}
			fv.Call(in)
		} else {
			//core.Respond(src, dest, rid, ""
			log.Error("function:%s not exist.", cmd)
		}
	}
}
func (hs *HallService) CallMSG(dest, src uint, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)
	core.Ret(src, dest, data...)
}
func (hs *HallService) RequestMSG(dest, src uint, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	cmd := data[0].(string)
	if cmd == "GameServerLogin" {
		acId := data[1].(int32)
		isNeedPlayerInfo := data[2].(bool)
		client, ok := hs.clientMap[acId]
		if ok {
			if isNeedPlayerInfo {
				jsonstr, err := client.saveInfoToString()
				if err != nil {
					core.Respond(src, dest, rid, false, int32(0), []byte{})
				} else {
					core.Respond(src, dest, rid, true, client.session, jsonstr)
				}
			} else {
				core.Respond(src, dest, rid, true, client.session, []byte{})
			}
		} else {
			core.Respond(src, dest, rid, false, int32(0), []byte{})
		}
		return
	}
	core.Respond(src, dest, rid, data...)
}

func (hs *HallService) socketMSG(src uint, cmd int, data []byte) {
	switch cmd {
	case tcp.AGENT_DATA:
		hs.socketData(src, data)
	}
}

func (hs *HallService) socketData(src uint, data []byte) {
	var basic btype.PHead
	hs.decoder.SetBuffer(data)
	hs.decoder.Decode(&basic)
	log.Debug("recv package: %v", basic)
	ctype := basic.Type
	switch ctype {
	case btype.C_MSG_CHECK_SESSION:
		hs.ccheckSession(src, &basic)
	case btype.C_MSG_QUERY_DC_INFO:
		hs.queryDCInfo(src, &basic)
	}
}

func (hs *HallService) queryDCInfo(src uint, basic *btype.PHead) {
	var param btype.CQueryDCInfo

	client := hs.getPlayer(src, basic)
	if client == nil {
		return
	}
	switch param.Type {
	case 0: //query
		client.queryDCInfo(basic)
	case 1: //apply
		client.applyDC(basic, &param)
	}
}
func (hs *HallService) getPlayer(src uint, basic *btype.PHead) *HallClient {
	client, ok := hs.clientMap[basic.AcId]
	if ok {
		if basic.Session == client.session {
			return client
		} else {
			basic.Type = btype.S_MSG_SESSION_ERROR
			hs.encoder.Reset()
			hs.encoder.Encode(*basic)
			hs.sendToAgent(src)
			log.Debug("session error")
			time.AfterFunc(time.Second*2, func() {
				core.Close(src, hs.Id())
			})
		}
	} else {
		basic.Type = btype.S_MSG_SESSION_ERROR
		hs.encoder.Reset()
		hs.encoder.Encode(*basic)
		hs.sendToAgent(src)
		log.Debug("session error")
		time.AfterFunc(time.Second*2, func() {
			core.Close(src, hs.Id())
		})
		core.Close(src, hs.Id())
	}
	return nil
}

func (hs *HallService) ccheckSession(src uint, basic *btype.PHead) {
	var param btype.CCheckSession
	hs.decoder.Decode(&param)
	log.Debug("check session param:%v", param)

	cb := func(ok bool) {
		log.Info("check session isok:%v", ok)
		if ok {

			loginCB := func(acId int32, nicknamestr, datastr, namestr string, playType, accountType, qdType int32, macstr string) {
				log.Info("acid:%d", acId)
				if acId == param.AcId {
					//init player data
					playerInfo := global.NewPropertySet()
					playerInfo.LoadJson(datastr)
					playerInfo.SetPropertyString(global.KPropertyTypeNickname, nicknamestr)
					playerInfo.SetPropertyString(global.KPropertyTypeAcName, namestr)
					playerInfo.SetPropertyString(global.KPropertyTypeMac, macstr)
					playerInfo.SetPropertyInt32(global.KPropertyTypePlayerType, playType)
					playerInfo.SetPropertyInt32(global.KPropertyTypeAccountType, accountType)
					playerInfo.SetPropertyInt32(global.KPropertyTypeQuDaoType, qdType)
					if !playerInfo.HasFlag(global.KPropertyTypeFlags1, global.KFlagsIsFirstChuangJian) {
						playerInfo.SetPropertyInt64(global.KPropertyTypeCoin, 10000)
						playerInfo.SetPropertyInt64(global.KPropertyTypeQueBiDC, 10000)
						playerInfo.SetPropertyInt64(global.KPropertyTypeJiangQuan, 10000)
						playerInfo.SetFlag(global.KPropertyTypeFlags1, global.KFlagsIsFirstChuangJian)
					}
					client := &HallClient{agentId: src, playerInfo: playerInfo, acId: acId, server: hs}
					client.init()
					client.replyLogin(basic)
					hs.clientMap[acId] = client
					hs.sendJsonToClient(src)
					hs.sendLobbyToClient(src)
				} else {
					hs.sendJsonToClient(src)
					core.Close(src, hs.Id())
				}
			}
			client, ok := hs.clientMap[param.AcId]
			if ok {
				client.reLogin(src)
				client.replyLogin(basic)
				hs.sendJsonToClient(src)
				hs.sendLobbyToClient(src)
			} else {
				core.Request(hs.dbId, hs, loginCB, "PlayerLogin", param.AcId)
			}
		} else {
			core.Close(src, hs.Id())
		}
	}
	timeout := func() {
		core.Close(src, hs.Id())
	}
	session, _ := strconv.ParseUint(param.Session, 10, 64)
	log.Debug("platid:%x, src id :%x", hs.platId, hs.Id())
	core.RequestTimeout(hs.platId, hs, cb, timeout, 20000, "CheckSeesion", int(param.AcId), session)
}

func (hs *HallService) initPlayerInfo() {
}

func (hs *HallService) sendToAgent(dest uint) {
	hs.encoder.UpdateLen()
	b := hs.encoder.Buffer()
	nb := make([]byte, len(b))
	copy(nb, b)
	core.Send(dest, hs.Id(), tcp.AGENT_CMD_SEND, nb)
}

func (hs *HallService) sendJsonToClient(dest uint) {
	t1 := `{"needMinCoin":20000,"buzuCoin":10000,"buzuNumMax":4}`
	t2 := `{"coin":1,"num":30}`
	t3 := `[1,2,5,10,20,50,100,200,500,1000]`
	t4 := `[{"jifenMax":80,"jifenMin":-10000,"jiangQuanMax":0,"jiangQuan":0,"rate":5180,"lv":0},{"jifenMax":160,"jifenMin":80,"jiangQuanMax":38000,"jiangQuan":4500,"rate":2800,"lv":5},{"jifenMax":240,"jifenMin":160,"jiangQuanMax":50000,"jiangQuan":6000,"rate":1000,"lv":4},{"jifenMax":320,"jifenMin":240,"jiangQuanMax":76000,"jiangQuan":9000,"rate":800,"lv":3},{"jifenMax":640,"jifenMin":320,"jiangQuanMax":102000,"jiangQuan":12000,"rate":200,"lv":2},{"jifenMax":2000000,"jifenMin":640,"jiangQuanMax":128000,"jiangQuan":15000,"rate":20,"lv":1}]`
	t5 := `{"daily":{"timeDesc":"00:00-24:00","title":"全天开放","timeValue":""},"final":{"timeDesc":"线下比赛","title":"11月31日","timeValue":1480564800},"month":{"timeDesc":"00:00-24:00","title":"9月6日","timeValue":1475251200}}`
	t6 := `[{"coin":300,"ExChangeRate":50,"isWeb":1,"hotSale":0,"GivenRate":0,"billType":1,"rmb":1,"billDesc":"巴适游戏","QuDaoType":1000},{"coin":1800,"ExChangeRate":60,"isWeb":1,"hotSale":0,"GivenRate":10,"billType":2,"rmb":30,"billDesc":"巴适游戏","QuDaoType":1000},{"coin":4420,"ExChangeRate":65,"isWeb":1,"hotSale":1,"GivenRate":0,"billType":3,"rmb":68,"billDesc":"巴适游戏","QuDaoType":1000},{"coin":12800,"ExChangeRate":100,"isWeb":1,"hotSale":0,"GivenRate":50,"billType":4,"rmb":128,"billDesc":"巴适游戏","QuDaoType":1000},{"coin":20790,"ExChangeRate":105,"isWeb":1,"hotSale":1,"GivenRate":0,"billType":5,"rmb":198,"billDesc":"巴适游戏","QuDaoType":1000},{"coin":36080,"ExChangeRate":110,"isWeb":1,"hotSale":1,"GivenRate":0,"billType":6,"rmb":328,"billDesc":"巴适游戏","QuDaoType":1000}]`
	hs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_SEND_JSON
	hs.encoder.Encode(head)
	hs.encoder.Encode(t1)
	hs.encoder.Encode(t2)
	hs.encoder.Encode(t3)
	hs.encoder.Encode(t4)
	hs.encoder.Encode(t5)
	hs.encoder.Encode(t6)
	hs.sendToAgent(dest)
}

var LobbySendOrder = [][]interface{}{
	[]interface{}{"roomType", global.KValueTypeInt32},
	[]interface{}{"roomId", global.KValueTypeInt32},
	[]interface{}{"coinMin", global.KValueTypeInt64},
	[]interface{}{"coinMax", global.KValueTypeInt64},
	[]interface{}{"diZhuCoin", global.KValueTypeInt64},
	[]interface{}{"maxFan", global.KValueTypeInt32},
	[]interface{}{"clientNum", global.KValueTypeInt16},
	[]interface{}{"clientMax", global.KValueTypeInt16},
	[]interface{}{"startMinNum", global.KValueTypeInt16},
	[]interface{}{"port", global.KValueTypeInt16},
	[]interface{}{"ip", global.KValueTypeString},
	[]interface{}{"coinSetp", global.KValueTypeInt32},
	[]interface{}{"baoMingCoin", global.KValueTypeInt64},
	[]interface{}{"duiHuaTime", global.KValueTypeInt32},
}

func (hs *HallService) sendLobbyToClient(dest uint) {
	t1 := `[{"dingqueTime":20000,"roomType":0,"coinMax":1000,"maxFan":8,"diZhuCoin":10,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":0,"daidaTime":10000,"coinMin":100,"jiesuanTime":0,"zhunbeiTime":10000,"roomId":1,"startMinNum":4},{"dingqueTime":20000,"roomType":0,"coinMax":10000,"maxFan":8,"diZhuCoin":100,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":1,"daidaTime":10000,"coinMin":1000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":2,"startMinNum":4},{"dingqueTime":30000,"roomType":1,"coinMax":1000,"maxFan":8,"diZhuCoin":10,"isNeedRobot":1,"coinType":1002,"duiHuaTime":15000,"firstFapaiTime":10000,"clientMax":5000,"clientNum":0,"jifenBeiLv":10,"dapaiTime":40000,"daidaTime":10000,"coinSetp":0,"baoMingCoin":100,"coinMin":100,"jiesuanTime":15000,"zhunbeiTime":10000,"roomId":3,"startMinNum":4},{"dingqueTime":20000,"roomType":0,"coinMax":100000,"maxFan":8,"diZhuCoin":800,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":2,"daidaTime":10000,"coinMin":50000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":6,"startMinNum":2},{"dingqueTime":20000,"roomType":0,"coinMax":200000,"maxFan":8,"diZhuCoin":2000,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":3,"daidaTime":10000,"coinMin":100000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":7,"startMinNum":2},{"dingqueTime":20000,"roomType":0,"coinMax":-1,"maxFan":8,"diZhuCoin":10000,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":4,"daidaTime":10000,"coinMin":200000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":8,"startMinNum":2}]`
	var pasre []map[string]interface{}
	json.Unmarshal([]byte(t1), &pasre)
	hs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_LOBBY_INFO
	hs.encoder.Encode(head)
	hs.encoder.Encode(len(pasre))
	for _, v := range pasre {
		v["ip"] = "127.0.0.1"
		v["port"] = float64(40001)
		for _, typeDefine := range LobbySendOrder {
			value, ok := v[typeDefine[0].(string)]
			if ok {
				switch typeDefine[1].(int) {
				case global.KValueTypeInt16:
					num := value.(float64)
					hs.encoder.Encode(int16(num))
				case global.KValueTypeInt32:
					num := value.(float64)
					hs.encoder.Encode(int32(num))
				case global.KValueTypeInt64:
					num := value.(float64)
					hs.encoder.Encode(int64(num))
				case global.KValueTypeString:
					str := value.(string)
					hs.encoder.Encode(str)
				}
			}
		}
	}
	hs.sendToAgent(dest)
}

func NewHS(platid, dbid uint) *HallService {
	hs := &HallService{Base: core.NewBaseLen(1024 * 1024)}
	hs.platId = platid
	hs.dbId = dbid
	hs.clientMap = make(map[int32]*HallClient)
	hs.decoder = binary.NewDecoder()
	hs.encoder = binary.NewEncoder()
	hs.SetDispatcher(hs)
	return hs
}

func (hs *HallService) Run() {
	core.RegisterService(hs)
	core.Name(hs.Id(), "hallService")
	hs.ticker = time.NewTicker(time.Millisecond * 10)
	go func() {
	EXIT:
		for {
			_, ok := <-hs.ticker.C
			if !ok {
				break
			}
			//loop for msg
		MSGLOOP:
			for {
				select {
				case msg, ok := <-hs.In():
					if !ok {
						break EXIT
					}
					hs.DispatchM(msg)
				default:
					break MSGLOOP
				}
			}
			//loop for client
			for _, cli := range hs.clientMap {
				cli.update()
			}
		}
	}()

	s := tcp.New("", "20001", hs.Id())
	s.Listen()
}
