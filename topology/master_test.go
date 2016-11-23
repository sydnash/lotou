package topology_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
)

type Game struct {
	*core.Base
}

func (g *Game) Close(dest, src uint, msgType string) {
	g.Base.Close()
	log.Info("close: %v, %v", src, dest)
}

func TestMaster(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	core.SetAsMaster()
	topology.StartMaster("127.0.0.1", "4000")

	game := &Game{core.NewBase()}
	core.RegisterService(game)
	core.Name(game.Id(), "game")
	c2 := func(dest, src uint, msgType string, data ...interface{}) {
		log.Info("%x, %x, %v", src, dest, data)
	}
	game.SetSelf(game)
	game.RegisterBaseCB(core.MSG_TYPE_CLOSE, (*Game).Close, true)
	game.RegisterBaseCB(core.MSG_TYPE_NORMAL, c2, false)
	go func() {
		for msg := range game.In() {
			game.DispatchM(msg)
		}
	}()

	core.Name(100, "我是服务2")
	core.Name(100, "service3")
	log.Info("test")
	for {
	}
}
