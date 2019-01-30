package cache

import (
	"encoding/json"
	"errors"
)

type CacheStore interface {
	Exists(key string) bool

	Get(key string) ([]byte, error)

	Set(key string, value interface{}, expire int64) error

	Delete(key string) error

	Flush() error
}

type SimpleMemCache map[string][]byte

func (c *SimpleMemCache) Exists(key string) bool {
	_, ok := (*c)[key]
	return ok
}

func (c *SimpleMemCache) Get(key string) ([]byte, error) {
	result, isok := (*c)[key]
	if isok {
		return result, nil
	} else {
		return []byte{}, errors.New("cannot found ")
	}
}

func (c *SimpleMemCache) Set(key string, value interface{}, expire int64) error {
	bs, _ := json.Marshal(value)
	(*c)[key] = bs
	return nil
}

func (c *SimpleMemCache) Delete(key string) error {
	delete(*c, key)
	return nil
}

func (c *SimpleMemCache) Flush() error {
	for key := range *c {
		delete(*c, key)
	}
	return nil
}
