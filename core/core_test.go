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
		for {
			select {
			case msg, ok := <-m.In():
				if ok {
					if msg.Type == core.MSG_TYPE_CLOSE {
						m.Close()
						log.Info("%v, %v, %v", msg.Src, msg.Dest, msg.Type)
						a <- 1
						break
					} else {
						log.Info("%v, %v, %v, %v", msg.Src, msg.Dest, msg.Type, msg.Data[0].(string))
					}
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
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
