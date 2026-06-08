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

	// Move legacy frontend-hosted media paths to backend-hosted media/download URLs.
	database.DB.Exec("ALTER TABLE games ADD COLUMN IF NOT EXISTS description TEXT DEFAULT ''")
	database.DB.Exec("UPDATE game_media SET file_url = ? WHERE game_id = 1 AND media_type = 'game_file'", "/downloads/cybercity-demo.txt")
	database.DB.Exec("UPDATE game_media SET file_url = ? WHERE game_id = 4 AND media_type = 'game_file'", "/downloads/farm-simulator-demo.txt")
	database.DB.Exec("UPDATE game_media SET file_url = regexp_replace(file_url, '^https?://[^/]+', '') WHERE file_url ~ '^https?://[^/]+/media/images/' AND media_type <> 'game_file'")
	database.DB.Exec("UPDATE game_media SET file_url = REPLACE(file_url, '/assets/images/', '/media/images/') WHERE file_url LIKE '/assets/images/%' AND media_type <> 'game_file'")
	database.DB.Exec("UPDATE game_media SET file_url = regexp_replace(file_url, '^https?://[^/]+', '') WHERE file_url ~ '^https?://[^/]+/downloads/' AND media_type = 'game_file'")
	database.DB.Exec("UPDATE game_media SET file_url = ? WHERE file_url LIKE ?", "/media/images/protoss_12+8.png", "%protoss_knife.png")
	database.DB.Exec("UPDATE game_media SET file_url = ? WHERE file_url LIKE ?", "/media/images/protoss_cross.png", "%protoss_best_16.png")

	// Fix existing TAKEN_DOWN games to have REVOKED licenses
	database.DB.Exec("UPDATE game_licenses SET status = 'REVOKED' FROM games WHERE game_licenses.game_id = games.game_id AND games.status = 'TAKEN_DOWN'")

	// Setup Gin router
	r := routes.SetupRouter()

	// Start server on port 8000
	log.Println("Starting server on port 8000...")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
