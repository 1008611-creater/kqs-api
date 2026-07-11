package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestConfigureTrustedClientIPCloudflareUsesCFConnectingIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	configureTrustedClientIP(router, &config.Config{
		Server: config.ServerConfig{
			Mode:            "release",
			TrustedPlatform: "cloudflare",
		},
	})
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("CF-Connecting-IP", "203.0.113.42")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "203.0.113.42", rec.Body.String())
}

func TestConfigureTrustedClientIPDefaultIgnoresSpoofedCFHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	configureTrustedClientIP(router, &config.Config{
		Server: config.ServerConfig{Mode: "release"},
	})
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.10:12345"
	req.Header.Set("CF-Connecting-IP", "203.0.113.42")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "198.51.100.10", rec.Body.String())
}
