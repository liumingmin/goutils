# net 网络库
## httpx 兼容http1.x和2.0的httpclient
### httpclientx_test.go
#### TestHttpXGet
```go

clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Get("http://127.0.0.1:8049")
	if err != nil {
		t.Fatal(fmt.Errorf("error making request: %v", err))
	}
	t.Log(resp.StatusCode)
	t.Log(resp.Proto)
}
```
#### TestHttpXPost
```go

clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Get("http://127.0.0.1:8881")
	if err != nil {
		t.Fatal(fmt.Errorf("error making request: %v", err))
	}
	t.Log(resp.StatusCode)
	t.Log(resp.Proto)
}
```
## ip
## packet tcp包model
## proxy ssh proxy
### ssh_client_test.go
#### TestSshClient
```go

client := getSshClient(t)
defer client.Close()

session, err := client.NewSession()
if err != nil {
	t.Fatalf("Create session failed %v", err)
}
defer session.Close()

// run command and capture stdout/stderr
output, err := session.CombinedOutput("ls -l /data")
if err != nil {
	t.Fatalf("CombinedOutput failed %v", err)
}
t.Log(string(output))
```
#### TestMysqlSshClient
```go

client := getSshClient(t)
defer client.Close()

//test时候，打开，会引入mysql包
//mysql.RegisterDialContext("tcp", func(ctx context.Context, addr string) (net.Conn, error) {
//	return client.Dial("tcp", addr)
//})

db, err := sql.Open("", "")
if err != nil {
	t.Fatalf("open db failed %v", err)
}
defer db.Close()

rs, err := db.Query("select  limit 10")
if err != nil {
	t.Fatalf("open db failed %v", err)
}
defer rs.Close()
for rs.Next() {

}
```
## serverx 兼容http1.x和2.0的http server
