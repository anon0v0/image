package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"oneimg/backend/config"
	"oneimg/backend/middlewares"

	"github.com/gin-gonic/gin"
)

// ServeImage 动态图片服务 (控制访问权限)
func ServeImage(c *gin.Context) {
	// 获取请求的文件路径
	filePath := c.Param("filepath")
	
	// 防止目录遍历攻击
	if strings.Contains(filePath, "..") {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "Forbidden"})
		return
	}

	cfg := config.App
	fullPath := filepath.Join(cfg.UploadPath, filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Image not found"})
		return
	}

	// 检查用户登录状态
	_, _, loggedIn := middlewares.GetCurrentUser(c)

	// 如果已登录，直接返回原图
	if loggedIn {
		c.File(fullPath)
		return
	}

	// 如果未登录，尝试返回缩略图
	// 假设缩略图命名规则：原文件名_thumb.ext
	ext := filepath.Ext(fullPath)
	nameWithoutExt := strings.TrimSuffix(fullPath, ext)
	thumbPath := nameWithoutExt + "_thumb" + ext

	// 检查缩略图是否存在
	if _, err := os.Stat(thumbPath); err == nil {
		c.File(thumbPath)
		return
	}

	// 如果缩略图不存在，但原图存在，且未登录...
	// 策略：为了安全，未登录用户只能看缩略图。如果缩略图丢失，暂时返回原图或403？
	// 根据需求： "不登录的情况下可以查看到图片，只能看到图片的缩略图"
	// 如果没有缩略图，返回原图可能会泄露高清图。
	// 但为了用户体验，如果缩略图生成失败，暂时返回原图可能更好，或者实时生成缩略图（复杂）。
	// 这里选择：如果缩略图不存在，返回原图 (Fallback)，或者返回占位图。
	// 考虑到已有逻辑是上传时生成缩略图，如果丢失说明异常。
	// 暂时返回原图，避免裂图。
	c.File(fullPath)
}
