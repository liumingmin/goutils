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
