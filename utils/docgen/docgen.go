package docgen

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liumingmin/goutils/middleware"
	"github.com/liumingmin/goutils/utils"
)

//optional
const (
	PARAM_MODULE_METHOD   = "moduleMethod"  // POST:application/json
	PARAM_REQ_MODEL_PATH  = "reqModelPath"  // request model file path
	PARAM_RESP_MODEL_PATH = "respModelPath" // response model file path
	PARAM_REQ_REMARK      = "reqRemark"     // request custom remark
	PARAM_RESP_REMARK     = "respRemark"    // response custom remark
)

const (
	MODEL_TYPE_REGX = `(?U)type\s+%s\s+struct\s+{`
)

var (
	modelFilePath = "."
)

func SetupModelPath(path string) {
	modelFilePath = path
}

func GenDoc(ctx context.Context, docName, docFilePath string, maxdepth int, moduleDocs string) {
	os.MkdirAll(filepath.Dir(docFilePath), 0700)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("# %s\n", docName))
	sb.WriteString("<!-- toc -->\n\n")

	sb.WriteString(moduleDocs)

	ioutil.WriteFile(docFilePath, []byte(sb.String()), 0666)

	if err := exec.Command("cmd", "/c", fmt.Sprintf("markdown-toc --maxdepth %d -i %s", maxdepth, docFilePath)).Run(); err != nil {
		fmt.Println(err)
	}
}

func GenModuleDoc(ctx context.Context, moduleName, moduleUri string, paramMap map[string]string, reqTypeNames []string, reqCase interface{},
	respTypeNames []string, respCases ...interface{}) string {
	if paramMap == nil {
		paramMap = make(map[string]string)
	}

	sb := strings.Builder{}

	sb.WriteString(genReqDoc(ctx, moduleName, moduleUri, paramMap, reqTypeNames, reqCase))
	sb.WriteString(genRespDoc(ctx, moduleName, moduleUri, paramMap, respTypeNames, respCases))
	sb.WriteString("\n")
	return sb.String()
}

func genReqDoc(ctx context.Context, moduleName, moduleUri string, paramMap map[string]string, reqStructNames []string,
	reqCase interface{}) string {
	method, _ := paramMap[PARAM_MODULE_METHOD]
	reqRemark, _ := paramMap[PARAM_REQ_REMARK]
	reqModelPath, _ := paramMap[PARAM_REQ_MODEL_PATH]

	if method == "" {
		method = "POST:application/json"
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("## %s\n", moduleName))
	sb.WriteString("### 请求说明\n")
	sb.WriteString(fmt.Sprintf("1. uri: %s\n", moduleUri))
	sb.WriteString(fmt.Sprintf("2. method: %s\n", method))

	if len(reqStructNames) != 0 {
		sb.WriteString("3. 参数说明：\n\n")
		for _, name := range reqStructNames {
			sb.WriteString(genDocTable(ctx, reqModelPath, name) + "\n\n")
		}
	}

	if reqRemark != "" {
		sb.WriteString("4. 补充说明：\n")
		sb.WriteString(reqRemark + "\n")
	}

	if reqCase != nil {
		sb.WriteString("### 请求样例\n")
		sb.WriteString(genJson(ctx, reqCase) + "\n")
	}
	return sb.String()
}

func genRespDoc(ctx context.Context, moduleName, moduleUri string, paramMap map[string]string, respStructNames []string,
	resps []interface{}) string {
	respRemark, _ := paramMap[PARAM_RESP_REMARK]
	respModelPath, _ := paramMap[PARAM_RESP_MODEL_PATH]

	sb := strings.Builder{}

	sb.WriteString("### 响应说明\n")

	for _, respStructName := range respStructNames {
		sb.WriteString(genDocTable(ctx, respModelPath, respStructName) + "\n\n")
	}

	if respRemark != "" {
		sb.WriteString("### 响应补充说明\n")
		sb.WriteString(respRemark + "\n")
	}

	for _, resp := range resps {
		if resp != nil {
			var respModel interface{}
			code := 0
			if commonRespPtr, ok := resp.(*middleware.DefaultServiceResponse); ok && commonRespPtr != nil {
				commonRespPtr.Tag = moduleUri
				respModel = commonRespPtr
				code = commonRespPtr.Code
			}

			if commonResp, ok := resp.(middleware.DefaultServiceResponse); ok {
				commonResp.Tag = moduleUri
				respModel = commonResp
				code = commonResp.Code
			}

			if commonResp, ok := resp.(middleware.ServiceResponse); ok {
				respModel = commonResp
				code = commonResp.GetCode()
			}

			if respModel != nil {
				msg := "成功响应样例"
				if code != 0 {
					msg = "失败响应样例"
				}
				sb.WriteString(fmt.Sprintf("### %v(code:%v)\n", msg, code))
				sb.WriteString(genJson(ctx, respModel) + "\n")
				continue
			}
		}

		respModel := middleware.DefaultServiceResponse{
			Code: 0,
			Msg:  "success",
			Tag:  moduleUri,
			Data: resp,
		}
		sb.WriteString("### 成功响应样例\n")
		sb.WriteString(genJson(ctx, respModel) + "\n")
	}

	return sb.String()
}

