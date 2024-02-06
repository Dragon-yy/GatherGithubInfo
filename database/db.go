package database

import (
	"GatherGithubInfo/config"
	"GatherGithubInfo/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(ip, port, usernmae, password, db string) error {
	// Replace "DB_USER", "DB_PASSWORD", "DB_NAME", and "DB_HOST" with your actual MySQL credentials
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", usernmae, password, ip, port, db)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// AutoMigrate will create the user table if it doesn't exist
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}

func SaveUser(user *models.User) error {
	// Save the user to the database
	// Connect to the database and set connection pool
	dsh := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DatabaseUser, config.DatabasePassword, config.DatabaseIP, config.DatabasePort, config.DatabaseName)
	DB, err := gorm.Open(mysql.Open(dsh), &gorm.Config{})
	if err != nil {
		return err
	}
	result := DB.Create(user)
	if result.Error != nil {
		return result.Error
	}

	fmt.Println("User saved successfully.")
	return nil
}
