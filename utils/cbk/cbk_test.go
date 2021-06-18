package cbk

import (
	"fmt"
	"testing"
	"time"
)

func TestCbkFailed(t *testing.T) {
	var ok bool
	var lastBreaked bool
	for j := 0; j < 200; j++ {
		i := j
		//safego.Go(func() {
		err := Impls[SIMPLE].Check("test") //30s 返回一次true尝试
		fmt.Println(i, "Check:", ok)

		if err == nil {
			time.Sleep(time.Millisecond * 10)
			Impls[SIMPLE].Failed("test")

			if i > 105 && lastBreaked {
				Impls[SIMPLE].Succeed("test")
				lastBreaked = false
				fmt.Println(i, "Succeed")
			}
		} else {
			if lastBreaked {
				time.Sleep(time.Second * 10)
			} else {
				lastBreaked = true
			}
		}
		//})
	}
}
