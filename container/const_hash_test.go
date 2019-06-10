package container

import (
	"fmt"
	"strconv"
	"testing"
)

type TestNode string

func (n TestNode) Id() string {
	return string(n)
}

func (n TestNode) Health() bool {
	//keyC32 := crc32.ChecksumIEEE([]byte(string(n)))
	return "node8" != string(n)
}

func TestConstHash(t *testing.T) {

	var ringchash CHashRing

	var configs []CHashNode
	for i := 0; i < 10; i++ {
		configs = append(configs, TestNode("node"+strconv.Itoa(i)))
	}

	ringchash.Adds(configs)

	fmt.Println(ringchash.Debug())

	fmt.Println("==================================", configs)

	fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

	fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
	//
	fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
	//
	fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

	fmt.Println(ringchash.Get("xvd:fs:xcvd"))

	var configs2 []CHashNode
	for i := 0; i < 2; i++ {
		configs2 = append(configs2, TestNode("node"+strconv.Itoa(10+i)))
	}
	ringchash.Adds(configs2)
	fmt.Println("==================================")
	fmt.Println(ringchash.Debug())
	fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

	fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
	//
	fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
	//
	fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

	fmt.Println(ringchash.Get("xvd:fs:xcvd"))

	ringchash.Del("node0")

	fmt.Println("==================================")
	fmt.Println(ringchash.Debug())
	fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

	fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
	//
	fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
	//
	fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

	fmt.Println(ringchash.Get("xvd:fs:xcvd"))
}
