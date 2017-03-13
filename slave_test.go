package lotou_test

import (
	"fmt"
	"github.com/sydnash/lotou"
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"testing"
)

func (g *Game) OnMainLoop(dt int) {
	log.Info("OnMainLoop %d", dt)
	if g.remoteId != 0 {
		log.Info("send remoteId: %v", g.remoteId)
		g.RawSend(g.remoteId, core.MSG_TYPE_NORMAL, 1, 2, 3, 4, "are you ok?")

		t := func(timeout bool, data string) {
			fmt.Println("request respond ", timeout, data)
		}
		g.Request(g.remoteId, 10, t, func() {
			fmt.Println("slave timeout")
		}, "hello")

		fmt.Println(g.Call(g.remoteId, "hello"))
	}
}

func TestSlave(t *testing.T) {
	conf.CoreIsStandalone = false
	conf.CoreIsMaster = false

	game := &Game{core.NewSkeleton(1000), 0}
	lotou.Start(nil, &lotou.ModuleParam{"game2", game})
}
