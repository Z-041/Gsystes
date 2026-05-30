package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/auth"
)

const (
	ContextKeyClaims = "claims"
)

func SetClaims(c *gin.Context, claims *auth.Claims) {
	c.Set(ContextKeyClaims, claims)
}

func GetClaims(c *gin.Context) *auth.Claims {
	value, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil
	}
	claims, ok := value.(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserID(c *gin.Context) uint {
	claims := GetClaims(c)
	if claims == nil {
		return 0
	}
	return claims.UserID
}

func GetUsername(c *gin.Context) string {
	claims := GetClaims(c)
	if claims == nil {
		return ""
	}
	return claims.Username
}
