package timer

import (
	"github.com/sydnash/lotou/vector"
	"sync"
)

type TimerSchedule struct {
	timers      *vector.Vector
	addedCache  *vector.Vector
	deleteCache *vector.Vector
	mutex       sync.Mutex
}

func NewTS() *TimerSchedule {
	ts := &TimerSchedule{}
	ts.timers = vector.New()
	ts.addedCache = vector.New()
	ts.deleteCache = vector.New()
	return ts
}

//Update update all timers
func (ts *TimerSchedule) Update(dt int) {
	ts.mutex.Lock()
	ts.timers.AppendVec(ts.addedCache)
	ts.addedCache.Clear()
	ts.mutex.Unlock()
	for i := 0; i < ts.timers.Len(); i++ {
		t := ts.timers.At(i).(*Timer)
		t.update(dt)
		if t.isComplete {
			ts.Unschedule(t)
		}
	}
	ts.mutex.Lock()
	for i := 0; i < ts.deleteCache.Len(); i++ {
		t := ts.deleteCache.At(i)
		for i := 0; i < ts.timers.Len(); i++ {
			if ts.timers.At(i) == t {
				ts.timers.Delete(i)
				break
			}
		}
	}
	ts.deleteCache.Clear()
	ts.mutex.Unlock()
}

//Schedule start a timer with interval and repeat.
//it's callback with each interval, and timer will delete after trigger repeat times
//if interval is small than schedule's interval
//it may trigger multitimes at a update.
func (ts *TimerSchedule) Schedule(interval, repeat int, cb TimerCallback) *Timer {
	t := NewTimer(interval, repeat, cb)
	ts.mutex.Lock()
	ts.addedCache.Push(t)
	ts.mutex.Unlock()
	return t
}

func (ts *TimerSchedule) Unschedule(t *Timer) {
	ts.mutex.Lock()
	ts.deleteCache.Push(t)
	ts.mutex.Unlock()
	t.cancel()
}
