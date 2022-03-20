package algorithm

import "testing"

func TestDescartes(t *testing.T) {
	result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})
	for _, item := range result {
		t.Log(item)
	}
}
