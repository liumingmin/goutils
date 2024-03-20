package container

import (
	"strconv"
	"testing"
)

type testConstHashNode string

func (n testConstHashNode) Id() string {
	return string(n)
}

func (n testConstHashNode) Health() bool {
	return "node8" != string(n)
}

func TestConstHash(t *testing.T) {

	var ringchash CHashRing

	var configs []CHashNode
	for i := 0; i < 10; i++ {
		configs = append(configs, testConstHashNode("node"+strconv.Itoa(i)))
	}

	ringchash.Adds(configs)

	t.Log("init:", ringchash.Debug())

	if ringchash.GetByC32(100, false).Id() != "node0" {
		t.Fail()
	}

	if ringchash.GetByC32(134217727, false).Id() != "node0" {
		t.Fail()
	}

	if ringchash.GetByC32(134217728, false).Id() != "node8" {
		t.Fail()
	}

	var configs2 []CHashNode
	for i := 0; i < 2; i++ {
		configs2 = append(configs2, testConstHashNode("node"+strconv.Itoa(10+i)))
	}
	ringchash.Adds(configs2)

	t.Log("add 2 nodes", ringchash.Debug())

	if ringchash.GetByC32(134217727, false).Id() != "node10" {
		t.Fail()
	}

	if ringchash.GetByC32(134217728, false).Id() != "node10" {
		t.Fail()
	}

	ringchash.Del("node0")
	t.Log("del 1 node", ringchash.Debug())

	if ringchash.GetByC32(100, false).Id() != "node10" {
		t.Fail()
	}

	t.Log(ringchash.GetByKey("goutils", false))
}
