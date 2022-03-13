package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		return
	}

	bs, _ := ioutil.ReadFile(GetCurrentDirectory() + "\\auto_doc_code.tpl")
	docStr := genDocCodeFromTpl(string(bs), os.Args[1], os.Args[2])
	ioutil.WriteFile(strings.ToLower(os.Args[1])+"_test.go", []byte(docStr), 0666)
}

func genDocCodeFromTpl(tplFileContent, moduleName, moduleTitle string) string {
	t, err := template.New(fmt.Sprint(time.Now().Unix())).Parse(tplFileContent)
	if err != nil {
		return ""
	}

	data := map[string]string{
		"moduleName":  moduleName,
		"moduleTitle": moduleTitle,
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return ""
	}

	return buf.String()
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
