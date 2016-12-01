package gameserver

const (
	KRoomTypeClassic = iota
	KRoomTypeDC
)

type RoomInfo struct {
	RoomType      int32
	RoomId        int32
	CoinMin       int64
	CoinMax       int64
	DiZhu         int64
	MaxFan        int32
	ClientNum     int16
	ClientMax     int16
	StartMinMax   int16
	Port          int16
	IP            string
	zbTime        int32
	DuanPaiTime   int32
	DingQueTime   int32
	ChuPaiTime    int32
	OpChooseTime  int32
	JieSuanTime   int32
	coinType      int32
	coinStep      int32
	isNeedRobot   bool
	SocreGainTime int32
	baoMingFei    int64
	jifenBeilv    int32
}

type RoomControl struct {
	roomInfo *RoomInfo
	desks    []*DeskConrol
}

func (rc *RoomControl) isCanEnter() bool {
	if rc.roomInfo.ClientNum+1 <= rc.roomInfo.ClientMax {
		return true
	}
	return false
}

func (rc *RoomControl) enter(client *GameClient) {
	dc := &DeskConrol{}
	dc.rc = rc
	client.roomId = rc.roomInfo.RoomId
	dc.enter(client)
	rc.desks = append(rc.desks, dc)
}

func NewRC(info *RoomInfo) *RoomControl {
	ret := &RoomControl{}
	ret.roomInfo = info
	ret.desks = make([]*DeskConrol, 0, 1000)
	return ret
}
