# d. 進階查詢 I (使用 EXISTS, NOT EXISTS, NULL, UNION, >=, LIKE 等)

本文件收錄系統中使用到的第一階層進階 SQL 語法，包含模糊搜尋、數值範圍區間比對、以及利用子查詢 (Subquery) 進行存在性 (`EXISTS` / `NOT EXISTS`) 判斷的實際應用範例。

---

### 1. 使用 NOT EXISTS 過濾已購買的遊戲
- **對應 API**：`GET /api/games?hide_owned=true`
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
- **對應 API**：`GET /api/games?q={keyword}`
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
- **對應 API**：`GET /api/games?min_price=100&max_price=500`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games 
  WHERE status = 'ACTIVE' 
    AND price >= 100.00 
    AND price <= 500.00;
  ```

### 4. 後台管理員使用 ILIKE 搜尋使用者
- **對應 API**：`GET /api/admin/users?q={keyword}`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users 
  WHERE username ILIKE '%test%' OR email ILIKE '%test%'
  ORDER BY created_at DESC;
  ```

### 5. 確保審核通過的遊戲才顯示 (常態條件限制)
- **對應 API**：(所有前台 `GET /api/games` 相關路由)
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE status != 'TAKEN_DOWN';
  ```
- **說明**：這是隱含的條件操作，確保被下架的遊戲不會被一般玩家搜尋到。
