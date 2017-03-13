package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	tick := time.NewTicker(time.Duration(100) * time.Millisecond)

	ts := NewTS()

	ts.Schedule(100, 10, func(dt int) {
		fmt.Println("time1")
	})
	ts.Schedule(10, 10, func(dt int) {
		fmt.Println("timet-----")
	})
	ts.Schedule(1000, 1, func(dt int) {
		fmt.Println("time2")
	})
	ts.Schedule(100, -1, func(dt int) {
		fmt.Println("time3")
	})

	go func() {
		for {
			<-tick.C
			ts.Update(10)
		}
	}()

	ch := make(chan int)
	<-ch
}
