package container

import (
	"fmt"
	"strconv"
	"testing"
)

func TestConstHash(t *testing.T) {

	var ringchash CHashRing

	var configs []interface{}
	for i := 0; i < 10; i++ {
		configs = append(configs, "node"+strconv.Itoa(i))
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

	var configs2 []interface{}
	for i := 0; i < 2; i++ {
		configs2 = append(configs2, "node"+strconv.Itoa(10+i))
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
