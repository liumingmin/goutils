package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/liumingmin/goutils/log"

	"github.com/google/uuid"
)

func NewUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func MD5(origStr string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(origStr))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//检查keyname的keyvalue是否符合预期值expectKeyValues，如果不存在keyvalue，使用defaultKeyValue判断
func CheckKeyValueExpected(keyValues map[string]string, keyName, defaultKeyValue string, expectKeyValues []string) bool {
	if keyValue, exist := keyValues[keyName]; exist {
		log.Debug(context.Background(), "Found keyName: %v keyValue: %v, expectValue: %+v",
			keyName, keyValue, expectKeyValues)

		if found, _ := StringsInArray(expectKeyValues, keyValue); found {
			return true
		}
	} else {
		log.Debug(context.Background(), "Not Found  keyName: %v, defaultValue: %v, expectValue: %+v",
			keyName, defaultKeyValue, expectKeyValues)

		if found, _ := StringsInArray(expectKeyValues, defaultKeyValue); found {
			return true
		}
	}

	return false
}

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return filepath.Dir(path)
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}
func FileExt(filePath string) string {
	idx := strings.LastIndex(filePath, ".")
	ext := ""
	if idx >= 0 {
		ext = strings.ToLower(filePath[idx:])
	}

	return ext
}

func ReadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	ext := FileExt(filePath)

	if ext == ".jpg" {
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()

		return img, nil
	} else if ext == ".png" {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()

		return img, nil
	}

	return nil, errors.New(ext)
}

func WriteImage(img image.Image, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, nil)
}
