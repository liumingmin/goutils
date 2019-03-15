package rpcbi

import (
	"fmt"
	"testing"
	"time"
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

func TestNewRpcServer(t *testing.T) {
	s, _ := NewRpcServer(1)
	s.Register(new(Arith))
	s.Start("tcp", "127.0.0.1:12345")
}

func TestNewRpcClient(t *testing.T) {
	c, err := NewRpcClient("127.0.0.1:12345", 10)
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
