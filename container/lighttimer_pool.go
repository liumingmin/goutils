package container

import (
	"hash/crc32"
	"time"
)

type LightTimerPool struct {
	lightTimers []*LightTimer
}

func NewLightTimerPool(size int, tickInterval time.Duration) *LightTimerPool {
	lightTimers := make([]*LightTimer, 0, size)

	for i := 0; i < size; i++ {
		lt := NewLightTimer()
		lt.StartTicks(tickInterval)
		lightTimers = append(lightTimers, lt)
	}

	return &LightTimerPool{
		lightTimers: lightTimers,
	}
}

func (p *LightTimerPool) AddTimer(key string, d time.Duration, callback CallbackFunc) {
	p.AddTimerWithDeadline(key, d, time.Time{}, callback, nil)
}

func (p *LightTimerPool) AddTimerWithDeadline(key string, d time.Duration, deadline time.Time,
	callback CallbackFunc, deadlineCallback CallbackFunc) {
	index := crc32.ChecksumIEEE([]byte(key)) % uint32(len(p.lightTimers))
	lt := p.lightTimers[index]
	lt.AddTimerWithDeadline(d, deadline, callback, deadlineCallback)
}

func (p *LightTimerPool) AddCallback(key string, d time.Duration, callback func()) {
	index := crc32.ChecksumIEEE([]byte(key)) % uint32(len(p.lightTimers))
	lt := p.lightTimers[index]
	lt.AddCallback(d, callback)
}
