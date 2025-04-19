package middleware

import (
	"fmt"
	"oauth-go/internal/store"

	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (*store.User, bool) {
	val, ok := c.Get(contextUserKey)
	if !ok {
		return nil, false
	}

	user, ok := val.(*store.User)
	return user, ok
}

func MustGetUserFromContext(c *gin.Context) (*store.User, error) {
	user, ok := GetUserFromContext(c)

	if !ok {
		return nil, fmt.Errorf("cannot get user from request")
	}

	return user, nil
}
