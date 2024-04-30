package http

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"net/http"
)

// setupRouter sets list of route to http handler.
func (h *Handler) setupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.MaxMultipartMemory = MaxMultipartMemory

	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	r.Use(sloggin.New(h.logger.Logger))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	return r
}
