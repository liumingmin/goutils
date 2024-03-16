package utils

import "strings"

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

// 检查keyname的keyvalue是否符合预期值expectKeyValues，如果不存在keyvalue，使用defaultKeyValue判断
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
