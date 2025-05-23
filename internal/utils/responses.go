package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorResponse defines the standard structure for error responses returned from the API.
//
// Fields:
//   - Code: an optional machine-readable error code for client-side handling
//   - Error: a human-readable error message
//   - Details: optional additional information helpful for debugging or context
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`    // Optional error code for frontend
	Error   string `json:"error"`             // Human-readable error message
	Details string `json:"details,omitempty"` // Optional detailed explanation
}

// ErrorJSON sends a structured JSON error response using the provided HTTP status code,
// error message, and optional details. It also logs the error server-side.
//
// Parameters:
//   - ctx: the Gin context used to send the response
//   - httpCode: the HTTP status code to return (e.g. 400, 404, 500)
//   - message: a human-readable error message
//   - details: optional additional context to include in the response and logs
func ErrorJSON(ctx *gin.Context, httpCode int, message string, details ...string) {
	resp := ErrorResponse{
		Error: message,
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}

	// Optional logging for server-side debugging
	log.Printf("[ERROR] %s - %s", message, resp.Details)

	ctx.JSON(httpCode, resp)
}

// ErrorWithCodeJSON sends a structured JSON error response including a custom error code,
// HTTP status code, and optional detailed context. Logs are emitted for traceability.
//
// Parameters:
//   - ctx: the Gin context used to send the response
//   - httpCode: the HTTP status code to return (e.g. 400, 404, 500)
//   - errorCode: a custom machine-readable error code (e.g. "VALIDATION_ERROR")
//   - message: a human-readable error message
//   - details: optional additional context to include in the response and logs
func ErrorWithCodeJSON(ctx *gin.Context, httpCode int, errorCode, message string, details ...string) {
	resp := ErrorResponse{
		Code:  errorCode,
		Error: message,
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}

	log.Printf("[ERROR] (%s) %s - %s", errorCode, message, resp.Details)

	ctx.JSON(httpCode, resp)
}
