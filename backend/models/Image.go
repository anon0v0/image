package models

import "time"

// 图片模型
type Image struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Url       string    `json:"url" gorm:"not null"`
	FileName  string    `json:"filename" gorm:"not null"`
	FileSize  int64     `json:"file_size" gorm:"not null"`
	MimeType  string    `json:"mimeType"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
    // ------------------ 新增 Hash 字段 ------------------
	Hash      string    `json:"hash" gorm:"size:64;index"` 
	Category  string    `json:"category" gorm:"size:50;index"` // 新增分类字段
	Tags      string    `json:"tags" gorm:"type:text"`         // 新增标签字段 (JSON array or comma-separated)
    // ---------------------------------------------------
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}