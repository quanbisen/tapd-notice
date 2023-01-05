package middleware

import "github.com/gin-gonic/gin"

func InitMiddleware(services map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = services
		c.Next()
	}
}
