package core

import (
	"fmt"
	"testing"
	"time"
)

type Game struct {
	*Skeleton
	Dst ServiceID
}

func (g *Game) OnRequestMSG(src ServiceID, rid uint64, data ...interface{}) {
	g.Respond(src, rid, "world")
}
func (g *Game) OnCallMSG(src ServiceID, cid uint64, data ...interface{}) {
	g.Ret(src, cid, "world")
}

func (g *Game) OnMainLoop(dt int) {
	g.Send(g.Dst, MSG_TYPE_NORMAL, g.Name, []byte{1, 2, 3, 4, 56})
	g.RawSend(g.Dst, MSG_TYPE_NORMAL, g.Name, g.Id)

	t := func(timeout bool, data ...interface{}) {
		fmt.Println("request respond ", timeout, data)
	}
	g.Request(g.Dst, 10, t, func() {
	}, "hello")

	fmt.Println(g.Call(g.Dst, "hello"))
}

func (g *Game) OnNormalMSG(src ServiceID, data ...interface{}) {
	fmt.Println(src, data)
}

func TestModule(t *testing.T) {
	id1 := StartService("g1", &Game{Skeleton: &Skeleton{D: 0}})
	StartService("g2", &Game{Skeleton: &Skeleton{D: 1000}, Dst: id1})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
