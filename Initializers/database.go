package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/Zenithive/it-crm-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	godotenv.Load()
	dsn := os.Getenv("DB_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if DB == nil {
		fmt.Println("DB is nil")
	}
	// DB.AutoMigrate(&models.Activity{})
	err = DB.AutoMigrate(&models.Lead{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	err = DB.AutoMigrate(&models.Activity{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}
