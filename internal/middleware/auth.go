package middleware

import (
	"fmt"
	"go-auth/internal/services"
	"go-auth/internal/store"
	"go-auth/pkg/response"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

const contextUserKey = "user"

func getTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization token")
	}

	// Check if the header is in the correct format
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("missing Authorization token")
	}

	tokenString := parts[1]

	return tokenString, nil
}

func AuthMiddleware(store *store.Store, services *services.Services, logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := getTokenFromHeader(ctx)
		if err != nil {
			logger.Debug("cannot get Authorization header", "error", err)
			ctx.AbortWithStatusJSON(response.ErrUnauthorized.Code, response.ErrUnauthorized)
			return
		}

		// Parse and validate the token
		token, err := services.Jwt.VerifyToken(tokenString)

		if err != nil || !token.Valid {
			logger.Debug("cannot validate token", "error", err)
			ctx.AbortWithStatusJSON(response.ErrUnauthorized.Code, response.ErrUnauthorized)
			return
		}

		// Optionally, you can extract claims and set them in the context
		claims, err := services.Jwt.GetClaims(token)

		if err != nil {
			logger.Debug("cannot get claims", "error", err)
			ctx.AbortWithStatusJSON(response.ErrUnauthorized.Code, response.ErrUnauthorized)
			return
		}

		filters := map[string]any{
			"id": claims.SessionID,
		}

		_, err = store.Session.GetSessionBy(ctx.Request.Context(), filters)
		if err != nil {
			logger.Debug("cannot get session", "error", err)
			ctx.AbortWithStatusJSON(response.ErrUnauthorized.Code, response.ErrUnauthorized)
			return
		}

		filters = map[string]any{
			"id": claims.UserID,
		}

		user, err := store.User.GetUserBy(ctx.Request.Context(), filters)

		if err != nil {
			logger.Debug("cannot get user", "error", err)
			ctx.AbortWithStatusJSON(response.ErrUnauthorized.Code, response.ErrUnauthorized)
			return
		}

		ctx.Set(contextUserKey, user)

		ctx.Next()
	}
}
