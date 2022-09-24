package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(origStr string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(origStr))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
