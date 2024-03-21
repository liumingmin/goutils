**其他语言版本: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [通用工具库](#%E9%80%9A%E7%94%A8%E5%B7%A5%E5%85%B7%E5%BA%93)
  * [async_test.go](#async_testgo)
  * [熔断器](#%E7%86%94%E6%96%AD%E5%99%A8)
  * [checksum](#checksum)
  * [CSV文件解析为MDB内存表](#csv%E6%96%87%E4%BB%B6%E8%A7%A3%E6%9E%90%E4%B8%BAmdb%E5%86%85%E5%AD%98%E8%A1%A8)
  * [分布式锁](#%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81)
  * [dll_mod_test.go](#dll_mod_testgo)
  * [文档自动生成](#%E6%96%87%E6%A1%A3%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90)
  * [encoding_test.go](#encoding_testgo)
  * [有限状态机](#%E6%9C%89%E9%99%90%E7%8A%B6%E6%80%81%E6%9C%BA)
  * [httpclient工具](#httpclient%E5%B7%A5%E5%85%B7)
  * [邮件工具](#%E9%82%AE%E4%BB%B6%E5%B7%A5%E5%85%B7)
  * [安全的go协程](#%E5%AE%89%E5%85%A8%E7%9A%84go%E5%8D%8F%E7%A8%8B)
  * [snowflake](#snowflake)
  * [stringutils_test.go](#stringutils_testgo)
  * [结构体TAG生成器](#%E7%BB%93%E6%9E%84%E4%BD%93tag%E7%94%9F%E6%88%90%E5%99%A8)

<!-- tocstop -->

# 通用工具库
## async_test.go
### TestAsyncInvokeWithTimeout
```go

f1 := false
f2 := false
result := AsyncInvokeWithTimeout(time.Second*1, func() {
	time.Sleep(time.Millisecond * 500)
	f1 = true
}, func() {
	time.Sleep(time.Millisecond * 500)
	f2 = true
})

if !result {
	t.FailNow()
}

if !f1 {
	t.FailNow()
}

if !f2 {
	t.FailNow()
}
```
### TestAsyncInvokeWithTimeouted
```go

f1 := false
f2 := false
result := AsyncInvokeWithTimeout(time.Second*1, func() {
	time.Sleep(time.Millisecond * 1500)
	f1 = true
}, func() {
	time.Sleep(time.Millisecond * 500)
	f2 = true
})

if result {
	t.FailNow()
}

if f1 {
	t.FailNow()
}

if !f2 {
	t.FailNow()
}
```
### TestAsyncInvokesWithTimeout
```go

f1 := false
f2 := false

fns := []func(){
	func() {
		time.Sleep(time.Millisecond * 500)
		f1 = true
	}, func() {
		time.Sleep(time.Millisecond * 500)
		f2 = true
	},
}
result := AsyncInvokesWithTimeout(time.Second*1, fns)

if !result {
	t.FailNow()
}

if !f1 {
	t.FailNow()
}

if !f2 {
	t.FailNow()
}
```
## 熔断器
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
#### goutils
##### 0
##### 1
##### 2
##### 3
##### 4
## CSV文件解析为MDB内存表
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
## 分布式锁
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
## 文档自动生成
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
## encoding_test.go
### TestGBK2UTF8
```go

src := []byte{206, 210, 202, 199, 103, 111, 117, 116, 105, 108, 115, 49}
utf8str, err := GBK2UTF8(src)
if err != nil {
	t.FailNow()
}

if string(utf8str) != "我是goutils1" {
	t.FailNow()
}
```
### TestUTF82GBK
```go

src := []byte{230, 136, 145, 230, 152, 175, 103, 111, 117, 116, 105, 108, 115, 49}
gbkStr, err := UTF82GBK(src)
if err != nil {
	t.FailNow()
}

if !reflect.DeepEqual(gbkStr, []byte{206, 210, 202, 199, 103, 111, 117, 116, 105, 108, 115, 49}) {
	t.FailNow()
}
```
### TestIsGBK
```go

if !IsGBK([]byte{206, 210}) {
	t.FailNow()
}
```
### TestIsUtf8
```go

if !IsUtf8([]byte{230, 136, 145}) {
	t.FailNow()
}
```
## 有限状态机
## httpclient工具
## 邮件工具
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
## 安全的go协程
## snowflake
### 雪花ID生成器
#### TestSnowflake
```go

n, _ := NewNode(1)
t.Log(n.Generate(), ",", n.Generate(), ",", n.Generate())
```
## stringutils_test.go
### TestStringsReverse
```go

var strs = []string{"1", "2", "3", "4"}
revStrs := StringsReverse(strs)

if !reflect.DeepEqual(revStrs, []string{"4", "3", "2", "1"}) {
	t.FailNow()
}
```
### TestStringsInArray
```go

var strs = []string{"1", "2", "3", "4"}
ok, index := StringsInArray(strs, "3")
if !ok {
	t.FailNow()
}

if index != 2 {
	t.FailNow()
}

ok, index = StringsInArray(strs, "5")
if ok {
	t.FailNow()
}

if index != -1 {
	t.FailNow()
}
```
### TestStringsExcept
```go

var strs1 = []string{"1", "2", "3", "4"}
var strs2 = []string{"3", "4", "5", "6"}

if !reflect.DeepEqual(StringsExcept(strs1, strs2), []string{"1", "2"}) {
	t.FailNow()
}

if !reflect.DeepEqual(StringsExcept(strs1, []string{}), []string{"1", "2", "3", "4"}) {
	t.FailNow()
}

if !reflect.DeepEqual(StringsExcept([]string{}, strs2), []string{}) {
	t.FailNow()
}
```
### TestStringsDistinct
```go

var strs1 = []string{"1", "2", "3", "4", "1", "3"}
distincted := StringsDistinct(strs1)
sort.Strings(distincted)
if !reflect.DeepEqual(distincted, []string{"1", "2", "3", "4"}) {
	t.FailNow()
}
```
## 结构体TAG生成器
### TestAutoGenTags
```go

fmt.Println(AutoGenTags(testUser{}, map[string]TAG_STYLE{
	"json": TAG_STYLE_SNAKE,
	"bson": TAG_STYLE_UNDERLINE,
	"form": TAG_STYLE_ORIG,
}))
```
