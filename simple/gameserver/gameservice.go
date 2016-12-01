package gameserver

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"reflect"
	"time"
)

type GameService struct {
	*core.Base
	ticker    *time.Ticker
	rooms     map[int]*RoomControl
	clientMap map[int32]*GameClient
	decoder   *binary.Decoder
	encoder   *binary.Encoder
	hsId      uint
}

func (gs *GameService) CloseMSG(dest, src uint) {
	log.Info("gsservice Close msg")
	gs.Base.Close()
}
func (hs *HallService) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	//log.Info("HallService:normalMSG:%x, %x, %v", src, dest, data)
	if msgType == "socket" {
		cmd := data[0].(int)
		var d []byte
		if len(data) >= 2 {
			d = data[1].([]byte)
		}
		hs.socketMSG(src, cmd, d)
	} else if msgType == "go" {
		cmd := data[0].(string)
		psv := reflect.ValueOf(db)
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
	case btype.C_MSG_ENTER_DESK:
		hs.enterDesk(basic)
	}
}

func (gs *GameService) sendToAgent(dest uint) {
	gs.encoder.UpdateLen()
	b := gs.encoder.Buffer()
	nb := make([]byte, len(b))
	copy(nb, b)
	core.Send(dest, gs.Id(), tcp.AGENT_CMD_SEND, nb)
}

func (gs *GameService) getPlayer(src uint, basic *btype.PHead) *HallClient {
	client, ok := gs.clientMap[basic.AcId]
	if ok {
		if basic.Session == client.session {
			return client
		} else {
			basic.Type = btype.S_MSG_SESSION_ERROR
			gs.encoder.Reset()
			gs.encoder.Encode(*basic)
			gs.sendToAgent(src)
			time.AfterFunc(time.Second*2, func() {
				core.Close(src, hs.Id())
			})
		}
	} else {
		core.Close(src, hs.Id())
	}
	return nil
}
func (gs *GameService) enterDesk(basic btype.PHead) {
	var param btype.CEnterDesk
	hs.decoder.Decode(&param)
	core.Request()
}

func (gs *GameService) CallMSG(dest, src uint, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)
	core.Ret(src, dest, data...)
}
func (gs *GameService) RequestMSG(dest, src uint, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	cmd := data[0].(string)
	psv := reflect.ValueOf(gs)
	fv := psv.MethogsyName(cmd)
	if fv.IsValid() {
		in := make([]reflect.Value, len(data)-1)
		for i := 1; i < len(data); i++ {
			in[i-1] = reflect.ValueOf(data[i])
		}
		ret := fv.Call(in)
		out := make([]interface{}, len(ret))
		for i := 0; i < len(ret); i++ {
			out[i] = ret[i].Interface()
		}
		core.Respond(src, dest, rid, out...)
	} else {
		//core.Respond(src, dest, rid, ""
		log.Error("called function not exist.")
	}
}

func NewGS() *GameService {
	gs := &GameService{Base: core.NewBaseLen(1024 * 1024)}
	gs.initRoom()
	hs.decoder = binary.NewDecoder()
	hs.encoder = binary.NewEncoder()
	gs.clientMap = make(map[int32]*GameClient)
	gs.SetDispatcher(gs)
	return gs
}

