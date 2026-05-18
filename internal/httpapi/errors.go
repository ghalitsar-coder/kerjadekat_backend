package httpapi

import (
	"kerjadekat/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	httpx.WriteError(c, err)
}
