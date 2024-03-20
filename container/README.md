**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [container](#container)
  * [bitmap_test.go](#bitmap_testgo)
  * [buff_pool_test.go](#buff_pool_testgo)
  * [const_hash_test.go](#const_hash_testgo)
  * [mdb_test.go](#mdb_testgo)
  * [queue_test.go](#queue_testgo)

<!-- tocstop -->

# container
## bitmap_test.go
### TestBitmapExists
```go

bitmap := initTestData()
t.Log(bitmap)

if !bitmap.Exists(122) {
	t.FailNow()
}

if bitmap.Exists(123) {
	t.FailNow()
}
```
### TestBitmapSet
```go

bitmap := initTestData()

if bitmap.Exists(1256) {
	t.FailNow()
}

bitmap.Set(1256)

if !bitmap.Exists(1256) {
	t.FailNow()
}
```
### TestBitmapUnionOr
```go

bitmap := initTestData()
bitmap2 := initTestData()
bitmap2.Set(256)

bitmap3 := bitmap.Union(&bitmap2)
if !bitmap3.Exists(256) {
	t.FailNow()
}

bitmap3.Set(562)

if !bitmap3.Exists(562) {
	t.FailNow()
}

if bitmap.Exists(562) {
	t.FailNow()
}
```
### TestBitmapBitInverse
```go

bitmap := initTestData()

if !bitmap.Exists(66) {
	t.FailNow()
}

bitmap.Inverse()

if bitmap.Exists(66) {
	t.FailNow()
}
```
## buff_pool_test.go
### TestBuffPool
```go

buf1 := GetPoolBuff(BUFF_128K)
ptr1 := uintptr(unsafe.Pointer(&buf1[0]))
len1 := len(buf1)
PutPoolBuff(BUFF_128K, buf1)

buf2 := GetPoolBuff(BUFF_128K)
ptr2 := uintptr(unsafe.Pointer(&buf2[0]))
len2 := len(buf2)
PutPoolBuff(BUFF_128K, buf2)

if len1 != BUFF_128K {
	t.Error("pool get BUFF_128K len failed")
}

if len1 != len2 {
	t.Error("pool get BUFF_128K len failed")
}

if ptr1 != ptr2 {
	t.Error("pool get BUFF_128K failed")
}

//4M
buf1 = GetPoolBuff(BUFF_4M)
ptr1 = uintptr(unsafe.Pointer(&buf1[0]))
len1 = len(buf1)
PutPoolBuff(BUFF_4M, buf1)

buf2 = GetPoolBuff(BUFF_4M)
ptr2 = uintptr(unsafe.Pointer(&buf2[0]))
len2 = len(buf2)
PutPoolBuff(BUFF_4M, buf2)

if len1 != BUFF_4M {
	t.Error("pool get BUFF_4M len failed")
}

if len1 != len2 {
	t.Error("pool get BUFF_4M len failed")
}

if ptr1 != ptr2 {
	t.Error("pool get BUFF_4M failed")
}
```
## const_hash_test.go
### TestConstHash
```go


var ringchash CHashRing

var configs []CHashNode
for i := 0; i < 10; i++ {
	configs = append(configs, testConstHashNode("node"+strconv.Itoa(i)))
}

ringchash.Adds(configs)

t.Log("init:", ringchash.Debug())

if ringchash.GetByC32(100, false).Id() != "node0" {
	t.Fail()
}

if ringchash.GetByC32(134217727, false).Id() != "node0" {
	t.Fail()
}

if ringchash.GetByC32(134217728, false).Id() != "node8" {
	t.Fail()
}

var configs2 []CHashNode
for i := 0; i < 2; i++ {
	configs2 = append(configs2, testConstHashNode("node"+strconv.Itoa(10+i)))
}
ringchash.Adds(configs2)

t.Log("add 2 nodes", ringchash.Debug())

if ringchash.GetByC32(134217727, false).Id() != "node10" {
	t.Fail()
}

if ringchash.GetByC32(134217728, false).Id() != "node10" {
	t.Fail()
}

ringchash.Del("node0")
t.Log("del 1 node", ringchash.Debug())

if ringchash.GetByC32(100, false).Id() != "node10" {
	t.Fail()
}

t.Log(ringchash.GetByKey("goutils", false))
```
## mdb_test.go
### TestDataTable
```go

if len(testDt.Rows()) != 10 {
	t.FailNow()
}

if testDt.PkString(testDt.Row("9")) != "9" {
	t.FailNow()
}

if testDt.PkInt(testDt.Row("8")) != 8 {
	t.FailNow()
}

if reflect.DeepEqual(testDt.Row("2"), []string{"2", "C2", "N2"}) {
	t.FailNow()
}

if reflect.DeepEqual(testDt.RowsBy("code", "C2")[0], []string{"2", "C2", "N2"}) {
	t.FailNow()
}

if reflect.DeepEqual(testDt.RowsByPredicate(func(dr *DataRow) bool { return dr.String("name") == "N4" })[0], []string{"4", "C4", "N4"}) {
	t.FailNow()
}

testDt.Push([]string{"2", "C2", "N3"})

if reflect.DeepEqual(testDt.RowsByIndexPredicate("code", "C2", func(dr *DataRow) bool { return dr.String("name") == "N3" })[0], []string{"2", "C2", "N2"}) {
	t.FailNow()
}
```
### TestDataSet
```go

if testDs.Table("testDt") != testDt {
	t.FailNow()
}
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
