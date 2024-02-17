

<!-- toc -->

- [container 容器模块](#container-%E5%AE%B9%E5%99%A8%E6%A8%A1%E5%9D%97)
  * [bitmap_test.go 比特位表](#bitmap_testgo-%E6%AF%94%E7%89%B9%E4%BD%8D%E8%A1%A8)
  * [const_hash_test.go 一致性HASH](#const_hash_testgo-%E4%B8%80%E8%87%B4%E6%80%A7hash)
  * [lighttimer_test.go 轻量级计时器](#lighttimer_testgo-%E8%BD%BB%E9%87%8F%E7%BA%A7%E8%AE%A1%E6%97%B6%E5%99%A8)
  * [queue_test.go](#queue_testgo)

<!-- tocstop -->

# container 容器模块
## bitmap_test.go 比特位表
### TestBitmapExists
```go

bitmap := initTestData()
t.Log(bitmap)

t.Log(bitmap.Exists(122))
t.Log(bitmap.Exists(123))

//data1 := []byte{1, 2, 4, 7}
//data2 := []byte{0, 1, 5}

```
### TestBitmapSet
```go

bitmap := initTestData()

t.Log(bitmap.Exists(1256))

bitmap.Set(1256)

t.Log(bitmap.Exists(1256))
```
### TestBitmapUnionOr
```go

bitmap := initTestData()
bitmap2 := initTestData()
bitmap2.Set(256)

bitmap3 := bitmap.Union(&bitmap2)
t.Log(bitmap3.Exists(256))

bitmap3.Set(562)
t.Log(bitmap3.Exists(562))

t.Log(bitmap.Exists(562))
```
### TestBitmapBitInverse
```go

bitmap := initTestData()

t.Log(bitmap.Exists(66))

bitmap.Inverse()

t.Log(bitmap.Exists(66))

```
## const_hash_test.go 一致性HASH
### TestConstHash
```go


var ringchash CHashRing

var configs []CHashNode
for i := 0; i < 10; i++ {
	configs = append(configs, TestNode("node"+strconv.Itoa(i)))
}

ringchash.Adds(configs)

fmt.Println(ringchash.Debug())

fmt.Println("==================================", configs)

fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
//
fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
//
fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

fmt.Println(ringchash.Get("xvd:fs:xcvd"))

var configs2 []CHashNode
for i := 0; i < 2; i++ {
	configs2 = append(configs2, TestNode("node"+strconv.Itoa(10+i)))
}
ringchash.Adds(configs2)
fmt.Println("==================================")
fmt.Println(ringchash.Debug())
fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
//
fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
//
fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

fmt.Println(ringchash.Get("xvd:fs:xcvd"))

ringchash.Del("node0")

fmt.Println("==================================")
fmt.Println(ringchash.Debug())
fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
//
fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
//
fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

fmt.Println(ringchash.Get("xvd:fs:xcvd"))
```
## lighttimer_test.go 轻量级计时器
### TestStartTicks
```go

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
```
### TestStartTicksDeadline
```go


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
```
### TestLtPool
```go

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
```
### TestStartTicks2
```go

lt := NewLightTimer()
lt.StartTicks(time.Millisecond)

lt.AddCallback(time.Second*time.Duration(3), func() {
	fmt.Println("invoke once")
})

time.Sleep(time.Hour)
```
## queue_test.go
### TestEnqueueBack
```go

q := InitQueue()
t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 25; i++ {
	deItem := q.EnqueueBack(fmt.Sprint(i))
	if deItem != "" {
		t.Log("EnqueueBack dequeue front: ", deItem)
	}
}
t.Log("queue is empty ", q.IsEmpty())
q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})

```
### TestDequeueFront
```go

q := InitQueue()
t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 25; i++ {
	deItem := q.EnqueueBack(fmt.Sprint(i))
	if deItem != "" {
		t.Log("EnqueueBack dequeue front: ", deItem)
	}
}

for i := 0; i < 5; i++ {
	deItem := q.DequeueFront()
	if deItem != "" {
		t.Log("DequeueFront dequeue front: ", deItem)
	}
}

q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})

t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 5; i++ {
	deItem := q.DequeueFront()
	if deItem != "" {
		t.Log("DequeueFront dequeue front: ", deItem)
	}
}

q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})

t.Log("queue is empty ", q.IsEmpty())
```
### TestEnqueueFront
```go

q := InitQueue()
t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 25; i++ {
	deItem := q.EnqueueFront(fmt.Sprint(i))
	if deItem != "" {
		t.Log("EnqueueFront dequeue back: ", deItem)
	}
}
t.Log("queue is empty ", q.IsEmpty())
q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})
```
### TestDequeueBack
```go

q := InitQueue()
t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 25; i++ {
	deItem := q.EnqueueFront(fmt.Sprint(i))
	if deItem != "" {
		t.Log("EnqueueFront dequeue back: ", deItem)
	}
}

for i := 0; i < 5; i++ {
	deItem := q.DequeueBack()
	if deItem != "" {
		t.Log("DequeueBack dequeue back: ", deItem)
	}
}
t.Log("queue is empty ", q.IsEmpty())

q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})

for i := 0; i < 5; i++ {
	deItem := q.DequeueBack()
	if deItem != "" {
		t.Log("DequeueBack dequeue back: ", deItem)
	}
}

t.Log("queue is empty ", q.IsEmpty())
```
### TestQueueClear
```go

q := InitQueue()
t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 25; i++ {
	deItem := q.EnqueueBack(fmt.Sprint(i))
	if deItem != "" {
		t.Log("EnqueueBack dequeue front: ", deItem)
	}
}
t.Log("queue is empty ", q.IsEmpty())
q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})

q.Clear()
t.Log("queue is empty ", q.IsEmpty())
q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})
```
### TestQueueFindBy
```go

q := InitQueue()
t.Log("queue is empty ", q.IsEmpty())

for i := 0; i < 25; i++ {
	deItem := q.EnqueueBack(fmt.Sprint(i))
	if deItem != "" {
		t.Log("EnqueueBack dequeue front: ", deItem)
	}
}

q.Range(func(i string) bool {
	t.Log("left:", i)
	return true
})

i20 := q.FindOneBy(func(i string) bool {
	ni, _ := strconv.Atoi(i)
	if ni == 20 {
		return true
	}
	return false
})
if i20 != "20" {
	t.FailNow()
} else {
	t.Log("FindOneBy: ", i20)
}

data := q.FindBy(func(i string) bool {
	ni, _ := strconv.Atoi(i)
	if ni%3 == 0 {
		return true
	}
	return false
})
t.Log("modern 3 :", data)
```
