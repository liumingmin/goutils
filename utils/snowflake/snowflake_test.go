package snowflake

import "testing"

func TestSnowflake(t *testing.T) {
	n, _ := NewNode(1)
	t.Log(n.Generate(), ",", n.Generate(), ",", n.Generate())
}
