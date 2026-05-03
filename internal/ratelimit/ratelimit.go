package ratelimit

import (
	"context"
	"net"
	"net/http"
	"time"

	"sync"

	"golang.org/x/time/rate"
)

type Limiter struct {
	limiters map[string]*entry
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func New(r rate.Limit, b int) *Limiter {
	return &Limiter{
		limiters: make(map[string]*entry),
		rate:     r,
		burst:    b,
	}
}

func (l *Limiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.limiters[ip]
	if !ok {
		e = &entry{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.limiters[ip] = e
	}
	e.lastSeen = time.Now()

	return e.limiter
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid remote address", http.StatusInternalServerError)
			return
		}

		limiter := l.getLimiter(host)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (l *Limiter) Cleanup(ctx context.Context, interval, ttl time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.evictOlderThan(ttl)
		}
	}
}

func (l *Limiter) evictOlderThan(ttl time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-ttl)
	for ip, e := range l.limiters {
		if e.lastSeen.Before(cutoff) {
			delete(l.limiters, ip)
		}
	}
}