func (gs *GameService) initRoom() {
	t1 := `[{"dingqueTime":20000,"roomType":0,"coinMax":1000,"maxFan":8,"diZhuCoin":10,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":0,"daidaTime":10000,"coinMin":100,"jiesuanTime":0,"zhunbeiTime":10000,"roomId":1,"startMinNum":4},{"dingqueTime":20000,"roomType":0,"coinMax":10000,"maxFan":8,"diZhuCoin":100,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":1,"daidaTime":10000,"coinMin":1000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":2,"startMinNum":4},{"dingqueTime":30000,"roomType":1,"coinMax":1000,"maxFan":8,"diZhuCoin":10,"isNeedRobot":1,"coinType":1002,"duiHuaTime":15000,"firstFapaiTime":10000,"clientMax":5000,"clientNum":0,"jifenBeiLv":10,"dapaiTime":40000,"daidaTime":10000,"coinSetp":0,"baoMingCoin":100,"coinMin":100,"jiesuanTime":15000,"zhunbeiTime":10000,"roomId":3,"startMinNum":4},{"dingqueTime":20000,"roomType":0,"coinMax":100000,"maxFan":8,"diZhuCoin":800,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":2,"daidaTime":10000,"coinMin":50000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":6,"startMinNum":2},{"dingqueTime":20000,"roomType":0,"coinMax":200000,"maxFan":8,"diZhuCoin":2000,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":3,"daidaTime":10000,"coinMin":100000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":7,"startMinNum":2},{"dingqueTime":20000,"roomType":0,"coinMax":-1,"maxFan":8,"diZhuCoin":10000,"isNeedRobot":1,"coinType":1001,"firstFapaiTime":10000,"clientMax":500,"clientNum":0,"dapaiTime":20000,"coinSetp":4,"daidaTime":10000,"coinMin":200000,"jiesuanTime":10000,"zhunbeiTime":10000,"roomId":8,"startMinNum":2}]`
	var pasre []map[string]int32
	json.Unmarshal([]byte(t1), &pasre)

	gs.roomControl = make([int]*RoomControl)
	for _, info := range pasre {
		roomType, ok := info["roomType"]
		if !ok {
			continue
		}
		switch roomType {
		case KRoomTypeClassic:
			roomInfo := &RoomInfo{}
			roomInfo.RoomType = roomType
			roomInfo.RoomId = info["room"]
			roomInfo.CoinMin = int64(info["coinMin"])
			roomInfo.CoinMax = int64(info["coinMax"])
			roomInfo.DiZhu = int64(info["diZhuCoin"])
			roomInfo.MaxFan = info["maxFan"]
			roomInfo.ClientNum = info["clientNum"]
			roomInfo.ClientMax = info["clientMax"]
			roomInfo.StartMinMax = int16(info["startMinNum"])
			roomInfo.coinType = info["coinType"]
			roomInfo.zbTime = info["zhunbeiTime"]
			roomInfo.DuanPaiTime = info["firstFapaiTime"]
			roomInfo.DingQueTime = info["dingqueTime"]
			roomInfo.ChuPaiTime = info["dapaiTime"]
			roomInfo.OpChooseTime = info["daidaTime"]
			roomInfo.JieSuanTime = info["jiesuanTime"]
			roomInfo.coinStep = info["coinSetp"]
			roomInfo.isNeedRobot = (info["isNeedRobot"] == 1)
			roomInfo.Port = 4001
			roomInfo.IP = "127.0.0.1"
			gs.addRoom(roomInfo)
		case KRoomTypeDC:
			roomInfo := &RoomInfo{}
			roomInfo.RoomType = roomType
			roomInfo.RoomId = info["room"]
			roomInfo.CoinMin = int64(info["coinMin"])
			roomInfo.CoinMax = int64(info["coinMax"])
			roomInfo.DiZhu = int64(info["diZhuCoin"])
			roomInfo.MaxFan = info["maxFan"]
			roomInfo.ClientNum = info["clientNum"]
			roomInfo.ClientMax = info["clientMax"]
			roomInfo.StartMinMax = int16(info["startMinNum"])
			roomInfo.coinType = info["coinType"]
			roomInfo.zbTime = info["zhunbeiTime"]
			roomInfo.DuanPaiTime = info["firstFapaiTime"]
			roomInfo.DingQueTime = info["dingqueTime"]
			roomInfo.ChuPaiTime = info["dapaiTime"]
			roomInfo.OpChooseTime = info["daidaTime"]
			roomInfo.JieSuanTime = info["jiesuanTime"]
			roomInfo.coinStep = info["coinSetp"]
			roomInfo.isNeedRobot = (info["isNeedRobot"] == 1)

			roomInfo.baoMingFei = int64(info["baoMingCoin"])
			roomInfo.SocreGainTime = info["duiHuaTime"]
			roomInfo.jifenBeilv = info["jifenBeiLv"]
			roomInfo.Port = 4001
			roomInfo.IP = "127.0.0.1"
			gs.addRoom(roomInfo)
		}
	}
}

func (gs *GameService) addRoom(info *RoomInfo) {
	rc := NewRC(info)
	_, ok := gs.rooms[rc.roomInfo.RoomId]
	if ok {
		log.Error("GameService:addRoom:%d is exist.", rc.roomInfo.RoomId)
	}
	gs.rooms[rc.roomInfo.RoomId] = rc
}

func (gs *GameService) Run() {
	core.RegisterService(gs)
	core.Name(gs.Id(), "gsserver")
	gs.hsId, _ = core.GetIdByName("hallService")
	gs.ticker = time.NewTicker(time.Millisecond * 10)
	go func() {
	EXIT:
		for {
			_, ok := <-gs.ticker.C
			if !ok {
				break
			}
			//loop for msg
		MSGLOOP:
			for {
				select {
				case msg, ok := <-gs.In():
					if !ok {
						break EXIT
					}
					gs.DispatchM(msg)
				default:
					break MSGLOOP
				}
			}
			//loop for desk
		}
	}()

	s := tcp.New("", "40001", gs.Id())
	s.Listen()
}
