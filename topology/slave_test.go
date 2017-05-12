package topology_test

import (
	"fmt"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
)

var (
	remoteId core.ServiceID = 100
)

func (g *Game) OnMainLoop(dt int) {
	if remoteId != 0 {
		log.Info("send")
		g.RawSend(remoteId, core.MSG_TYPE_NORMAL, "haha", 1, 2, 3, 4)

		t := func(timeout bool, data ...interface{}) {
			fmt.Println("request respond ", timeout, data)
		}
		g.Request(remoteId, core.MSG_ENC_TYPE_GO, 10, t, "hello")

		fmt.Println(g.Call(remoteId, core.MSG_ENC_TYPE_GO, "hello"))
	}
}

func TestSlavea(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)

	remoteId = 0
	core.InitNode(false, false)
	topology.StartSlave("127.0.0.1", "4000")

	log.Info("start register node")
	core.RegisterNode()
	log.Info("start create service")
	game := &Game{core.NewSkeleton(1000)}
	core.StartService(&core.ModuleParam{
		N: "game2",
		M: game,
		L: 0,
	})
	log.Info("game2's id: %v", game.Id)

	var err error
	remoteId, err = core.NameToId("game1")
	log.Info("NameToId: %v, %v %v", "game1", remoteId, err)

	ch := make(chan int)
	<-ch
}
