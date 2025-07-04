package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"service-signalering/pkg/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestAuthenticateKey_ValidKey(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "dev-key-123")
	c.Request = req

	middleware := auth.AuthenticateKey()
	middleware(c)

	assert.False(t, c.IsAborted())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthenticateKey_MissingKey(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	middleware := auth.AuthenticateKey()
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticateKey_InvalidKey(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "invalid-key")
	c.Request = req

	middleware := auth.AuthenticateKey()
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticateKey_CustomEnvironmentKeys(t *testing.T) {
	originalKeys := os.Getenv("API_KEYS")
	os.Setenv("API_KEYS", "custom-key-1,custom-key-2")
	defer os.Setenv("API_KEYS", originalKeys)

	tests := []struct {
		name        string
		apiKey      string
		shouldAbort bool
	}{
		{"Valid custom key 1", "custom-key-1", false},
		{"Valid custom key 2", "custom-key-2", false},
		{"Invalid key", "invalid-key", true},
		{"Default key should not work", "dev-key-123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-API-Key", tt.apiKey)
			c.Request = req

			middleware := auth.AuthenticateKey()
			middleware(c)

			assert.Equal(t, tt.shouldAbort, c.IsAborted())
		})
	}
}

func BenchmarkAuthenticateKey(b *testing.B) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "dev-key-123")
	c.Request = req

	middleware := auth.AuthenticateKey()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Abort()
		middleware(c)
	}
}
