package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Deps holds dependencies injected into the HTTP server.
// Extend this struct as new services are implemented.
type Deps struct{}

// NewRouter wires the HTTP router with common endpoints.
// This keeps bootstrap logic in one place for tests and main.
func NewRouter(_ context.Context, _ Deps) *gin.Engine {
	r := gin.New()

	// Default middleware: logging and recovery. Can be swapped if needed.
	r.Use(gin.Logger(), gin.Recovery())

	// Health endpoint for readiness checks.
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
