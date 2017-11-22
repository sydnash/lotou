package lotou

import (
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"math/rand"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type CloseFunc func()

//Start start lotou with given modules which in is in data.
//initialize log by config param
//if lotou's network is standalone, then only start master service.
//if lotou's network is not standalone,
// 		if node is master, then start master service
//		if node is slave, then start slave service, and register node to get a nodeid from master which will block until it success
//capture system's SIGKILL SIGTERM signal
//and wait until all service are closed.
//f will be called when SIGKILL or SIGTERM is received.
func Start(f CloseFunc, data ...*core.ModuleParam) {
	StartWithName("", f, nil, data...)
}

func StartWithName(nodeName string, f CloseFunc, customLogger log.Logger, data ...*core.ModuleParam) {
	rand.Seed(time.Now().UnixNano())

	if customLogger == nil {
		logFilePath := nodeName + "_" + conf.LogFilePath
		logger := log.Init(logFilePath, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)
		logger.SetColored(conf.LogHasColor)
	} else {
		log.SetLogger(customLogger)
	}
	core.InitNode(conf.CoreIsStandalone, conf.CoreIsMaster)

	log.Info("starting node: {%v}", nodeName)
	if !conf.CoreIsStandalone {
		if conf.CoreIsMaster {
			topology.StartMaster(conf.MasterListenIp, conf.MultiNodePort)
		} else {
			topology.StartSlave(conf.SlaveConnectIp, conf.MultiNodePort)
		}
		core.RegisterNode(nodeName)
	} else {
		topology.StartMaster(conf.MasterListenIp, conf.MultiNodePort)
	}

	for _, m := range data {
		core.StartService(m)
	}

	/*err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Error("ListenAndServe: %v", err)
	}*/

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	if f == nil {
		f = core.SendCloseToAll
	}
	go func() {
		for {
			sig := <-c
			log.Info("lotou closing down (signal: %v)", sig)
			f()
		}
	}()

	core.Wait()
}

//RawStart start lotou, with no wait
func RawStart(nodeName string, data ...*core.ModuleParam) {
	log.Init(conf.LogFilePath, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)

	core.InitNode(conf.CoreIsStandalone, conf.CoreIsMaster)

	log.Info("start node: {%v}", nodeName)
	if !conf.CoreIsStandalone {
		if conf.CoreIsMaster {
			topology.StartMaster(conf.MasterListenIp, conf.MultiNodePort)
		} else {
			topology.StartSlave(conf.SlaveConnectIp, conf.MultiNodePort)
		}
		core.RegisterNode(nodeName)
	}

	for _, m := range data {
		core.StartService(m)
	}
}
