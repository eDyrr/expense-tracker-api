package database

import (
	"fmt"
	"os"

	"github.com/eDyrr/expense-tracker-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {

	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/expensesDB?parseTime=true", os.Getenv("DB_CRED"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	DB = db
	db.AutoMigrate(&models.User{})

	return db, nil
}
