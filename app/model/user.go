package model

import "gorm.io/gorm"

// User モデル固有の定数
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

type User struct {
	gorm.Model
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"uniqueIndex"`
	Status string `json:"status"`
}
