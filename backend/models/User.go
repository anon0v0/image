package models

import "time"

// 用户模型
type User struct {
	Id       int       `json:"id" gorm:"primaryKey"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
