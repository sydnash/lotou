package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

const (
	ACTION, AC, TYPE, PW, QU_DAO_TYPE, DEVICE_ID, ACCOUNT_TYPE, CRC = "action", "ac", "type", "pw", "qudaoType", "deviceId", "accountType", "crc"
)

func intToByteSlice(v uint32) []byte {
	a := make([]byte, 4)
	a[3] = byte((v >> 24) & 0xFF)
	a[2] = byte((v >> 16) & 0XFF)
	a[1] = byte((v >> 8) & 0XFF)
	a[0] = byte(v & 0XFF)
	return a
}
func byteSliceToInt(s []byte) (v uint32) {
	v = uint32(s[3])<<24 | uint32(s[2])<<16 | uint32(s[1])<<8 | uint32(s[0])
	return v
}

func int16ToByteSlice(v uint16) []byte {
	a := make([]byte, 2)
	a[1] = byte((v >> 8) & 0xFF)
	a[0] = byte(v & 0XFF)
	return a
}
func byteSliceToInt16(s []byte) (v uint16) {
	v = uint16(s[1])<<8 | uint16(s[0])
	return v
}

func pack(js []byte, msgType uint32) []byte {
	a := make([]byte, 0, 1000)
	a = append(a, intToByteSlice(uint32(18+len(js)))...)
	a = append(a, intToByteSlice(0)...)
	a = append(a, int16ToByteSlice(0)...)
	a = append(a, intToByteSlice(msgType)...)
	a = append(a, intToByteSlice(0)...)
	a = append(a, js...)
	return a
}
func unpack(in []byte, rlen int, t interface{}) {
	tmt := in[18:rlen]
	err := json.Unmarshal(tmt, t)
	fmt.Println("recv json:", string(tmt))
	if err != nil {
		fmt.Println("error type:", reflect.TypeOf(err))
		fmt.Println("parse err:", err)
		return
	}
	fmt.Println(t)
}

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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	action := r.FormValue(ACTION)
	ac := r.FormValue(AC)
	ltype := r.FormValue(TYPE)
	pw := r.FormValue(PW)
	qudaoType := r.FormValue(QU_DAO_TYPE)
	deviceId := r.FormValue(DEVICE_ID)
	accountType := r.FormValue(ACCOUNT_TYPE)
	crc := r.FormValue(CRC)
	fmt.Println(action, ac, ltype, pw, qudaoType, deviceId, accountType, crc)

	t1, _ := strconv.Atoi(accountType)
	t2, _ := strconv.Atoi(qudaoType)
	loginSend := LoginSend{AcName: ac, AccountType: t1, QudaoType: t2, MacAddress: deviceId, LoginType: 2}

	b, err := json.Marshal(loginSend)
	if err != nil {
		fmt.Println("json error:", err)
		w.Write([]byte(`{"status":1000001}`))
		return
	}
	fmt.Println(string(b))

	raddress, _ := net.ResolveTCPAddr("tcp", "192.168.23.7:55000")
	tcpCon, err := net.DialTCP("tcp", nil, raddress)
	if err != nil {
		fmt.Println("connect to game server faield:", err)
		w.Write([]byte(`{"status":1000001}`))
		return
	}
	defer tcpCon.Close()
	sendPacket := pack(b, 2013)
	tcpCon.Write(sendPacket)

	buf := make([]byte, 1024)
	recvLen, err := tcpCon.Read(buf)
	if err != nil {
		fmt.Println("recv from server faield:", err)
		w.Write([]byte(`{"status":1000001}`))
		return
	}
	a := LoginRecv{}
	unpack(buf, recvLen, &a)
	if a.IRet != 1 {
		w.Write([]byte(`{"status":1000001}`))
		return
	}
	a.Status = 1
	sendClient, err := json.Marshal(a)
	w.Write(sendClient)
}
func main() {
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8080", nil)
}
