# e. 進階查詢 II (使用 ORDER BY, IN, MAX/MIN/AVG/SUM/COUNT, GROUP BY, HAVING 等)

本文件收錄系統中使用到的第二階層進階 SQL 語法，包含群組聚合運算 (Aggregation)、排序、以及清單包含判斷 (`IN`)。

---

### 1. 遊戲評分統計 (使用 COUNT, AVG 與 CASE WHEN)
- **對應功能**：新增/刪除評論時觸發的重算邏輯 (內部呼叫)
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
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games 
  WHERE game_id IN (14, 25, 42, 58);
  ```

### 4. 開發者總銷售額與銷量統計 (使用 SUM, COUNT)
- **對應 API**：`GET /api/developer/games/:id/stats`
- **原生 SQL 語法**：
  ```sql
  SELECT 
    COUNT(ti.transaction_item_id) as total_sales_count,
    COALESCE(SUM(ti.purchase_price), 0) as total_revenue
  FROM transaction_items ti
  JOIN games g ON ti.game_id = g.game_id
  WHERE g.developer_id = 5 AND g.game_id = 42;
  ```
- **說明**：為開發者儀表板提供銷售數據，使用 `SUM` 將該遊戲所有歷史銷售單價加總計算營收。

### 5. 後台使用者列表排序 (ORDER BY)
- **對應 API**：`GET /api/admin/users`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users ORDER BY created_at DESC;
  ```

### 6. 退款紀錄列表排序 (ORDER BY)
- **對應 API**：`GET /api/csr/refunds`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests ORDER BY status DESC, created_at ASC;
  ```
- **說明**：將待處理 (`PENDING`) 的退款單排在前面，已處理的排在後面。
