package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HotlinkProtectionMiddleware 防盗链中间件
func HotlinkProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对图片资源做防盗链检查
		if !strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
			c.Next()
			return
		}

		referer := c.Request.Referer()
		
		// 如果没有 referer（直接访问），允许通过
		if referer == "" {
			c.Next()
			return
		}

		// 白名单检查
		allowedDomains := []string{
			c.Request.Host, // 当前域名
			"localhost",
			"127.0.0.1",
		}

		// 检查 referer 是否在白名单中
		allowed := false
		for _, domain := range allowedDomains {
			if strings.Contains(referer, domain) {
				allowed = true
				break
			}
		}

		if !allowed {
			// 盗链请求，返回403
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
