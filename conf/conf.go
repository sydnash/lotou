package conf

import (
	"github.com/sydnash/lotou/log"
)

var (
	LogFilePath      string = "loger"         //log's file path
	LogFileLevel     int    = log.LEVEL_MAX   //log's file level
	LogShellLevel    int    = log.DEBUG_LEVEL //log's shell level
	LogMaxLine       int    = 10000           //log's max line per file
	LogBufferSize    int    = 1000            //
	LogHasColor      bool   = true
	CoreIsStandalone bool   = true        //set system is a standalone or multinode
	CoreIsMaster     bool   = true        //set node is master
	MasterListenIp   string = "127.0.0.1" //master listen ip
	SlaveConnectIp   string = "127.0.0.1" //master ip
	MultiNodePort    string = "4000"      //master listen port
	CallTimeOut      int    = 10000       //global timeout for Call fucntion
)
