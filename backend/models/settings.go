package models

import "time"

// Settings 系统设置
type Settings struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"unique;not null"`
	Value     string    `json:"value" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIConfig AI配置结构体 (用于JSON序列化存储在Settings中)
type AIConfig struct {
	ApiUrl string `json:"api_url"`
	ApiKey string `json:"api_key"`
	Model  string `json:"model"`
}
