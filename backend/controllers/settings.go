package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"oneimg/backend/database"
	"oneimg/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAISettings 获取AI配置
func GetAISettings(c *gin.Context) {
	db := database.GetDB().DB
	var setting models.Settings

	// 查找 key 为 "ai_config" 的设置
	if err := db.Where("`key` = ?", "ai_config").First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，返回默认空配置
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": models.AIConfig{
					ApiUrl: "https://api.openai.com",
					ApiKey: "",
					Model:  "gpt-3.5-turbo",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取设置失败"})
		return
	}

	var aiConfig models.AIConfig
	if err := json.Unmarshal([]byte(setting.Value), &aiConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "解析配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": aiConfig,
	})
}

// SaveAISettings 保存AI配置
func SaveAISettings(c *gin.Context) {
	var aiConfig models.AIConfig
	if err := c.ShouldBindJSON(&aiConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	configJson, err := json.Marshal(aiConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "配置序列化失败"})
		return
	}

	db := database.GetDB().DB
	var setting models.Settings

	// 查找或创建
	err = db.Where("`key` = ?", "ai_config").First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		setting = models.Settings{
			Key:       "ai_config",
			Value:     string(configJson),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存设置失败"})
			return
		}
	} else if err == nil {
		setting.Value = string(configJson)
		setting.UpdatedAt = time.Now()
		if err := db.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新设置失败"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功"})
}

// GetAIModels 获取可用模型列表
func GetAIModels(c *gin.Context) {
	// 先获取当前配置
	db := database.GetDB().DB
	var setting models.Settings
	var aiConfig models.AIConfig

	if err := db.Where("`key` = ?", "ai_config").First(&setting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请先保存 AI 配置"})
		return
	}

	if err := json.Unmarshal([]byte(setting.Value), &aiConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "配置解析失败"})
		return
	}

	if aiConfig.ApiUrl == "" || aiConfig.ApiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "API配置不完整"})
		return
	}

	// 构建请求
	// 处理 API URL，去掉 /chat/completions 或 /v1 后缀，指向根路径
	// 通常 OpenAI 兼容接口的模型列表在 /v1/models
	// 如果用户配置的是 https://api.openai.com/v1，我们需要请求 https://api.openai.com/v1/models
	
	// 简单处理：如果以 /chat/completions 结尾，去掉它
	// 目标是构造出 /v1/models 的 URL
	
	// 假设用户填写的 api_url 是 https://api.openai.com/v1
	// 或者 https://api.openai.com/v1/chat/completions
	
	baseUrl := aiConfig.ApiUrl
	// 移除可能存在的 /chat/completions
	if len(baseUrl) > 17 && baseUrl[len(baseUrl)-17:] == "/chat/completions" {
		baseUrl = baseUrl[:len(baseUrl)-17]
	}
	
	// 确保以 /v1 结尾（如果不是，尝试追加，或者假设用户填的就是 base url）
	// 这里做一个简单的尝试：直接请求 baseUrl + "/models" (如果 baseUrl 包含 /v1)
	// 或者 baseUrl + "/v1/models" (如果 baseUrl 不包含 /v1)
	
	// 更稳妥的方式：
	// 如果 baseUrl 包含 /v1，则请求 /models
	// 否则请求 /v1/models
	
	var modelsUrl string
	if len(baseUrl) > 3 && baseUrl[len(baseUrl)-3:] == "/v1" {
		modelsUrl = baseUrl + "/models"
	} else if len(baseUrl) > 4 && baseUrl[len(baseUrl)-4:] == "/v1/" {
		modelsUrl = baseUrl + "models"
	} else {
		// 尝试追加 /v1/models
		if baseUrl[len(baseUrl)-1] == '/' {
			modelsUrl = baseUrl + "v1/models"
		} else {
			modelsUrl = baseUrl + "/v1/models"
		}
	}

	req, err := http.NewRequest("GET", modelsUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建请求失败"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+aiConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "请求模型列表失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取模型失败，请检查 URL 和 Key"})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "解析响应失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": result["data"], // 通常模型列表在 data 字段
	})
}
