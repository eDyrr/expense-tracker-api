package models

import (
	"time"

	"github.com/eDyrr/expense-tracker-api/models"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string        `json:"name"`
	Email     string        `gorm:"unique" json:"email"`
	Password  []byte        `json:"-"`
	Items     []models.Item `gorm:"many2many:purchases"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
	DeletdAt  gorm.DeletedAt
}
