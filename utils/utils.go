package utils

import (
	"math/rand"
	"strconv"
	"time"
)

const tsBase36 = 36

func NanoTsBase36() string {
	return strconv.FormatInt(time.Now().UnixNano(), tsBase36)
}

func RandBase36() string {
	return strconv.FormatInt(rand.Int63n(tsBase36), tsBase36)
}
