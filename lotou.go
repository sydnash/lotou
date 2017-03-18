package lotou

import (
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

type ModuleParam struct {
	N string
	M core.Module
}

type CloseFunc func()

func Start(f CloseFunc, data ...*ModuleParam) {
	log.Init(conf.LogFilePath, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)

	core.InitNode(conf.CoreIsStandalone, conf.CoreIsMaster)

	if !conf.CoreIsStandalone {
		if conf.CoreIsMaster {
			topology.StartMaster(conf.MasterListenIp, conf.MultiNodePort)
		} else {
			topology.StartSlave(conf.SlaveConnectIp, conf.MultiNodePort)
		}
		core.RegisterNode()
	}

	for _, m := range data {
		core.StartService(m.N, m.M)
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

//RawStart start lotou, with no block.
func RawStart(data ...*ModuleParam) {
	log.Init(conf.LogFilePath, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)

	core.InitNode(conf.CoreIsStandalone, conf.CoreIsMaster)

	if !conf.CoreIsStandalone {
		if conf.CoreIsMaster {
			topology.StartMaster(conf.MasterListenIp, conf.MultiNodePort)
		} else {
			topology.StartSlave(conf.SlaveConnectIp, conf.MultiNodePort)
		}
		core.RegisterNode()
	}

	for _, m := range data {
		core.StartService(m.N, m.M)
	}
}
