package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type KeyLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor

	rate     rate.Limit
	burst    int
	idleTTL  time.Duration
	maxKeys  int
	stopOnce sync.Once
	stop     chan struct{}
}

type Options struct {
	Rate            rate.Limit
	Burst           int
	IdleTTL         time.Duration
	CleanupInterval time.Duration
	MaxKeys         int
}

func (o Options) withDefaults() Options {
	if o.Rate <= 0 {
		o.Rate = 10
	}
	if o.Burst <= 0 {
		o.Burst = 20
	}
	if o.IdleTTL <= 0 {
		o.IdleTTL = 10 * time.Minute
	}
	if o.CleanupInterval <= 0 {
		o.CleanupInterval = time.Minute
	}
	if o.MaxKeys <= 0 {
		o.MaxKeys = 100_000
	}
	return o
}

func New(opts Options) *KeyLimiter {
	opts = opts.withDefaults()

	l := &KeyLimiter{
		visitors: make(map[string]*visitor),
		rate:     opts.Rate,
		burst:    opts.Burst,
		idleTTL:  opts.IdleTTL,
		maxKeys:  opts.MaxKeys,
		stop:     make(chan struct{}),
	}

	go l.cleanupLoop(opts.CleanupInterval)
	return l
}

func (l *KeyLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, ok := l.visitors[key]
	if !ok {
		if len(l.visitors) >= l.maxKeys {
			l.evictOldestLocked()
		}
		v = &visitor{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.visitors[key] = v
	}
	v.lastSeen = time.Now()

	return v.limiter.Allow()
}

func (l *KeyLimiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.visitors)
}

func (l *KeyLimiter) Close() {
	l.stopOnce.Do(func() { close(l.stop) })
}

func (l *KeyLimiter) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			l.cleanup()
		}
	}
}

func (l *KeyLimiter) cleanup() {
	cutoff := time.Now().Add(-l.idleTTL)

	l.mu.Lock()
	defer l.mu.Unlock()

	for key, v := range l.visitors {
		if v.lastSeen.Before(cutoff) {
			delete(l.visitors, key)
		}
	}
}

func (l *KeyLimiter) evictOldestLocked() {
	var oldestKey string
	var oldestTime time.Time

	for key, v := range l.visitors {
		if oldestKey == "" || v.lastSeen.Before(oldestTime) {
			oldestKey, oldestTime = key, v.lastSeen
		}
	}
	if oldestKey != "" {
		delete(l.visitors, oldestKey)
	}
}