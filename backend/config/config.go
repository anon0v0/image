package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// 服务器配置
	Port string

	// 数据库配置
	SqlitePath string
	IsMysql    bool
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string

	// 上传文件配置
	MaxFileSize  int64
	AllowedTypes []string
	UploadPath   string

	// 默认用户
	DefaultUser string
	DefaultPass string

	// 安全配置
	JWTSecret     string
	SessionSecret string

	// AI 配置
	AiApiUrl string
	AiApiKey string
	AiModel  string
	AiPrompt string

	// App URL
	AppUrl string

	// GitHub 配置
	GitHubConfig GitHubConfig
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// 设置全局
var App *Config

func NewConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64)
	allowedTypes := strings.Split(getEnv("ALLOWED_TYPES", "image/jpeg,image/png,image/gif"), ",")
	port := getEnv("SERVER_PORT", getEnv("PORT", "8080"))

	sqlitePath := getEnv("SQLITE_PATH", "./data/data.db")
	isMysql := getEnv("IS_MYSQL", "false") == "true"
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "oneimgxru")

	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	defaultUser := getEnv("DEFAULT_USER", "admin")
	defaultPass := getEnv("DEFAULT_PASS", "123456")

	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-this-in-production")
	sessionSecret := getEnv("SESSION_SECRET", "your-session-secret-key-change-this-in-production")

	// AI 配置读取
	aiApiUrl := getEnv("AI_API_URL", "")
	aiApiKey := getEnv("AI_API_KEY", "")
	aiModel := getEnv("AI_MODEL", "gemini-1.5-flash")
	
	// 默认提示词
	defaultPrompt := "Please describe the content of this image in English using 3 to 5 words, joined by hyphens (e.g. cat-eating-fish). Only return the filename string, no extension, no other text."
	aiPrompt := getEnv("AI_PROMPT", defaultPrompt)

	appUrl := getEnv("APP_URL", "http://localhost:8080")

	App = &Config{
		Port:          port,
		SqlitePath:    sqlitePath,
		IsMysql:       isMysql,
		DbHost:        dbHost,
		DbPort:        dbPort,
		DbUser:        dbUser,
		DbPassword:    dbPassword,
		DbName:        dbName,
		UploadPath:    uploadPath,
		MaxFileSize:   maxFileSize,
		AllowedTypes:  allowedTypes,
		DefaultUser:   defaultUser,
		DefaultPass:   defaultPass,
		JWTSecret:     jwtSecret,
		SessionSecret: sessionSecret,
		AiApiUrl:      aiApiUrl,
		AiApiKey:      aiApiKey,
		AiModel:       aiModel,
		AiPrompt:      aiPrompt,
		AppUrl:        appUrl,
		GitHubConfig: GitHubConfig{
			ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GITHUB_REDIRECT_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}