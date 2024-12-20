package services

import (
	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/models"
)

func GetPurchases(id interface{}) []models.Purchase {
	var purchases []models.Purchase
	err := database.DB.Where("user_id = ?", id).Find(&purchases).Error
	if err != nil {
		panic("failed to load user purchases")
	}
	return purchases
}
