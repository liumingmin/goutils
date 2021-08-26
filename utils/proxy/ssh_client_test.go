package proxy

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func getSshClient(t *testing.T) *SshClient {
	pemBytes, _ := ioutil.ReadFile("")
	client, err := NewSshClient("127.0.0.1:22", "app", pemBytes, "")
	if err != nil {
		t.Fatalf("NewSshClient failed %v", err)
	}

	err = client.Connect()
	if err != nil {
		t.Fatalf("SshClient connect failed %v", err)
	}
	return client
}

func TestSshClient(t *testing.T) {
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
}

func TestMysqlSshClient(t *testing.T) {
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
}
