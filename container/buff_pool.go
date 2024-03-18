package container

import (
	"sync"
)

var buffPoolMap map[int32]*sync.Pool

const (
	BUFF_128K = 128 * 1024
	BUFF_4M   = 4 * 1024 * 1024
)

type PoolBuffer128K []byte
type PoolBuffer4M []byte

func init() {
	buffPoolMap = make(map[int32]*sync.Pool)

	buffPoolMap[BUFF_128K] = &sync.Pool{
		New: func() interface{} {
			return PoolBuffer128K(make([]byte, BUFF_128K))
		},
	}

	buffPoolMap[BUFF_4M] = &sync.Pool{
		New: func() interface{} {
			return PoolBuffer4M(make([]byte, BUFF_4M))
		},
	}
}

func GetPoolBuff(size int32) []byte {
	val := buffPoolMap[size].Get()
	switch buf := val.(type) {
	case PoolBuffer128K:
		return buf
	case PoolBuffer4M:
		return buf
	}

	return nil
}

func PutPoolBuff(size int32, buff []byte) {
	switch size {
	case BUFF_128K:
		buffPoolMap[size].Put(PoolBuffer128K(buff))
	case BUFF_4M:
		buffPoolMap[size].Put(PoolBuffer4M(buff))
	}
}
