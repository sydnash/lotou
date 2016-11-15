package core_test

import "testing"
import "github.com/sydnash/lotou/core"
import "github.com/sydnash/lotou/log"
import "time"

func f(m *core.Message) {
	_ = m
}

type Game struct {
	*core.Base
}

func (g *Game) Close(dest, src uint) {
	g.Base.Close()
	log.Info("close: %v, %v", src, dest)
}

type Game2 struct {
	Game
}

var m *Game
var m2 *Game2

func init() {
	log.Init("log", log.FATAL_LEVEL, log.INFO_LEVEL, 10000, 1000)
	m = &Game{core.NewBase()}
	core.RegisterService(m)
	core.Name(m.Id(), "m1")

	m2 = &Game2{Game: Game{Base: core.NewBase()}}
	core.RegisterService(m2)
}

func TestCore(t *testing.T) {
	a := make(chan int)
	go func() {
		c2 := func(dest, src uint, data ...interface{}) {
			log.Info("%v, %v, %v", src, dest, data[0].(string))
		}
		m.SetSelf(m)
		m.RegisterBaseCB(core.MSG_TYPE_CLOSE, (*Game).Close, true)
		m.RegisterBaseCB(core.MSG_TYPE_NORMAL, c2, false)
	OUTER_FOR:
		for {
			select {
			case msg, ok := <-m.In():
				if !ok {
					log.Info("m.In is closed.")
					a <- 1
					break OUTER_FOR
				}
				m.DispatchM(msg)
			}
		}
	}()

	go func() {
		for {
			if !(core.SendName("m1", m2.Id(), "kdjfajdfkdf", "aksjdflkajsdf")) {
				a <- 1
				break
			}
			time.Sleep(time.Second)
			time.AfterFunc(time.Second*10, func() {
				core.Close(m.Id(), m2.Id())
			})
		}
	}()

	<-a
	<-a
}
