package container

import (
	"sync"
	"time"

	"github.com/liumingmin/goutils/safego"
)

type BufferInvoker struct {
	sync.Map
	ChanSize int
	Func     func(interface{})
}

func (b *BufferInvoker) Invoke(key, item interface{}) {
	if value, exists := (*b).Load(key); exists {
		if itemChan, ok := value.(chan interface{}); ok {
			select {
			case itemChan <- item:
			default:
			}
		}
	} else {
		itemChan := make(chan interface{}, b.ChanSize)

		act, loaded := (*b).LoadOrStore(key, itemChan)
		if loaded { // 被其他线程创建
			itemChan = nil
			if itemChan2, ok := act.(chan interface{}); ok {
				select {
				case itemChan2 <- item:
				default:
				}
			}
		} else {
			itemChan <- item
			safego.Go(func() {
				defer func() {
					(*b).Delete(key)
					//fmt.Println("exit gorruntine")
				}()

				var waitCnt = 0
				for {
					select {
					case procItem, ok2 := <-itemChan:
						if !ok2 {
							return
						}
						b.Func(procItem)
					default:
						time.Sleep(time.Second)
						waitCnt++

						if waitCnt > 5 {
							return
						}
					}
				}
			})
		}
	}
}
