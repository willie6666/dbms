# a. 僅使用單一資料表的查詢 (Single Table Queries)

本文件列出系統中單純對單一資料表進行檢索，不牽涉 `JOIN` 的所有 API 查詢。這些查詢通常用於取得單一資源、全域列表，或進行簡單的狀態確認。

---

### 1. 取得全域標籤列表
- **對應 API**：`GET /api/tags`
- **Go 實作 (GORM)**：
  ```go
  var tags []models.Tag
  database.DB.Find(&tags)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM tags;
  ```

### 2. 使用者登入檢查 (依 Email 查詢)
- **對應 API**：`POST /api/auth/login`
- **Go 實作 (GORM)**：
  ```go
  var user models.User
  database.DB.Where("email = ?", input.Email).First(&user)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users WHERE email = 'user@example.com' ORDER BY user_id LIMIT 1;
  ```

### 3. 使用者登出 (狀態確認)
- **對應 API**：`POST /api/auth/logout`
- **Go 實作 (GORM)**：
  *(無資料庫查詢，單純清除客戶端 Token)*

### 4. 取得單一遊戲基本資訊
- **對應 API**：`GET /api/games/{id}`
- **Go 實作 (GORM)**：
  ```go
  var game models.Game
  database.DB.Where("game_id = ? AND status != 'TAKEN_DOWN'", id).First(&game)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE game_id = 42 AND status != 'TAKEN_DOWN' ORDER BY game_id LIMIT 1;
  ```

### 5. 查看自己的好友名單
- **對應 API**：`GET /api/social/friends`
- **Go 實作 (GORM)**：
  ```go
  var friendships []models.Friendship
  database.DB.Where("user_id1 = ? OR user_id2 = ?", userID, userID).Find(&friendships)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM friendships WHERE user_id1 = 5 OR user_id2 = 5;
  ```

### 6. 查看未處理的好友邀請
- **對應 API**：`GET /api/social/friends/requests`
- **Go 實作 (GORM)**：
  ```go
  var requests []models.FriendRequest
  database.DB.Where("receiver_id = ? AND status = 'PENDING'", userID).Find(&requests)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM friend_requests WHERE receiver_id = 5 AND status = 'PENDING';
  ```

### 7. 查看自己的黑名單
- **對應 API**：`GET /api/social/blacklist`
- **Go 實作 (GORM)**：
  ```go
  var blacklist []models.Blocklist
  database.DB.Where("blocker_id = ?", userID).Find(&blacklist)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM blocklist WHERE blocker_id = 5;
  ```

### 8. 管理員獲取全站使用者名單
- **對應 API**：`GET /api/admin/users`
- **Go 實作 (GORM)**：
  ```go
  var users []models.User
  database.DB.Order("created_at DESC").Find(&users)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users ORDER BY created_at DESC;
  ```

### 9. 客服獲取所有待處理退款單
- **對應 API**：`GET /api/csr/refunds`
- **Go 實作 (GORM)**：
  ```go
  var requests []models.RefundRequest
  database.DB.Order("created_at DESC").Find(&requests)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests ORDER BY created_at DESC;
  ```

### 10. 開發者獲取自己發布的所有遊戲
- **對應 API**：`GET /api/developer/games`
- **Go 實作 (GORM)**：
  ```go
  var games []models.Game
  database.DB.Where("developer_id = ?", userID).Order("created_at DESC").Find(&games)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE developer_id = 5 ORDER BY created_at DESC;
  ```

### 11. 檢查交易明細是否存在 (申請退款前)
- **對應 API**：`POST /api/social/refunds`
- **Go 實作 (GORM)**：
  ```go
  var item models.TransactionItem
  database.DB.First(&item, input.TransactionItemID)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM transaction_items WHERE transaction_item_id = 123 ORDER BY transaction_item_id LIMIT 1;
  ```

### 12. 取得個人退款歷史紀錄
- **對應 API**：`GET /api/protected/refunds`
- **Go 實作 (GORM)**：
  ```go
  var refunds []models.RefundRequest
  database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&refunds)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests WHERE user_id = 5 ORDER BY created_at DESC;
  ```

### 13. 驗證玩家是否擁有該遊戲遊玩權限
- **對應 API**：`GET /api/protected/library/{game_id}/play`
- **Go 實作 (GORM)**：
  ```go
  var license models.GameLicense
  database.DB.Where("user_id = ? AND game_id = ? AND status = 'ACTIVE'", userID, gameID).First(&license)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM game_licenses WHERE user_id = 5 AND game_id = 42 AND status = 'ACTIVE' ORDER BY license_id LIMIT 1;
  ```

### 14. 顯示與某使用者的對話紀錄
- **對應 API**：`GET /api/social/messages/{user_id}`
- **Go 實作 (GORM)**：
  ```go
  var messages []models.Message
  database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", myID, peerID, peerID, myID).Order("created_at ASC").Find(&messages)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM messages WHERE (sender_id = 5 AND receiver_id = 10) OR (sender_id = 10 AND receiver_id = 5) ORDER BY created_at ASC;
  ```
