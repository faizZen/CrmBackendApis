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
	DB.Exec(`CREATE TYPE resource_type AS ENUM ('CONSULTANT', 'FREELANCER', 'CONTRACTOR', 'EMPLOYEE');`)
	DB.Exec(`CREATE TYPE resource_status AS ENUM ('ACTIVE', 'INACTIVE', 'ON_BENCH');`)
	DB.Exec(`CREATE TYPE vendor_status AS ENUM ('ACTIVE', 'INACTIVE', 'PREFERRED');`)
	DB.Exec(`CREATE TYPE payment_terms AS ENUM ('NET_30', 'NET_60', 'NET_90');`)
	DB.Exec(`CREATE TYPE task_status AS ENUM ('TODO', 'IN_PROGRESS', 'COMPLETED', 'ON_HOLD');`)
	DB.Exec(`CREATE TYPE task_priority AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'URGENT');`)
	// DB.AutoMigrate(&models.Activity{})
	err = DB.AutoMigrate(
		&models.User{},
		&models.Campaign{},
		&models.Organization{},
		&models.Lead{},
		&models.Activity{},
		&models.Deals{},
		&models.ResourceProfile{},   // New Model
		&models.Vendor{},            // New Model
		&models.Skill{},             // Supporting model
		&models.PastProject{},       // Supporting model
		&models.Contact{},           // Supporting model
		&models.PerformanceRating{}, // Supporting model
		&models.Task{},              // Supporting model
		&models.ResourceSkill{},     // Supporting model
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}
