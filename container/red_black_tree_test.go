package container

import (
	"fmt"
	"testing"
)

func TestReaBlackTree(t *testing.T) {

	type personT struct {
		name   string
		age    int
		gender bool
		score  int64
	}

	tree := RedBlackTree{}

	for i := 0; i < 1000000; i++ {

		name := fmt.Sprintf("rongo%d", i)

		tree.Put(name, &personT{
			name:   name,
			age:    i,
			gender: true,
			score:  int64(i) + 100,
		})
	}

	nodeVal := tree.Get("rongo999999")

	personVal, ok := nodeVal.(*personT)
	if !ok {
		t.FailNow()
	}

	if personVal.name != "rongo999999" {
		t.FailNow()
	}

	if personVal.age != 999999 {
		t.FailNow()
	}

	if personVal.score != 1000099 {
		t.FailNow()
	}

}
