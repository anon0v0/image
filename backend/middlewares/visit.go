package middlewares

import (
	"strings"
	"time"

	"oneimg/backend/database"
	"oneimg/backend/models"

	"github.com/gin-gonic/gin"
)

// VisitMiddleware 访客记录中间件
func VisitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// 排除前端静态资源请求（但保留uploads图片访问的记录）
		if strings.HasPrefix(path, "/assets") ||
			strings.HasPrefix(path, "/static") ||
			path == "/favicon.ico" ||
			path == "/logo.png" ||
			path == "/" ||
			strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		// 只记录图片外链调用（/uploads路径）
		if strings.HasPrefix(path, "/uploads") {
			// 创建访客记录
			visit := models.Visit{
				IP:        c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Referer:   c.Request.Referer(),
				Path:      path,
				Method:    c.Request.Method,
				CreatedAt: time.Now(),
			}

			// 异步保存到数据库，不阻塞请求
			go func() {
				database.GetDB().DB.Create(&visit)
			}()
		}

		c.Next()
	}
}

