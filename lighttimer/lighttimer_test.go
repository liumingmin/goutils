package lighttimer

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestStartTicks(t *testing.T) {
	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddTimer(time.Second*time.Duration(2), func(fireSeqNo uint) bool {
		fmt.Println("callback", fireSeqNo, "-")
		if fireSeqNo == 4 {
			return true
		}
		return false
	})

	time.Sleep(time.Hour)
}

func TestStartTicksDeadline(t *testing.T) {

	//NewLightTimerPool

	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddTimerWithDeadline(time.Second*time.Duration(2), time.Now().Add(time.Second*5), func(seqNo uint) bool {
		fmt.Println("callback", seqNo, "-")
		if seqNo == 4 {
			return true
		}
		return false
	}, func(seqNo uint) bool {
		fmt.Println("end callback", seqNo, "-")
		return true
	})

	time.Sleep(time.Hour)
}

func TestLtPool(t *testing.T) {
	pool := NewLightTimerPool(10, time.Millisecond)

	for i := 0; i < 100000; i++ {
		tmp := i
		pool.AddTimerWithDeadline(strconv.Itoa(tmp), time.Second*time.Duration(2), time.Now().Add(time.Second*5), func(seqNo uint) bool {
			fmt.Println("callback", tmp, "-", seqNo, "-")
			if seqNo == 4 {
				return true
			}
			return false
		}, func(seqNo uint) bool {
			fmt.Println("end callback", tmp, "-", seqNo, "-")
			return true
		})
	}

	time.Sleep(time.Second * 20)

	fmt.Println(runtime.NumGoroutine())

	time.Sleep(time.Hour)
}

func TestStartTicks2(t *testing.T) {
	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddCallback(time.Second*time.Duration(3), func() {
		fmt.Println("invoke once")
	})

	time.Sleep(time.Hour)
}

func BenchmarkStartTicks(b *testing.B) {
	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	for i := 0; i < 100000; i++ {
		tmp := i
		timeout := 1 + rand.Intn(20)
		lt.AddTimer(time.Second*time.Duration(timeout), func(fireSeqNo uint) bool {
			fmt.Println("callback", tmp, "-", timeout)
			return true
		})
	}

	time.Sleep(time.Hour)
}
