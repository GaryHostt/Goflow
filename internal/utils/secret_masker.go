package utils

import (
	"regexp"
	"strings"
)

// SecretMasker handles sanitization of sensitive data in logs
// CRITICAL: Prevents API keys, tokens, and credentials from appearing in logs
type SecretMasker struct {
	patterns []*regexp.Regexp
}

// NewSecretMasker creates a new secret masker
func NewSecretMasker() *SecretMasker {
	return &SecretMasker{
		patterns: []*regexp.Regexp{
			// Slack webhook URLs
			regexp.MustCompile(`https://hooks\.slack\.com/services/[A-Z0-9]+/[A-Z0-9]+/[A-Za-z0-9]+`),
			
			// Discord webhook URLs
			regexp.MustCompile(`https://discord\.com/api/webhooks/[0-9]+/[A-Za-z0-9_-]+`),
			
			// Generic API keys (various formats)
			regexp.MustCompile(`[aA][pP][iI]_?[kK][eE][yY][\s:=]+['"]*([A-Za-z0-9_\-]{20,})['"]*`),
			regexp.MustCompile(`[aA][pP][iI][-_]?[tT][oO][kK][eE][nN][\s:=]+['"]*([A-Za-z0-9_\-]{20,})['"]*`),
			
			// Bearer tokens
			regexp.MustCompile(`[bB]earer\s+([A-Za-z0-9_\-\.]{20,})`),
			
			// Authorization headers
			regexp.MustCompile(`[aA]uthorization[\s:=]+['"]*([A-Za-z0-9_\-\.]{20,})['"]*`),
			
			// Password patterns
			regexp.MustCompile(`[pP]assword[\s:=]+['"]*([^'"\s]{6,})['"]*`),
			regexp.MustCompile(`[pP]ass[\s:=]+['"]*([^'"\s]{6,})['"]*`),
			
			// Secret patterns
			regexp.MustCompile(`[sS]ecret[\s:=]+['"]*([A-Za-z0-9_\-]{20,})['"]*`),
			
			// AWS keys
			regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
			
			// JWT tokens (basic pattern)
			regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`),
			
			// Email in JSON (for PII protection)
			regexp.MustCompile(`"email"[\s:=]+"[^"]+@[^"]+"`),
			
			// Credit card numbers (basic pattern)
			regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
		},
	}
}

// Mask replaces sensitive data with [REDACTED]
func (sm *SecretMasker) Mask(input string) string {
	masked := input

	for _, pattern := range sm.patterns {
		masked = pattern.ReplaceAllStringFunc(masked, func(match string) string {
			// Keep first few characters for debugging context
			if len(match) > 10 {
				return match[:4] + "***REDACTED***"
			}
			return "***REDACTED***"
		})
	}

	return masked
}

// MaskMap sanitizes a map of data (useful for JSON payloads)
func (sm *SecretMasker) MaskMap(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range data {
		keyLower := strings.ToLower(key)

		// Check if key name indicates sensitive data
		if sm.isSensitiveKey(keyLower) {
			sanitized[key] = "***REDACTED***"
			continue
		}

		// Recursively sanitize nested maps
		switch v := value.(type) {
		case string:
			sanitized[key] = sm.Mask(v)
		case map[string]interface{}:
			sanitized[key] = sm.MaskMap(v)
		case []interface{}:
			sanitized[key] = sm.maskArray(v)
		default:
			sanitized[key] = value
		}
	}

	return sanitized
}

// maskArray sanitizes an array of data
func (sm *SecretMasker) maskArray(arr []interface{}) []interface{} {
	sanitized := make([]interface{}, len(arr))

	for i, item := range arr {
		switch v := item.(type) {
		case string:
			sanitized[i] = sm.Mask(v)
		case map[string]interface{}:
			sanitized[i] = sm.MaskMap(v)
		case []interface{}:
			sanitized[i] = sm.maskArray(v)
		default:
			sanitized[i] = item
		}
	}

	return sanitized
}

// isSensitiveKey checks if a key name indicates sensitive data
func (sm *SecretMasker) isSensitiveKey(key string) bool {
	sensitiveKeys := []string{
		"password",
		"passwd",
		"pass",
		"secret",
		"api_key",
		"apikey",
		"api-key",
		"token",
		"access_token",
		"refresh_token",
		"auth",
		"authorization",
		"credential",
		"encrypted_key",
		"webhook_url",
		"private_key",
		"client_secret",
		"session_id",
		"session",
		"cookie",
	}

	for _, sensitive := range sensitiveKeys {
		if strings.Contains(key, sensitive) {
			return true
		}
	}

	return false
}

// MaskURL removes sensitive parts from URLs (query params, credentials)
func (sm *SecretMasker) MaskURL(url string) string {
	// Mask credentials in URLs like https://user:pass@example.com
	credPattern := regexp.MustCompile(`://([^:]+):([^@]+)@`)
	url = credPattern.ReplaceAllString(url, "://***:***@")

	// Mask query parameters that look sensitive
	queryPattern := regexp.MustCompile(`([?&])(api_key|token|secret|password|key)=([^&]+)`)
	url = queryPattern.ReplaceAllString(url, "$1$2=***REDACTED***")

	return url
}

// Global secret masker instance
var globalMasker = NewSecretMasker()

// Mask is a convenience function using the global masker
func Mask(input string) string {
	return globalMasker.Mask(input)
}

// MaskMap is a convenience function using the global masker
func MaskMap(data map[string]interface{}) map[string]interface{} {
	return globalMasker.MaskMap(data)
}

// MaskURL is a convenience function using the global masker
func MaskURL(url string) string {
	return globalMasker.MaskURL(url)
}

