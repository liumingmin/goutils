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
  * [file_test.go](#file_testgo)
  * [有限状态机](#%E6%9C%89%E9%99%90%E7%8A%B6%E6%80%81%E6%9C%BA)
  * [httpclient工具](#httpclient%E5%B7%A5%E5%85%B7)
  * [邮件工具](#%E9%82%AE%E4%BB%B6%E5%B7%A5%E5%85%B7)
  * [math_test.go](#math_testgo)
  * [reflectutils_test.go](#reflectutils_testgo)
  * [安全的go协程](#%E5%AE%89%E5%85%A8%E7%9A%84go%E5%8D%8F%E7%A8%8B)
  * [snowflake](#snowflake)
  * [stringutils_test.go](#stringutils_testgo)
  * [struct_test.go](#struct_testgo)
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

checkSumPath, err := GenerateChecksumFile(context.Background(), testTempDirPath, testChecksumName)
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

err = CompareChecksumFiles(context.Background(), testTempDirPath, checkSumPath)
if err != nil {
	t.Error(err)
	return
}
```
## CSV文件解析为MDB内存表
### csv_parse_test.go
#### TestReadCsvToDataTable
```go

dt, err := ReadCsvToDataTable(context.Background(), filepath.Join(testTempDirPath, testCsvFilePath), ',',
	[]string{"id", "name", "age", "remark"}, "id", []string{"name"})
if err != nil {
	t.Error(err)
}

if !reflect.DeepEqual(dt.Row("10").Data(), []string{"10", "name10", "10", "remark10"}) {
	t.FailNow()
}

if !reflect.DeepEqual(dt.RowsBy("name", "name10")[0].Data(), []string{"10", "name10", "10", "remark10"}) {
	t.FailNow()
}
```
#### TestParseCsvRaw
```go

records := ParseCsvRaw(context.Background(),
	`id	name	age	remark
0	name0	0	remark0
1	name1	1	remark1
2	name2	2	remark2
3	name3	3	remark3
4	name4	4	remark4
5	name5	5	remark5
6	name6	6	remark6
7	name7	7	remark7
8	name8	8	remark8
9	name9	9	remark9
10	name10	10	remark10
11	name11	11	remark11
12	name12	12	remark12
13	name13	13	remark13
14	name14	14	remark14
15	name15	15	remark15
16	name16	16	remark16
17	name17	17	remark17
18	name18	18	remark18
19	name19	19	remark19`)

dt := container.NewDataTable([]string{"id", "name", "age", "remark"}, "id", []string{"name"}, 20)
dt.PushAll(records)

if !reflect.DeepEqual(dt.Row("10").Data(), []string{"10", "name10", "10", "remark10"}) {
	t.FailNow()
}

if !reflect.DeepEqual(dt.RowsBy("name", "name10")[0].Data(), []string{"10", "name10", "10", "remark10"}) {
	t.FailNow()
}
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

mod := NewDllMod("ntdll.dll")

info := &struct {
	osVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformId        uint32
	CsdVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	_                 byte
}{}

info.osVersionInfoSize = uint32(unsafe.Sizeof(*info))
retCode, err := mod.Call("RtlGetVersion", uintptr(unsafe.Pointer(info)))
if err != nil {
	t.Error(err)
}

if retCode != 0 {
	t.Error(retCode)
}

if info.MajorVersion == 0 {
	t.Error(info.MajorVersion)
}

retCode, err = mod.Call("RtlGetVersion", uintptr(unsafe.Pointer(info)))
if err != nil {
	t.Error(err)
}
if err != nil {
	t.Error(err)
}

if retCode != 0 {
	t.Error(retCode)
}

if info.MajorVersion == 0 {
	t.Error(info.MajorVersion)
}
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

testDllConvertNum(t, mod, int(-1080))
testDllConvertNum(t, mod, uint(1080))
testDllConvertNum(t, mod, int8(-128))
testDllConvertNum(t, mod, uint8(255))
testDllConvertNum(t, mod, int16(-30000))
testDllConvertNum(t, mod, uint16(30000))
testDllConvertNum(t, mod, int32(-3000000))
testDllConvertNum(t, mod, uint32(3000000))
testDllConvertNum(t, mod, int64(-3000000))
testDllConvertNum(t, mod, uint64(3000000))
testDllConvertNum(t, mod, uintptr(11080))

testData := 123
up := unsafe.Pointer(&testData)
testDllConvertNum(t, mod, up)

testDllConvertNumPtr(t, mod, int(-1080))
testDllConvertNumPtr(t, mod, uint(1080))
testDllConvertNumPtr(t, mod, int8(-128))
testDllConvertNumPtr(t, mod, uint8(255))
testDllConvertNumPtr(t, mod, int16(-30000))
testDllConvertNumPtr(t, mod, uint16(30000))
testDllConvertNumPtr(t, mod, int32(-3000000))
testDllConvertNumPtr(t, mod, uint32(3000000))
testDllConvertNumPtr(t, mod, int64(-3000000))
testDllConvertNumPtr(t, mod, uint64(3000000))
testDllConvertNumPtr(t, mod, uintptr(11080))

testDllConvertNumPtr(t, mod, float32(100.12))
testDllConvertNumPtr(t, mod, float64(100.12))
testDllConvertNumPtr(t, mod, complex64(100.12))
testDllConvertNumPtr(t, mod, complex128(100.12))
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
### TestDllConvertUnsupport
```go

mod := NewDllMod("test.dll")

_, err := mod.convertArg(float32(11.12))
if err != ErrUnsupportArg {
	t.Error(err)
}

_, err = mod.convertArg(float64(11.12))
if err != ErrUnsupportArg {
	t.Error(err)
}

_, err = mod.convertArg(complex64(11.12))
if err != ErrUnsupportArg {
	t.Error(err)
}

_, err = mod.convertArg(complex128(11.12))
if err != ErrUnsupportArg {
	t.Error(err)
}

m := make(map[string]string)
_, err = mod.convertArg(m)
if err != ErrUnsupportArg {
	t.Error(err)
}

c := make(chan struct{})
_, err = mod.convertArg(c)
if err != ErrUnsupportArg {
	t.Error(err)
}

s := struct{}{}
_, err = mod.convertArg(s)
if err != ErrUnsupportArg {
	t.Error(err)
}

_, err = mod.convertArg(interface{}(s))
if err != ErrUnsupportArg {
	t.Error(err)
}
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
## file_test.go
### TestGetCurrPath
```go

path := GetCurrPath()
t.Log(path)
```
### TestFileExist
```go

runFile := os.Args[0]

if !FileExist(runFile) {
	t.Error(runFile)
}
```
### TestFileExt
```go

if FileExt("aaa.txt") != ".txt" {
	t.Error(FileExt("aaa.txt"))
}

if FileExt("aaa.txt.zip") != ".zip" {
	t.Error(FileExt("aaa.txt.zip"))
}

if FileExt("aaa.txt.") != "." {
	t.Error(FileExt("aaa.txt."))
}

if FileExt("aaa") != "" {
	t.Error(FileExt("aaa"))
}
```
### TestFileCopy
```go

runFile := os.Args[0]
err := FileCopy(runFile, filepath.Join(testTempDirPath, "test_file"))
if err != nil {
	t.Error()
}

err = FileCopy(filepath.Join(testTempDirPath, "test_file"), filepath.Join(testTempDirPath, "test_file"))
if err != nil {
	t.Error()
}

err = FileCopy(filepath.Join(testTempDirPath, "test_file"), ".")
if err == nil {
	t.Error()
}
```
### TestIsPathTravOut
```go

if IsPathTravOut(`C:\a\b`, `C:\a`) {
	t.FailNow()
}

if IsPathTravOut(`C:\A\B`, `C:\a`) {
	t.FailNow()
}

if IsPathTravOut(`C:\a\b\..`, `C:\a`) {
	t.FailNow()
}

if !IsPathTravOut(`C:\a\b\..\..`, `C:\a`) {
	t.FailNow()
}

if !IsPathTravOut(`C:\A\B\..\..`, `C:\a`) {
	t.FailNow()
}

if !IsPathTravOut(`C:\a\b`, `C:\c`) {
	t.FailNow()
}
```
### TestUniformPathStyle
```go

if UniformPathStyle(`C:\a\b`) != `C:/a/b` {
	t.FailNow()
}

if UniformPathStyleCase(`C:\A\B`) != `c:/a/b` {
	t.FailNow()
}

if !reflect.DeepEqual(UniformPathListStyleCase([]string{`C:\A\B`}), []string{`c:/a/b`}) {
	t.FailNow()
}
```
### TestIsSameFilePath
```go

if !IsSameFilePath(`C:\a\b`, `C:/a/b`) {
	t.FailNow()
}

if !IsSameFilePath(`C:\a\..\a\b`, `C:/a/b`) {
	t.FailNow()
}

if IsSameFilePath(`C:\a\..\a\b\c`, `C:/a/b`) {
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
## math_test.go
### TestMathMin
```go

if Min(10, 9) != 9 {
	t.FailNow()
}

if Min(-1, -2) != -2 {
	t.FailNow()
}

if Min(3.1, 4.02) != 3.1 {
	t.FailNow()
}
```
### TestMathMax
```go

if Max(10, 9) != 10 {
	t.FailNow()
}

if Max(-1, -2) != -1 {
	t.FailNow()
}

if Max(3.1, 4.02) != 4.02 {
	t.FailNow()
}
```
### TestMathAbs
```go

if Abs(-1) != 1 {
	t.FailNow()
}

if Abs(1) != 1 {
	t.FailNow()
}
```
## reflectutils_test.go
### TestAnyIndirect
```go

val := reflect.ValueOf(10)
if AnyIndirect(val) != val {
	t.Error(val)
}

x := 10
val2 := reflect.ValueOf(&x)
if AnyIndirect(val2) == val2 {
	t.Error(val2)
}

if AnyIndirect(val2).Int() != int64(x) {
	t.Error(val2)
}
```
### TestIsNil
```go

var m map[string]string
if !IsNil(m) {
	t.Error(m)
}

var c chan string
if !IsNil(c) {
	t.Error(c)
}

var fun func()
if !IsNil(fun) {
	t.Error("func not nil")
}

var s []string
if !IsNil(s) {
	t.Error(s)
}

var sp *string
if !IsNil(sp) {
	t.Error(sp)
}

var up unsafe.Pointer
if !IsNil(up) {
	t.Error(up)
}

testIsNil[map[string]string](t)
testIsNil[chan string](t)
testIsNil[func()](t)
testIsNil[[]string](t)
testIsNil[*string](t)
testIsNil[unsafe.Pointer](t)
```
### testIsNil[T any]
```go

value := testWrapperNil[T]()

if value == nil {
	t.Error(value)
}

if !IsNil(value) {
	t.Error(value)
}
```
## 安全的go协程
## snowflake
### 雪花ID生成器
#### TestSnowflake
```go

n, err := NewNode(-1)
if err == nil {
	t.FailNow()
}

n, err = NewNode(1024)
if err == nil {
	t.FailNow()
}

n, err = NewNode(2)
if err != nil {
	t.Fatal(err)
}

id1 := n.Generate()
id2 := n.Generate()
if id1 == id2 {
	t.FailNow()
}

if ParseInt64(id1.Int64()) != id1 {
	t.FailNow()
}

idtemp, err := ParseString(id1.String())
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp, err = ParseBase2(id1.Base2())
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp, err = ParseBase32([]byte(id1.Base32()))
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp, err = ParseBase36(id1.Base36())
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp, err = ParseBase58([]byte(id1.Base58()))
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp, err = ParseBase64(id1.Base64())
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp, err = ParseBytes(id1.Bytes())
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}

idtemp = ParseIntBytes(id1.IntBytes())
if idtemp != id1 {
	t.FailNow()
}

bs, err := id1.MarshalJSON()
if err != nil {
	t.Fatal(err)
}

idtemp = ID(0)
err = idtemp.UnmarshalJSON(bs)
if err != nil {
	t.Fatal(err)
}
if idtemp != id1 {
	t.FailNow()
}
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
## struct_test.go
### TestCopyStruct
```go

type SrcFoo struct {
	A                int
	B                []*string
	C                map[string]*int
	SrcUnique        string
	SameNameDiffType time.Time
}
type DstFoo struct {
	A                int
	B                []*string
	C                map[string]*int
	DstUnique        int
	SameNameDiffType string
}

// Create the initial value
str1 := "hello"
str2 := "bye bye"
int1 := 1
int2 := 2
f1 := &SrcFoo{
	A: 1,
	B: []*string{&str1, &str2},
	C: map[string]*int{
		"A": &int1,
		"B": &int2,
	},
	SrcUnique:        "unique",
	SameNameDiffType: time.Now(),
}
var f2 DstFoo

CopyStruct(f1, &f2, BaseConvert)

if !reflect.DeepEqual(f1.A, f2.A) {
	t.Error(f2)
}

if !reflect.DeepEqual(f1.B, f2.B) {
	t.Error(f2)
}

if !reflect.DeepEqual(f1.C, f2.C) {
	t.Error(f2)
}

if !reflect.DeepEqual(BaseConvert(f1.SameNameDiffType, reflect.TypeOf("")), f2.SameNameDiffType) {
	t.Error(f2)
}
```
### TestCopyStructs
```go

type SrcFoo struct {
	A                int
	B                []*string
	C                map[string]*int
	SrcUnique        string
	SameNameDiffType time.Time
}
type DstFoo struct {
	A                int
	B                []*string
	C                map[string]*int
	DstUnique        int
	SameNameDiffType string
}

// Create the initial value
str1 := "hello"
str2 := "bye bye"
int1 := 1
int2 := 2
f1 := []SrcFoo{{
	A: 1,
	B: []*string{&str1, &str2},
	C: map[string]*int{
		"A": &int1,
		"B": &int2,
	},
	SrcUnique:        "unique",
	SameNameDiffType: time.Now(),
}}
var f2 []DstFoo
CopyStructs(f1, &f2, BaseConvert)

if !reflect.DeepEqual(f1[0].A, f2[0].A) {
	t.Error(f2)
}

if !reflect.DeepEqual(f1[0].B, f2[0].B) {
	t.Error(f2)
}

if !reflect.DeepEqual(f1[0].C, f2[0].C) {
	t.Error(f2)
}

if !reflect.DeepEqual(BaseConvert(f1[0].SameNameDiffType, reflect.TypeOf("")), f2[0].SameNameDiffType) {
	t.Error(f2)
}
```
## 结构体TAG生成器
### TestAutoGenTags
```go

structStrWithTag := AutoGenTags(testUser{}, map[string]TAG_STYLE{
	"json":      TAG_STYLE_SNAKE,
	"bson":      TAG_STYLE_UNDERLINE,
	"form":      TAG_STYLE_ORIG,
	"nonestyle": TAG_STYLE_NONE,
})

if !strings.Contains(structStrWithTag, `bson:"user_id"`) {
	t.FailNow()
}

if !strings.Contains(structStrWithTag, `form:"UserId"`) {
	t.FailNow()
}

if !strings.Contains(structStrWithTag, `json:"userId"`) {
	t.FailNow()
}

if !strings.Contains(structStrWithTag, `json:"status"`) {
	t.FailNow()
}

if strings.Contains(structStrWithTag, `nonestyle:`) {
	t.FailNow()
}

//t.Log(structStrWithTag)
```
