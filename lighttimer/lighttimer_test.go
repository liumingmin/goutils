package lighttimer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestStartTicks(t *testing.T) {
	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddTimer(time.Second*time.Duration(2), time.Now().Add(time.Second*5), func(fireSeqNo uint) bool {
		fmt.Println("callback", fireSeqNo, "-")
		if fireSeqNo == 4 {
			return true
		}
		return false
	})

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
		lt.AddTimer(time.Second*time.Duration(timeout), time.Now().Add(time.Hour), func(fireSeqNo uint) bool {
			fmt.Println("callback", tmp, "-", timeout)
			return true
		})
	}

	time.Sleep(time.Hour)
}
