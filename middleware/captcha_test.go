package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVerifyCaptcha(t *testing.T) {
	router := gin.New()
	router.GET("/cimage", GetCaptchaImage)

	g := router.Group("/", VerifyCaptcha(func(c *gin.Context) (string, string) {
		return c.DefaultPostForm("cid", ""), c.DefaultPostForm("ccode", "")
	}))
	g.POST("/submit", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	var tplStr = `
<!doctype html>
<html>
 <body>
  <form method="post" action="/submit">
		<div><input type="hidden" name="cid" value="%s"></div>
		<div><input type="image" src="/cimage?id=%s"></div>
		<div><input type="text" name="ccode" value=""></div>
		<div><input type="submit" value="submit"></div>
  </form>
 </body>
</html>
`
	router.GET("/", func(c *gin.Context) {
		cid := GetCaptchaId()
		c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf(tplStr, cid, cid)))
	})
	router.Run(":8080")
}
