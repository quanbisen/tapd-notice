package middleware

import (
	"github.com/gin-gonic/gin"
	"tapd-notice/common/orm"
)

func ContextDB() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("db", orm.DB)
		context.Next()
	}
}
