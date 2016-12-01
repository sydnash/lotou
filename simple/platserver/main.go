package main

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/simple/dbserver"
	"github.com/sydnash/lotou/topology"
)

func main() {
	//init master ndoe
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	core.SetAsMaster()
	topology.StartMaster("127.0.0.1", "4000")

	//init service
	OpenDB()
	ps := NewPS()
	ps.Run()

	db := dbserver.NewDB()
	db.Run()
	ch := make(chan int)
	<-ch
}
