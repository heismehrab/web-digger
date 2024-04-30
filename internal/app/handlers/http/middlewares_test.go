package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddlewareWithOptionsMethod(t *testing.T) {
	// Create a new test server.
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// Create a test request to the server without authentication.
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodOptions,
		"/",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check if the response is unauthorized.
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d but got %d", http.StatusNoContent, w.Code)
	}
}

func TestCORSMiddlewareWithGetMethod(t *testing.T) {
	// Create a new test server.
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// Create a test request to the server without authentication.
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD", w.Header().Get("Access-Control-Allow-Methods"))
}
