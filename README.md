# gin-limiter

This is to add a middleware to [Gin framework](https://github.com/gin-gonic/gin) to support rate limiting. It wraps [Juju's ratelimit](https://github.com/juju/ratelimit) implemetation as a Gin middleware

## Usage

```go
lm := limiter.NewRateLimiter(time.Minute, 10, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-API-KEY")
		if key != "" {
			return key, nil
		}
		return "", errors.New("API key is missing")
	})

r.GET("/ping", lm.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
```

This means the URI "/ping" only allows 10 requests per minutes per X-API-KEY. The key can be not only header but also with your own rules. You can decide what you try to limit with by returning the key.For example, it's also ok to use cookie or client ip as the key. 
