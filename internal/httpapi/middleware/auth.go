package middleware

import (
	"net/http"
	"strings"

	"kerjadekat/backend/pkg/token"

	"github.com/gin-gonic/gin"
)

const ContextClaimsKey = "auth_claims"

func JWTAuth(issuer *token.Issuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := strings.TrimSpace(c.GetHeader("Authorization"))
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimSpace(parts[1])
		claims, err := issuer.ParseAccess(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ContextClaimsKey, claims)
		c.Next()
	}
}

func Claims(c *gin.Context) (token.Claims, bool) {
	v, ok := c.Get(ContextClaimsKey)
	if !ok {
		return token.Claims{}, false
	}
	cl, ok := v.(token.Claims)
	return cl, ok
}
