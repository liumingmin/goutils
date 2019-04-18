package middleware

import (
	"bytes"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetCaptchaImage(c *gin.Context) {
	id := c.Query("id")

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	var content bytes.Buffer
	c.Header("Content-Type", "image/png")

	captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)

	c.Data(http.StatusOK, "image/png", content.Bytes())
}

func GetCaptchaId() string {
	return captcha.New()
}

func VerifyCaptcha(paramFunc func(*gin.Context) (string, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaId, captchaCode := paramFunc(c)
		if !captcha.VerifyString(captchaId, captchaCode) {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		c.Next()
		return
	}
}
