package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"oneimg/backend/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GithubLogin 跳转到 GitHub 登录页面
func GithubLogin(c *gin.Context) {
	cfg := config.App.GitHubConfig
	if cfg.ClientID == "" || cfg.RedirectURL == "" {
		c.String(http.StatusInternalServerError, "GitHub configuration missing")
		return
	}

	// 生成随机 State 防止 CSRF (这里简单处理，生产环境应使用随机字符串并存入 Session)
	state := "random_state_string" 

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user&state=%s",
		cfg.ClientID, cfg.RedirectURL, state,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GithubCallback 处理 GitHub 回调
func GithubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "No code provided")
		return
	}

	cfg := config.App.GitHubConfig

	// 1. 使用 Code 换取 Access Token
	tokenReqBody := map[string]string{
		"client_id":     cfg.ClientID,
		"client_secret": cfg.ClientSecret,
		"code":          code,
	}
	jsonBody, _ := json.Marshal(tokenReqBody)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonBody))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange token")
		return
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.String(http.StatusInternalServerError, "Failed to decode token response")
		return
	}

	if tokenResp.Error != "" {
		c.String(http.StatusUnauthorized, "GitHub Error: "+tokenResp.Error)
		return
	}

	// 2. 使用 Access Token 获取用户信息
	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	
	userResp, err := client.Do(userReq)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch user info")
		return
	}
	defer userResp.Body.Close()

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&githubUser); err != nil {
		c.String(http.StatusInternalServerError, "Failed to decode user info")
		return
	}

	// 3. 设置 Session
	session := sessions.Default(c)
	session.Set("user_id", githubUser.ID) // 使用 GitHub ID 作为用户 ID
	session.Set("username", githubUser.Login)
	session.Set("role", "user") // GitHub 用户默认为普通用户
	session.Set("logged_in", true)

	// 设置 Cookie 属性
	session.Options(sessions.Options{
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	if err := session.Save(); err != nil {
		fmt.Printf("GithubCallback: Failed to save session: %v\n", err)
		c.String(http.StatusInternalServerError, "Failed to save session")
		return
	}

	fmt.Printf("GithubCallback: Login successful for user %s (ID: %d)\n", githubUser.Login, githubUser.ID)

	// 重定向回首页
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
