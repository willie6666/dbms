# c. 使用三個資料表的查詢 (Three-Table JOIN Queries)

本文件列出系統中橫跨三個或以上資料表的複雜關聯查詢 (JOIN)。這些查詢通常發生在多對多關聯、或是一對多對多的深層結構中。

---

### 1. 透過特定標籤搜尋遊戲 (多對多關聯)
- **說明**：玩家點擊某個標籤 (例如 RPG) 來找遊戲時，由於這是一個多對多關係，必須先從 `games` 表 JOIN 中介表 `game_tags`，再 JOIN 標籤主檔 `tags` 來做名稱比對。
- **對應 API**：`GET /api/games?tag={name}`
- **Go 實作 (GORM)**：
  ```go
  query = query.
      Joins("JOIN game_tags filter_game_tags ON filter_game_tags.game_id = games.game_id").
      Joins("JOIN tags filter_tags ON filter_tags.tag_id = filter_game_tags.tag_id").
      Where("filter_tags.tag_name ILIKE ?", tag)
  ```
- **原生 SQL 語法 (連續 INNER JOIN)**：
  ```sql
  SELECT games.* 
  FROM games
  JOIN game_tags ON game_tags.game_id = games.game_id
  JOIN tags ON tags.tag_id = game_tags.tag_id
  WHERE games.status = 'ACTIVE' 
    AND tags.tag_name ILIKE 'RPG';
  ```

### 2. 查詢玩家的遊戲庫並包含遊戲封面圖
- **說明**：玩家的遊戲庫紀錄了買過的遊戲 (`game_licenses`)。為了在前端渲染出遊戲圖示，必須先串接 `games` 主檔，再串接 `game_media` 取出附屬的圖片，總共橫跨三個表。
- **對應 API**：`GET /api/protected/library`
- **Go 實作 (GORM)**：
  ```go
  var licenses []models.GameLicense
  database.DB.Preload("Game.Media").Where("user_id = ? AND status = 'ACTIVE'", userID).Find(&licenses)
  ```
- **原生 SQL 語法 (JOIN + LEFT JOIN 等效邏輯)**：
  ```sql
  SELECT game_licenses.license_id, games.title, game_media.file_url 
  FROM game_licenses 
  JOIN games ON game_licenses.game_id = games.game_id 
  LEFT JOIN game_media ON games.game_id = game_media.game_id 
  WHERE game_licenses.user_id = 5 
    AND game_licenses.status = 'ACTIVE';
  ```

### 3. 查詢歷史訂單明細與對應的遊戲名稱
- **說明**：這是標準的電商一對多對一三表關聯。從主訂單表 (`transactions`) 關聯出底下的購物明細項目 (`transaction_items`)，再關聯到遊戲主檔 (`games`) 來顯示玩家到底買了哪些遊戲。
- **對應 API**：`GET /api/protected/transactions`
- **Go 實作 (GORM)**：
  ```go
  var transactions []models.Transaction
  database.DB.Preload("Items.Game").Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions)
  ```
- **原生 SQL 語法 (連續 JOIN 等效邏輯)**：
  ```sql
  SELECT transactions.transaction_id, transactions.created_at, 
         transaction_items.purchase_price, games.title
  FROM transactions
  JOIN transaction_items ON transactions.transaction_id = transaction_items.transaction_id
  JOIN games ON transaction_items.game_id = games.game_id
  WHERE transactions.user_id = 5
  ORDER BY transactions.created_at DESC;
  ```

### 4. 下載遊戲檔案前之授權與實體檔案驗證
- **說明**：當玩家點擊下載遊戲時，系統必須橫跨三個表進行確認：1) 玩家真的擁有這款遊戲 (`game_licenses`)，2) 該遊戲沒被硬刪除 (`games`)，3) 該遊戲有上傳對應的實體安裝檔 (`game_media`) 供下載。
- **對應 API**：`GET /api/protected/library/:game_id/download`
- **Go 實作 (GORM)**：
  ```go
  var license models.GameLicense
  // 驗證授權 (Table 1: game_licenses, Table 2: games)
  database.DB.Joins("Game").Where("user_id = ? AND game_licenses.game_id = ? AND game_licenses.status = 'ACTIVE'", userID, gameID).First(&license)

  var media models.GameMedia
  // 尋找實體檔案 (Table 3: game_media)
  database.DB.Where("game_id = ? AND media_type = 'GAME_FILE'", gameID).First(&media)
  ```
- **原生 SQL 語法 (邏輯上等同三表驗證 JOIN)**：
  ```sql
  SELECT game_licenses.license_id, games.game_id, game_media.file_url 
  FROM game_licenses
  JOIN games ON game_licenses.game_id = games.game_id
  JOIN game_media ON games.game_id = game_media.game_id
  WHERE game_licenses.user_id = 5 
    AND game_licenses.game_id = 42 
    AND game_licenses.status = 'ACTIVE'
    AND game_media.media_type = 'GAME_FILE';
  ```
