package controllers

import (
	"net/http"

	"oneimg/backend/database"
	"oneimg/backend/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// User 用户信息结构（不包含密码）
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
			Success: false,
		})
		return
	}

	// 获取数据库实例
	db := database.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Code:    500,
			Message: "数据库连接失败",
			Success: false,
		})
		return
	}

	// 查找用户
	var user models.User
	result := db.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Code:    401,
			Message: "用户名或密码错误",
			Success: false,
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Code:    401,
			Message: "用户名或密码错误",
			Success: false,
		})
		return
	}

	// 获取session
	session := sessions.Default(c)

	// 设置session数据
	session.Set("user_id", user.Id)
	session.Set("username", user.Username)
	session.Set("role", "admin") // 设置角色为管理员
	session.Set("logged_in", true)

	// 设置session选项
	session.Options(sessions.Options{
		MaxAge:   24 * 60 * 60,            // 24小时，单位秒
		HttpOnly: true,                    // 防止XSS攻击
		Secure:   false,                   // 生产环境应设为true（需要HTTPS）
		SameSite: http.SameSiteLaxMode, // 防止CSRF攻击
		Path:     "/",                     // cookie路径
	})

	// 保存session
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Code:    500,
			Message: "保存会话失败",
			Success: false,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, LoginResponse{
		Code:    200,
		Message: "登录成功",
		Success: true,
		User: &User{
			ID:       user.Id,
			Username: user.Username,
		},
	})
}

// CheckLogin 检查登录状态
func CheckLogin(c *gin.Context) {
	session := sessions.Default(c)
	loggedIn := session.Get("logged_in")
	if loggedIn == nil || !loggedIn.(bool) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"logged_in": false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"logged_in": true,
			"username":  session.Get("username"),
			"role":      session.Get("role"),
			"user_id":   session.Get("user_id"),
		},
	})
}
