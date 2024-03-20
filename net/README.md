**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [net](#net)
  * [bwlimit](#bwlimit)
  * [httpdownloader](#httpdownloader)
  * [httpx](#httpx)
  * [ip](#ip)
  * [packet](#packet)
  * [proxy](#proxy)
  * [serverx](#serverx)

<!-- tocstop -->

# net
## bwlimit
## httpdownloader
### httpdownloader_test.go
#### TestHttpDownloaderDownload
```go

dialer := bwlimit.NewDialer()
dialer.RxBwLimit().SetBwLimit(20 * 1024 * 1024)
hc := &http.Client{
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

downloader := &HttpDownloader{
	HttpClient:   hc,
	ConBlockChan: make(chan struct{}, 10),
	BlockSize:    1024 * 1024,
	RetryCnt:     1,
}

err := downloader.Download(context.Background(), "https://golang.google.cn/dl/go1.21.7.windows-amd64.zip", http.Header{
	"User-Agent": []string{"Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Mobile Safari/537.36"},
}, "./go1.21.7.windows-amd64.zip")

if err != nil {
	t.Error(err)
}
```
## httpx
### httpclientx_test.go
#### TestHttpXGet
```go

clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Get("http://golang.google.cn")
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}
	t.Log(resp.StatusCode)
	t.Log(resp.Proto)
}
```
#### TestHttpXPost
```go

clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Get("http://golang.google.cn")
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}
	t.Log(resp.StatusCode)
	t.Log(resp.Proto)
}
```
## ip
## packet
## proxy
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
## serverx
