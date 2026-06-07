package main

import (
	"log"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"
	"vapor_auror_backend/routes"
	"vapor_auror_backend/utils"
)

func main() {
	// Initialize database connection
	database.ConnectDB()

	// 臨時修復：將資料庫內明碼的假密碼轉換為真實的 Bcrypt 雜湊密碼
	// 讓預設測試帳號可以使用密碼 'admin' 登入
	var users []models.User
	database.DB.Find(&users)
	for _, u := range users {
		// 如果密碼不是以 $2a$ 開頭 (Bcrypt 格式)，就強制更新為 'admin' 的雜湊
		if len(u.PasswordHash) < 4 || u.PasswordHash[:4] != "$2a$" {
			hashed, _ := utils.HashPassword("admin")
			database.DB.Model(&u).Update("password_hash", hashed)
		}
	}

	// Setup Gin router
	r := routes.SetupRouter()

	// Start server on port 8000
	log.Println("Starting server on port 8000...")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
