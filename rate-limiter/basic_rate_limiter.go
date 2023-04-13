package ratelimiter

import (
	"sync"
	"time"
)

// Leaky bucket with request limit per period
type BasicRateLimiter struct {
	mu                         *sync.Mutex
	duration                   time.Duration
	maxRequestsPerPeriod       int
	numRunningRequests         int
	numRequestsInCurrentPeriod int
}

func NewBasicRateLimiter(duration time.Duration, maxRequestsPerPeriod int) *BasicRateLimiter {
	rl := &BasicRateLimiter{new(sync.Mutex), duration, maxRequestsPerPeriod, 0, 0}

	go func() {
		ticker := time.NewTicker(rl.duration)
		for range ticker.C {
			rl.mu.Lock()
			rl.numRequestsInCurrentPeriod = 0
			rl.mu.Unlock()
		}

	}()

	return rl
}

func (limiter *BasicRateLimiter) allow(id int) bool {
	if limiter.numRequestsInCurrentPeriod >= limiter.maxRequestsPerPeriod || limiter.numRunningRequests >= limiter.maxRequestsPerPeriod {
		return false
	} else {
		limiter.numRunningRequests += 1
		limiter.numRequestsInCurrentPeriod += 1
		return true
	}
}

func (limiter *BasicRateLimiter) done(id int) {
	limiter.mu.Lock()
	limiter.numRunningRequests -= 1
	limiter.mu.Unlock()
}
