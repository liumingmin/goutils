**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [utils](#utils)
  * [cbk](#cbk)
  * [checksum](#checksum)
  * [context_test.go](#context_testgo)
  * [csv](#csv)
  * [distlock](#distlock)
  * [dll_mod_test.go](#dll_mod_testgo)
  * [docgen](#docgen)
  * [fsm](#fsm)
  * [hc](#hc)
  * [ismtp](#ismtp)
  * [safego](#safego)
  * [snowflake](#snowflake)
  * [tags_test.go](#tags_testgo)

<!-- tocstop -->

# utils
## cbk
### cbk_test.go
#### TestCbkFailed
```go

InitCbk()

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
```
## checksum
### crc32_test.go
#### TestCompareChecksumFiles
```go

checkSumPath, err := GenerateChecksumFile(context.Background(), testChecksumDirPath, testChecksumName)
if err != nil {
	t.Error(err)
	return
}
t.Log(checkSumPath)

checksumMd5Path, err := GenerateChecksumMd5File(context.Background(), checkSumPath)
if err != nil {
	t.Error(err)
	return
}

valid := IsChecksumFileValid(context.Background(), checkSumPath, checksumMd5Path)
if !valid {
	t.Error(valid)
	return
}

err = CompareChecksumFiles(context.Background(), testChecksumDirPath, checkSumPath)
if err != nil {
	t.Error(err)
	return
}
```
### temp
## context_test.go
### TestContextWithTsTrace
```go

t.Log(ContextWithTrace())

t.Log(time.Now().UnixNano())
time.Sleep(time.Second)
t.Log(time.Now().UnixNano())
t.Log(NanoTsBase32())
t.Log(ContextWithTsTrace())
t.Log(ContextWithTsTrace())
t.Log(ContextWithTsTrace())
```
## csv
### csv_parse_test.go
#### TestReadCsvToDataTable
```go

dt, err := ReadCsvToDataTable(context.Background(), `goutils.log`, '\t',
	[]string{"xx", "xx", "xx", "xx"}, "xxx", []string{"xxx"})
if err != nil {
	t.Log(err)
	return
}
for _, r := range dt.Rows() {
	t.Log(r.Data())
}

rs := dt.RowsBy("xxx", "869")
for _, r := range rs {
	t.Log(r.Data())
}

t.Log(dt.Row("17"))
```
## distlock
### consullock_test.go
#### TestAquireConsulLock
```go

l, _ := NewConsulLock("accountId", 10)
//l.Lock(15)
//l.Unlock()
ctx := context.Background()
fmt.Println("try lock 1")

fmt.Println(l.Lock(ctx, 5))
//time.Sleep(time.Second * 6)

//fmt.Println("try lock 2")
//fmt.Println(l.Lock(3))

l2, _ := NewConsulLock("accountId", 10)
fmt.Println("try lock 3")
fmt.Println(l2.Lock(ctx, 15))

l3, _ := NewConsulLock("accountId", 10)
fmt.Println("try lock 4")
fmt.Println(l3.Lock(ctx, 15))

time.Sleep(time.Minute)
```
### filelock_test.go
#### TestFileLock
```go

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

```
### rdslock_test.go
#### TestRdsLock
```go

redis.InitRedises()
l, err := NewRdsLuaLock("rdscdb", "accoutId", 4)
if err != nil {
	t.Error(err)
}

l2, err := NewRdsLuaLock("rdscdb", "accoutId", 4)
if err != nil {
	t.Error(err)
}

ctx := context.Background()

if !l.Lock(ctx, 1) {
	t.Error("can not get lock")
}

time.Sleep(time.Millisecond * 300)
if l2.Lock(ctx, 1) {
	t.Error("except get lock")
}
l.Unlock(ctx)

time.Sleep(time.Millisecond * 100)

if !l2.Lock(ctx, 1) {
	t.Error("can not get lock")
}
```
## dll_mod_test.go
### TestDllCall
```go

// mod := NewDllMod("machineinfo.dll")

// result := int32(0)

// retCode, err := mod.Call("GetDiskType", "C:", &result)
// if err != nil {
// 	t.Fatal(err)
// }

// if retCode != 0 {
// 	t.FailNow()
// }

// if result != 4 {
// 	t.FailNow()
// }
```
### TestDllConvertString
```go

mod := NewDllMod("test.dll")

testStr := "abcde很棒"
var arg uintptr
var err error
arg, err = mod.convertArg(testStr)
if err != nil {
	t.FailNow()
}

var slice []byte
header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
header.Data = arg
header.Len = len(testStr)
header.Cap = header.Len

if string(slice) != testStr {
	t.FailNow()
}
```
### TestDllConvertInt
```go

mod := NewDllMod("test.dll")

var arg uintptr
var err error
arg, err = mod.convertArg(12345)
if err != nil {
	t.FailNow()
}

if arg != 12345 {
	t.FailNow()
}

intptr := int(1080)
arg, err = mod.convertArg(&intptr)
if err != nil {
	t.FailNow()
}

if *(*int)(unsafe.Pointer(arg)) != intptr {
	t.FailNow()
}

uintptr1 := uintptr(11080)
arg, err = mod.convertArg(&uintptr1)
if err != nil {
	t.FailNow()
}

if *(*uintptr)(unsafe.Pointer(arg)) != uintptr1 {
	t.FailNow()
}
```
### TestDllConvertBool
```go

mod := NewDllMod("test.dll")

var arg uintptr
var err error
arg, err = mod.convertArg(true)
if err != nil {
	t.FailNow()
}

if arg != 1 {
	t.FailNow()
}
```
### TestDllConvertSlice
```go

mod := NewDllMod("test.dll")

origSlice := []byte("testslicecvt")

var arg uintptr
var err error
arg, err = mod.convertArg(origSlice)
if err != nil {
	t.FailNow()
}

var slice []byte
header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
header.Data = arg
header.Len = len(origSlice)
header.Cap = header.Len

if bytes.Compare(origSlice, slice) != 0 {
	t.FailNow()
}
```
### TestDllConvertStructPtr
```go

mod := NewDllMod("test.dll")

s := testDllModStruct{100, 200, 300}

var arg uintptr
var err error
arg, err = mod.convertArg(&s)
if err != nil {
	t.FailNow()
}

s2 := *(*testDllModStruct)(unsafe.Pointer(arg))
if s2.x1 != s.x1 || s2.x2 != s.x2 || s2.x4 != s.x4 {
	t.FailNow()
}
```
### TestGetCStrFromUintptr
```go

mod := NewDllMod("test.dll")

testStr := "abcde很棒"
var arg uintptr
var err error
arg, err = mod.convertArg(testStr)
if err != nil {
	t.FailNow()
}

origStr := mod.GetCStrFromUintptr(arg)

if testStr != origStr {
	t.FailNow()
}
```
### TestDllConvertFunc
```go

//cannot convert back
// mod := NewDllMod("test.dll")

// var testCallback = func(s uintptr) uintptr {
// 	fmt.Println("test callback")
// 	return s + 900000
// }

// var arg uintptr
// var err error
// arg, err = mod.convertArg(testCallback)
// if err != nil {
// 	t.FailNow()
// }

// callback := *(*(func(s uintptr) uintptr))(unsafe.Pointer(arg))

// t.Log(callback(12345))
```
## docgen
### cmd
### doc
### docgen_test.go
#### TestGenDocTestUser
```go

sb := strings.Builder{}
sb.WriteString(genDocTestUserQuery())
sb.WriteString(genDocTestUserCreate())
sb.WriteString(genDocTestUserUpdate())
sb.WriteString(genDocTestUserDelete())

GenDoc(context.Background(), "用户管理", "doc/testuser.md", 2, sb.String())
```
## fsm
## hc
## ismtp
### ismtp_test.go
#### TestSendEmail
```go

emailauth := LoginAuth(
	"from",
	"xxxxxx",
	"mailhost.com",
)

ctype := fmt.Sprintf("Content-Type: %s; charset=%s", "text/plain", "utf-8")

msg := fmt.Sprintf("To: %s\r\nCc: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
	strings.Join([]string{"target@mailhost.com"}, ";"),
	"",
	"from@mailhost.com",
	"测试",
	ctype,
	"测试")

err := SendMail("mailhost.com:port", //convert port number from int to string
	emailauth,
	"from@mailhost.com",
	[]string{"target@mailhost.com"},
	[]byte(msg),
)

if err != nil {
	t.Log(err)
	return
}

return
```
## safego
## snowflake
### snowflake_test.go
#### TestSnowflake
```go

n, _ := NewNode(1)
t.Log(n.Generate(), ",", n.Generate(), ",", n.Generate())
```
## tags_test.go
### TestAutoGenTags
```go

fmt.Println(AutoGenTags(testUser{}, map[string]TAG_STYLE{
	"json": TAG_STYLE_SNAKE,
	"bson": TAG_STYLE_UNDERLINE,
	"form": TAG_STYLE_ORIG,
}))
```
