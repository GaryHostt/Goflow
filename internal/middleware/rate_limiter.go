package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter manages rate limits per tenant
// MULTI-TENANT: Different tiers get different limits
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	
	// Configuration
	freeLimit  rate.Limit // requests per second (e.g., 5)
	paidLimit  rate.Limit // requests per second (e.g., 50)
	burstSize  int        // burst capacity
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(freeLimit, paidLimit float64, burstSize int) *RateLimiter {
	return &RateLimiter{
		limiters:  make(map[string]*rate.Limiter),
		freeLimit: rate.Limit(freeLimit),
		paidLimit: rate.Limit(paidLimit),
		burstSize: burstSize,
	}
}

// getLimiter returns or creates a rate limiter for a tenant
func (rl *RateLimiter) getLimiter(tenantID string, isPaid bool) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[tenantID]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	// Create new limiter
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists := rl.limiters[tenantID]; exists {
		return limiter
	}

	// Determine limit based on tier
	limit := rl.freeLimit
	if isPaid {
		limit = rl.paidLimit
	}

	limiter = rate.NewLimiter(limit, rl.burstSize)
	rl.limiters[tenantID] = limiter
	return limiter
}

// RateLimitMiddleware enforces rate limits per tenant
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract tenant ID from context (set by AuthMiddleware)
		tenantID, ok := GetTenantIDFromContext(r.Context())
		if !ok {
			// No tenant ID, allow (public endpoints)
			next.ServeHTTP(w, r)
			return
		}

		// TODO: MULTI-TENANT - Query tenant tier from database
		// For now, assume all tenants are free tier
		isPaid := false

		// Get limiter for this tenant
		limiter := rl.getLimiter(tenantID, isPaid)

		// Check if request is allowed
		if !limiter.Allow() {
			// Rate limit exceeded
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-RateLimit-Limit", "5")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"success":false,"error":"Rate limit exceeded. Please try again later."}`))
			return
		}

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", "5")
		// Note: Getting remaining tokens requires additional logic

		next.ServeHTTP(w, r)
	})
}

// CleanupOldLimiters removes inactive limiters (memory optimization)
func (rl *RateLimiter) CleanupOldLimiters() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// In production, track last access time and remove inactive limiters
			// For simplicity, we keep all limiters (small memory footprint)
			rl.mu.Unlock()
		}
	}()
}

