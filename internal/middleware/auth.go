package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexmacdonald/simple-ipass/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is a custom type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// TenantIDKey is the context key for tenant ID (MULTI-TENANT READY!)
	TenantIDKey ContextKey = "tenant_id"
)

var jwtSecret = []byte("ipaas-jwt-secret-change-in-production")

// SetJWTSecret sets the JWT secret (should be called on startup)
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// GetJWTSecret returns the JWT secret
func GetJWTSecret() []byte {
	return jwtSecret
}

// AuthMiddleware validates JWT tokens and extracts user_id and tenant_id
func AuthMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("Missing Authorization header", map[string]interface{}{
					"path":   r.URL.Path,
					"method": r.Method,
				})
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Expected format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Warn("Invalid Authorization header format", map[string]interface{}{
					"path":   r.URL.Path,
					"header": authHeader,
				})
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				log.Warn("Invalid or expired token", map[string]interface{}{
					"path":  r.URL.Path,
					"error": err.Error(),
				})
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Error("Invalid token claims", map[string]interface{}{
					"path": r.URL.Path,
				})
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Extract user_id (required)
			userID, ok := claims["user_id"].(string)
			if !ok {
				log.Error("Missing user_id in token", map[string]interface{}{
					"path": r.URL.Path,
				})
				http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
				return
			}

			// Extract tenant_id (optional for now, required in multi-tenant phase)
			tenantID, _ := claims["tenant_id"].(string)
			if tenantID == "" {
				// MIGRATION PHASE: For backwards compatibility, derive tenant from user
				// In Phase 1 (current): Each user is their own tenant
				// In Phase 2 (multi-tenant): This would come from JWT
				tenantID = "tenant_" + userID
			}

			// Log successful authentication with context
			log.InfoWithContext("Request authenticated", userID, tenantID, map[string]interface{}{
				"path":   r.URL.Path,
				"method": r.Method,
			})

			// Add both user_id and tenant_id to request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, TenantIDKey, tenantID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts user_id from request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetTenantIDFromContext extracts tenant_id from request context (MULTI-TENANT READY!)
func GetTenantIDFromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	return tenantID, ok
}

// GetUserAndTenantFromContext extracts both IDs (convenience method)
func GetUserAndTenantFromContext(ctx context.Context) (userID, tenantID string, ok bool) {
	userID, ok1 := GetUserIDFromContext(ctx)
	tenantID, ok2 := GetTenantIDFromContext(ctx)
	return userID, tenantID, ok1 && ok2
}
