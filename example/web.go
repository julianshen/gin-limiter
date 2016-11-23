package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	limiter "github.com/julianshen/gin-limiter"
	"time"
)

func main() {
	r := gin.Default()

	//Allow only 10 requests per minute per API-Key
	lm := limiter.NewRateLimiter(time.Minute, 10, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-API-KEY")
		if key != "" {
			return key, nil
		}
		return "", errors.New("API key is missing")
	})
	//Apply only to /ping
	r.GET("/ping", lm.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Allow only 5 requests per second per user
	lm2 := limiter.NewRateLimiter(time.Second, 5, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-USER-TOKEN")
		if key != "" {
			return key, nil
		}
		return "", errors.New("User is not authorized")
	})

	//Apply to a group
	x := r.Group("/v2")
	x.Use(lm2.Middleware())
	{
		x.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		x.GET("/another_ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong pong",
			})
		})
	}

	r.Run(":8080") // listen and server on 0.0.0.0:8080
}
