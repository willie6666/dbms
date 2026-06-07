package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// The port is 5432 and user/password are admin/admin based on docker-compose.yml
	dsn := "host=localhost user=admin password=admin dbname=vapor_auror port=5432 sslmode=disable TimeZone=Asia/Taipei"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = database

	fmt.Println("Database connection successfully opened")
}
