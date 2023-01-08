package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// StrictHeader 添加安全请求头
func StrictHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Download-Options", "noopen")
		c.Header("Strict-Transport-Security", "max-age=5184000")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-DNS-Prefetch-Control", "off")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "True")
		// 放行索引options
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
		return
	}
}
