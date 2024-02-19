package utils

import (
	"testing"
	"time"
)

func TestContextWithTsTrace(t *testing.T) {
	t.Log(ContextWithTrace())

	t.Log(time.Now().UnixNano())
	time.Sleep(time.Second)
	t.Log(time.Now().UnixNano())
	t.Log(NanoTsBase32())
	t.Log(ContextWithTsTrace())
	t.Log(ContextWithTsTrace())
	t.Log(ContextWithTsTrace())
}
