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

	t4, t5 := core.GetIdByName("game1")
	log.Debug("get id by name:%v, %v, %v", "game1", t4, t5)
	if t5 {
		for {
			core.Send(t4, 0XFF<<24, "你好，我现在进行测试")
			time.Sleep(time.Second)
		}
	}
	for {
	}
}
