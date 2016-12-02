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
	ZbTime        int32
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

func (rc *RoomControl) getDesk() *DeskConrol {
	for _, dc := range rc.desks {
		if dc.isCaneEnter() {
			return dc
		}
	}
	dc := NewDC(rc)
	rc.desks = append(rc.desks, dc)
	return dc
}

func (rc *RoomControl) enter(client *GameClient) {
	rc.roomInfo.ClientNum++
	dc := rc.getDesk()
	client.roomId = rc.roomInfo.RoomId
	dc.enter(client)
}
func (rc *RoomControl) decreaseClientNum() {
	rc.roomInfo.ClientNum--
	if rc.roomInfo.ClientNum < 0 {
		rc.roomInfo.ClientNum = 0
	}
}

func (rc *RoomControl) exit(client *GameClient) {
	dc := client.dc
	dc.exit(client)
}

func NewRC(info *RoomInfo) *RoomControl {
	ret := &RoomControl{}
	ret.roomInfo = info
	ret.desks = make([]*DeskConrol, 0, 1000)
	return ret
}
