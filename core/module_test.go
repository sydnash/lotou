package core

import (
	"fmt"
	"github.com/sydnash/lotou/conf"
	"github.com/sydnash/lotou/log"
	"testing"
	"time"
)

type Game struct {
	*Skeleton
	Dst ServiceID
}

type XMsg struct {
	A int32
	B string
	C int64
}

/*
func (g *Game) OnRequestMSG(src ServiceID, encType int32, rid uint64, data ...interface{}) {
	g.Respond(src, encType, rid, "world")
}
func (g *Game) OnCallMSG(src ServiceID, encType int32, cid uint64, data ...interface{}) {
	g.Ret(src, encType, cid, "world")
}
func (g *Game) OnNormalMSG(src ServiceID, encType int32, data ...interface{}) {
	fmt.Println(src, encType, data)
}
*/

func (g *Game) OnMainLoop(dt int) {
	g.Send(g.Dst, MSG_TYPE_NORMAL, MSG_ENC_TYPE_GO, "testNormal", g.Name, []byte{1, 2, 3, 4, 56})
	g.RawSend(g.Dst, MSG_TYPE_NORMAL, "testNormal", g.Name, g.Id)

	t := func(timeout bool, data ...interface{}) {
		fmt.Println("request respond ", timeout, data)
	}
	g.Request(g.Dst, MSG_ENC_TYPE_GO, 10, t, "testRequest", "hello")

	fmt.Println(g.Call(g.Dst, MSG_ENC_TYPE_GO, "testCall", "hello"))
}

func (g *Game) OnInit() {
	//test for go and no enc
	g.SubscribeFunc(MSG_TYPE_NORMAL, "testNormal", func(src ServiceID, data ...interface{}) {
		log.Info("%v, %v", src, data)
	})
	g.SubscribeFunc(MSG_TYPE_REQUEST, "testRequest", func(src ServiceID, data ...interface{}) string {
		return "world"
	})
	g.SubscribeFunc(MSG_TYPE_CALL, "testCall", func(src ServiceID, data ...interface{}) (string, string) {
		return "hello", "world"
	})
}

func TestModule(t *testing.T) {
	log.Init(conf.LogFilePath, conf.LogFileLevel, conf.LogShellLevel, conf.LogMaxLine, conf.LogBufferSize)
	id1 := StartService("g1", &Game{Skeleton: NewSkeleton(0)})
	StartService("g2", &Game{Skeleton: NewSkeleton(1000), Dst: id1})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
