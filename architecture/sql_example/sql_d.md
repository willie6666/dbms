# d. 進階查詢 I (使用 EXISTS, NOT EXISTS, NULL, UNION, >=, LIKE 等)

本文件收錄系統中使用到的第一階層進階 SQL 語法，包含模糊搜尋、數值範圍區間比對、以及利用子查詢 (Subquery) 進行存在性 (`EXISTS` / `NOT EXISTS`) 判斷的實際應用範例。

---

### 1. 使用 NOT EXISTS 過濾已購買的遊戲
- **說明**：玩家在瀏覽商店首頁時，可以勾選「隱藏已擁有的遊戲」。此時後端會利用 `NOT EXISTS` 子查詢去檢查 `game_licenses` 中是否已經有這款遊戲的購買紀錄，將買過的濾除。
- **對應 API**：`GET /api/games?hide_owned=true`
- **Go 實作 (GORM)**：
  ```go
  query = query.Where("NOT EXISTS (SELECT 1 FROM game_licenses WHERE game_licenses.game_id = games.game_id AND game_licenses.user_id = ? AND game_licenses.status = 'ACTIVE')", userID)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games 
  WHERE games.status = 'ACTIVE' 
    AND NOT EXISTS (
      SELECT 1 FROM game_licenses 
      WHERE game_licenses.game_id = games.game_id 
        AND game_licenses.user_id = 5 
        AND game_licenses.status = 'ACTIVE'
    );
  ```

### 2. 使用 ILIKE 進行多欄位模糊搜尋 (Keyword Search)
- **說明**：當玩家在搜尋框輸入關鍵字時，系統必須同時去比對遊戲標題、遊戲描述、標籤名稱以及開發者名稱。透過 `ILIKE` 能達到不分大小寫的模糊比對效果。
- **對應 API**：`GET /api/games?q={keyword}`
- **Go 實作 (GORM)**：
  ```go
  keyword := "%" + q + "%"
  query = query.Where("games.title ILIKE ? OR games.description ILIKE ? OR filter_tags.tag_name ILIKE ? OR filter_developers.username ILIKE ?", keyword, keyword, keyword, keyword)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT games.* FROM games
  LEFT JOIN game_tags ON game_tags.game_id = games.game_id
  LEFT JOIN tags ON tags.tag_id = game_tags.tag_id
  LEFT JOIN users ON users.user_id = games.developer_id
  WHERE games.status = 'ACTIVE' AND (
    games.title ILIKE '%戰神%' OR 
    games.description ILIKE '%戰神%' OR 
    tags.tag_name ILIKE '%戰神%' OR 
    users.username ILIKE '%戰神%'
  );
  ```

### 3. 使用 >= 與 <= 進行價格區間過濾
- **說明**：商店支援透過設定最低與最高價格區間來過濾遊戲清單，幫助玩家找到符合預算的遊戲。
- **對應 API**：`GET /api/games?min_price=100&max_price=500`
- **Go 實作 (GORM)**：
  ```go
  query = query.Where("games.price >= ?", minPrice)
  query = query.Where("games.price <= ?", maxPrice)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games 
  WHERE status = 'ACTIVE' 
    AND price >= 100.00 
    AND price <= 500.00;
  ```

### 4. 後台管理員使用 ILIKE 搜尋使用者
- **說明**：管理員在後台想要尋找特定玩家時，可以使用模糊搜尋同時比對 `username` 與 `email` 欄位。
- **對應 API**：`GET /api/admin/users?q={keyword}`
- **Go 實作 (GORM)**：
  ```go
  keyword := "%" + q + "%"
  database.DB.Where("username ILIKE ? OR email ILIKE ?", keyword, keyword).Order("registration_date DESC").Find(&users)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users 
  WHERE username ILIKE '%test%' OR email ILIKE '%test%'
  ORDER BY registration_date DESC;
  ```

### 5. 確保審核通過的遊戲才顯示 (使用 != 不等於)
- **說明**：這是隱含在所有前台顯示邏輯中的條件操作，確保被下架 (`TAKEN_DOWN`) 或仍是草稿 (`DRAFT`) 的遊戲不會被一般玩家搜尋到。
- **對應 API**：(所有前台 `GET /api/games` 相關路由)
- **Go 實作 (GORM)**：
  ```go
  query = query.Where("games.status != 'TAKEN_DOWN'")
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE status != 'TAKEN_DOWN';
  ```
