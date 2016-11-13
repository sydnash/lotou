package core_test

import "testing"
import "fmt"
import "github.com/sydnash/majiang/core"
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
	m = &Game{Base: core.NewBase()}
	core.RegisterService(m)

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
						fmt.Println(msg.Src, msg.Dest, msg.Type)
						m.Close()
						a <- 1
						break
					} else {
						fmt.Println(msg.Src, msg.Dest, msg.Type, msg.Data[0].(string))
					}
				}
			}
		}
	}()

	go func() {
		for {
			if !(core.Send(m.Id(), m2.Id(), "kdjfajdfkdf", "aksjdflkajsdf")) {
				m2.Close()
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
