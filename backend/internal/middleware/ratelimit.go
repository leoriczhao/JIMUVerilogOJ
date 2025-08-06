package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶结构
type TokenBucket struct {
	capacity    int64     // 桶容量
	tokens      int64     // 当前令牌数
	refillRate  int64     // 每秒补充令牌数
	lastRefill  time.Time // 上次补充时间
	mu          sync.Mutex
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(capacity, refillRate int64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// TryConsume 尝试消费一个令牌
func (tb *TokenBucket) TryConsume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	// 计算需要补充的令牌数
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int64(elapsed.Seconds()) * tb.refillRate

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}

	// 尝试消费令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// RateLimiter 速率限制器
type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
	capacity    int64 // 桶容量
	refillRate  int64 // 每秒补充令牌数
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(capacity, refillRate int64) *RateLimiter {
	return &RateLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

// GetBucket 获取或创建令牌桶
func (rl *RateLimiter) GetBucket(key string) *TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.buckets[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// 双重检查
		if bucket, exists = rl.buckets[key]; !exists {
			bucket = NewTokenBucket(rl.capacity, rl.refillRate)
			rl.buckets[key] = bucket
		}
		rl.mu.Unlock()
	}
	return bucket
}

// RateLimit 速率限制中间件
// capacity: 令牌桶容量（最大突发请求数）
// refillRate: 每秒补充令牌数（稳定请求速率）
func RateLimit(capacity, refillRate int64) gin.HandlerFunc {
	limiter := NewRateLimiter(capacity, refillRate)

	return gin.HandlerFunc(func(c *gin.Context) {
		// 使用客户端IP作为限制键
		clientIP := c.ClientIP()
		bucket := limiter.GetBucket(clientIP)

		if !bucket.TryConsume() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "请求过于频繁，请稍后再试",
				"retry_after": "1s",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// APIRateLimit API专用速率限制中间件
// 针对API接口的默认配置：每个IP每秒最多10个请求，突发最多20个请求
func APIRateLimit() gin.HandlerFunc {
	return RateLimit(20, 10) // 容量20，每秒补充10个令牌
}

// StrictRateLimit 严格速率限制中间件
// 针对敏感接口的严格配置：每个IP每秒最多2个请求，突发最多5个请求
func StrictRateLimit() gin.HandlerFunc {
	return RateLimit(5, 2) // 容量5，每秒补充2个令牌
}