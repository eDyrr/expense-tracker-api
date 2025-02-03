package controllers

import (
	"fmt"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/models"
)

func DeletePurchase(purchaseID, userID uint) error {
	result := database.DB.Where("ID = ? AND user_id = ?", purchaseID, userID).Delete(&models.Purchase{})

	if result.RowsAffected == 0 {
		fmt.Println("no record found")
	}

	return result.Error
}
