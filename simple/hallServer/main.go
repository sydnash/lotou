package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/simple/json_type"
	"github.com/sydnash/lotou/topology"
)

func main() {
	//init slave node
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	topology.StartSlave("127.0.0.1", "4000")
	core.RegisterNode()
	platid, _ = core.GetIdByName("platservice")

	hs := NewHS()
	hs.Run()

	ch = make(chan int)
	<-ch
}
