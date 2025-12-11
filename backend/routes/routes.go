package routes

import (
	"net/http"
	"time"

	"oneimg/backend/config"
	"oneimg/backend/controllers"
	"oneimg/backend/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 设置路由
func SetupRoutes() *gin.Engine {
	cfg := config.App

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 基础中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.ConfigMiddleware(cfg))
	r.Use(middlewares.SessionMiddleware(cfg))
	// 访客记录中间件
	r.Use(middlewares.VisitMiddleware())

	// 跨域配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 静态资源
	r.Static("/static", "./static/frontend")
	// 对uploads目录启用防盗链保护
	uploadsGroup := r.Group("/uploads")
	uploadsGroup.Use(middlewares.HotlinkProtectionMiddleware())
	uploadsGroup.Static("/", cfg.UploadPath)
	
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	r.StaticFile("/logo.png", "./frontend/dist/logo.png")

	// API路由分组
	api := r.Group("/api")
	{
		// 公开接口（无需认证）
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)
		api.GET("/logout", controllers.Logout)
		api.GET("/user/status", controllers.CheckLoginStatus)

		// GitHub 认证接口（公开）
		authGroup := api.Group("/auth")
		{
			authGroup.GET("/github/login", controllers.GithubLogin)
			authGroup.GET("/github/callback", controllers.GithubCallback)
		}

		// 可选认证接口（不强制要求登录）
		optional := api.Group("")
		optional.Use(middlewares.OptionalAuthMiddleware())
		{
			optional.GET("/images", controllers.GetImageList)
			optional.GET("/images/:id", controllers.GetImageDetail)
		}

		// 需要认证的接口分组（应用AuthMiddleware）
		auth := api.Group("")
		auth.Use(middlewares.AuthMiddleware())
		{
			// 统计数据 (普通用户也可查看统计，或者你可以决定是否限制)
			auth.GET("/stats/dashboard", controllers.GetDashboardStats)
			auth.GET("/stats/images", controllers.GetImageStats)
			auth.GET("/stats/visits", controllers.GetVisitStats)
			
			// 账户管理接口
			auth.POST("/account/change", controllers.ChangeAccountInfo)
			auth.POST("/sessions/clear", controllers.ClearAllSessions)

			// 管理员接口分组
			admin := auth.Group("")
			admin.Use(middlewares.AdminMiddleware())
			{
				// 图片相关接口 (上传、删除)
				admin.POST("/upload", controllers.UploadImage)
				admin.POST("/upload/images", controllers.UploadImages)
				admin.DELETE("/images/:id", controllers.DeleteImage)

				// AI 设置
				admin.GET("/settings/ai", controllers.GetAISettings)
				admin.POST("/settings/ai", controllers.SaveAISettings)
				admin.GET("/ai/models", controllers.GetAIModels)

				// AI 任务
				admin.POST("/batch-tag", controllers.BatchTagImages)
				admin.GET("/ai/progress", controllers.GetAIProgress)
				// 图片去重
				admin.POST("/deduplicate", controllers.BatchDeduplicate)
			}
		}
	}

	// 前端SPA路由支持
	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "API Not Found"})
			return
		}
		c.File("./frontend/dist/index.html")
	})

	return r
}
