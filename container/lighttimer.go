package container

import (
	"container/heap"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

const (
	MIN_TIMER_INTERVAL = 1 * time.Millisecond
)

// 定时执行函数
// 入参回调执行序号(次数)
// 回参是否完成，true停止执行，false继续执行，在dealine的callback中无效
type CallbackFunc func(uint) (done bool)

type LightTimer struct {
	timerHeap     lightTimerHeap
	timerHeapLock sync.Mutex
	nextAddSeq    uint
}

func NewLightTimer() *LightTimer {
	lt := &LightTimer{nextAddSeq: 1}
	heap.Init(&lt.timerHeap)

	return lt
}

// Add a callback which will be called after specified duration
func (lt *LightTimer) AddCallback(d time.Duration, callback func()) *Timer {
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: func(u uint) (done bool) {
			callback()
			return true
		},
		repeat: false,
	}
	lt.timerHeapLock.Lock()
	t.addseq = lt.nextAddSeq // set addseq when locked
	lt.nextAddSeq += 1

	heap.Push(&lt.timerHeap, t)
	lt.timerHeapLock.Unlock()
	return t
}

func (lt *LightTimer) AddTimer(d time.Duration, callback CallbackFunc) *Timer {
	return lt.AddTimerWithDeadline(d, time.Time{}, callback, nil)
}

// Add a timer which calls callback periodly
func (lt *LightTimer) AddTimerWithDeadline(d time.Duration, deadline time.Time, callback CallbackFunc, deadlineCallback CallbackFunc) *Timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}

	t := &Timer{
		fireTime:         time.Now().Add(d),
		deadline:         deadline,
		interval:         d,
		callback:         callback,
		deadlineCallback: deadlineCallback,
		repeat:           true,
	}
	lt.timerHeapLock.Lock()
	t.addseq = lt.nextAddSeq // set addseq when locked
	lt.nextAddSeq += 1

	heap.Push(&lt.timerHeap, t)
	lt.timerHeapLock.Unlock()
	return t
}

// Tick once for timers
func (lt *LightTimer) tick() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "tick paniced: %v\n", err)
			debug.PrintStack()
		}
	}()

	now := time.Now()
	lt.timerHeapLock.Lock()

	for {
		if lt.timerHeap.Len() <= 0 {
			break
		}

		nextFireTime := lt.timerHeap.timers[0].fireTime
		//fmt.Printf(">>> nextFireTime %s, now is %s\n", nextFireTime, now)
		if nextFireTime.After(now) {
			break
		}

		t := heap.Pop(&lt.timerHeap).(*Timer)

		callback := t.callback
		if callback == nil {
			continue
		}

		if !t.repeat {
			t.callback = nil
		}

		t.fireSeqNo++

		fireSeqNo := t.fireSeqNo
		// unlock the lock to run callback, because callback may add more callbacks / timers
		lt.timerHeapLock.Unlock()
		done := runCallback(callback, fireSeqNo)
		lt.timerHeapLock.Lock()

		tNow := time.Now()
		if t.repeat && !done && t.deadline.After(tNow) {
			// add Timer back to heap
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) { // might happen when interval is very small
				t.fireTime = now.Add(t.interval)
			}
			t.addseq = lt.nextAddSeq
			lt.nextAddSeq += 1
			heap.Push(&lt.timerHeap, t)
		}

		deadlineCallback := t.deadlineCallback
		if deadlineCallback != nil && (!t.deadline.After(tNow)) {
			lt.timerHeapLock.Unlock()
			runCallback(deadlineCallback, fireSeqNo)
			lt.timerHeapLock.Lock()
		}
	}
	lt.timerHeapLock.Unlock()
}

// Start the self-ticking routine, which ticks per tickInterval
func (lt *LightTimer) StartTicks(tickInterval time.Duration) {
	go lt.selfTickRoutine(tickInterval)
}

func (lt *LightTimer) selfTickRoutine(tickInterval time.Duration) {
	for {
		time.Sleep(tickInterval)
		lt.tick()
	}
}

func runCallback(callback CallbackFunc, fireSeqNo uint) (done bool) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback %v paniced: %v\n", callback, err)
			debug.PrintStack()
		}
	}()
	return callback(fireSeqNo)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Timer struct {
	fireTime         time.Time
	interval         time.Duration
	deadline         time.Time
	callback         CallbackFunc
	deadlineCallback CallbackFunc
	repeat           bool
	addseq           uint

	fireSeqNo uint
}

func (t *Timer) Cancel() {
	t.callback = nil
}

func (t *Timer) IsActive() bool {
	return t.callback != nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type lightTimerHeap struct {
	timers []*Timer
}

func (h *lightTimerHeap) Len() int {
	return len(h.timers)
}

func (h *lightTimerHeap) Less(i, j int) bool {
	//log.Println(h.timers[i].fireTime, h.timers[j].fireTime)
	t1, t2 := h.timers[i].fireTime, h.timers[j].fireTime
	if t1.Before(t2) {
		return true
	}

	if t1.After(t2) {
		return false
	}
	// t1 == t2, making sure Timer with same deadline is fired according to their add order
	return h.timers[i].addseq < h.timers[j].addseq
}

func (h *lightTimerHeap) Swap(i, j int) {
	var tmp *Timer
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *lightTimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(*Timer))
}

func (h *lightTimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
