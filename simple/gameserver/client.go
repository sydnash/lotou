package gameserver

import (
	"encoding/json"
	"github.com/sydnash/lotou/simple/global"
	"strconv"
)

type GameClient struct {
	playerInfo *global.PropertySet
	acId       int32
	session    int32
	roomId     int32
	deskId     int32
	deskPos    int32
	gs         *GameService
	agentId    uint
	dc         *DeskConrol
}

func (gc *GameClient) saveInfoToString() ([]byte, error) {
	sendMap := make(map[string]string)
	for k, v := range gc.playerInfo.Property {
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
