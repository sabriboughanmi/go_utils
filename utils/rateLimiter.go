package utils

import (
	"time"
)

type RateLimiter struct {
	chanBlocker chan interface{}
	ticker      *time.Ticker
	running     bool
}

//Start must be called before any API call to start the Rate Limiter.
func (r *RateLimiter) Start() {
	r.running = true
	go func() {
		for range r.ticker.C {
			//r.chanBlocker <- nil
			if r.running == false {
				return
			}
		}
	}()
}

//Stop must be called after all api calls has finished
func (r *RateLimiter) Stop() {
	r.running = false
	close(r.chanBlocker)
}

//Wait must be called before any rate limited API Call
func (r *RateLimiter) Wait() {
	if !r.running {
		panic("RateLimiter need to be running in order to wait for it")
	}
	<-r.chanBlocker
}

//CreateLimiter returns a RateLimiter which can be used to rate limit any process operations.
//	requestInterval: is the interval duration between Operations
//	burstCount:  allow short bursts of requests in the rate limiting scheme while preserving the overall rate limit.
//	requestsCount: is the number of Operations you are intending to Rate Limit. (this value is optional but can optimize data initialization)
func CreateLimiter(requestInterval time.Duration, burstCount int, requestsCount ...int) *RateLimiter {
	var ticker = time.NewTicker(requestInterval)
	var rateLimiter = RateLimiter{
		chanBlocker: make(chan interface{}, func() int {
			if len(requestsCount) > 0 {
				return requestsCount[0]
			}
			return 0
		}()),
		ticker: ticker,
	}
	//fill chanBlocker with free operations to prevent waiting for burst calls.
	if burstCount > 0 {
		go func() {
			for i := 0; i < burstCount; i++ {
				rateLimiter.chanBlocker <- nil
			}
		}()
	}
	return &rateLimiter
}
