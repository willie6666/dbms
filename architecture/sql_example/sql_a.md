# a. 僅使用單一資料表的查詢 (Single Table Queries)

本文件列出系統中單純對單一資料表進行檢索，不牽涉 `JOIN` 的所有 API 查詢。這些查詢通常用於取得單一資源、全域列表，或進行簡單的狀態確認。

---

### 1. 取得全域標籤列表
- **對應 API**：`GET /api/tags`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM tags;
  ```

### 2. 使用者登入檢查 (依 Email 查詢)
- **對應 API**：`POST /api/auth/login`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users WHERE email = 'user@example.com' ORDER BY user_id LIMIT 1;
  ```

### 3. 獲取當前登入者資訊
- **對應 API**：`GET /api/auth/me`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users WHERE user_id = 5 ORDER BY user_id LIMIT 1;
  ```

### 4. 查看公開的玩家個人檔案
- **對應 API**：`GET /api/users/:id`
- **原生 SQL 語法**：
  ```sql
  SELECT user_id, username, avatar_url, bio, role, created_at 
  FROM users WHERE user_id = 10 ORDER BY user_id LIMIT 1;
  ```

### 5. 取得單一遊戲基本資訊
- **對應 API**：`GET /api/games/:id`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE game_id = 42 AND status = 'ACTIVE' ORDER BY game_id LIMIT 1;
  ```

### 6. 查看自己的好友名單
- **對應 API**：`GET /api/social/friends`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM friendships WHERE user_id1 = 5 OR user_id2 = 5;
  ```
- **說明**：由於好友關係在資料庫中是無方向性的，因此查詢時必須使用 `OR` 判斷當前使用者是 `user_id1` 或 `user_id2`。

### 7. 查看未處理的好友邀請
- **對應 API**：`GET /api/social/friends/requests`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM friend_requests WHERE receiver_id = 5 AND status = 'PENDING';
  ```

### 8. 查看自己的黑名單
- **對應 API**：`GET /api/social/blocklist`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM blocklist WHERE blocker_id = 5;
  ```

### 9. 管理員獲取全站使用者名單
- **對應 API**：`GET /api/admin/users`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM users ORDER BY created_at DESC;
  ```

### 10. 客服獲取所有待處理退款單
- **對應 API**：`GET /api/csr/refunds`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM refund_requests ORDER BY created_at DESC;
  ```

### 11. 開發者獲取自己發布的所有遊戲
- **對應 API**：`GET /api/developer/games`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM games WHERE developer_id = 5 ORDER BY created_at DESC;
  ```

### 12. 檢查交易明細是否存在 (退款前置檢查)
- **對應 API**：`POST /api/social/refunds`
- **原生 SQL 語法**：
  ```sql
  SELECT * FROM transaction_items WHERE transaction_item_id = 123 ORDER BY transaction_item_id LIMIT 1;
  ```
