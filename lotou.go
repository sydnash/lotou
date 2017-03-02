package lotou

import (
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
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
}
