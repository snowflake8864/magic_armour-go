package rate_limit

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/token_bucket"
	"time"
)

func LimitHandler(rate float64, maxConn int64) gin.HandlerFunc {
	tb := &token_bucket.TokenBucket{
		Capacity:  maxConn, //100
		Rate:      rate,    //1.0,
		Tokens:    0,
		LastToken: time.Now(),
	}
	return func(c *gin.Context) {
		if !tb.Allow() {
			c.String(503, "Too many request")
			c.Abort()
			return
		}
		c.Next()
	}
}
