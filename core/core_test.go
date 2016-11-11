package core_test

import "testing"
import "fmt"
import "github.com/sydnash/majiang/core"
import "time"

func f(m *core.Message) {
	_ = m
}

type Game struct {
	in chan *core.Message
	id uint
}

type Game2 struct {
	Game
}

var m *Game
var m2 *Game2

func init() {
	m = &Game{}
	m.in = make(chan *core.Message, 1000)
	m.id = core.RegisterService(m)

	m2 = &Game2{}
	m2.in = make(chan *core.Message, 1000)
	m2.id = core.RegisterService(m2)
}
func (self *Game) Send(m *core.Message) {
	self.in <- m
}

func TestCore(t *testing.T) {
	a := make(chan int)
	go func() {
		for {
			select {
			case msg, ok := <-m.in:
				if ok {
					if msg.Type == core.MSG_TYPE_CLOSE {
						fmt.Println(msg.Src, msg.Dest, msg.Type)
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
			time.Sleep(time.Second)
			if !(core.Send(m.id, m2.id, "kdjfajdfkdf", "aksjdflkajsdf")) {
				a <- 1
				break
			}
			time.AfterFunc(time.Second*10, func() {
				core.Close(m.id, m2.id)
			})
		}
	}()

	<-a
	<-a
}
