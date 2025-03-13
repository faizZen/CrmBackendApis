package main

import (
	"log"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/auth"
	"github.com/Zenithive/it-crm-backend/internal/graphql"
	"github.com/joho/godotenv"
)

func init() {
	initializers.ConnectToDatabase()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	auth.InitGoogleStore()
}
func main() {
	graphql.Handler()
}
