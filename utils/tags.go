package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func AutoGenTags(vo interface{}, tagDefs map[string]TAG_STYLE) string {
	voType := reflect.TypeOf(vo)
	if voType.Kind() == reflect.Ptr {
		voType = voType.Elem()
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("type %s struct{\n", voType.Name()))

	var sortedTagName []string
	for k := range tagDefs {
		sortedTagName = append(sortedTagName, k)
	}
	sort.Strings(sortedTagName)

	for i := 0; i < voType.NumField(); i++ {
		dstField := voType.Field(i)

		var tagBuffer bytes.Buffer
		for _, tagName := range sortedTagName {
			tagStr := ConvertFieldStyle(dstField.Name, tagDefs[tagName])
			if tagStr != "" {
				tagBuffer.WriteString(fmt.Sprintf(`%s:"%s" `, tagName, tagStr))
			}
		}

		buffer.WriteString(fmt.Sprintf("%s %s `%s`\n", dstField.Name,
			dstField.Type.String(), tagBuffer.String())) //pkgStr+dstField.Type.Name()
	}

	buffer.WriteString("}")

	return buffer.String()
}

type TAG_STYLE = int

var (
	TAG_STYLE_NONE      = 0
	TAG_STYLE_ORIG      = 1
	TAG_STYLE_SNAKE     = 2
	TAG_STYLE_UNDERLINE = 3
)

func ConvertFieldStyle(str string, style TAG_STYLE) string {
	if len(str) == 0 {
		return str
	}

	switch style {
	case TAG_STYLE_NONE:
		return ""
	case TAG_STYLE_ORIG:
		return str
	case TAG_STYLE_SNAKE:
		if len(str) == 1 {
			return strings.ToLower(str)
		}

		return strings.ToLower(str[:1]) + str[1:]
	case TAG_STYLE_UNDERLINE:
		if len(str) == 1 {
			return strings.ToLower(str)
		}

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
