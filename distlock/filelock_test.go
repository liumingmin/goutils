package distlock

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestFileLock(t *testing.T) {
	test_file_path, _ := os.Getwd()
	locked_file := test_file_path

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			flock := NewFileLock(locked_file, false)
			err := flock.Lock()
			if err != nil {
				wg.Done()
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("output : %d\n", num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(2 * time.Second)

}
