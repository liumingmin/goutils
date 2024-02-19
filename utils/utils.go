package utils

import (
	"encoding/base32"
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"
)

const tsBase36 = 36

func NanoTsBase36() string {
	return strconv.FormatInt(time.Now().UnixNano(), tsBase36)
}

func RandBase36() string {
	return strconv.FormatInt(tsRander.Int63n(tsBase36), tsBase36)
}

var tsRander *rand.Rand

func NanoTsBase32() string {
	var bs [8]byte
	binary.LittleEndian.PutUint64(bs[:], uint64(time.Now().UnixNano()))
	return Base32EncodeToString(bs[:])
}

func RandBase32() string {
	return strconv.FormatInt(tsRander.Int63n(32), 32)
}

func Base32EncodeToString(bs []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bs)
}

func init() {
	tsRander = rand.New(rand.NewSource(time.Now().UnixNano()))
}
