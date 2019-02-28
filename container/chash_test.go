package container

import "testing"

func TestNewChash(t *testing.T) {
	strs := []string{"111", "222", "3333", "4444", "5555"}
	var configs []interface{}
	for _, str := range strs {
		configs = append(configs, str)
	}
	c := NewChash(configs)
	//t.Log(c.ring.Value)

	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode("66666")
	t.Log(c.GetNode("fdsafdwfe"))
	//
	c.AddNode("777777")
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode("8888")
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode("99999")
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode("aaaaa")
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode("bbbbb")
	t.Log(c.GetNode("fdsafdwfe"))

	c.AddNode("ccccc")
	t.Log(c.GetNode("fdsafdwfe"))
}
