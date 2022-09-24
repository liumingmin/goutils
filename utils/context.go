package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/liumingmin/goutils/log"
)

func ContextWithTrace() context.Context {
	traceId := strings.Replace(uuid.New().String(), "-", "", -1)
	return context.WithValue(context.Background(), log.LOG_TRADE_ID, traceId)
}

func ContextWithTsTrace() context.Context {
	return context.WithValue(context.Background(), log.LOG_TRADE_ID, NanoTsBase36()+RandBase36())
}
