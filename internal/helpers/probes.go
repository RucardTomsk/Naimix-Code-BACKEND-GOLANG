package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// WithProbes attaches liveness and readiness probes to a router.
func WithProbes(r *gin.Engine) *gin.Engine {
	r.GET("/system/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ready"})
	})
	r.GET("/system/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "healthy"})
	})

	return r
}
