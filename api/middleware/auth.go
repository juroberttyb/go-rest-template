package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/service"
	"github.com/gin-gonic/gin"
)

func AuthUser(a service.Auth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		if len(header) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(header, " ", 2)
		token := parts[0] // for backward compatible since old apps use 'Authorization: {{token}}'
		if len(parts) > 1 {
			switch parts[0] {
			case "Bearer":
				token = parts[1]
			default:
				// we don't support other authorization schemes
				ctx.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
		}

		c := ctx.Request.Context()

		// validate if this token is issued by us
		claims, err := a.ValidateToken(c, token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("userID", claims.Subject)
		ctx.Set("user_id", claims.Subject)
		ctx.Set("scope", claims.Scope)
		ctx.Set("aud", []string(claims.Audience))
		c = context.WithValue(c, "userID", ctx.GetString("userID"))
		c = context.WithValue(c, "user_id", ctx.GetString("user_id"))
		c = context.WithValue(c, "scope", ctx.Value("scope").(models.UserType).String())
		c = context.WithValue(c, "aud", ctx.Value("aud"))
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

// NeedPermission specifies required permission level of user
// should be called after AuthUser()
func NeedPermission(userType models.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, exists := ctx.Get("scope")
		if !exists {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		scope, ok := v.(models.UserType)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if scope < userType {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Next()
	}
}
