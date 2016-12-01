package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/simple/gameserver"
	"github.com/sydnash/lotou/topology"
	"math/rand"
	"time"
)

func main() {
	//init slave node
	rand.Seed(time.Now().UnixNano())
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	topology.StartSlave("127.0.0.1", "4000")
	core.RegisterNode()
	platid, _ := core.GetIdByName("platservice")
	dbid, ok := core.GetIdByName("dbserver")
	if !ok {
		panic("could not find dbserver.")
	}

	log.Info("platid: %x", platid)
	log.Info("dbid: %x", dbid)
	hs := NewHS(platid, dbid)
	hs.Run()

	gs := gameserver.NewGS()
	gs.Run()

	ch := make(chan int)
	<-ch
}
