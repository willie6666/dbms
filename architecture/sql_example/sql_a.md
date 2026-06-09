# a. 僅使用單一資料表的查詢 (Single Table Queries)

本文件列出系統中單純對單一資料表進行檢索，不牽涉 `JOIN` 的所有 API 查詢。這些查詢通常用於取得單一資源、全域列表，或進行簡單的狀態確認。

---

### 1. 取得全域標籤列表
- **說明**：列出系統中所有可用的遊戲分類標籤，供開發者上架或玩家搜尋時使用。
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
- **說明**：使用者登入時，系統會先根據輸入的 Email 到資料庫中尋找對應的帳號紀錄，再進行密碼比對。
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
- **說明**：由於系統採用無狀態 JWT Token 架構，登出動作僅由前端負責清除 Token，後端不涉及資料庫的查詢與寫入。
- **對應 API**：`POST /api/auth/logout`
- **Go 實作 (GORM)**：
  *(無資料庫查詢，單純清除客戶端 Token)*

### 4. 取得單一遊戲基本資訊
- **說明**：玩家點擊進入遊戲介紹頁面時，讀取該遊戲的基本資料 (名稱、價格、介紹等)。同時過濾掉已被下架的遊戲。
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
- **說明**：查詢與自己有關，且狀態為已接受的好友關係紀錄。
- **對應 API**：`GET /api/social/friends`
- **Go 實作 (GORM)**：
  ```go
  var friendships []models.Friendship
  database.DB.Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "ACCEPTED").Order("created_at desc").Find(&friendships)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM friendships WHERE (sender_id = 5 OR receiver_id = 5) AND status = 'ACCEPTED' ORDER BY created_at DESC;
  ```

### 6. 查看未處理的好友邀請
- **說明**：列出別人發送給自己或是自己發送出去，且尚未同意或拒絕的好友邀請紀錄。
- **對應 API**：`GET /api/social/friends/requests`
- **Go 實作 (GORM)**：
  ```go
  var requests []models.Friendship
  database.DB.Where("(receiver_id = ? OR sender_id = ?) AND status = ?", userID, userID, "PENDING").Order("created_at desc").Find(&requests)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM friendships WHERE (receiver_id = 5 OR sender_id = 5) AND status = 'PENDING' ORDER BY created_at DESC;
  ```

### 7. 查看自己的黑名單
- **說明**：取得自己曾經封鎖過的使用者清單，供未來想要解除封鎖時參考。
- **對應 API**：`GET /api/social/blacklist`
- **Go 實作 (GORM)**：
  ```go
  var blacklist []models.Blacklist
  database.DB.Where("blocker_id = ?", userID).Order("created_at desc").Find(&blacklist)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM blacklists WHERE blocker_id = 5 ORDER BY created_at DESC;
  ```

### 8. 管理員獲取全站使用者名單
- **說明**：系統管理員在後台查看所有的註冊玩家與開發者名單，以時間倒序排列。
- **對應 API**：`GET /api/admin/users`
- **Go 實作 (GORM)**：
  ```go
  var users []models.User
  database.DB.Order("registration_date DESC").Find(&users)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users ORDER BY registration_date DESC;
  ```

### 9. 客服獲取所有待處理退款單
- **說明**：客服人員 (CSR) 登入後台時，列出所有玩家提出的退款申請，以便進行後續審核。
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
- **說明**：開發者進入自己的儀表板時，只會看到屬於自己開發的遊戲清單。
- **對應 API**：`GET /api/developer/games`
- **Go 實作 (GORM)**：
  ```go
  var games []models.Game
  database.DB.Where("developer_id = ?", userID).Order("created_at DESC").Find(&games)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE developer_id = 5 ORDER BY game_id DESC;
  ```

### 11. 檢查交易明細是否存在 (申請退款前)
- **說明**：玩家申請退款時，後端必須先驗證該筆交易明細是否真的存在且屬於該玩家。
- **對應 API**：`POST /api/social/refunds`
- **Go 實作 (GORM)**：
  ```go
  var item models.TransactionItem
  database.DB.First(&item, input.TransactionItemID)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM transaction_items WHERE item_id = 123 ORDER BY item_id LIMIT 1;
  ```

### 12. 取得個人退款歷史紀錄
- **說明**：玩家可以查看自己過去所有的退款申請紀錄與目前的審核狀態。
- **對應 API**：`GET /api/protected/refunds`
- **Go 實作 (GORM)**：
  ```go
  var refunds []models.RefundRequest
  database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&refunds)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests WHERE buyer_id = 5 ORDER BY created_at DESC;
  ```

### 13. 驗證玩家是否擁有該遊戲遊玩權限
- **說明**：玩家點擊「遊玩」時，後端必須檢查 `game_licenses` 授權表，確認玩家是否購買過且未被退款或撤銷。
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
- **說明**：查詢自己與另一名玩家的歷史私訊內容，並按照時間先後順序排列。
- **對應 API**：`GET /api/social/messages/{user_id}`
- **Go 實作 (GORM)**：
  ```go
  var messages []models.Message
  database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", myID, peerID, peerID, myID).Order("sent_at asc, message_id asc").Find(&messages)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM messages WHERE (sender_id = 5 AND receiver_id = 10) OR (sender_id = 10 AND receiver_id = 5) ORDER BY sent_at ASC, message_id ASC;
  ```
