package topology_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
)

func TestSlaveb(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	log.Debug("start slave")
	topology.StartSlave("127.0.0.1", "4000")
	core.RegisterNode()
	core.Name(1, "我是服务1")

	game := &Game{core.NewBase()}
	core.RegisterService(game)
	core.Name(game.Id(), "game1")
	c2 := func(dest, src uint, msgType string, data ...interface{}) {
		log.Info("%x, %x, %v", src, dest, data)
		core.Send(src, dest, "copy1", "copy2")
	}
	c3 := func(dest, src uint, msgType string, rid int, data ...interface{}) {
		log.Info("request: %x, %x, %v, %v", src, dest, rid, data)
		core.Respond(src, dest, rid, data...)
	}
	c4 := func(dest, src uint, msgType string, data ...interface{}) {
		log.Info("call: %x, %x, %v", src, dest, data)
		core.Ret(src, dest, data...)
	}
	game.SetSelf(game)
	game.RegisterBaseCB(core.MSG_TYPE_CLOSE, (*Game).Close, true)
	game.RegisterBaseCB(core.MSG_TYPE_NORMAL, c2, false)
	game.RegisterBaseCB(core.MSG_TYPE_REQUEST, c3, false)
	game.RegisterBaseCB(core.MSG_TYPE_CALL, c4, false)
	go func() {
		for msg := range game.In() {
			game.DispatchM(msg)
		}
	}()

	for {
	}
}
