package lighttimer

import (
	"time"
	"fmt"
	"math/rand"
)

func TStartTicks() {
	StartTicks(time.Millisecond)

	AddTimer(time.Second* time.Duration(2), func(fireSeqNo uint) bool {
		fmt.Println("callback",fireSeqNo,"-")
		if fireSeqNo == 4{
			return true
		}
		return false
	})

	time.Sleep(time.Hour)
}

func TStartTicks2() {
	StartTicks(time.Millisecond)

	AddCallback(time.Second* time.Duration(3), func() {
		fmt.Println("invoke once")
	})

	time.Sleep(time.Hour)
}

func BStartTicks(){
	StartTicks(time.Millisecond)

	for i:=0;i<100000;i++{
		tmp := i
		timeout:=1+rand.Intn(20)
		AddTimer(time.Second* time.Duration(timeout), func(fireSeqNo uint) bool {
			fmt.Println("callback",tmp,"-",timeout)
			return true
		})
	}

	time.Sleep(time.Hour)
}
