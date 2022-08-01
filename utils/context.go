package utils

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/liumingmin/goutils/log"
)

func ContextWithTrace() context.Context {
	traceId := strings.Replace(uuid.New().String(), "-", "", -1)
	return context.WithValue(context.Background(), log.LOG_TRADE_ID, traceId)
}

const ctxTsBase = 36

func ContextWithTsTrace() context.Context {
	return context.WithValue(context.Background(), log.LOG_TRADE_ID, strconv.FormatInt(time.Now().UnixNano(), ctxTsBase)+
		strconv.FormatInt(rand.Int63n(ctxTsBase), ctxTsBase))
}
