package cache

import "errors"

type CacheStore interface {
	Get(key string) (string, error)

	Set(key string, value string, expire int64) error

	Delete(key string) error

	Flush() error
}

type SimpleMemCache map[string]string

func (c *SimpleMemCache) Get(key string) (string, error) {
	result, isok := (*c)[key]
	if isok {
		return result, nil
	} else {
		return "", errors.New("cannot found ")
	}
}

func (c *SimpleMemCache) Set(key string, value string, expire int64) error {
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
