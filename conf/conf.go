package conf

import (
	"github.com/sydnash/lotou/log"
)

var (
	LogFilePath      string = "loger"
	LogFileLevel     int    = log.LEVEL_MAX
	LogShellLevel    int    = log.DEBUG_LEVEL
	LogMaxLine       int    = 10000
	LogBufferSize    int    = 1000
	CoreIsStandalone bool   = true
	CoreIsMaster     bool   = true
	MasterListenIp   string = "127.0.0.1"
	SlaveConnectIp   string = "127.0.0.1"
	MultiNodePort    string = "4000"
	CallTimeOut      int    = 10000
)
