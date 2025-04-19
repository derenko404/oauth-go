package middleware

import (
	"fmt"
	"go-auth/internal/services"
	"go-auth/internal/store"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserKey = "user"

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

func AuthMiddleware(store *store.Store, services *services.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the Authorization header
		tokenString, err := getTokenFromHeader(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Parse and validate the token
		token, err := services.Jwt.VerifyToken(tokenString)

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is not valid"})
			return
		}

		// Optionally, you can extract claims and set them in the context
		claims, err := services.Jwt.GetClaims(token)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to cast claims to CustomClaims"})
			return
		}

		filters := map[string]any{
			"id": claims.SessionID,
		}

		_, err = store.Session.GetSessionBy(ctx.Request.Context(), filters)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot find session"})
			return
		}

		filters = map[string]any{
			"id": claims.UserID,
		}

		user, err := store.User.GetUserBy(ctx.Request.Context(), filters)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot find user"})
			return
		}

		ctx.Set(ContextUserKey, user)

		ctx.Next()
	}
}
