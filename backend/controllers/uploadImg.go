package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"oneimg/backend/config"
	"oneimg/backend/database"
	"oneimg/backend/models"
	"oneimg/backend/services"

	"github.com/gin-gonic/gin"
)

// UploadResponse 上传响应结构
type UploadResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []ImageResult `json:"data"`
}

// ImageResult 单个图片上传结果
type ImageResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ID        int    `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	FileName  string `json:"filename,omitempty"`
	FileSize  int64  `json:"file_size,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Category  string `json:"category,omitempty"`
	Tags      string `json:"tags,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// OpenAIRequest OpenAI API 请求结构
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageUrl *ImageUrl `json:"image_url,omitempty"`
}

type ImageUrl struct {
	Url string `json:"url"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AIProgress 进度结构
type AIProgressStruct struct {
	Total     int  `json:"total"`
	Current   int  `json:"current"`
	IsRunning bool `json:"is_running"`
}

var AIProgress = AIProgressStruct{}

func ensureUploadDir(uploadPath string) error {
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		return os.MkdirAll(uploadPath, 0755)
	}
	return nil
}

func generateUniqueFileName(ext string) string {
	timestamp := time.Now().UnixNano()
	hash := fmt.Sprintf("%x", timestamp)
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(900) + 100
	return fmt.Sprintf("%s%d%s", hash, randomNum, ext)
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func determineOutputFormat(contentType, originalExt string) string {
	specialFormats := map[string]string{"image/gif": ".gif", "image/svg+xml": ".svg"}
	if ext, exists := specialFormats[contentType]; exists {
		return ext
	}
	switch strings.ToLower(originalExt) {
	case ".gif":
		return ".gif"
	case ".svg":
		return ".svg"
	default:
		return ".webp"
	}
}

func calculateFileHash(file multipart.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// getAIInfo 调用 AI 获取标签和分类
func getAIInfo(imageBytes []byte, cfg *config.Config) (string, string) {
	// 从数据库获取 AI 配置
	db := database.GetDB().DB
	var setting models.Settings
	var aiConfig models.AIConfig

	if err := db.Where("`key` = ?", "ai_config").First(&setting).Error; err == nil {
		json.Unmarshal([]byte(setting.Value), &aiConfig)
	}

	// 如果数据库没有配置，尝试使用配置文件（兼容旧逻辑），或者直接返回
	if aiConfig.ApiUrl == "" {
		aiConfig.ApiUrl = cfg.AiApiUrl
		aiConfig.ApiKey = cfg.AiApiKey
		aiConfig.Model = cfg.AiModel
	}

	if aiConfig.ApiUrl == "" || aiConfig.ApiKey == "" {
		return "", ""
	}

	// Base64 编码
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	imgDataUrl := fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)

	// 构建 Prompt - 强调中文输出
	prompt := `你是一个图片分析助手。请分析这张图片，返回JSON格式结果。

要求：
1. tags字段：包含3-8个简体中文关键词的数组，描述图片内容。如果能识别出具体角色（如动漫人物、游戏角色、明星等），把角色名作为第一个标签。
2. category字段：从以下分类中选择一个最匹配的：动漫、人物、风景、影视、游戏、美食、动物、艺术、宇宙、科技、简约、机车、其他

注意：所有标签必须是简体中文，不要使用英文！

示例输出：{"tags": ["初音未来", "双马尾", "蓝色头发", "舞台", "演唱会"], "category": "动漫"}

只返回JSON字符串，不要添加任何markdown格式或其他说明文字。`

	reqBody := OpenAIRequest{
		Model: aiConfig.Model,
		Messages: []Message{{Role: "user", Content: []Content{
			{Type: "text", Text: prompt},
			{Type: "image_url", ImageUrl: &ImageUrl{Url: imgDataUrl}},
		}}},
	}

	jsonData, _ := json.Marshal(reqBody)

	// 处理 API URL，确保正确拼接 /chat/completions
	apiURL := aiConfig.ApiUrl
	if strings.HasSuffix(apiURL, "/v1") {
		apiURL = apiURL + "/chat/completions"
	} else if !strings.HasSuffix(apiURL, "/chat/completions") {
		// 如果没有 /v1 也没有 /chat/completions，假设是根地址，补全标准路径
		// 但为了兼容性，最好检查是否以 / 结尾
		if strings.HasSuffix(apiURL, "/") {
			apiURL = apiURL + "v1/chat/completions"
		} else {
			apiURL = apiURL + "/v1/chat/completions"
		}
	}

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+aiConfig.ApiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("⚠️ [Debug] AI请求失败:", err)
		return "", ""
	}
	defer resp.Body.Close()

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		fmt.Println("⚠️ [Debug] AI响应解析失败:", err)
		return "", ""
	}

	if len(openAIResp.Choices) > 0 {
		content := openAIResp.Choices[0].Message.Content
		// 清理 Markdown 标记
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)

		type AIResult struct {
			Tags     []string `json:"tags"`
			Category string   `json:"category"`
		}
		var result AIResult
		if err := json.Unmarshal([]byte(content), &result); err == nil {
			// 将 tags 数组转换为逗号分隔的字符串
			tagsStr := strings.Join(result.Tags, ",")
			fmt.Printf("✅ [Debug] AI识别成功: Tags=%s, Category=%s\n", tagsStr, result.Category)
			return tagsStr, result.Category
		} else {
			fmt.Println("⚠️ [Debug] AI返回非JSON格式:", content)
		}
	}

	return "", ""
}

