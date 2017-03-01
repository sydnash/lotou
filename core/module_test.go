package core

import (
	"fmt"
	"testing"
	"time"
)

type Game struct {
	*Skeleton
	Dst uint
}

func (g *Game) OnMainLoop(dt int) {
	g.Send(g.Dst, MSG_TYPE_NORMAL, g.Name, []byte{1, 2, 3, 4, 56})
	g.RawSend(g.Dst, MSG_TYPE_NORMAL, g.Name, g.Id)
}

func (g *Game) OnNormalMSG(src uint, data ...interface{}) {
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
