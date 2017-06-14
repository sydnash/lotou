package timer

import (
	"errors"
	"github.com/sydnash/lotou/helper"
	"github.com/sydnash/lotou/log"
)

var (
	TimerIsComplete = errors.New("timer is complete")
)

type TimerCallback func(int)
type Timer struct {
	cb         TimerCallback
	interval   int //interval time of milloseconds per trigger
	elapsed    int //time elapsed
	repeat     int //repeat times, <= 0 forever
	repeated   int //allready repeated times
	isComplete bool
	isForever  bool
}

func NewTimer(interval, repeat int, cb TimerCallback) *Timer {
	if interval <= 0 {
		panic("NewTimer: interval is negative or zero.")
	}
	t := &Timer{}
	t.interval = interval
	t.cb = cb
	t.repeat = repeat
	t.isForever = (t.repeat <= 0)
	return t
}

func (t *Timer) update(dt int) {
	if t.isComplete {
		return
	}

	t.elapsed += dt
	if t.elapsed < t.interval {
		return
	}

	for t.elapsed >= t.interval {
		t.elapsed -= t.interval
		t.repeated += 1

		t.trigger()

		if !t.isForever {
			if t.repeated >= t.repeat {
				t.isComplete = true
				return
			}
		}
	}
}

func (t *Timer) trigger() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Timer:trigger stack: %v\n, %v", helper.GetStack(), err)
		}
	}()
	t.cb(t.interval)
}

//Reset reset timer's time elapsed and repeated times.
func (t *Timer) Reset() error {
	if t.isComplete {
		return TimerIsComplete
	}
	t.elapsed = 0
	t.repeated = 0
	return nil
}

func (t *Timer) cancel() {
	t.isComplete = true
}