func processUploadFile(fileHeader *multipart.FileHeader, cfg *config.Config, db *database.Database) ImageResult {
	// 1. 基础验证
	file, err := fileHeader.Open()
	if err != nil {
		return ImageResult{Success: false, Message: "无法打开文件"}
	}
	defer file.Close()

	// 计算哈希
	fileHash, err := calculateFileHash(file)
	if err != nil {
		return ImageResult{Success: false, Message: "哈希计算失败"}
	}
	file.Seek(0, 0) // 重置指针

	// 查重
	var existingImage models.Image
	if err := db.DB.Where("hash = ?", fileHash).First(&existingImage).Error; err == nil {
		return ImageResult{
			Success:   true,
			Message:   "图片已存在",
			ID:        existingImage.Id,
			URL:       existingImage.Url,
			FileName:  existingImage.FileName,
			FileSize:  existingImage.FileSize,
			MimeType:  existingImage.MimeType,
			Width:     existingImage.Width,
			Height:    existingImage.Height,
			Category:  existingImage.Category,
			Tags:      existingImage.Tags,
			CreatedAt: existingImage.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// 图片处理服务
	imgService := services.NewImageService()
	processedImage, err := imgService.ProcessImage(file, fileHeader)
	if err != nil {
		return ImageResult{Success: false, Message: "图片处理失败: " + err.Error()}
	}

	// 2. 命名逻辑 (随机) & AI 标签
	originalExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	outputExt := originalExt

	// 根据解码后的格式规范化扩展名
	switch processedImage.Format {
	case "jpeg":
		outputExt = ".jpg"
	case "png":
		outputExt = ".png"
	case "gif":
		outputExt = ".gif"
	case "webp":
		outputExt = ".webp"
	}
	
	// 如果 outputExt 仍然为空，使用原始扩展名
	if outputExt == "" {
		outputExt = originalExt
	}

	uniqueFileName := generateUniqueFileName(outputExt)
	var category string = "其他"
	var tags string = ""

	// 调用 AI
	aiTags, aiCategory := getAIInfo(processedImage.ThumbnailBytes, cfg)
	if aiCategory != "" {
		category = aiCategory
	}
	if aiTags != "" {
		tags = aiTags
	}

	// 3. 保存文件
	today := time.Now().Format("20060102")
	saveDir := filepath.Join(cfg.UploadPath, today)
	if err := ensureUploadDir(saveDir); err != nil {
		return ImageResult{Success: false, Message: "创建目录失败"}
	}

	savePath := filepath.Join(saveDir, uniqueFileName)
	fileUrl := fmt.Sprintf("/uploads/%s/%s", today, uniqueFileName)

	// 保存主文件
	out, err := os.Create(savePath)
	if err != nil {
		return ImageResult{Success: false, Message: "保存文件失败"}
	}
	defer out.Close()
	out.Write(processedImage.CompressedBytes)

	// 保存缩略图
	thumbName := strings.TrimSuffix(uniqueFileName, outputExt) + "_thumb" + outputExt
	thumbPath := filepath.Join(saveDir, thumbName)
	thumbOut, err := os.Create(thumbPath)
	if err == nil {
		defer thumbOut.Close()
		thumbOut.Write(processedImage.ThumbnailBytes)
	}

	// 保存预览图 (用于前端展示)
	previewName := strings.TrimSuffix(uniqueFileName, outputExt) + "_preview.webp"
	previewPath := filepath.Join(saveDir, previewName)
	previewOut, err := os.Create(previewPath)
	if err == nil {
		defer previewOut.Close()
		previewOut.Write(processedImage.PreviewBytes)
	}

	// 4. 数据库记录
	imageModel := models.Image{
		Url:       fileUrl,
		FileName:  uniqueFileName,
		FileSize:  int64(len(processedImage.CompressedBytes)),
		MimeType:  processedImage.MimeType,
		Width:     processedImage.Width,
		Height:    processedImage.Height,
		Hash:      fileHash,
		Category:  category,
		Tags:      tags,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&imageModel).Error; err != nil {
		return ImageResult{Success: false, Message: "数据库保存失败"}
	}

	return ImageResult{
		Success:   true,
		ID:        imageModel.Id,
		URL:       imageModel.Url,
		FileName:  imageModel.FileName,
		FileSize:  imageModel.FileSize,
		MimeType:  imageModel.MimeType,
		Width:     imageModel.Width,
		Height:    imageModel.Height,
		Category:  imageModel.Category,
		Tags:      imageModel.Tags,
		CreatedAt: imageModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// UploadImages 批量上传 (Frontend uses this)
func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "获取表单数据失败", "data": []string{}})
		return
	}

	files := form.File["images[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "未检测到文件", "data": []string{}})
		return
	}

	cfg := c.MustGet("config").(*config.Config)
	db := database.GetDB()

	var results []ImageResult
	successCount := 0

	for _, file := range files {
		result := processUploadFile(file, cfg, db)
		results = append(results, result)
		if result.Success {
			successCount++
		}
	}

	if successCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": fmt.Sprintf("成功上传 %d 张图片", successCount),
			"data":    results,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "所有图片上传失败",
			"data":    results,
		})
	}
}

