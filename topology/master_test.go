package topology_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
	"time"
)

type Game struct {
	*core.Base
}

func (g *Game) CloseMSG(dest, src uint) {
	g.Base.Close()
	log.Info("close: %v, %v", src, dest)
}
func (g *Game) NormalMSG(dest, src uint, msgType string, data ...interface{}) {
	log.Info("%x, %x, %v", src, dest, data)
	//core.Send(src, dest, a)
}

func (g *Game) RequestMSG(dest, src uint, rid int, data ...interface{}) {
	log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
	core.Respond(src, dest, rid, data...)
}
func (g *Game) CallMSG(dest, src uint, data ...interface{}) {
	log.Info("call: %x, %x, %v", src, dest, data)
	core.Ret(src, dest, data...)
}
func TestMaster(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	core.SetAsMaster()
	topology.StartMaster("127.0.0.1", "4000")

	game := &Game{core.NewBase()}
	core.RegisterService(game)
	core.Name(game.Id(), "game")
	game.SetDispatcher(game)
	go func() {
		for msg := range game.In() {
			game.DispatchM(msg)
		}
	}()

	core.Name(100, "我是服务2")
	core.Name(100, "service3")
	log.Info("test")
	for {
		time.Sleep(time.Minute * 10)
	}
}
