package rpcbi

import (
	"fmt"
	"testing"
	"time"

	"crypto/rand"
	"crypto/tls"

	"github.com/liumingmin/goutils/safego"
)

type Arith int

type Args struct {
	A, B int
}

var test int

func (t *Arith) Multiply(args *Args, reply *int) error {
	test++

	//if test%1000 == 0 {
	//	return errors.New("normal error ")
	//}
	fmt.Println("aaaaaaaaaaaaaa")
	*reply = args.A*args.B + 10000
	//time.Sleep(time.Millisecond * 100)
	return nil
}

type connCallback struct {
	server *RpcServer
}

func (c *connCallback) ConnFinished(id string) {
	safego.Go(func() {
		time.Sleep(time.Second * 2)

		s := c.server.GetSession(id)
		args := &Args{9, 100}
		var reply int
		err := s.Call("SArith.Multiply", args, &reply)
		if err != nil {
			fmt.Println("arith error:", err)
		}

		fmt.Println(reply)
	})
}

func (c *connCallback) DisconnFinished(id string) {

}

func TestNewRpcServer(t *testing.T) {
	c := &connCallback{}
	s := NewRpcServer(PROTOCOL_FORMAT_JSON, c)
	c.server = s
	s.RegisterService("Arith", new(Arith))

	crt, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		t.Fatal(err.Error())
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader

	lis, _ := tls.Listen("tcp", ":12345", tlsConfig) //tlsConfig
	s.Serve(lis)
}

func TestNewRpcClient(t *testing.T) {
	//crt, err := tls.LoadX509KeyPair("server.crt", "server.key")
	//if err != nil {
	//	t.Fatal(err.Error())
	//}
	tlsConfig := &tls.Config{}
	//tlsConfig.Certificates = []tls.Certificate{crt}
	//tlsConfig.Time = time.Now
	//tlsConfig.Rand = rand.Reader
	tlsConfig.InsecureSkipVerify = true

	conn, err := tls.Dial("tcp", "10.11.134.34:12345", tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	c := NewRpcClient(PROTOCOL_FORMAT_JSON, "22345")
	err = c.Start(conn)
	if err != nil {
		t.Log(err)
		return
	}

	c.RegisterService("SArith", new(Arith))
	//fmt.Println(c.Handshake("0001", 2))
	//fmt.Println(err)
	//time.Sleep(time.Second * 3)

	args := &Args{7, 100}
	var reply int
	err = c.Call("Arith.Multiply", args, &reply)
	if err != nil {
		fmt.Println("arith error:", err)
	}
	fmt.Println(reply)

	args2 := &Args{6, 100}
	var reply2 int
	err = c.Call("Arith.Multiply", args2, &reply2)
	if err != nil {
		fmt.Println("arith2 error:", err)
	}
	fmt.Println(reply2)

	time.Sleep(time.Hour)
}
