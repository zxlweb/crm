package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuperAdminMiddleware_AllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_super_admin", true)
		c.Next()
	})
	r.Use(SuperAdminMiddleware())
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ok", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestSuperAdminMiddleware_BlocksNonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_super_admin", false)
		c.Next()
	})
	r.Use(SuperAdminMiddleware())
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ok", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status %d", w.Code)
	}
}
