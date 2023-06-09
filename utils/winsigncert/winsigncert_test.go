package winsigncert

import (
	"fmt"
	"testing"
)

func TestValidateSignCert(t *testing.T) {
	result := ValidateSignCert("xxxx")
	fmt.Println(result)

	result = ValidateSignCert("xxxx")
	fmt.Println(result)

	result = ValidateSignCert("D:\\7-Zip-19\\7z.dll")
	fmt.Println(result)
}

func TestGetSignCertInfo(t *testing.T) {
	info := GetSignCertInfo("xxxxx")

	fmt.Println("ProgramName:", *info.ProgramName)
	fmt.Println("Subject:", *info.Subject)
	fmt.Println("MoreInfo:", *info.MoreInfo)
	fmt.Println("Timestamp:", *info.Timestamp)

	info = GetSignCertInfo("xxxxx")

	fmt.Println("ProgramName:", *info.ProgramName)
	fmt.Println("Subject:", *info.Subject)
	fmt.Println("Timestamp:", *info.Timestamp)

}
