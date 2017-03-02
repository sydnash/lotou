package topology_test

import (
	"fmt"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
)

var (
	remoteId uint = 100
)

func (g *Game) OnMainLoop(dt int) {
	if remoteId != 0 {
		log.Info("send")
		g.RawSend(remoteId, core.MSG_TYPE_NORMAL, 1, 2, 3, 4)

		t := func(timeout bool, data ...interface{}) {
			fmt.Println("request respond ", timeout, data)
		}
		g.Request(remoteId, 10, t, func() {
			log.Info("timeout")
		}, "hello")

		fmt.Println(g.Call(remoteId, "hello"))
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
	core.StartService("game2", game)
	log.Info("game2's id: %v", game.Id)

	var err error
	remoteId, err = core.NameToId("game1")
	log.Info("NameToId: %v, %v %v", "game1", remoteId, err)

	ch := make(chan int)
	<-ch

	/*
		response := func(data ...interface{}) {
			log.Info("respond: %v", data)
		}
		_ = response
		go func() {
			for msg := range game.In() {
				game.DispatchM(msg)
			}
		}()
		core.Send(t4, game.Id(), "测试中", "同时测试中文")
		if t5 {
			for {
				core.Request(t4, game, response, "response1", "response2")
				ret, err := core.Call(t4, game, "call1", "call2")
				log.Info("ret: %v", ret)
			}
		}
		for {
		}
	*/
}
