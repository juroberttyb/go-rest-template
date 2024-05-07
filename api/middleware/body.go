package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxBodySizeInMB(maxMBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, maxMBytes*1024*1024)
		c.Next()
	}
}
