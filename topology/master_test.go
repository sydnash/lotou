package topology_test

import (
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/topology"
	"testing"
	"time"
)

type Game struct {
	*core.Skeleton
}

func (g *Game) OnRequestMSG(src core.ServiceID, rid int, data ...interface{}) {
	g.Respond(src, rid, "world")
}
func (g *Game) OnCallMSG(src core.ServiceID, cid int, data ...interface{}) {
	g.Ret(src, cid, "world")
}

func (g *Game) OnNormalMSG(src core.ServiceID, data ...interface{}) {
	log.Info("%v, %v", src, data)
	//g.RawSend(src, core.MSG_TYPE_NORMAL, "222")
}
func (g *Game) OnDistributeMSG(data ...interface{}) {
	log.Info("%v", data)
}
func TestMaster(t *testing.T) {
	log.Init("test", log.FATAL_LEVEL, log.DEBUG_LEVEL, 10000, 1000)

	core.InitNode(false, true)
	topology.StartMaster("127.0.0.1", "4000")
	core.RegisterNode()

	game := &Game{core.NewSkeleton(0)}
	id := core.StartService("game1", game)
	log.Info("game1's id :%v", id)

	log.Info("test")
	for {
		time.Sleep(time.Minute * 10)
	}
}
