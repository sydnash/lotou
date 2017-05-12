package core_test

import (
	"fmt"
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/core"
	"github.com/sydnash/lotou/log"
	"testing"
	"time"
)

type Game struct {
	*core.Skeleton
	Dst core.ServiceID
}

type XMsg struct {
	A int32
	B string
	C int64
}

func (g *Game) OnMainLoop(dt int) {
	g.Send(g.Dst, core.MSG_TYPE_NORMAL, core.MSG_ENC_TYPE_GO, "testNormal", g.Name, []byte{1, 2, 3, 4, 56})
	g.RawSend(g.Dst, core.MSG_TYPE_NORMAL, "testNormal", g.Name, g.Id)

	t := func(timeout bool, data ...interface{}) {
		fmt.Println("request respond ", timeout, data)
	}
	g.Request(g.Dst, core.MSG_ENC_TYPE_GO, 10, t, "testRequest", "hello")

	fmt.Println(g.Call(g.Dst, core.MSG_ENC_TYPE_GO, "testCall", "hello"))
}

func (g *Game) OnInit() {
	//test for go and no enc
	g.RegisterHandlerFunc(core.MSG_TYPE_NORMAL, "testNormal", func(src core.ServiceID, data ...interface{}) {
		log.Info("%v, %v", src, data)
	}, true)
	g.RegisterHandlerFunc(core.MSG_TYPE_REQUEST, "testRequest", func(src core.ServiceID, data ...interface{}) string {
		return "world"
	}, true)
	g.RegisterHandlerFunc(core.MSG_TYPE_CALL, "testCall", func(src core.ServiceID, data ...interface{}) (string, string) {
		return "hello", "world"
	}, true)
}

func TestModule(t *testing.T) {
	log.Init(conf.LogFilePath, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)
	id1 := core.StartService(&core.ModuleParam{
		N: "g1",
		M: &Game{Skeleton: core.NewSkeleton(0)},
		L: 0,
	})
	core.StartService(&core.ModuleParam{
		N: "g2",
		M: &Game{Skeleton: core.NewSkeleton(1000), Dst: id1},
		L: 0,
	})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
