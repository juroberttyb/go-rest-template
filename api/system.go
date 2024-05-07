package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/A-pen-app/kickstart/global"
)

func addSystemRoutes(root *gin.RouterGroup) {
	// Create router group for system module.
	systemGroup := root.Group("system")

	// Register the system module handlers.
	systemGroup.GET("version", version)
	systemGroup.GET("time", timenow)

}

// version is the handler for responding system version requests.
func version(ctx *gin.Context) {
	// Respond with the commit hash of this code and its build time.
	ctx.JSON(http.StatusOK, gin.H{
		"service": global.ServiceName,
		"commit":  global.GitCommitHash,
		"time":    global.BuildTime,
	})
}

// timenow is the handler for responding the current system time.
func timenow(ctx *gin.Context) {
	// Respond with the current system timestamp in milliseconds.
	timestamp := time.Now().UnixNano() / 10e6
	ctx.JSON(http.StatusOK, gin.H{
		"time": timestamp,
	})
}
