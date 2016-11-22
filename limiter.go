package ginlimiter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
	"fmt"
)

type RateKeyFunc func(ctx *gin.Context) (string, error)

type RateLimiterMiddleware struct {
	fillInterval time.Duration
	capacity     int64
	ratekeygen   RateKeyFunc
	limiters     map[string]*ratelimit.Bucket
}

func (r *RateLimiterMiddleware) get(ctx *gin.Context) (*ratelimit.Bucket, error) {
	key, err := r.ratekeygen(ctx)

	if err != nil {
		return nil, err
	}

	if limiter, existed := r.limiters[key]; existed {
		return limiter, nil
	}

	limiter := ratelimit.NewBucket(r.fillInterval, r.capacity)
	r.limiters[key] = limiter
	return limiter, nil
}

func (r *RateLimiterMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter, err := r.get(ctx)
		if err != nil || limiter.TakeAvailable(1) == 0 {
			if err == nil {
				err = errors.New("Too many requests")
			}
			ctx.AbortWithError(429, err)
		} else {
			ctx.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
			ctx.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
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
