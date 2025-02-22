package bcgin

import (
	"cgin"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkCginRouting(b *testing.B) {
	router := cgin.New() // 初始化 cgin 路由
	// middleware := func(c *cgin.Context) {
	// 	c.Next()
	// }
	// router.Use(middleware)
	router.GET("/user/:id", func(c *cgin.Context) {
		c.String(200, "Hello World")
	})

	req := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGinRouting(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New() // 初始化 gin 路由
	// middleware := func(c *gin.Context) {
	// 	c.Next()
	// }
	// router.Use(middleware)
	router.GET("/user/:id", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	req := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
