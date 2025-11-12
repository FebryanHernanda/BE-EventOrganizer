package middleware

import (
	"sync"
	"time"

	"github.com/FebryanHernanda/BE-EventOrganizer/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	visitors map[string]*visitor
	mu       *sync.Mutex
	rate     rate.Limit
	burst    int
}

type MiddlewareDeps struct {
	GlobalLimiter *RateLimiterStore
	AuthLimiter   *RateLimiterStore
	UserLimiter   *RateLimiterStore
}

func NewRateLimiterStore(r rate.Limit, b int) *RateLimiterStore {
	return &RateLimiterStore{
		visitors: make(map[string]*visitor),
		mu:       &sync.Mutex{},
		rate:     r,
		burst:    b,
	}
}

func (s *RateLimiterStore) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		s.mu.Lock()
		for ip, v := range s.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(s.visitors, ip)
			}
		}
		s.mu.Unlock()
	}
}

func (s *RateLimiterStore) getVisitor(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(s.rate, s.burst)
		v = &visitor{limiter, time.Now()}
		s.visitors[ip] = v
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (s *RateLimiterStore) RateLimiterMiddleware() gin.HandlerFunc {
	go s.cleanupVisitors()

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := s.getVisitor(ip)

		if !limiter.Allow() {
			ctx.Header("Retry-After", "60")
			logrus.WithFields(logrus.Fields{
				"ip":   ip,
				"path": ctx.FullPath(),
			}).Warn("Rate limit exceeded")
			response.Error(ctx, "too many request, try again later", 429, nil)
			ctx.Abort()
			return
		}

		logrus.WithFields(logrus.Fields{
			"ip":   ip,
			"path": ctx.FullPath(),
		}).Info("Request allowed")

		ctx.Next()
	}
}
