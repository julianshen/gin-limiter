package ginlimiter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

type RateKeyFunc func(ctx *gin.Context) string

type RateLimiterMiddleware struct {
	fillInterval time.Duration
	capacity     int64
	ratekeygen   RateKeyFunc
	limiters     map[string]*ratelimit.Bucket
}

func (r *RateLimiterMiddleware) get(ctx *gin.Context) *ratelimit.Bucket {
	key := r.ratekeygen(ctx)
	if limiter, existed := r.limiters[key]; existed {
		return limiter
	}

	limiter := ratelimit.NewBucket(r.fillInterval, r.capacity)
	r.limiters[key] = limiter
	return limiter
}

func (r *RateLimiterMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter := r.get(ctx)
		if limiter.TakeAvailable(1) == 0 {
			ctx.AbortWithError(429, errors.New("Too many requests"))
		}
	}
}

func NewRateLimiter(fillInterval time.Duration, capacity int64, keyGen RateKeyFunc) *RateLimiterMiddleware {
	limiters := make(map[string]*ratelimit.Bucket)
	return &RateLimiterMiddleware{
		fillInterval,
		capacity,
		keyGen,
		limiters,
	}
}
