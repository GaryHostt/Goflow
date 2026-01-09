package engine

import (
	"sync"
	"time"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState string

const (
	StateClosed   CircuitBreakerState = "closed"   // Normal operation
	StateOpen     CircuitBreakerState = "open"     // Failing, reject requests
	StateHalfOpen CircuitBreakerState = "half_open" // Testing if service recovered
)

// CircuitBreaker implements the Circuit Breaker pattern for connectors
// Prevents cascading failures when third-party APIs are down
type CircuitBreaker struct {
	mu sync.RWMutex

	// Configuration
	maxFailures   int           // Failures before opening circuit
	timeout       time.Duration // How long to wait before attempting recovery
	halfOpenMax   int           // Max attempts in half-open state

	// State
	state           CircuitBreakerState
	failures        int
	lastFailureTime time.Time
	halfOpenAttempts int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
		halfOpenMax: 3,
		state:       StateClosed,
	}
}

// Call executes a function with circuit breaker protection
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if circuit is open
	if cb.state == StateOpen {
		// Check if timeout has elapsed
		if time.Since(cb.lastFailureTime) > cb.timeout {
			// Move to half-open state (testing recovery)
			cb.state = StateHalfOpen
			cb.halfOpenAttempts = 0
		} else {
			// Circuit still open, reject immediately
			return &CircuitBreakerError{
				State:   StateOpen,
				Message: "Circuit breaker is open, service unavailable",
			}
		}
	}

	// Execute the function
	err := fn()

	if err != nil {
		// Record failure
		cb.recordFailure()
		return err
	}

	// Record success
	cb.recordSuccess()
	return nil
}

// recordFailure handles a failed request
func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.state == StateHalfOpen {
		// Failed during recovery test, reopen circuit
		cb.state = StateOpen
		cb.halfOpenAttempts = 0
	} else if cb.failures >= cb.maxFailures {
		// Too many failures, open circuit
		cb.state = StateOpen
	}
}

// recordSuccess handles a successful request
func (cb *CircuitBreaker) recordSuccess() {
	if cb.state == StateHalfOpen {
		cb.halfOpenAttempts++
		// If enough successful attempts, close circuit
		if cb.halfOpenAttempts >= cb.halfOpenMax {
			cb.state = StateClosed
			cb.failures = 0
			cb.halfOpenAttempts = 0
		}
	} else {
		// Reset failure count on success
		cb.failures = 0
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailures returns the current failure count
func (cb *CircuitBreaker) GetFailures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// Reset manually resets the circuit breaker (admin action)
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenAttempts = 0
}

// CircuitBreakerError represents a circuit breaker error
type CircuitBreakerError struct {
	State   CircuitBreakerState
	Message string
}

func (e *CircuitBreakerError) Error() string {
	return e.Message
}

// CircuitBreakerManager manages circuit breakers for all connectors
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetBreaker returns or creates a circuit breaker for a connector
func (m *CircuitBreakerManager) GetBreaker(connectorKey string) *CircuitBreaker {
	m.mu.RLock()
	breaker, exists := m.breakers[connectorKey]
	m.mu.RUnlock()

	if exists {
		return breaker
	}

	// Create new breaker
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if breaker, exists := m.breakers[connectorKey]; exists {
		return breaker
	}

	// Create new circuit breaker
	// 5 failures → open circuit for 60 seconds
	breaker = NewCircuitBreaker(5, 60*time.Second)
	m.breakers[connectorKey] = breaker
	return breaker
}

// GetAllStates returns the state of all circuit breakers
func (m *CircuitBreakerManager) GetAllStates() map[string]CircuitBreakerState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	states := make(map[string]CircuitBreakerState)
	for key, breaker := range m.breakers {
		states[key] = breaker.GetState()
	}
	return states
}

// ResetBreaker manually resets a specific circuit breaker
func (m *CircuitBreakerManager) ResetBreaker(connectorKey string) {
	m.mu.RLock()
	breaker, exists := m.breakers[connectorKey]
	m.mu.RUnlock()

	if exists {
		breaker.Reset()
	}
}

