package httpx

import (
	"errors"
	"net/http"

	"kerjadekat/backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	case errors.Is(err, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrInvalidOTP):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	case errors.Is(err, domain.ErrInvalidTransition):
		c.JSON(http.StatusConflict, gin.H{"error": "invalid state transition"})
	case errors.Is(err, domain.ErrPaymentFailed):
		c.JSON(http.StatusBadGateway, gin.H{"error": "payment error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
