package lotou_test

import (
	"github.com/sydnash/lotou"
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"testing"
)

type Game struct {
	*core.Skeleton
	remoteId uint
}

func (g *Game) OnRequestMSG(src uint, rid int, data ...interface{}) {
	g.Respond(src, rid, "world")
}
func (g *Game) OnCallMSG(src uint, cid int, data ...interface{}) {
	g.Ret(src, cid, "world")
}

func (g *Game) OnNormalMSG(src uint, data ...interface{}) {
	log.Info("%v, %v", src, data)
	//g.RawSend(src, core.MSG_TYPE_NORMAL, "222")
}
func (g *Game) OnDistributeMSG(data ...interface{}) {
	log.Info("%v", data)
}
func (g *Game) OnInit() {
	log.Info("OnInit: name:%v  id:%v", g.Name, g.Id)
	g.remoteId, _ = core.NameToId("game1")
	log.Info("name2id: game1:%v", g.remoteId)
}

func TestMaster(t *testing.T) {
	conf.CoreIsStandalone = false
	conf.CoreIsMaster = true
	game := &Game{core.NewSkeleton(0), 0}
	lotou.Start(nil, &lotou.ModuleParam{"game1", game})
}
