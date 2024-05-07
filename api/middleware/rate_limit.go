package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func NewRateLimiter(t time.Duration, limit int64) gin.HandlerFunc {
	// Create a rate with the given limit (number of requests) for the given
	// period (a time.Duration of your choice).
	rate := limiter.Rate{
		Period: t,
		Limit:  limit,
	}

	store := memory.NewStore()

	// Then, create the limiter instance which takes the store and the rate as arguments.
	// Now, you can give this instance to any supported middleware.
	instance := limiter.New(store, rate)

	return mgin.NewMiddleware(instance)
}
