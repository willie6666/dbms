# e. 進階查詢 II (使用 ORDER BY, IN, MAX/MIN/AVG/SUM/COUNT, GROUP BY, HAVING 等)

本文件收錄系統中使用到的第二階層進階 SQL 語法，包含群組聚合運算 (Aggregation)、排序、以及清單包含判斷 (`IN`)。

---

### 1. 遊戲評分統計 (使用 COUNT, AVG 與 CASE WHEN)
- **說明**：當玩家新增或刪除遊戲評論時，系統會動態觸發此查詢，針對該款遊戲的所有可見評論 (`VISIBLE`) 進行統計。利用 `COUNT` 算出總評論數，再利用 `AVG` 搭配 `CASE WHEN` 將正負評轉換為具體的 0.0 ~ 5.0 平均分數，最後將這個分數存回 `games.overall_rating`。
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
- **說明**：在商店查詢時，為了避免因為 JOIN 多個標籤導致同一款遊戲重複出現，必須使用 `GROUP BY games.game_id` 將結果合併。同時也利用 `ORDER BY` 實現依照價格由高至低或發售日期排列的功能。
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
- **說明**：當玩家按下結帳按鈕，系統會將購物車內的多個 `game_id` 集合成一個陣列，並使用 `IN` 語法一次性向資料庫拉出所有遊戲最新的價格與狀態，以防在結帳瞬間有遊戲被下架或改價。
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
- **說明**：開發者在查看某款遊戲的統計數據時，系統會 JOIN 交易明細表，並利用聚合函數 `SUM` 計算總營業額、`COUNT` 統計總共賣出的套數。
- **對應 API**：`GET /api/developer/games/:id/stats`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Table("transaction_items ti").
      Select("COUNT(ti.item_id) as total_sales_count, COALESCE(SUM(ti.purchase_price), 0) as total_revenue").
      Joins("JOIN games g ON ti.game_id = g.game_id").
      Where("g.developer_id = ? AND g.game_id = ?", developerID, gameID).
      Row().Scan(&stats.TotalSalesCount, &stats.TotalRevenue)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT 
    COUNT(ti.item_id) as total_sales_count,
    COALESCE(SUM(ti.purchase_price), 0) as total_revenue
  FROM transaction_items ti
  JOIN games g ON ti.game_id = g.game_id
  WHERE g.developer_id = 5 AND g.game_id = 42;
  ```

### 5. 後台使用者列表排序 (ORDER BY)
- **說明**：管理員檢視全站使用者名單時，系統預設將最新註冊的帳號排在最前面，方便管理與監控近期湧入的會員。
- **對應 API**：`GET /api/admin/users`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Order("registration_date DESC").Find(&users)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users ORDER BY registration_date DESC;
  ```

### 6. 退款紀錄列表排序 (ORDER BY 多欄位)
- **說明**：客服 (CSR) 介面需要列出退款申請，為確保客服第一時間看到亟需處理的項目，會進行多重欄位排序。此語法會先將 `status` 為 `PENDING` (待處理) 的單子排前面，並且依照申請時間正序排列 (最老的待處理單優先)。
- **對應 API**：`GET /api/csr/refunds`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Order("status DESC").Order("created_at ASC").Find(&refunds)
  ```
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests ORDER BY status DESC, created_at ASC;
  ```
