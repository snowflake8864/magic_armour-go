package token_bucket

import (
	"sync"
	"time"
)

// TokenBucket 令牌桶
type TokenBucket struct {
	Capacity  int64      // 桶的容量，即最多能存储多少令牌
	Rate      float64    // 令牌放入速率，单位为每秒
	Tokens    float64    // 当前桶中令牌的数量
	LastToken time.Time  // 上一次向桶中放令牌的时间
	Mtx       sync.Mutex // 互斥锁，确保并发安全
}

func (tb *TokenBucket) Allow() bool {
	tb.Mtx.Lock()
	defer tb.Mtx.Unlock()
	now := time.Now()
	// 计算需要放的令牌数量
	tb.Tokens = tb.Tokens + tb.Rate*now.Sub(tb.LastToken).Seconds()
	if tb.Tokens > float64(tb.Capacity) {
		tb.Tokens = float64(tb.Capacity)
	}
	// 判断是否允许请求
	if tb.Tokens >= 1 {
		tb.Tokens--
		tb.LastToken = now
		return true
	} else {
		return false
	}
}
