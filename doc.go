package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liumingmin/goutils/utils"
)

func main() {
	genDocByTestFile(utils.GetCurrPath(), 1, &strings.Builder{}, false)
	genDocByTestFile(utils.GetCurrPath(), 1, &strings.Builder{}, true)
}

var moduleCnName = map[string]string{
	"algorithm":          "算法模块",
	"cache":              "缓存模块",
	"mem_cache_test.go":  "内存缓存",
	"rds_cache_test.go":  "Redis缓存",
	"conf":               "yaml配置模块",
	"container":          "容器模块",
	"bitmap_test.go":     "比特位表",
	"const_hash_test.go": "一致性HASH",
	"lighttimer_test.go": "轻量级计时器",
	"db":                 "数据库",
	"elasticsearch":      "ES搜索引擎",
	"es6":                "ES6版本API",
	"es7":                "ES7版本API",
	"kafka":              "kafka消息队列",
	"mongo":              "mongo数据库",
	"redis":              "go-redis",
	"log":                "zap日志库",
	"net":                "网络库",
	"httpx":              "兼容http1.x和2.0的httpclient",
	"packet":             "tcp包model",
	"proxy":              "ssh proxy",
	"serverx":            "兼容http1.x和2.0的http server",
	"utils":              "通用工具库",
	"cbk":                "熔断器",
	"csv":                "CSV文件解析为MDB内存表",
	"distlock":           "分布式锁",
	"fsm":                "有限状态机",
	"hc":                 "httpclient工具",
	"ismtp":              "邮件工具",
	"safego":             "安全的go协程",
	"ws":                 "websocket客户端和服务端库",
	"docgen":             "文档自动生成",
	"crc16_test.go":      "crc16算法",
	"descartes_test.go":  "笛卡尔组合",
	"list_test.go":       "Redis List工具库",
	"zset_test.go":       "Redis ZSet工具库",
	"mq_test.go":         "Redis PubSub工具库",
	"lock_test.go":       "Redis 锁工具库",
	"tags_test.go":       "结构体TAG生成器",
	"snowflake_test.go":  "雪花ID生成器",
}

var enSwitchStr = "**Read this in other languages: [English](README.md), [中文](README_zh.md).**\n\n"
var cnSwitchStr = "**其他语言版本: [English](README.md), [中文](README_zh.md).**\n\n"

// dir := filepath.Dir(filePath)
func genDocByTestFile(dir string, level int, sb *strings.Builder, isCn bool) map[string]string {
	files, _ := ioutil.ReadDir(dir)

	nextLevel := level + 1

	outlines := make(map[string]string)

	for _, file := range files {
		if file.IsDir() {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			genDirLevel(file.Name(), nextLevel, sb, isCn)
			genDocByTestFile(filepath.Join(dir, file.Name()), nextLevel, sb, isCn)

			if level == 1 {
				readmeName := "README.md"
				if isCn {
					readmeName = "README_zh.md"
				}

				readmePath := strings.TrimRight(file.Name(), "/\\") + "/" + readmeName
				outlines[file.Name()] = readmePath

				if file.Name() == "ws" {
					sb.Reset()
					continue
				}
				content := sb.String()
				if strings.TrimSpace(content) == "" {
					sb.Reset()
					continue
				}
				ioutil.WriteFile(readmePath, []byte("<!-- toc -->\n"+content), 0666)
				sb.Reset()

				if err := exec.Command("cmd", "/c", "markdown-toc --maxdepth 2 -i "+readmePath).Run(); err != nil {
					fmt.Println(err)
				}

				bs1, _ := ioutil.ReadFile(readmePath)
				fileContent := string(bs1)
				if isCn {
					fileContent = cnSwitchStr + fileContent
				} else {
					fileContent = enSwitchStr + fileContent
				}
				ioutil.WriteFile(readmePath, []byte(fileContent), 0666)

			}
			continue
		}

		if strings.HasSuffix(file.Name(), "_test.go") {
			codeFilePath := dir + "/" + file.Name()
			bs, err := ioutil.ReadFile(codeFilePath)
			if err != nil {
				continue
			}
			content := string(bs)

			genDirLevel(file.Name(), nextLevel, sb, isCn)
			parseTestCode(nextLevel, content, sb, isCn)
		}
	}

	return outlines
}

func genDirLevel(dirName string, level int, sb *strings.Builder, isCn bool) {
	prefixSymbol := ""
	for i := 0; i < level-1; i++ {
		prefixSymbol += "#"
	}

	dirActName := dirName

	if isCn {
		if cnName, ok := moduleCnName[dirName]; ok {
			dirActName = cnName
		}
	}

	fmt.Println(dirName)
	sb.WriteString(fmt.Sprintf("%s %s\n", prefixSymbol, dirActName))
}

func parseTestCode(level int, content string, sb *strings.Builder, isCn bool) {
	reg, _ := regexp.Compile(`(?U)func (?P<fname>.*)\(t \*testing\.T\) *\{(?P<body>(.|\n)*)\n\}`)
	match := reg.FindAllStringSubmatch(content, -1)

	for _, item := range match {
		genDirLevel(item[1], level+1, sb, isCn)

		sb.WriteString("```go\n")
		sb.WriteString(removePrefixTab(item[2]) + "\n")
		sb.WriteString("```\n")
	}
}

func removePrefixTab(code string) string {
	lines := strings.Split(code, "\n")
	newLines := make([]string, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "\t") {
			line = strings.Replace(line, "\t", "", 1)
		}
		newLines = append(newLines, line)
	}
	return strings.Join(newLines, "\n")
}
