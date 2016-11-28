package json_type

import (
	"fmt"
)

type LoginSend struct {
	AcName      string `json:"acName"`
	AcPWD       string `json:"acPwd"`
	AccountType int    `json:"accountType"`
	QudaoType   int    `json:"qudaoType"`
	MacAddress  string `json:"mac"`
	LoginType   int    `json:"loginType"`
	TitleURL    string `json:"titleUrl"`
	Coin        string `json:"coin"`
	Nickname    string `json:"nickname"`
	Status      int    `json:"status"`
}
type LoginRecv struct {
	IRet        int    `json:"iRet"`
	Msg         string `json:"msg"`
	AcName      string `json:"acName"`
	Session     string `json:"session"`
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	AccountType int    `json:"accountType"`
	AccountId   int    `json:"accountId"`
	Pwd_second  string `json:"pwd_second"`
	Status      int    `json:"status"`
}

func (self *LoginRecv) String() string {
	s := fmt.Sprintf("iRet: %d, msg: %s, acName:%s, session:%s, ip:%s, port:%d, accountType:%d, accountId:%d", self.IRet, self.Msg, self.AcName, self.Session, self.Ip, self.Port, self.AccountType, self.AccountId)
	return s
}
