package http

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/slog-gin"
	"os"
)

const resourcesDir = "/resources/static"

// setupRouter sets list of route to http handler.
func (h *Handler) setupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.MaxMultipartMemory = MaxMultipartMemory

	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	r.Use(sloggin.New(h.logger.Logger))

	// Define static pages.
	h.staticRoutes(r)

	// Define API routes.
	h.APIRoutes(r)

	return r
}

// staticRoutes defines static routes including web pages.
func (h *Handler) staticRoutes(r *gin.Engine) {
	pwd, _ := os.Getwd()
	staticWebPageAddr := pwd + resourcesDir

	r.Static("/", staticWebPageAddr)
}

// APIRoutes defines HTTP groups endpoints.
func (h *Handler) APIRoutes(r *gin.Engine) {
	// Defining groups.
	api := r.Group("api")
	v1 := api.Group("v1")

	// Defining endpoints.
	v1.POST("analyze-html", h.AnalyzeWebPage)
}
