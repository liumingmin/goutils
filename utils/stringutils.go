package utils

import (
	"strings"
)

func StringsReverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func ComposeKey(keys ...string) string {
	return strings.Join(keys, "#")
}

func StringsInArray(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, -1
}

func StringsExcept(ss1 []string, ss2 []string) (se []string) {
	if len(ss1) == 0 {
		return
	}

	if len(ss2) == 0 {
		return ss1
	}

	for _, s1 := range ss1 {
		found := false
		for _, s2 := range ss2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			se = append(se, s1)
		}
	}
	return
}
