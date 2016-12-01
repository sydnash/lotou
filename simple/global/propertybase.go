package global

import (
	"errors"
)

const PROPERTY_JIANGEVALUE = 1000
const PROPERTYBASE_JIANGGEVALUE = 10000

const (
	KValueTypeInt8 = iota + 1
	KValueTypeInt16
	KValueTypeInt32
	KValueTypeInt64
	KValueTypeString
	KValueTypeMax
)
const (
	KPropertyBaseNone = iota
	KPropertyBaseMax
)
const (
	KPropertyType32Start = KPropertyBaseNone*PROPERTYBASE_JIANGGEVALUE + iota
	KPropertyTypeAcId
	KPropertyTypeYuanBao
	KPropertyTypePlayerType
	KPropertyTypeSex
	KPropertyTypeVipLv
	KPropertyTypeFlags1
	KPropertyTypeFlags2
	KPropertyTypeLv
	KPropertyTypeExp
	KPropertyTypeTiliDC
	KPropertyTypeScoreDC
	KPropertyTypeRoomIdDC
	KPropertyTypeMaxScoreDC
	KPropertyTypeGameNuJbs
	KPropertyTypeDayBuZuNum
	KPropertyTypeAccountType
	KPropertyTypeQuDaoType
	KPropertyTypeRobotType
	KPropertyTypeJBSRewardBeiShu
	KPropertyTypeRandomLV_dc
	KPropertyType32End
)
const (
	KPropertyType64Start = KPropertyBaseNone*PROPERTYBASE_JIANGGEVALUE + PROPERTY_JIANGEVALUE + iota
	KPropertyTypeCoin
	KPropertyTypeQueBiDC
	KPropertyTypeJiangQuan
	KPropertyTypeLastLoginTime
	KPropertyTypeLastOutTime
	KPropertyType64End
)
const (
	KPropertyTypeStringStart = KPropertyBaseNone*PROPERTYBASE_JIANGGEVALUE + PROPERTY_JIANGEVALUE*2 + iota
	KPropertyTypeNickname
	KPropertyTypeTitleUrl
	KPropertyTypeVideoUrl
	KPropertyTypeSign
	KPropertyTypeAcName
	KPropertyTypeLastIp
	KPropertyTypeFanPai
	KPropertyTypeSim
	KPropertyTypeMac
	KPropertyTypeStringEnd
)
const (
	KFlagsIsFirstChuangJian = 1 << 0
	KFlagsIsHaveMjDC        = 1 << 1
)
const (
	KPlayerTypeCommon = iota
	KPlayerTypeRobot
	KPlayerTypeMax
)

type PropertyBase struct {
	PropertyName   string
	IsSendToClient bool
	IsSqlPlayerDup bool
	IsHaveDbLog    bool
	ValueType      int
	Value          interface{}
}

var (
	keyToProperty  map[string]int
	propertyToBase map[int]*PropertyBase
)

var KeyIsNotExist = errors.New("key is not exist")

func TypeToKey(ptype int) (*PropertyBase, error) {
	base, ok := propertyToBase[ptype]
	if !ok {
		return nil, KeyIsNotExist
	}
	return base, nil
}

func init() {
	keyToProperty = make(map[string]int)
	propertyToBase = make(map[int]*PropertyBase)
	type tmp struct {
		PropertyName   string
		IsSendToClient bool
		IsSqlPlayerDup bool
		IsHaveDbLog    bool
		def            interface{}
	}
	addProperty := func(ptype int, t *tmp, vtype int) {
		keyToProperty[t.PropertyName] = ptype
		propertyToBase[ptype] = &PropertyBase{t.PropertyName, t.IsSendToClient, t.IsSqlPlayerDup, t.IsHaveDbLog, vtype, t.def}
	}
	addPropertyArray := func(start int, array []tmp, vtype int) {
		for i, v := range array {
			ptype := i + start + 1
			addProperty(ptype, &v, vtype)
		}
	}
	array1 := []tmp{
		{"acId", false, false, false, int32(0)},           // 玩家账户Id
		{"yuanBao", true, true, true, int32(0)},           // 玩家元宝
		{"playerType", false, true, true, int32(0)},       // 玩家类型
		{"sex", true, true, true, int32(0)},               // 玩家性别
		{"vipLv", true, true, true, int32(0)},             // vip等级
		{"flags1", true, false, true, int32(0)},           // 标志位1
		{"flags2", true, false, true, int32(0)},           // 标志位2
		{"lv", true, true, true, int32(0)},                // 等级
		{"exp", true, true, true, int32(0)},               // 经验
		{"tiLi_dc", true, false, false, int32(0)},         // 日赛体力
		{"jifen_dc", true, false, true, int32(0)},         // 日赛积分
		{"roomId_dc", true, false, false, int32(0)},       // 日赛场 房间号
		{"jifen_Max_dc", true, false, true, int32(0)},     // 日赛场积分最大值
		{"game_nu_jbs", true, false, false, int32(0)},     // 当前轮数的局数
		{"DayBuzuNum", true, false, false, int32(0)},      // 当前每日补助的数量
		{"AccountType", true, false, false, int32(0)},     // 账户类型QudaoType
		{"QudaoType", true, false, false, int32(0)},       // 渠道类型
		{"RobotType", true, false, false, int32(0)},       // 账户类型
		{"JBSRewardBeishu", true, false, false, int32(0)}, // 账户类型
		{"randomLv_dc", true, true, true, int32(0)},       // 账户类型
	}
	addPropertyArray(KPropertyType32Start, array1, KValueTypeInt32)
	array2 := []tmp{
		{"coin", true, true, true, int64(0)},           // 金币
		{"quebi_dc", true, true, true, int64(0)},       // 日赛雀币
		{"jiangQuan", true, true, true, int64(0)},      // 奖券
		{"lastLoginTime", false, true, true, int64(0)}, // 上次登录的时间
		{"lastOutTime", false, true, true, int64(0)},   // 上次登出的时间
	}
	addPropertyArray(KPropertyType64Start, array2, KValueTypeInt64)
	array3 := []tmp{
		{"name", true, false, true, ""},      // 玩家名字
		{"titleUrl", true, false, false, ""}, // 头像url
		{"videoUrl", true, false, false, ""}, // 视频url
		{"sign", true, false, false, ""},     // 签名
		{"acName", true, false, false, ""},   // 账户名字
		{"lastIp", true, true, false, ""},    // 最后登录的Ip
		{"fanPai", false, false, false, ""},  // 翻牌相关的
		{"Sim", false, false, false, ""},     // 翻牌相关的
		{"Mac", false, false, false, ""},     // 翻牌相关的
	}
	addPropertyArray(KPropertyTypeStringStart, array3, KValueTypeString)
}
