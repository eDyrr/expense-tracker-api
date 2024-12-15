package models

import "gorm.io/gorm"

type Purchase struct {
	gorm.Model
	Name     string  `json:"name"`
	Cost     float32 `json:"cost"`
	Category string  `json:"category"`
	UserID   uint    `json:"user_id"`
}
