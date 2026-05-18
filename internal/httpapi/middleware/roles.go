package middleware

import (
	"net/http"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/token"

	"github.com/gin-gonic/gin"
)

// RequireRoles allows the request only when the JWT role matches one of the allowed roles.
func RequireRoles(allowed ...string) gin.HandlerFunc {
	set := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		set[r] = struct{}{}
	}
	return func(c *gin.Context) {
		cl, ok := Claims(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if _, ok := set[cl.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

// MustClaims returns JWT claims set by JWTAuth middleware.
func MustClaims(c *gin.Context) token.Claims {
	cl, ok := Claims(c)
	if !ok {
		panic("middleware: claims missing — route must use JWTAuth")
	}
	return cl
}

// IsAgentOrAdmin is a convenience guard for agent dashboard routes.
func IsAgentOrAdmin(role string) bool {
	return role == domain.RoleAgent || role == domain.RoleAdmin
}