func genJson(ctx context.Context, instance interface{}) string {
	sb := strings.Builder{}

	sb.WriteString("```json\n")
	bs, _ := json.MarshalIndent(instance, "", "\t")
	sb.WriteString(string(bs) + "\n")
	sb.WriteString("```\n")
	return sb.String()
}

func genDocTable(ctx context.Context, filePath, typeName string) string {
	if typeName == "" {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("|参数名     |类型|是否必须|默认值  |说明    |\n")
	sb.WriteString("|----------|----|-------|-------|--------|\n")
	sb.WriteString(genDocTableContent(ctx, filePath, typeName))
	return sb.String()
}

func genDocTableContent(ctx context.Context, filePath, typeName string) string {
	if filePath == "" {
		filePath = modelFilePath
	}
	structFilePath := findTypeStructByName(ctx, filePath, typeName)
	if structFilePath == "" {
		return ""
	}

	sb := strings.Builder{}

	bs, _ := ioutil.ReadFile(structFilePath)
	content := string(bs)

	reg, _ := regexp.Compile(fmt.Sprintf(MODEL_TYPE_REGX, typeName))
	typeStart := reg.FindString(content)

	typeContent, _ := utils.ParseContentByTag(content, typeStart, "\n}")
	lines := strings.Split(typeContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		clause := strings.Split(line, " ")

		if len(clause) < 3 {
			if len(clause) > 0 {
				sb.WriteString(genDocTableContent(ctx, filePath, strings.TrimLeft(strings.TrimSpace(clause[0]), "*")))
			}
			continue
		}

		sb.WriteString(genDocTableLine(ctx, line))
	}

	return sb.String()
}

func findTypeStructByName(ctx context.Context, filePath, typeName string) string {
	reg, _ := regexp.Compile(fmt.Sprintf(MODEL_TYPE_REGX, typeName))

	absPath := filePath
	if !filepath.IsAbs(filePath) {
		absPath, _ = filepath.Abs(filePath)
	}

	dir := absPath
	fi, _ := os.Stat(absPath)
	if !fi.IsDir() {
		dir = filepath.Dir(absPath)
	}

	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		codeFilePath := dir + "/" + file.Name()
		bs, err := ioutil.ReadFile(codeFilePath)
		if err != nil {
			continue
		}
		content := string(bs)

		typeStart := reg.FindString(content)
		if typeStart != "" {
			return codeFilePath
		}
	}

	return ""
}

func genDocTableLine(ctx context.Context, line string) string {
	fieldName, _ := utils.ParseContentByTag(line, "json:\"", "\"")
	if specCharIdx := strings.Index(fieldName, ","); specCharIdx >= 0 { //去掉,omitempty 之类
		fieldName = fieldName[:specCharIdx]
	}

	fieldType, _ := utils.ParseContentByTag(line, " ", "`")
	remark, _ := utils.ParseContentByTag(line+"\n", "//", "\n")
	required, _ := utils.ParseContentByTag(line, "binding:\"", "\"")

	requiredCn := "否"
	if strings.Contains(required, "required") {
		requiredCn = "是"
	}

	docLine := fmt.Sprintf("|%s|%s|%s|-|%s|", fieldName, strings.TrimSpace(fieldType), requiredCn, remark) + "\n"
	return docLine
}
