package lotou

import (
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"os"
	"os/signal"
)

type ModuleParam struct {
	N string
	M core.Module
}

func Start(data ...*ModuleParam) {
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Info("lotou closing down (signal: %v)", sig)
}
