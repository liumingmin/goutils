package middleware

import (
	"sync"
	"time"

	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/utils"
)

type cbkParam struct {
	isBreaker  bool
	errCount   int64
	totalCount int64

	accessLast int64
	simpleLast int64
}

type CircuitBreaker struct {
	Name        string
	cbkParamMap map[string]*cbkParam
	lock        sync.RWMutex

	isTurnOn            bool
	simpleInterval      time.Duration
	testRecoverInterval time.Duration
	totalThreshold      int64
	errorRateThreshold  float64
}

func (c *CircuitBreaker) Init() {
	c.cbkParamMap = make(map[string]*cbkParam)

	confPrefix := "cbk"
	if c.Name != "" {
		confPrefix += "." + c.Name
	}

	c.isTurnOn = conf.ExtBool(confPrefix+".isTurnOn", true)
	c.simpleInterval = conf.ExtDuration(confPrefix+".simpleInterval", time.Second*10)
	c.testRecoverInterval = conf.ExtDuration(confPrefix+".simpleInterval", time.Second*30)
	c.totalThreshold = conf.ExtInt64(confPrefix+".totalThreshold", 100)
	c.errorRateThreshold = conf.ExtFloat64(confPrefix+".errorRateThreshold", 0.9)
}

func (c *CircuitBreaker) Check(key string) bool {
	if !c.isTurnOn {
		return true
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	if param, ok := c.cbkParamMap[key]; ok {
		if param.isBreaker {
			if utils.Abs64(time.Now().UnixNano()-param.accessLast) < int64(c.testRecoverInterval) {
				return false
			}
		}
	}

	return true
}

func (c *CircuitBreaker) accessed(param *cbkParam) {
	now := time.Now().UnixNano()
	if utils.Abs64(now-param.simpleLast) > int64(c.simpleInterval) {
		param.errCount = 0
		param.totalCount = 0
		param.simpleLast = now
	}
	param.totalCount++
	param.accessLast = now
}

func (c *CircuitBreaker) Succeed(key string) {
	if !c.isTurnOn {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if param, ok := c.cbkParamMap[key]; ok {
		c.accessed(param)

		if param.isBreaker {
			param.isBreaker = false
		}
	}
}

func (c *CircuitBreaker) Failed(key string) {
	if !c.isTurnOn {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if param, ok := c.cbkParamMap[key]; ok {
		c.accessed(param)
		param.errCount++

		if param.totalCount > c.totalThreshold && float64(param.errCount)/float64(param.totalCount) > c.errorRateThreshold {
			param.isBreaker = true
		}
	} else {
		param := &cbkParam{}
		c.accessed(param)
		param.errCount++
		c.cbkParamMap[key] = param
	}
}
