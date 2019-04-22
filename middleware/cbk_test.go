package middleware

import (
	"fmt"
	"testing"
	"time"
)

func TestCbkFailed(t *testing.T) {
	cbkone := &CircuitBreaker{}
	cbkone.Init()
	cbkone.isTurnOn = false

	var ok bool
	var lastBreaked bool
	for i := 0; i < 200; i++ {
		ok = cbkone.Check("test") //30s 返回一次true尝试
		fmt.Println(i, "Check:", ok)

		if ok {
			time.Sleep(time.Millisecond * 10)
			cbkone.Failed("test")

			if i > 105 && lastBreaked {
				cbkone.Succeed("test")
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
	}
}
