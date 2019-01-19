package utils

import "strings"

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
