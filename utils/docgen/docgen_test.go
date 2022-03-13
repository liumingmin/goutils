package docgen

import (
	"context"
	"strings"
	"testing"

	"github.com/liumingmin/goutils/middleware"
)

type testUser struct {
	UserId   string `json:"userId" binding:"required"`   //用户ID
	Nickname string `json:"nickname" binding:"required"` //用户名称
	Status   string `json:"status"`                      //用户状态
	Type     string `json:"pType"`                       //用户类型
}

type Page struct {
	Cursor int `json:"cursor"`
	Size   int `json:"size"`
}

type queryByNameReq struct {
	Page
	Nickname string `json:"nickname" binding:"required"` //用户名称
}

type CommonPageData struct {
	Cursor int                    `json:"cursor"`
	Size   int                    `json:"size"`
	Total  int                    `json:"total"`
	More   bool                   `json:"more"`
	Data   interface{}            `json:"data"`
	Extra  map[string]interface{} `json:"extra"`
}

type repDeleteId struct {
	UserId string `json:"userId" binding:"required"` //用户ID
}

//usage： docgen_cmd.exe 模块名和模块中文描述，自动生成模块的接口文档基础代码
//docgen_cmd.exe TestModule 测试模块

func TestGenDocTestUser(t *testing.T) {
	sb := strings.Builder{}
	sb.WriteString(genDocTestUserQuery())
	sb.WriteString(genDocTestUserCreate())
	sb.WriteString(genDocTestUserUpdate())
	sb.WriteString(genDocTestUserDelete())

	GenDoc(context.Background(), "用户管理", "doc/testuser.md", 2, sb.String())
}

func genDocTestUserQuery() string {
	moduleName := "用户管理列表"
	moduleUri := "/auto_doc/testuser/query"

	paramMap := map[string]string{
		PARAM_REQ_MODEL_PATH:  ".",
		PARAM_RESP_MODEL_PATH: ".",
		PARAM_REQ_REMARK:      "",
		PARAM_RESP_REMARK:     "",
	}

	dts := CommonPageData{
		Cursor: 0,
		Size:   10,
		Total:  1,
		Data:   []interface{}{},
		Extra:  nil,
	}

	docStr := GenModuleDoc(context.Background(), moduleName, moduleUri, paramMap, []string{"queryByNameReq"}, queryByNameReq{},
		[]string{"testUser"}, dts)
	return docStr
}

func genDocTestUserCreate() string {
	moduleName := "用户管理创建"
	moduleUri := "/auto_doc/testuser/create"

	paramMap := map[string]string{
		PARAM_REQ_MODEL_PATH:  ".",
		PARAM_RESP_MODEL_PATH: ".",
		PARAM_REQ_REMARK:      "",
		PARAM_RESP_REMARK:     "",
	}

	reqCase := testUser{
		UserId:   "1",
		Nickname: "超级棒棒",
		Status:   "0",
		Type:     "1",
	}
	docStr := GenModuleDoc(context.Background(), moduleName, moduleUri, paramMap, []string{"testUser"}, reqCase,
		[]string{}, nil)
	return docStr
}

func genDocTestUserUpdate() string {
	moduleName := "用户管理修改"
	moduleUri := "/auto_doc/testuser/update"

	paramMap := map[string]string{
		PARAM_REQ_MODEL_PATH:  ".",
		PARAM_RESP_MODEL_PATH: ".",
		PARAM_REQ_REMARK:      "",
		PARAM_RESP_REMARK:     "",
	}

	reqCase := testUser{
		UserId:   "1",
		Nickname: "超级人",
		Status:   "0",
		Type:     "1",
	}
	docStr := GenModuleDoc(context.Background(), moduleName, moduleUri, paramMap, []string{"testUser"}, reqCase,
		[]string{}, nil)
	return docStr
}

func genDocTestUserDelete() string {
	moduleName := "用户管理删除"
	moduleUri := "/auto_doc/testuser/delete"
	paramMap := map[string]string{
		PARAM_REQ_MODEL_PATH:  ".",
		PARAM_RESP_MODEL_PATH: ".",
		PARAM_REQ_REMARK:      "删除记录必须存在",
		PARAM_RESP_REMARK:     "删除记录必须存在",
	}

	reqCase := repDeleteId{
		UserId: "1",
	}

	respCase := middleware.DefaultServiceResponse{
		Code: -1,
		Msg:  "删除失败，记录不存在",
	}
	docStr := GenModuleDoc(context.Background(), moduleName, moduleUri, paramMap, []string{"repDeleteId"}, reqCase,
		[]string{}, nil, respCase)
	return docStr
}
