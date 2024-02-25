package proxy

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func getSshClient(t *testing.T) *ssh.Client {
	pemBytes, _ := ioutil.ReadFile("")
	config, err := NewSshClient("127.0.0.1:22", "app", pemBytes, "")
	if err != nil {
		t.Fatalf("NewSshClient failed %v", err)
	}

	sshClient, err := config.SshDial()
	if err != nil {
		t.Fatalf("SshClientConfig connect failed %v", err)
	}
	return sshClient
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

func TestMain(m *testing.M) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:22", time.Second*2)
	if err != nil {
		fmt.Println("Please install ssh on local and start at port: 22, then run test.")
		return
	}
	conn.Close()

	m.Run()
}
