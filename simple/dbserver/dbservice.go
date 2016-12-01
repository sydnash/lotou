package dbserver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"reflect"
)

type DBService struct {
	*core.Base
	db *sql.DB
}

func (db *DBService) CloseMSG(dest, src uint) {
	log.Info("dbservice Close msg")
	db.db.Close()
	db.Base.Close()
}
func (db *DBService) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	log.Info("%x, %x, %v", src, dest, data)
	if msgType == "go" {
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
	} else if msgType == "socket" {
	}
}
func (db *DBService) UpdatePlayerData(acId int32, nickname string, jsonStr []byte) bool {
	_, err := db.db.Exec("replace into player(accountId, name, data) values (?,?,?) ", acId, nickname, string(jsonStr))
	if err != nil {
		log.Error("DBService:UpdatePlayerData faield %s", err)
		return false
	}
	return true
}

func (db *DBService) CallMSG(dest, src uint, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)
	core.Ret(src, dest, data...)
}
func (db *DBService) RequestMSG(dest, src uint, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	cmd := data[0].(string)
	psv := reflect.ValueOf(db)
	fv := psv.MethodByName(cmd)
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
func (db *DBService) PlayerLogin(acid int32) (acId int32, nicknamestr, datastr, namestr string, playType, accountType, qdType int32, macstr string) {
	row := db.db.QueryRow("call SelectLoginOne(?)", acid)
	var name sql.NullString
	var nickname sql.NullString
	var data sql.NullString
	var mac sql.NullString
	err := row.Scan(&acId, &nickname, &data, &playType, &name, &accountType, &qdType, &mac)
	if err != nil {
		log.Error("playerlogin is failed, scan db:%s", err)
	}

	if name.Valid {
		namestr = name.String
	}
	if nickname.Valid {
		nicknamestr = nickname.String
	}
	if data.Valid {
		datastr = data.String
	}
	if mac.Valid {
		macstr = mac.String
	}
	return
}

func NewDB() *DBService {
	db := &DBService{Base: core.NewBaseLen(1024 * 1024)}
	var err error
	db.db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/niuniu")
	if err != nil {
		panic("connect database failed.")
		return nil
	}
	db.SetDispatcher(db)
	return db
}

func (db *DBService) Run() {
	core.RegisterService(db)
	core.Name(db.Id(), "dbserver")
	go func() {
		for msg := range db.In() {
			db.DispatchM(msg)
		}
	}()
}
