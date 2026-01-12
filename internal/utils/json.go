package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	// ErrRequestBodyTooLarge is returned when request body exceeds max size
	ErrRequestBodyTooLarge = errors.New("request body too large")
	
	// ErrMalformedJSON is returned when JSON is invalid
	ErrMalformedJSON = errors.New("malformed JSON")
	
	// ErrUnknownFields is returned when unknown fields are present
	ErrUnknownFields = errors.New("unknown fields in request")
)

const (
	// MaxRequestBodySize is the maximum allowed request body size (1MB)
	MaxRequestBodySize = 1_048_576
)

// DecodeJSONStrict decodes JSON from request body with strict validation
// This prevents:
// - Large payloads that could exhaust memory
// - Unknown fields that might indicate API misuse
// - Malformed JSON that could cause panics
func DecodeJSONStrict(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Limit request body size to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)
	
	// Create decoder with strict settings
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Reject unknown fields
	
	// Attempt to decode
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		
		switch {
		// Catch syntax errors
		case errors.As(err, &syntaxError):
			return fmt.Errorf("%w at byte offset %d", ErrMalformedJSON, syntaxError.Offset)
		
		// Catch type errors
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("%w: field %q has wrong type", ErrMalformedJSON, unmarshalTypeError.Field)
			}
			return fmt.Errorf("%w: body contains incorrect JSON type", ErrMalformedJSON)
		
		// Catch EOF errors
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("%w: incomplete JSON", ErrMalformedJSON)
		
		// Catch unknown field errors
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("%w: %s", ErrUnknownFields, fieldName)
		
		// Catch body too large
		case err.Error() == "http: request body too large":
			return ErrRequestBodyTooLarge
		
		// Catch invalid unmarshal error (programming error)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		
		// Generic error
		default:
			return err
		}
	}
	
	// Ensure only one JSON value in body
	if dec.More() {
		return fmt.Errorf("%w: body must only contain a single JSON value", ErrMalformedJSON)
	}
	
	return nil
}

// WriteJSONError writes a JSON error response
func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// WriteJSON writes a JSON success response
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

