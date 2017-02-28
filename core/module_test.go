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
	g.Send(g.Dst, MSG_TYPE_NORMAL, g.Name)
}
func (g *Game) OnNormalMSG(src uint, data ...interface{}) {
	fmt.Println(src, data)
}

func TestModule(t *testing.T) {
	id1 := StartService("g1", 0, &Game{Skeleton: &Skeleton{}})
	StartService("g2", 1000, &Game{Skeleton: &Skeleton{}, Dst: id1})

	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()

	<-ch
}
