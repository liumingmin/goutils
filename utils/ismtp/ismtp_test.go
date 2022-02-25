package ismtp

import (
	"fmt"
	"strings"
	"testing"
)

func TestSendEmail(t *testing.T) {
	emailauth := LoginAuth(
		"from",
		"xxxxxx",
		"mailhost.com",
	)

	ctype := fmt.Sprintf("Content-Type: %s; charset=%s", "text/plain", "utf-8")

	msg := fmt.Sprintf("To: %s\r\nCc: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
		strings.Join([]string{"target@mailhost.com"}, ";"),
		"",
		"from@mailhost.com",
		"测试",
		ctype,
		"测试")

	err := SendMail("mailhost.com:port", //convert port number from int to string
		emailauth,
		"from@mailhost.com",
		[]string{"target@mailhost.com"},
		[]byte(msg),
	)

	if err != nil {
		t.Log(err)
		return
	}

	return
}
