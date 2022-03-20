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

func ParseContentByTag(content, tagStart, tagEnd string) (string, int) {
	if sIdx := strings.Index(content, tagStart); sIdx >= 0 {
		pos := sIdx + len(tagStart)
		content = content[pos:]
		if eIdx := strings.Index(content, tagEnd); eIdx >= 0 {
			tagContent := content[:eIdx]
			return tagContent, pos + eIdx + len(tagEnd)
		}
	}
	return "", 0
}

//检查keyname的keyvalue是否符合预期值expectKeyValues，如果不存在keyvalue，使用defaultKeyValue判断
func CheckKeyValueExpected(keyValues map[string]string, keyName, defaultKeyValue string, expectKeyValues []string) bool {
	if keyValue, exist := keyValues[keyName]; exist {
		if found, _ := StringsInArray(expectKeyValues, keyValue); found {
			return true
		}
	} else {
		if found, _ := StringsInArray(expectKeyValues, defaultKeyValue); found {
			return true
		}
	}

	return false
}
