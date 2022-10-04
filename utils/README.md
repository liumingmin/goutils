

<!-- toc -->

- [utils 通用工具库](#utils-%E9%80%9A%E7%94%A8%E5%B7%A5%E5%85%B7%E5%BA%93)
  * [cbk 熔断器](#cbk-%E7%86%94%E6%96%AD%E5%99%A8)
    + [cbk_test.go](#cbk_testgo)
  * [checksum](#checksum)
    + [crc32_test.go](#crc32_testgo)
    + [testdir](#testdir)
  * [csv CSV文件解析为MDB内存表](#csv-csv%E6%96%87%E4%BB%B6%E8%A7%A3%E6%9E%90%E4%B8%BAmdb%E5%86%85%E5%AD%98%E8%A1%A8)
    + [csv_parse_test.go](#csv_parse_testgo)
  * [distlock 分布式锁](#distlock-%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81)
    + [consullock_test.go](#consullock_testgo)
    + [filelock_test.go](#filelock_testgo)
    + [rdslock_test.go](#rdslock_testgo)
  * [docgen 文档自动生成](#docgen-%E6%96%87%E6%A1%A3%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90)
    + [cmd](#cmd)
    + [docgen_test.go](#docgen_testgo)
  * [fsm 有限状态机](#fsm-%E6%9C%89%E9%99%90%E7%8A%B6%E6%80%81%E6%9C%BA)
  * [hc httpclient工具](#hc-httpclient%E5%B7%A5%E5%85%B7)
  * [ismtp 邮件工具](#ismtp-%E9%82%AE%E4%BB%B6%E5%B7%A5%E5%85%B7)
    + [ismtp_test.go](#ismtp_testgo)
  * [safego 安全的go协程](#safego-%E5%AE%89%E5%85%A8%E7%9A%84go%E5%8D%8F%E7%A8%8B)
  * [snowflake](#snowflake)
    + [snowflake_test.go 雪花ID生成器](#snowflake_testgo-%E9%9B%AA%E8%8A%B1id%E7%94%9F%E6%88%90%E5%99%A8)
  * [tags_test.go 结构体TAG生成器](#tags_testgo-%E7%BB%93%E6%9E%84%E4%BD%93tag%E7%94%9F%E6%88%90%E5%99%A8)
    + [TestAutoGenTags](#testautogentags)

<!-- tocstop -->

# utils 通用工具库
## cbk 熔断器
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
#### TestGenerateCheckSumFile
```go

src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product"
checksumName := "nwjs"
checkSumPath, err := GenerateChecksumFile(context.Background(), src, checksumName)
if err != nil {
	t.Error(err)
	return
}
t.Log(checkSumPath)
```
#### TestGenerateChecksumMd5File
```go

src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum"
checksumMd5Path, err := GenerateChecksumMd5File(context.Background(), src)
if err != nil {
	t.Error(err)
	return
}
t.Log(checksumMd5Path)
```
#### TestGenerateChecksumFileWithIgnore
```go

dirMap := make(map[string][]string)
dirMap["jsex"] = []string{"E:\\game\\jsex\\base"}
for code, dirs := range dirMap {
	t.Log("game: ", code)
	for _, dir := range dirs {
		t.Log("dir: ", dir)
		start := time.Now() // 获取当前时间
		checksumName := fmt.Sprintf("%v", code)
		checksumMd5Path, err := GenerateChecksumFileWithIgnore(context.Background(), dir, checksumName, []string{fmt.Sprintf("%v.checksum", code), "pak", "locales\\pak"})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(checksumMd5Path)
		elapsed := time.Since(start)
		t.Log("time：", elapsed)
	}
}
```
#### TestIsChecksumFileValid
```go

src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum"
md5Path := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum.md5"
valid := IsChecksumFileValid(context.Background(), src, md5Path)
if !valid {
	t.Error(valid)
	return
}
t.Log(valid)
```
#### TestRelPath
```go

repos := []string{"", "a", "b", "a\\b", "a/c", "a\\b/c", "a/d/c", "d/a", "d/c"}

for _, repo1 := range repos {
	t.Log(">>>", repo1)
	for _, repo2 := range repos {
		rel, _ := filepath.Rel(repo1, repo2)
		if !strings.Contains(rel, ".") {
			t.Log(repo2, ":", rel)
		}
	}
}

```
#### TestCompareChecksumFiles
```go

src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum"
root := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product"
err := CompareChecksumFiles(context.Background(), root, src)
if err != nil {
	t.Error(err)
	return
}
```
### testdir
## csv CSV文件解析为MDB内存表
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
## distlock 分布式锁
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
l, _ := NewRdsLuaLock("rdscdb", "accoutId", 4)
l2, _ := NewRdsLuaLock("rdscdb", "accoutId", 4)
//l.Lock(15)
//l.Unlock()
ctx := context.Background()
fmt.Println(l.Lock(ctx, 5))
fmt.Println("1getlock")
fmt.Println(l2.Lock(ctx, 5))
fmt.Println("2getlock")
time.Sleep(time.Second * 15)

//l2, _ := NewRdsLuaLock("accoutId", 15)

//t.Log(l2.Lock(5))
```
## docgen 文档自动生成
### cmd
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
## fsm 有限状态机
## hc httpclient工具
## ismtp 邮件工具
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
## safego 安全的go协程
## snowflake
### snowflake_test.go 雪花ID生成器
#### TestSnowflake
```go

n, _ := NewNode(1)
t.Log(n.Generate(), ",", n.Generate(), ",", n.Generate())
```
## tags_test.go 结构体TAG生成器
### TestAutoGenTags
```go

fmt.Println(AutoGenTags(testUser{}, map[string]TAG_STYLE{
	"json": TAG_STYLE_SNAKE,
	"bson": TAG_STYLE_UNDERLINE,
	"form": TAG_STYLE_ORIG,
}))
```
