package container

import (
	"testing"

	"time"

	"github.com/liumingmin/goutils/safego"
)

type Node string

func (n Node) Health() bool {
	if n == "3333" || n == "4444" || n == "5555" {
		return false
	}
	return true
}

func TestNewChash(t *testing.T) {
	strs := []NodeHealth{Node("111"), Node("222"), Node("3333"), Node("4444"), Node("5555")}
	var configs []interface{}
	for _, str := range strs {
		configs = append(configs, str)
	}
	c := NewChash(strs)
	//t.Log(c.ring.Value)

	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode(Node("66666"))
	t.Log(c.GetNode("fdsafdwfe"))
	//
	c.AddNode(Node("777777"))
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode(Node("8888"))
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode(Node("99999"))
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode(Node("aaaaa"))
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode(Node("bbbbb"))
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode(Node("ccccc"))
	t.Log(c.GetNode("fdsafdwfe"))
}

func TestNewChash2(t *testing.T) {
	strs := []NodeHealth{Node("111"), Node("222"), Node("3333"), Node("4444"), Node("5555")}
	var configs []interface{}
	for _, str := range strs {
		configs = append(configs, str)
	}
	c := NewChash(strs)
	//t.Log(c.ring.Value)

	//safego.Go(func() {
	t.Log(c.GetNode("fdsafdwfe"))
	//})

	//safego.Go(func() {
	c.AddNode(Node("66666"))
	//})
	//safego.Go(func() {
	c.AddNode(Node("777777"))
	//})
	//safego.Go(func() {
	//c.AddNode(Node("8888"))
	//})
	//safego.Go(func() {
	//c.AddNode(Node("99999"))
	//})

	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})
	//

	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})

	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})

	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})

	safego.Go(func() {
		c.AddNode(Node("aaaaa"))
	})

	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})
	safego.Go(func() {
		c.AddNode(Node("bbbbb"))
	})

	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})
	safego.Go(func() {
		c.AddNode(Node("ccccc"))
	})
	safego.Go(func() {
		t.Log(c.GetNode("fdsafdwfe"))
	})

	time.Sleep(time.Second)
}
