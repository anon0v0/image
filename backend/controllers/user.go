package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserInfoResponse 用户信息响应结构
type UserInfoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func CheckLoginStatus(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	role := session.Get("role")
	loggedIn := session.Get("logged_in")
	
	fmt.Printf("CheckLoginStatus: UserID=%v, Role=%v, LoggedIn=%v\n", userID, role, loggedIn)
	
	// 检查是否真正登录
	isLoggedIn := loggedIn != nil && loggedIn == true
	
	c.JSON(http.StatusOK, UserInfoResponse{
		Code:    200,
		Message: "success",
		Data: map[string]any{
			"user_id":   userID,
			"role":      role,
			"logged_in": isLoggedIn,
		},
	})
}
