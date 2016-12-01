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
}

func NewRC(info *RoomInfo) *RoomControl {
	ret := &RoomControl{}
	ret.roomInfo = info
	return ret
}
