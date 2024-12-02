package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name      string
	Price     float32
	Categorie string
	CreatedAt time.Time `gorm:"autoCreatedAt"`
	UpdatedAt time.Time `gorm:"autoUpdatedAt"`
	DeletedAt gorm.DeletedAt
}
