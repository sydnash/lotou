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
		go func() {
			log.Info("send remoteId: %v", g.remoteId)
			g.RawSend(g.remoteId, core.MSG_TYPE_NORMAL, "testNormal", 1, 2, 3, 4, "are you ok?")

			t := func(timeout bool, data string) {
				fmt.Println("request respond ", timeout, data)
			}
			g.Request(g.remoteId, core.MSG_ENC_TYPE_GO, 10, t, "testRequest", "hello")

			fmt.Println(g.Call(g.remoteId, core.MSG_ENC_TYPE_GO, "testCall", "hello"))
		}()
	}
}

func (g *Game) OnCloseNotify() {
	log.Info("recieve close notify")
	g.SendClose(g.Id, true)
}

func TestSlave(t *testing.T) {
	conf.CoreIsStandalone = false
	conf.CoreIsMaster = false

	game := &Game{core.NewSkeleton(1000), 0}
	lotou.Start(nil, &lotou.ModuleParam{"game2", game})
}
