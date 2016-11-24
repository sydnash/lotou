package topology_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
	"time"
)

func TestSlavea(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	log.Debug("start slave")
	topology.StartSlave("127.0.0.1", "4000")
	core.RegisterNode()
	core.Name(1, "我是服务1")
	t1, t2 := core.GetIdByName("service3")
	log.Debug("get id by name:%v, %v, %v", "service3", t1, t2)
	t1, t2 = core.GetIdByName("service4")
	log.Debug("get id by name:%v, %v, %v", "service4", t1, t2)
	game := &Game{core.NewBase()}
	core.RegisterService(game)

	t4, t5 := core.GetIdByName("game1")
	log.Debug("get id by name:%v, %v, %v", "game1", t4, t5)
	c2 := func(dest, src uint, msgType string, data ...interface{}) {
		log.Info("%x, %x, %v", src, dest, data)
	}
	game.RegisterBaseCB(core.MSG_TYPE_NORMAL, c2, false)
	response := func(entype string, data ...interface{}) {
		log.Info("respond: %v", data)
	}
	go func() {
		for msg := range game.In() {
			game.DispatchM(msg)
		}
	}()
	if t5 {
		for {
			core.Send(t4, game.Id(), "你好，我现在进行测试", "测试一下")
			time.Sleep(time.Second)

			core.Request(t4, game, response, "response1", "response2")

			ret := core.Call(t4, game, "call1", "call2")
			log.Info("ret: %v", ret)
		}
	}
	for {
	}
}
