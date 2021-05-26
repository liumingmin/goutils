package container

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/liumingmin/goutils/safego"
)

var fb = BufferInvoker{ChanSize: 100, Func: processItem}

func TestFuncBuffer(t *testing.T) {
	for i := 0; i < 100; i++ {
		item := strconv.Itoa(i)
		safego.Go(func() {
			fb.Invoke("1234", item)
		})
	}

	fmt.Println("for end1")

	time.Sleep(time.Second * 10)

	for i := 0; i < 100; i++ {
		item := strconv.Itoa(i)
		safego.Go(func() {
			fb.Invoke("1234", item)
		})
	}

	fmt.Println("for end2")

	time.Sleep(time.Second * 60)
}

func processItem(item interface{}) {
	i := item.(string)
	fmt.Println("process", i)

	fmt.Println("processend", i)
}
