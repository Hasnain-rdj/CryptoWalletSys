package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

// Visitor tracks requests from a single IP
type Visitor struct {
	count      int
	lastReset  time.Time
	lastAccess time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     requestsPerMinute,
		window:   time.Minute,
	}

	// Cleanup old visitors every 5 minutes
	go rl.cleanupVisitors()

	return rl
}

// RateLimit middleware
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		visitor, exists := rl.visitors[ip]

		if !exists {
			visitor = &Visitor{
				count:      1,
				lastReset:  time.Now(),
				lastAccess: time.Now(),
			}
			rl.visitors[ip] = visitor
			rl.mu.Unlock()
			c.Next()
			return
		}

		// Reset count if window has passed
		if time.Since(visitor.lastReset) > rl.window {
			visitor.count = 1
			visitor.lastReset = time.Now()
			visitor.lastAccess = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}

		// Check if rate limit exceeded
		if visitor.count >= rl.rate {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		visitor.count++
		visitor.lastAccess = time.Now()
		rl.mu.Unlock()

		c.Next()
	}
}

// cleanupVisitors removes inactive visitors
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			// Remove visitors inactive for more than 10 minutes
			if time.Since(visitor.lastAccess) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
