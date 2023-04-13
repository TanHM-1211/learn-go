package ratelimiter

type RateLimiter interface {
	allow(id int) bool
	done(id int)
}
