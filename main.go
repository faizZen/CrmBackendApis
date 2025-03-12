package main

import (
	"log"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/internal/graphql"
	"github.com/joho/godotenv"
)

func init() {
	initializers.ConnectToDatabase()
}
func main() {
	graphql.Handler()
}
