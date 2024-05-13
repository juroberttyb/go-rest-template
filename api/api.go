// Package api defines all routes and handlers
package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"github.com/A-pen-app/kickstart/api/middleware"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/database"
	"github.com/A-pen-app/kickstart/global"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/service"
	"github.com/A-pen-app/kickstart/store"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/mq"
)

// NewRouter returns the global HTTP router instance.
func NewRouter() *gin.Engine {
	return initializeRouter()
}

// initializeSingletons is the function called by sync.Once to intialize the
// HTTP engine and router group singleton instances.
func initializeRouter() *gin.Engine {
	var router *gin.Engine
	// Create router and group instances. Check whether we should use the
	// microservice name as root router group URL prefix. This depends on
	// whether or our Kubernetes ingress is configured to use path-based
	// routing or name-based virtual hosting.
	if config.GetBool("SERVICE_NAME_AS_ROOT") {
		router = createRouterAndGroup(global.ServiceName)
	} else {
		router = createRouterAndGroup("")
	}
	return router
}

// @title					Order (aka Broadcast Service) API
// @description				This service provides message broadcasting service for official accounts,
// @description				it also supports audience filtering and performance reporting.
// @version					v1
// @host					localhost:8000
// @basePath					/
// @schemes					http
// @securityDefinitions.apiKey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT.
func createRouterAndGroup(prefix string) *gin.Engine {
	ctx := context.Background()

	// Create a clean HTTP router engine.
	engine := gin.New()

	// Configure HTTP router engine settings.
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = false
	engine.HandleMethodNotAllowed = true
	engine.ForwardedByClientIP = true
	// engine.MaxMultipartMemory = maxUploadFileSize

	// Create from the engine a router group with the given prefix.
	root := engine.Group(prefix)

	// Install common middleware to the router group.
	installCommonMiddleware(root)

	// initialize dependencies for injection
	db := database.GetPostgres()

	pubsub := mq.GetPubsub()
	cryptoStore := store.NewCrypto(ctx)
	orderStore := store.NewOrder(db)

	authSvc := service.NewAuth(ctx, cryptoStore)
	orderSvc := service.NewOrder(orderStore, pubsub)

	// register routes
	addDocRoutes(root)
	addProbesRoutes(root)
	addSystemRoutes(root)
	addOrderRoutes(root, orderSvc, authSvc)

	return engine
}

// installCommonMiddleware installs common middleware to the router group.
func installCommonMiddleware(root *gin.RouterGroup) {
	// support open tracing
	root.Use(otelgin.Middleware("kickstart-api"))

	// Install logger middleware, a middleware to log requests.
	root.Use(logging.RequestLogger([]string{"/alive", "/ready"}))

	// Install recovery middleware, a middleware to recover & log panics.
	// NOTE: The recovery middleware should always be the last one installed.
	root.Use(middleware.Recovery())
}

type errorResp struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
}

type pageReq struct {
	Next  string `form:"next" default:"" validate:"optional"`    // next cursor value, use it when requesting next page
	Count int    `form:"count" default:"10" validate:"optional"` // number of elements requested
}

type pageResp struct {
	Data interface{} `json:"data"`
	Next string      `json:"next" example:"next cursor value"`
}

// helper function for translating internal errors to http status code
func handleError(ctx *gin.Context, err error) {
	logging.Error(ctx.Request.Context(), err.Error())

	// we use trace id as request id
	requestID := trace.SpanContextFromContext(ctx.Request.Context()).TraceID().String()
	resp := &errorResp{
		Error:     err.Error(),
		RequestID: requestID,
	}
	switch err {
	case models.ErrorNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, resp)
	case models.ErrorWrongParams:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
	case models.ErrorUnsupported:
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, resp)
	case models.ErrorNotAllowed:
		ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
	}
}
