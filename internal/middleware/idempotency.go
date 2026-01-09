package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

// IdempotencyResult represents a cached result
type IdempotencyResult struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Timestamp  time.Time
}

// IdempotencyManager manages idempotency keys to prevent duplicate operations
// Solves the "double-click" problem in distributed systems
type IdempotencyManager struct {
	cache map[string]*IdempotencyResult
	mu    sync.RWMutex
	ttl   time.Duration // How long to cache results
}

// NewIdempotencyManager creates a new idempotency manager
func NewIdempotencyManager(ttl time.Duration) *IdempotencyManager {
	im := &IdempotencyManager{
		cache: make(map[string]*IdempotencyResult),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go im.cleanup()
	
	return im
}

// cleanup removes expired entries
func (im *IdempotencyManager) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		im.mu.Lock()
		now := time.Now()
		for key, result := range im.cache {
			if now.Sub(result.Timestamp) > im.ttl {
				delete(im.cache, key)
			}
		}
		im.mu.Unlock()
	}
}

// Get retrieves a cached result
func (im *IdempotencyManager) Get(key string) (*IdempotencyResult, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	
	result, exists := im.cache[key]
	if !exists {
		return nil, false
	}
	
	// Check if expired
	if time.Since(result.Timestamp) > im.ttl {
		return nil, false
	}
	
	return result, true
}

// Set caches a result
func (im *IdempotencyManager) Set(key string, result *IdempotencyResult) {
	im.mu.Lock()
	defer im.mu.Unlock()
	
	im.cache[key] = result
}

// GenerateKey generates an idempotency key from request details
func (im *IdempotencyManager) GenerateKey(method, path, body string) string {
	h := sha256.New()
	h.Write([]byte(method + path + body))
	return hex.EncodeToString(h.Sum(nil))
}

// IdempotencyMiddleware provides idempotency for POST/PUT/PATCH requests
func (im *IdempotencyManager) IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to mutating methods
		if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
			next.ServeHTTP(w, r)
			return
		}

		// Check for idempotency key header
		idempotencyKey := r.Header.Get("X-Idempotency-Key")
		if idempotencyKey == "" {
			// No idempotency key provided, process normally
			next.ServeHTTP(w, r)
			return
		}

		// Check if we've seen this key before
		if result, exists := im.Get(idempotencyKey); exists {
			// Return cached result
			for key, values := range result.Headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Idempotency-Replay", "true")
			w.WriteHeader(result.StatusCode)
			w.Write(result.Body)
			return
		}

		// Create a response recorder to capture the result
		recorder := &ResponseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           []byte{},
		}

		// Process the request
		next.ServeHTTP(recorder, r)

		// Cache the result
		result := &IdempotencyResult{
			StatusCode: recorder.statusCode,
			Body:       recorder.body,
			Headers:    recorder.Header().Clone(),
			Timestamp:  time.Now(),
		}
		im.Set(idempotencyKey, result)
	})
}

// ResponseRecorder captures the response for caching
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rr *ResponseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *ResponseRecorder) Write(b []byte) (int, error) {
	rr.body = append(rr.body, b...)
	return rr.ResponseWriter.Write(b)
}

