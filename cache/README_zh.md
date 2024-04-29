**其他语言版本: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [缓存模块](#%E7%BC%93%E5%AD%98%E6%A8%A1%E5%9D%97)
  * [cache_test.go](#cache_testgo)
  * [内存缓存](#%E5%86%85%E5%AD%98%E7%BC%93%E5%AD%98)
  * [Redis缓存](#redis%E7%BC%93%E5%AD%98)

<!-- tocstop -->

# 缓存模块
## cache_test.go
### TestCache
```go

tc := New(0, 0)

a, found := tc.Get("a")
if found || a != nil {
	t.Error("Getting A found value that shouldn't exist:", a)
}

b, found := tc.Get("b")
if found || b != nil {
	t.Error("Getting B found value that shouldn't exist:", b)
}

c, found := tc.Get("c")
if found || c != nil {
	t.Error("Getting C found value that shouldn't exist:", c)
}

tc.Set("a", 1, 0)
tc.Set("b", "b", 0)
tc.Set("c", 3.5, 0)

x, found := tc.Get("a")
if !found {
	t.Error("a was not found while getting a2")
}
if x == nil {
	t.Error("x for a is nil")
} else if a2 := x.(int); a2+2 != 3 {
	t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
}

x, found = tc.Get("b")
if !found {
	t.Error("b was not found while getting b2")
}
if x == nil {
	t.Error("x for b is nil")
} else if b2 := x.(string); b2+"B" != "bB" {
	t.Error("b2 (which should be b) plus B does not equal bB; value:", b2)
}

x, found = tc.Get("c")
if !found {
	t.Error("c was not found while getting c2")
}
if x == nil {
	t.Error("x for c is nil")
} else if c2 := x.(float64); c2+1.2 != 4.7 {
	t.Error("c2 (which should be 3.5) plus 1.2 does not equal 4.7; value:", c2)
}
```
### TestCacheTimes
```go

var found bool

tc := New(50*time.Millisecond, 1*time.Millisecond)
tc.Set("a", 1, 0)
tc.Set("b", 2, -1)
tc.Set("c", 3, 20*time.Millisecond)
tc.Set("d", 4, 70*time.Millisecond)

<-time.After(25 * time.Millisecond)
_, found = tc.Get("c")
if found {
	t.Error("Found c when it should have been automatically deleted")
}

<-time.After(30 * time.Millisecond)
_, found = tc.Get("a")
if found {
	t.Error("Found a when it should have been automatically deleted")
}

_, found = tc.Get("b")
if !found {
	t.Error("Did not find b even though it was set to never expire")
}

_, found = tc.Get("d")
if !found {
	t.Error("Did not find d even though it was set to expire later than the default")
}

<-time.After(20 * time.Millisecond)
_, found = tc.Get("d")
if found {
	t.Error("Found d when it should have been automatically deleted (later than the default)")
}
```
### TestStorePointerToStruct
```go

tc := New(0, 0)
tc.Set("foo", &TestStruct{Num: 1}, 0)
x, found := tc.Get("foo")
if !found {
	t.Fatal("*TestStruct was not found for foo")
}
foo := x.(*TestStruct)
foo.Num++

y, found := tc.Get("foo")
if !found {
	t.Fatal("*TestStruct was not found for foo (second time)")
}
bar := y.(*TestStruct)
if bar.Num != 2 {
	t.Fatal("TestStruct.Num is not 2")
}
```
### TestIncrementUint
```go

tc := New(0, 0)
tc.Set("tuint", uint(1), 0)
_, err := tc.Increment("tuint", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}

x, found := tc.Get("tuint")
if !found {
	t.Error("tuint was not found")
}
if x.(uint) != 3 {
	t.Error("tuint is not 3:", x)
}
```
### TestIncrementUintptr
```go

tc := New(0, 0)
tc.Set("tuintptr", uintptr(1), 0)
_, err := tc.Increment("tuintptr", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}

x, found := tc.Get("tuintptr")
if !found {
	t.Error("tuintptr was not found")
}
if x.(uintptr) != 3 {
	t.Error("tuintptr is not 3:", x)
}
```
### TestIncrementUint8
```go

tc := New(0, 0)
tc.Set("tuint8", uint8(1), 0)
_, err := tc.Increment("tuint8", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}

x, found := tc.Get("tuint8")
if !found {
	t.Error("tuint8 was not found")
}
if x.(uint8) != 3 {
	t.Error("tuint8 is not 3:", x)
}
```
### TestIncrementUint16
```go

tc := New(0, 0)
tc.Set("tuint16", uint16(1), 0)
_, err := tc.Increment("tuint16", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}

x, found := tc.Get("tuint16")
if !found {
	t.Error("tuint16 was not found")
}
if x.(uint16) != 3 {
	t.Error("tuint16 is not 3:", x)
}
```
### TestIncrementUint32
```go

tc := New(0, 0)
tc.Set("tuint32", uint32(1), 0)
_, err := tc.Increment("tuint32", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}

x, found := tc.Get("tuint32")
if !found {
	t.Error("tuint32 was not found")
}
if x.(uint32) != 3 {
	t.Error("tuint32 is not 3:", x)
}
```
### TestIncrementUint64
```go

tc := New(0, 0)
tc.Set("tuint64", uint64(1), 0)
_, err := tc.Increment("tuint64", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}

x, found := tc.Get("tuint64")
if !found {
	t.Error("tuint64 was not found")
}
if x.(uint64) != 3 {
	t.Error("tuint64 is not 3:", x)
}
```
### TestIncrementInt
```go

tc := New(0, 0)
tc.Set("tint", 1, 0)
_, err := tc.Increment("tint", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}
x, found := tc.Get("tint")
if !found {
	t.Error("tint was not found")
}
if x.(int) != 3 {
	t.Error("tint is not 3:", x)
}
```
### TestIncrementInt8
```go

tc := New(0, 0)
tc.Set("tint8", int8(1), 0)
_, err := tc.Increment("tint8", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}
x, found := tc.Get("tint8")
if !found {
	t.Error("tint8 was not found")
}
if x.(int8) != 3 {
	t.Error("tint8 is not 3:", x)
}
```
### TestIncrementInt16
```go

tc := New(0, 0)
tc.Set("tint16", int16(1), 0)
_, err := tc.Increment("tint16", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}
x, found := tc.Get("tint16")
if !found {
	t.Error("tint16 was not found")
}
if x.(int16) != 3 {
	t.Error("tint16 is not 3:", x)
}
```
### TestIncrementInt32
```go

tc := New(0, 0)
tc.Set("tint32", int32(1), 0)
_, err := tc.Increment("tint32", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}
x, found := tc.Get("tint32")
if !found {
	t.Error("tint32 was not found")
}
if x.(int32) != 3 {
	t.Error("tint32 is not 3:", x)
}
```
### TestIncrementInt64
```go

tc := New(0, 0)
tc.Set("tint64", int64(1), 0)
_, err := tc.Increment("tint64", 2)
if err != nil {
	t.Error("Error incrementing:", err)
}
x, found := tc.Get("tint64")
if !found {
	t.Error("tint64 was not found")
}
if x.(int64) != 3 {
	t.Error("tint64 is not 3:", x)
}
```
### TestDecrementInt64
```go

tc := New(0, 0)
tc.Set("int64", int64(5), 0)
_, err := tc.Decrement("int64", 2)
if err != nil {
	t.Error("Error decrementing:", err)
}
x, found := tc.Get("int64")
if !found {
	t.Error("int64 was not found")
}
if x.(int64) != 3 {
	t.Error("int64 is not 3:", x)
}
```
### TestAdd
```go

tc := New(0, 0)
err := tc.Add("foo", "bar", 0)
if err != nil {
	t.Error("Couldn't add foo even though it shouldn't exist")
}
err = tc.Add("foo", "baz", 0)
if err == nil {
	t.Error("Successfully added another foo when it should have returned an error")
}
```
### TestReplace
```go

tc := New(0, 0)
err := tc.Replace("foo", "bar", 0)
if err == nil {
	t.Error("Replaced foo when it shouldn't exist")
}
tc.Set("foo", "bar", 0)
err = tc.Replace("foo", "bar", 0)
if err != nil {
	t.Error("Couldn't replace existing key foo")
}
```
### TestDelete
```go

tc := New(0, 0)
tc.Set("foo", "bar", 0)
tc.Delete("foo")
x, found := tc.Get("foo")
if found {
	t.Error("foo was found, but it should have been deleted")
}
if x != nil {
	t.Error("x is not nil:", x)
}
```
### TestFlush
```go

tc := New(0, 0)
tc.Set("foo", "bar", 0)
tc.Set("baz", "yes", 0)
tc.Flush()
x, found := tc.Get("foo")
if found {
	t.Error("foo was found, but it should have been deleted")
}
if x != nil {
	t.Error("x is not nil:", x)
}
x, found = tc.Get("baz")
if found {
	t.Error("baz was found, but it should have been deleted")
}
if x != nil {
	t.Error("x is not nil:", x)
}
```
### TestIncrementOverflowInt
```go

tc := New(0, 0)
tc.Set("int8", int8(127), 0)
_, err := tc.Increment("int8", 1)
if err != nil {
	t.Error("Error incrementing int8:", err)
}
x, _ := tc.Get("int8")
int8 := x.(int8)
if int8 != -128 {
	t.Error("int8 did not overflow as expected; value:", int8)
}

```
### TestIncrementOverflowUint
```go

tc := New(0, 0)
tc.Set("uint8", uint8(255), 0)
_, err := tc.Increment("uint8", 1)
if err != nil {
	t.Error("Error incrementing int8:", err)
}
x, _ := tc.Get("uint8")
uint8 := x.(uint8)
if uint8 != 0 {
	t.Error("uint8 did not overflow as expected; value:", uint8)
}
```
### TestDecrementUnderflowUint
```go

tc := New(0, 0)
tc.Set("uint8", uint8(0), 0)
_, err := tc.Decrement("uint8", 1)
if err != nil {
	t.Error("Error decrementing int8:", err)
}
x, _ := tc.Get("uint8")
uint8 := x.(uint8)
if uint8 != 0 {
	t.Error("uint8 was not capped at 0; value:", uint8)
}
```
### TestCacheSerialization
```go

tc := New(0, 0)
testFillAndSerialize(t, tc)

// Check if gob.Register behaves properly even after multiple gob.Register
// on c.Items (many of which will be the same type)
testFillAndSerialize(t, tc)
```
### TestFileSerialization
```go

tc := New(0, 0)
tc.Add("a", "a", 0)
tc.Add("b", "b", 0)
f, err := ioutil.TempFile("", "go-cache-cache.dat")
if err != nil {
	t.Fatal("Couldn't create cache file:", err)
}
fname := f.Name()
f.Close()
tc.SaveFile(fname)

oc := New(0, 0)
oc.Add("a", "aa", 0) // this should not be overwritten
err = oc.LoadFile(fname)
if err != nil {
	t.Error(err)
}
a, found := oc.Get("a")
if !found {
	t.Error("a was not found")
}
astr := a.(string)
if astr != "aa" {
	if astr == "a" {
		t.Error("a was overwritten")
	} else {
		t.Error("a is not aa")
	}
}
b, found := oc.Get("b")
if !found {
	t.Error("b was not found")
}
if b.(string) != "b" {
	t.Error("b is not b")
}
```
### TestSerializeUnserializable
```go

tc := New(0, 0)
ch := make(chan bool, 1)
ch <- true
tc.Set("chan", ch, 0)
fp := &bytes.Buffer{}
err := tc.Save(fp) // this should fail gracefully
if err.Error() != "gob NewTypeObject can't handle type: chan bool" {
	t.Error("Error from Save was not gob NewTypeObject can't handle type chan bool:", err)
}
```
## 内存缓存
### TestMemCacheFunc
```go

ctx := context.Background()

const cacheKey = "UT:%v:%v"

var lCache = New(5*time.Minute, 5*time.Minute)
result, err := MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc0, cacheKey, "p1", "p2")
log.Info(ctx, "%v %v %v", result, err, printKind(result))

_memCacheFuncTestMore(ctx, lCache, cacheKey)
```
## Redis缓存
### TestRdscCacheFunc
```go

if !isRdsRun() {
	return
}

redisDao.InitRedises()
ctx := context.Background()

const cacheKey = "UT:%v:%v"
const RDSC_DB = "rdscdb"

rds := redisDao.Get(RDSC_DB)

result, err := RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

log.Info(ctx, "%v %v %v", result, err, printKind(result))
```
### TestRdsDeleteCacheTestMore
```go

if !isRdsRun() {
	return
}

redisDao.InitRedises()
ctx := context.Background()

const cacheKey = "UT:%v:%v"
const RDSC_DB = "rdscdb"

rds := redisDao.Get(RDSC_DB)

var result interface{}
var err error

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc1, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc1, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc2, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc2, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc3, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc3, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc6, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc6, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", drainToArray(result), err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc7, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc7, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", drainToMap(result), err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc8, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc8, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc9, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))

result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc9, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result, err, printKind(result))
err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}

//result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc10, cacheKey, "p1", "p2")
//log.Info(ctx, "%v %v %v", result, err, printKind(result))
//
//result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc10, cacheKey, "p1", "p2")
//log.Info(ctx, "%v %v %v", result, err, printKind(result))

err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
if err != nil {
	t.Error(err)
}
```
### TestRdsCacheMultiFunc
```go

if !isRdsRun() {
	return
}

redisDao.InitRedises()
ctx := context.Background()
const RDSC_DB = "rdscdb"

rds := redisDao.Get(RDSC_DB)
result, err := RdsCacheMultiFunc(ctx, rds, 30, getThingsByIds, "multikey:%s", []string{"1", "2", "5", "3", "4", "10"})
if err == nil && result != nil {
	mapValue, ok := result.(map[string]*Thing)
	if ok {
		for key, value := range mapValue {
			log.Info(ctx, "%v===%v", key, value)
		}
	}
}
```