// UploadImage 单个上传
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "上传失败: " + err.Error(), "data": []string{}})
		return
	}

	cfg := c.MustGet("config").(*config.Config)
	db := database.GetDB()

	result := processUploadFile(file, cfg, db)

	if result.Success {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "上传成功", "data": result})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": result.Message, "data": []string{}})
	}
}

// GetAIProgress 获取 AI 进度
func GetAIProgress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": AIProgress,
	})
}

// BatchTagImages 批量打标签 (原 BatchRenameOldImages)
func BatchTagImages(c *gin.Context) {
	if AIProgress.IsRunning {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "任务已在运行中"})
		return
	}

	db := database.GetDB().DB
	var images []models.Image

	// 查找所有 Tags 为空或 Category 为 "其他" 的图片
	if err := db.Where("tags = ? OR tags IS NULL", "").Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询图片失败"})
		return
	}

	if len(images) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "没有需要处理的图片"})
		return
	}

	cfg := c.MustGet("config").(*config.Config)
	
	// 初始化进度
	AIProgress.Total = len(images)
	AIProgress.Current = 0
	AIProgress.IsRunning = true

	go func() {
		defer func() {
			AIProgress.IsRunning = false
		}()

		// 限制并发数为 3
		semaphore := make(chan struct{}, 3)
		var wg sync.WaitGroup

		for _, img := range images {
			wg.Add(1)
			semaphore <- struct{}{} // 获取令牌

			go func(image models.Image) {
				defer wg.Done()
				defer func() { <-semaphore }() // 释放令牌

				// 读取图片文件
				relPath := strings.TrimPrefix(image.Url, "/uploads/")
				fullPath := filepath.Join(cfg.UploadPath, relPath)

				fileBytes, err := os.ReadFile(fullPath)
				if err != nil {
					fmt.Printf("读取文件失败: %s\n", fullPath)
					AIProgress.Current++
					return
				}

				// 尝试读取缩略图以节省token
				ext := filepath.Ext(fullPath)
				thumbPath := strings.TrimSuffix(fullPath, ext) + "_thumb" + ext
				if thumbBytes, err := os.ReadFile(thumbPath); err == nil {
					fileBytes = thumbBytes
				}

				tags, category := getAIInfo(fileBytes, cfg)
				if tags != "" {
					image.Tags = tags
					image.Category = category
					db.Save(&image)
				}
				AIProgress.Current++
			}(img)
		}
		
		wg.Wait()
		fmt.Printf("批量打标签完成，更新了 %d 张图片\n", AIProgress.Current)
	}()

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "AI 打标签任务已在后台启动"})
}

// BatchDeduplicate 批量去重
func BatchDeduplicate(c *gin.Context) {
	db := database.GetDB().DB

	// 查找所有重复的 hash
	type Result struct {
		Hash  string
		Count int
	}
	var results []Result

	if err := db.Model(&models.Image{}).Select("hash, count(*) as count").Group("hash").Having("count > 1").Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	deletedCount := 0
	cfg := c.MustGet("config").(*config.Config)

	for _, r := range results {
		var images []models.Image
		db.Where("hash = ?", r.Hash).Order("id asc").Find(&images)

		// 保留第一个，删除其余的
		if len(images) > 1 {
			for i := 1; i < len(images); i++ {
				img := images[i]
				
				// 删除物理文件
				relPath := strings.TrimPrefix(img.Url, "/uploads/")
				fullPath := filepath.Join(cfg.UploadPath, relPath)
				os.Remove(fullPath)
				
				// 删除缩略图
				ext := filepath.Ext(fullPath)
				thumbPath := strings.TrimSuffix(fullPath, ext) + "_thumb" + ext
				os.Remove(thumbPath)

				// 删除预览图
				previewPath := strings.TrimSuffix(fullPath, ext) + "_preview.webp"
				os.Remove(previewPath)

				// 删除数据库记录
				db.Delete(&img)
				deletedCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("去重完成，删除了 %d 张重复图片", deletedCount),
	})
}