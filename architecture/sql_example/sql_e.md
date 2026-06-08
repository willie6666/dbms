# e. 進階查詢 II (使用 ORDER BY, IN, MAX/MIN/AVG/SUM/COUNT, GROUP BY, HAVING 等)

本文件收錄系統中使用到的第二階層進階 SQL 語法，包含群組聚合運算 (Aggregation)、排序、以及清單包含判斷 (`IN`)。

---

### 1. 遊戲評分統計 (使用 COUNT, AVG 與 CASE WHEN)
- **對應功能**：新增/刪除評論時觸發的重算邏輯 (內部呼叫 `refreshStoredGameRating`)
- **Go 實作 (GORM)**：
  ```go
  database.DB.Model(&models.Review{}).
      Where("game_id = ? AND status = 'VISIBLE'", gameID).
      Select("COUNT(*) as total_reviews, COALESCE(AVG(CASE WHEN attitude = 'POSITIVE' THEN 5.0 ELSE 1.0 END), 0) as average_rating").
      Row().Scan(&totalReviews, &averageRating)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT 
    COUNT(*) as total_reviews, 
    COALESCE(AVG(CASE WHEN attitude = 'POSITIVE' THEN 5.0 ELSE 1.0 END), 0) as average_rating 
  FROM reviews 
  WHERE game_id = 42 AND status = 'VISIBLE';
  ```

### 2. 商店首頁的分群與自訂排序 (使用 GROUP BY 與 ORDER BY)
- **對應 API**：`GET /api/games?sort={sort_type}`
- **Go 實作 (GORM)**：
  ```go
  query = query.Group("games.game_id")
  if sort == "price_desc" {
      query = query.Order("games.price DESC")
  } else {
      query = query.Order("games.release_date DESC")
  }
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT games.* FROM games 
  LEFT JOIN game_tags ON game_tags.game_id = games.game_id
  WHERE games.status = 'ACTIVE' 
  GROUP BY games.game_id
  ORDER BY games.price DESC;
  ```

### 3. 結帳前的購物車商品清單校驗 (使用 IN)
- **對應 API**：`POST /api/shopping/checkout`
- **Go 實作 (GORM)**：
  ```go
  var games []models.Game
  database.DB.Where("game_id IN ?", gameIDs).Find(&games)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games 
  WHERE game_id IN (14, 25, 42, 58);
  ```

### 4. 開發者總銷售額與銷量統計 (使用 SUM, COUNT)
- **對應 API**：`GET /api/developer/games/:id/stats`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Table("transaction_items ti").
      Select("COUNT(ti.transaction_item_id) as total_sales_count, COALESCE(SUM(ti.purchase_price), 0) as total_revenue").
      Joins("JOIN games g ON ti.game_id = g.game_id").
      Where("g.developer_id = ? AND g.game_id = ?", developerID, gameID).
      Row().Scan(&stats.TotalSalesCount, &stats.TotalRevenue)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT 
    COUNT(ti.transaction_item_id) as total_sales_count,
    COALESCE(SUM(ti.purchase_price), 0) as total_revenue
  FROM transaction_items ti
  JOIN games g ON ti.game_id = g.game_id
  WHERE g.developer_id = 5 AND g.game_id = 42;
  ```

### 5. 後台使用者列表排序 (ORDER BY)
- **對應 API**：`GET /api/admin/users`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Order("created_at DESC").Find(&users)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users ORDER BY created_at DESC;
  ```

### 6. 退款紀錄列表排序 (ORDER BY 多欄位)
- **對應 API**：`GET /api/csr/refunds`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Order("status DESC").Order("created_at ASC").Find(&refunds)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests ORDER BY status DESC, created_at ASC;
  ```
- **說明**：將待處理 (`PENDING`) 的退款單排在前面，已處理的排在後面。
