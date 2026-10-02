package handler

import (
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	version string
	timeout time.Duration
}

func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
		timeout: 3 * time.Second,
	}
}

func (h *HealthHandler) Live(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "alive",
		"version": h.version,
	})
}
