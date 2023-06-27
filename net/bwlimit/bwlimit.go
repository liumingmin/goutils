package bwlimit

import (
	"context"
	"sync"

	"github.com/liumingmin/goutils/log"
	"golang.org/x/time/rate"
)

// holds info about the rate limiters in use
type BwLimit struct {
	mu        sync.RWMutex
	limiter   *rate.Limiter
	bandwidth int64
}

func (t *BwLimit) LimitBandwidth(n int) {
	if t.limiter == nil {
		return
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.limiter == nil {
		return
	}

	// Limit the transfer speed if required
	err := t.limiter.WaitN(context.Background(), n)
	if err != nil {
		log.Error(context.Background(), "limiter.WaitN error: %v", err)
	}
}

// SetBwLimit sets the current bandwidth limit
func (t *BwLimit) SetBwLimit(bandwidth int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if bandwidth == t.bandwidth {
		return
	}

	t.bandwidth = bandwidth

	if t.bandwidth > 0 {
		t.limiter = newRateLimiter(t.bandwidth)
		log.Debug(context.Background(), "Bandwidth limit set to %v", t.bandwidth)
	} else {
		t.limiter = nil
		log.Debug(context.Background(), "Bandwidth limit reset to unlimited")
	}
}

const defaultMaxBurstSize = 4 * 1024 * 1024 // must be bigger than the biggest request

func newRateLimiter(bandwidth int64) *rate.Limiter {
	// Relate maxBurstSize to bandwidth limit
	// 4M gives 2.5 Gb/s on Windows
	// Use defaultMaxBurstSize up to 2GBit/s (256MiB/s) then scale
	maxBurstSize := (bandwidth * defaultMaxBurstSize) / (256 * 1024 * 1024)
	if maxBurstSize < defaultMaxBurstSize {
		maxBurstSize = defaultMaxBurstSize
	}
	log.Debug(context.Background(), "bandwidth=%v maxBurstSize=%v", bandwidth, maxBurstSize)

	limiter := rate.NewLimiter(rate.Limit(bandwidth), int(maxBurstSize))
	if limiter != nil {
		err := limiter.WaitN(context.Background(), int(maxBurstSize))
		if err != nil {
			log.Error(nil, "Failed to NewLimiter: %v", err)
		}
	}
	return limiter
}
