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

type TAG_STYLE = int

var (
	TAG_STYLE_NONE      = 0
	TAG_STYLE_ORIG      = 1
	TAG_STYLE_SNAKE     = 2
	TAG_STYLE_UNDERLINE = 3
)

func ConvertFieldStyle(str string, style TAG_STYLE) string {
	switch style {
	case TAG_STYLE_NONE:
		return ""
	case TAG_STYLE_ORIG:
		return str
	case TAG_STYLE_SNAKE:
		return strings.ToLower(str[:1]) + str[1:]
	case TAG_STYLE_UNDERLINE:
		tmpStr := str[1:]
		resultStr := make([]rune, 0, len(tmpStr))
		for _, r := range tmpStr {
			if r >= 65 && r <= 90 {
				resultStr = append(resultStr, '_', r+32)
			} else {
				resultStr = append(resultStr, r)
			}
		}
		return strings.ToLower(str[:1]) + string(resultStr)
	}

	return ""
}
