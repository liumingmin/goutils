package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

const (
	TRACE_ID = "__traceId__"
)

func ContextWithTrace() context.Context {
	traceId := strings.Replace(uuid.New().String(), "-", "", -1)
	return context.WithValue(context.Background(), TRACE_ID, traceId)
}
