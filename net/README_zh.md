**其他语言版本: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [网络库](#%E7%BD%91%E7%BB%9C%E5%BA%93)
  * [bwlimit](#bwlimit)
  * [httpdownloader](#httpdownloader)
  * [兼容http1.x和2.0的httpclient](#%E5%85%BC%E5%AE%B9http1x%E5%92%8C20%E7%9A%84httpclient)
  * [ip](#ip)
  * [tcp包model](#tcp%E5%8C%85model)
  * [ssh proxy](#ssh-proxy)
  * [兼容http1.x和2.0的http server](#%E5%85%BC%E5%AE%B9http1x%E5%92%8C20%E7%9A%84http-server)

<!-- tocstop -->

# 网络库
## bwlimit
## httpdownloader
### file_offset_writer_test.go
#### TestFileOffsetWriter
```go

os.MkdirAll(testTempDirPath, 0666)

savePath := filepath.Join(testTempDirPath, "testFileOffsetWriter")
file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
if err != nil {
	t.Error(err)
}

var mutex sync.Mutex
writer := NewOffsetWriter(file, &mutex, 0, 1024)
writer.Write([]byte("1234567890"))
writer.ResetOffset()
writer.Write([]byte("1234567890"))
file.Close()

bs, err := os.ReadFile(savePath)
if err != nil {
	t.Error(err)
}

if string(bs) != "1234567890" {
	t.Error(string(bs))
}
```
### httpdownloader_test.go
#### TestHttpDownloaderDownload
```go

os.MkdirAll(testTempDirPath, 0666)

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

savePath := filepath.Join(testTempDirPath, "vc_redist")

url := "https://aka.ms/vs/17/release/vc_redist.arm64.exe"
resp, err := http.Head(url)
if err != nil {
	t.Error(err)
}
if resp != nil {
	resp.Body.Close()
}

contentLen := resp.Header.Get("Content-Length")

fileSize, _ := strconv.ParseInt(contentLen, 10, 64)

err = downloader.Download(context.Background(), url, http.Header{
	"User-Agent": []string{"Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Mobile Safari/537.36"},
}, savePath)

if err != nil {
	t.Error(err)
}

fileInfo, err := os.Stat(savePath)
if err != nil {
	t.Error(err)
}

if fileInfo.Size() != fileSize {
	t.Error(fileSize)
}

```
## 兼容http1.x和2.0的httpclient
### httpclientx_test.go
#### TestHttpXGet
```go


clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Get("http://127.0.0.1:" + testHttpxPortH1 + "/goutils/httpx")
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/1.1" {
		t.Error(resp.Proto)
	}
}

for i := 0; i < 3; i++ {
	resp, err := clientX.Get("http://127.0.0.1:" + testHttpxPortH2 + "/goutils/httpx")
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/2.0" {
		t.Error(resp.Proto)
	}
}
```
#### TestHttpXPost
```go


clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Post("http://127.0.0.1:"+testHttpxPortH1+"/goutils/httpx", "application/json", strings.NewReader(""))
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/1.1" {
		t.Error(resp.Proto)
	}
}

for i := 0; i < 3; i++ {
	resp, err := clientX.Post("http://127.0.0.1:"+testHttpxPortH2+"/goutils/httpx", "application/json", strings.NewReader(""))
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/2.0" {
		t.Error(resp.Proto)
	}
}
```
#### TestHttpXHead
```go


clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.Head("http://127.0.0.1:" + testHttpxPortH1 + "/goutils/httpx")
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/1.1" {
		t.Error(resp.Proto)
	}
}

for i := 0; i < 3; i++ {
	resp, err := clientX.Head("http://127.0.0.1:" + testHttpxPortH2 + "/goutils/httpx")
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/2.0" {
		t.Error(resp.Proto)
	}
}
```
#### TestHttpXPostForm
```go


clientX := getHcx()

for i := 0; i < 3; i++ {
	resp, err := clientX.PostForm("http://127.0.0.1:"+testHttpxPortH1+"/goutils/httpx", url.Values{"key": []string{"value"}})
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/1.1" {
		t.Error(resp.Proto)
	}
}

for i := 0; i < 3; i++ {
	resp, err := clientX.PostForm("http://127.0.0.1:"+testHttpxPortH2+"/goutils/httpx", url.Values{"key": []string{"value"}})
	if err != nil {
		t.Error(fmt.Errorf("error making request: %v", err))
	}

	if resp.StatusCode >= 400 {
		t.Error(resp.StatusCode)
	}

	if resp.Proto != "HTTP/2.0" {
		t.Error(resp.Proto)
	}
}
```
## ip
## tcp包model
## ssh proxy
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
## 兼容http1.x和2.0的http server
