package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		context.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		context.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		context.Next()
	}
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.Method != "OPTIONS" {
			context.Next()
		} else {
			context.Header("Access-Control-Allow-Origin", "*")
			context.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			context.Header("Access-Control-Allow-Headers", "authorization, origin, content-types, accept")
			context.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			context.Header("Content-Type", "application/json")
			context.AbortWithStatus(200)
		}
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		//c.Header("X-Frame-Options", "DENY")
		context.Header("X-Content-Type-Options", "nosniff")
		context.Header("X-XSS-Protection", "1; mode=block")
		if context.Request.TLS != nil {
			context.Header("Strict-Transport-Security", "max-age=31536000")
		}
	}
}
