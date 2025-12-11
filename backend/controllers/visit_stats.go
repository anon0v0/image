package controllers

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"oneimg/backend/database"
	"oneimg/backend/models"
	"oneimg/backend/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetVisitStats 获取访客统计数据
func GetVisitStats(c *gin.Context) {
	db := database.GetDB().DB

	// 检查是否是管理员
	session := sessions.Default(c)
	role := session.Get("role")
	isAdmin := role == "admin"

	var stats struct {
		TotalVisits    int64                    `json:"total_visits"`
		TodayVisits    int64                    `json:"today_visits"`
		UniqueVisitors int64                    `json:"unique_visitors"`
		TotalUsers     int64                    `json:"total_users"`
		ActiveUsers    int64                    `json:"active_users"`
		VisitTrend     []map[string]interface{} `json:"visit_trend"`
		TopIPs         []map[string]interface{} `json:"top_ips"`
		ImageCallRank  []map[string]interface{} `json:"image_call_rank"`
		IsAdmin        bool                     `json:"is_admin"`
	}

	stats.IsAdmin = isAdmin

	// 1. 总访问量
	db.Model(&models.Visit{}).Count(&stats.TotalVisits)

	// 2. 今日访问量
	today := time.Now().Format("2006-01-02")
	db.Model(&models.Visit{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.TodayVisits)

	// 3. 唯一访客数（基于IP）
	db.Model(&models.Visit{}).
		Distinct("ip").
		Count(&stats.UniqueVisitors)

	// 4. 总用户数
	db.Model(&models.User{}).Count(&stats.TotalUsers)

	// 5. 活跃用户数（最近7天有更新的用户）
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	db.Model(&models.User{}).
		Where("updated_at >= ?", sevenDaysAgo).
		Count(&stats.ActiveUsers)

	// 6. 访问趋势（最近7天）
	stats.VisitTrend = make([]map[string]interface{}, 0)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		db.Model(&models.Visit{}).
			Where("DATE(created_at) = ?", date).
			Count(&count)
		stats.VisitTrend = append(stats.VisitTrend, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}

	// 7. 来源IP排行（Top 15，含归属地）
	type IPStats struct {
		IP    string
		Count int64
	}
	var ipStats []IPStats
	db.Model(&models.Visit{}).
		Select("ip, COUNT(*) as count").
		Group("ip").
		Order("count DESC").
		Limit(15).
		Scan(&ipStats)

	stats.TopIPs = make([]map[string]interface{}, 0)
	for _, is := range ipStats {
		// 查询IP归属地
		location := services.GetIPLocation(is.IP)

		// 根据是否是管理员决定是否隐藏IP
		displayIP := is.IP
		if !isAdmin {
			displayIP = maskIP(is.IP)
		}

		stats.TopIPs = append(stats.TopIPs, map[string]interface{}{
			"ip":       displayIP,
			"full_ip":  is.IP, // 完整IP，前端根据isAdmin决定是否使用
			"count":    is.Count,
			"country":  location.Country,
			"region":   location.Region,
			"city":     location.City,
			"isp":      location.ISP,
			"is_admin": isAdmin,
		})
	}

	// 8. 图片调用排行（统计外链访问次数）
	// 统计访问 /uploads/ 的记录（真正的图片外链调用）
	type ImageCallStats struct {
		Path  string
		Count int64
	}
	var imageCallStats []ImageCallStats
	db.Model(&models.Visit{}).
		Select("path, COUNT(*) as count").
		Where("path LIKE '/uploads/%'").
		Group("path").
		Order("count DESC").
		Limit(20).
		Scan(&imageCallStats)

	stats.ImageCallRank = make([]map[string]interface{}, 0)
	for _, ics := range imageCallStats {
		// 从路径中提取文件名
		parts := strings.Split(ics.Path, "/")
		filename := ""
		if len(parts) > 0 {
			filename = parts[len(parts)-1]
		}

		// 去掉可能的查询参数
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}

		// 提取原始文件名（去除 _thumb 和 _preview 后缀）
		// 例如: abstract-flowing-gradients-gi8z_thumb.webp -> abstract-flowing-gradients-gi8z.webp
		originalFilename := filename
		ext := filepath.Ext(filename)
		baseName := strings.TrimSuffix(filename, ext)
		if strings.HasSuffix(baseName, "_thumb") {
			baseName = strings.TrimSuffix(baseName, "_thumb")
			originalFilename = baseName + ext
		} else if strings.HasSuffix(baseName, "_preview") {
			baseName = strings.TrimSuffix(baseName, "_preview")
			originalFilename = baseName + ".webp" // _preview 是 webp 格式的，但原图可能是其他格式
		}

		// 尝试根据原始文件名查找图片信息
		var image models.Image
		if originalFilename != "" {
			// 用原始文件名匹配
			db.Where("url LIKE ?", "%"+originalFilename).First(&image)
			// 如果没找到，尝试用 baseName 模糊匹配
			if image.Id == 0 && baseName != "" {
				db.Where("url LIKE ?", "%"+baseName+"%").First(&image)
			}
		}

		item := map[string]interface{}{
			"path":     ics.Path,
			"calls":    ics.Count,
			"filename": filename,
		}

		if image.Id > 0 {
			item["id"] = image.Id
			// 生成缩略图路径：从完整图片 URL 生成缩略图 URL
			// 例如: /uploads/2025/11/xxx.webp -> /uploads/2025/11/xxx_thumb.webp
			ext := filepath.Ext(image.Url)
			thumbnailUrl := strings.TrimSuffix(image.Url, ext) + "_thumb" + ext
			item["thumbnail"] = thumbnailUrl
			item["original_name"] = image.FileName
			item["width"] = image.Width
			item["height"] = image.Height
		}

		stats.ImageCallRank = append(stats.ImageCallRank, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
	})
}

// maskIP 隐藏IP地址的中间部分
func maskIP(ip string) string {
	if ip == "" {
		return "未知"
	}
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".***.***"
	}
	// IPv6 或其他格式，隐藏中间部分
	if len(ip) > 8 {
		return ip[:4] + "****" + ip[len(ip)-4:]
	}
	return ip
}
