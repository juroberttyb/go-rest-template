package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/A-pen-app/kickstart/global"
)

func addProbesRoutes(root *gin.RouterGroup) {
	// Register the liveness/readiness probe handlers.
	root.GET("alive", alive)
	root.GET("ready", ready)

}

// alive is the handler for Kubernetes liveness probes.
func alive(ctx *gin.Context) {
	// Set status code based on liveness indication flag.
	statusCode := http.StatusServiceUnavailable
	if global.Alive {
		statusCode = http.StatusOK
	}

	// Respond to probe according to current liveness status.
	ctx.JSON(statusCode, gin.H{
		"alive": global.Alive,
	})
}

// ready is the handler for Kubernetes readiness probes.
func ready(ctx *gin.Context) {
	// Set status code based on readiness indication flag.
	statusCode := http.StatusServiceUnavailable
	if global.Ready {
		statusCode = http.StatusOK
	}

	// Respond to probe according to current readiness status.
	ctx.JSON(statusCode, gin.H{
		"ready": global.Ready,
	})
}
