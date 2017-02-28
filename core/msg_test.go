package core

import (
	"fmt"
	"testing"
	"time"
)

func TestMSG(t *testing.T) {
	s1 := newService("s1")
	registerService(s1)

	s2 := newService("s2")
	id2 := registerService(s2)
	fmt.Println(id2)
	s2.run()

	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			send(s1, id2, MSG_TYPE_NORMAL, "test string", 1)
			send(s1, id2, MSG_TYPE_NORMAL)
			sendNoEnc(s1, id2, MSG_TYPE_NORMAL, 1.5, []byte{1, 2, 4})
		}
		time.Sleep(time.Duration(10 * time.Second))
		ch <- 1
	}()

	<-ch
}
