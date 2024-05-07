package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS ...
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if ctx.Request.Method != "OPTIONS" {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
