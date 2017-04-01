package conf

import (
	"encoding/json"
	"fmt"
	"github.com/sydnash/lotou/log"
	"io/ioutil"
	"os"
	"path"
	"reflect"
)

var (
	LogFilePath      string = "loger"         //log's file path
	LogFileLevel     int    = log.LEVEL_MAX   //log's file level
	LogShellLevel    int    = log.DEBUG_LEVEL //log's shell level
	LogMaxLine       int    = 10000           //log's max line per file
	LogBufferSize    int    = 1000            //
	LogHasColor      bool   = true
	CoreIsStandalone bool   = false       //set system is a standalone or multinode
	CoreIsMaster     bool   = true        //set node is master
	MasterListenIp   string = "127.0.0.1" //master listen ip
	SlaveConnectIp   string = "127.0.0.1" //master ip
	MultiNodePort    string = "4000"      //master listen port
	CallTimeOut      int    = 10000       //global timeout for Call fucntion
)

func PrintCurrentConfSetToStd() {
	fmt.Printf("LogFilePath = %v\n", LogFilePath)
	fmt.Printf("LogFileLevel = %v\n", LogFileLevel)
	fmt.Printf("LogShellLevel = %v\n", LogShellLevel)
	fmt.Printf("LogMaxLine = %v\n", LogMaxLine)
	fmt.Printf("LogBufferSize = %v\n", LogBufferSize)
	fmt.Printf("LogHasColor = %v\n", LogHasColor)
	fmt.Printf("CoreIsStandalone = %v\n", CoreIsStandalone)
	fmt.Printf("CoreIsMaster = %v\n", CoreIsMaster)
	fmt.Printf("MasterListenIp = %v\n", MasterListenIp)
	fmt.Printf("SlaveConnectIp = %v\n", SlaveConnectIp)
	fmt.Printf("MultiNodePort = %v\n", MultiNodePort)
	fmt.Printf("CallTimeOut = %v\n", CallTimeOut)
}

//~
func assignTo(r map[string]interface{}, target interface{}, name string) {
	if t := reflect.TypeOf(target); t.Kind() != reflect.Ptr {
		fmt.Println("Kind is", t.Kind())
		panic("Not a proper type")
	}
	if s, ok := r[name]; ok {
		t0 := reflect.TypeOf(s)
		t1 := reflect.TypeOf(target).Elem()
		// var doOK bool = true
		tSet := reflect.ValueOf(target).Elem()
		if t0.AssignableTo(t1) {
			tSet.Set(reflect.ValueOf(s))
		} else if t1.Kind() == reflect.Int {
			if t0.Kind() == reflect.Float64 {
				f := reflect.ValueOf(s).Float()
				tSet.Set(reflect.ValueOf(int(f)))
			} else {
				// doOK = false
			}
		} else {
			// doOK = false
		}
		// if doOK {
		// 	fmt.Printf("Set of <%v> to %v\n", name, s)
		// } else {
		// 	fmt.Printf("Cannot assign %v to %v\n", t0.Name(), t1.Name())
		// }
	}
}

const (
	privateConfiguraPath = ".private/svrconf.json"
)

//~ If you wish to alter the configuration filepath
//~ Overwrite the configures by configuration file.
func init() {
	goPath := os.ExpandEnv("$GOPATH")
	if len(goPath) <= 0 {
		return
	}
	fname := path.Join(goPath, privateConfiguraPath)
	fin, err := os.Open(fname)
	if err != nil {
		//The configure file may not exist
		return
	}
	defer fin.Close()
	chunk, r_err := ioutil.ReadAll(fin)
	if r_err != nil {
		return
	}
	var intfs interface{}
	jErr := json.Unmarshal(chunk, &intfs)
	if jErr != nil {
		return
	}

	if mIntfs, ok := intfs.(map[string]interface{}); ok {
		assignTo(mIntfs, &LogFilePath, `LogFilePath`)
		assignTo(mIntfs, &LogHasColor, `LogHasCoor`)
		assignTo(mIntfs, &CoreIsStandalone, `IsStandalone`)
		assignTo(mIntfs, &CoreIsMaster, `IsMaster`)
		assignTo(mIntfs, &MasterListenIp, `MasterIP`)
		assignTo(mIntfs, &SlaveConnectIp, `SlaveIP`)
		assignTo(mIntfs, &MultiNodePort, `MasterPort`)
		assignTo(mIntfs, &CallTimeOut, `CallTimeout`)
	}
}

func SetMasterMode() {
	CoreIsMaster = true
	CoreIsStandalone = false
}

func SetStandaloneMode() {
	CoreIsMaster = false
	CoreIsStandalone = true
}

func SetSlaveMode() {
	CoreIsMaster = false
	CoreIsStandalone = false
}
