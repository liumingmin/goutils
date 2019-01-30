package cache

import "errors"

type CacheStore interface {
	Exists(key string) bool

	Get(key string) (interface{}, error)

	Set(key string, value interface{}, expire int64) error

	Delete(key string) error

	Flush() error
}

type SimpleMemCache map[string]interface{}

func (c *SimpleMemCache) Exists(key string) bool {
	_, ok := (*c)[key]
	return ok
}

func (c *SimpleMemCache) Get(key string) (interface{}, error) {
	result, isok := (*c)[key]
	if isok {
		return result, nil
	} else {
		return "", errors.New("cannot found ")
	}
}

func (c *SimpleMemCache) Set(key string, value interface{}, expire int64) error {
	(*c)[key] = value
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
