package algorithm

import (
	"fmt"
	"testing"
)

func testDescartesStrContains(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, -1
}

func TestDescartes(t *testing.T) {
	result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})

	descartMap := make(map[string]bool)
	for _, item := range result {
		if ok, _ := testDescartesStrContains([]string{"A", "B"}, item[0]); !ok {
			t.FailNow()
		}

		if ok, _ := testDescartesStrContains([]string{"1", "2", "3"}, item[1]); !ok {
			t.FailNow()
		}

		if ok, _ := testDescartesStrContains([]string{"a", "b", "c", "d"}, item[2]); !ok {
			t.FailNow()
		}
		descartMap[fmt.Sprint(item)] = true
	}

	if len(descartMap) != 24 {
		t.FailNow()
	}
}
