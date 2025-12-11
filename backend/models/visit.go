package models

import (
	"time"
)

// Visit 访客记录模型
type Visit struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	IP        string    `gorm:"index;size:45" json:"ip"`         // 支持IPv6
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Referer   string    `gorm:"size:500" json:"referer"`
	Path      string    `gorm:"index;size:255" json:"path"`
	Method    string    `gorm:"size:10" json:"method"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (Visit) TableName() string {
	return "visits"
}
