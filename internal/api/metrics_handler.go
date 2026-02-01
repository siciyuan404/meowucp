package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	token string
}

func NewMetricsHandler(token string) *MetricsHandler {
	return &MetricsHandler{token: token}
}

func (h *MetricsHandler) Serve(c *gin.Context) {
	if h.token != "" && c.GetHeader("X-Metrics-Token") != h.token {
		respondError(c, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	c.String(http.StatusOK, "# metrics placeholder\n")
}
