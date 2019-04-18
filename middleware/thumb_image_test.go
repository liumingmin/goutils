package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestThumbImageServe(t *testing.T) {
	router := gin.New()
	router.Use(ThumbImageServe("/images", GinHttpFs("G:/images", false)))
	router.Run(":8080")
}
