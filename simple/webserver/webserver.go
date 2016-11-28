package main

import (
	"encoding/json"
	"fmt"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/simple/json_type"
	"github.com/sydnash/lotou/topology"
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

	s := &Service{core.NewBase()}
	core.RegisterService(s)
	ret := core.Call(platid, s, "Login", ac, t1, t2, deviceId, 2)
	fmt.Println(ret)
	//iret int, msg, ip, port, pwd_sec string, session uint64, acid int
	a := json_type.LoginRecv{}
	a.IRet = ret[0].(int)
	if a.IRet != 1 {
		w.Write([]byte(`{"status":1000001}`))
		return
	} else {
		a.Msg = ret[1].(string)
		a.Ip = ret[2].(string)
		var err error
		a.Port, err = strconv.Atoi(ret[3].(string))
		_ = err
		a.Pwd_second = ""
		a.Session = strconv.FormatUint(ret[5].(uint64), 10)
		a.AccountId = ret[6].(int)
		a.AccountType = t1
		a.Status = 1
		sendClient, err := json.Marshal(a)
		w.Write(sendClient)
		return
	}

	loginSend := json_type.LoginSend{AcName: ac, AccountType: t1, QudaoType: t2, MacAddress: deviceId, LoginType: 2}

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
	a = json_type.LoginRecv{}
	unpack(buf, recvLen, &a)
	if a.IRet != 1 {
		w.Write([]byte(`{"status":1000001}`))
		return
	}
	a.Status = 1
	sendClient, err := json.Marshal(a)
	w.Write(sendClient)
}

type Service struct {
	*core.Base
}

var platid uint

func main() {
	//init slave node
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	log.Debug("start slave")
	topology.StartSlave("127.0.0.1", "4000")
	core.RegisterNode()
	platid, _ = core.GetIdByName("platservice")

	//start http
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8080", nil)
}
