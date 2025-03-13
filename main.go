package main

import (
	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/auth"
	"github.com/Zenithive/it-crm-backend/internal/graphql"
)

func init() {
	initializers.ConnectToDatabase()

	auth.InitGoogleStore()
}
func main() {
	graphql.Handler()
}
