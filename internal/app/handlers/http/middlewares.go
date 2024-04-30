package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORSMiddleware provides CORS ability for the router.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}
