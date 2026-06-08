# b. 使用兩個資料表的查詢 (Two-Table JOIN Queries)

本文件列出所有需要跨越兩個關聯資料表 (Table) 來取得資料的 API 查詢。透過 JOIN，我們能將代碼 (如 ID) 轉換為人類可讀的名稱或帶出關聯狀態。

---

### 1. 查詢購物車內容與遊戲資訊
- **對應 API**：`GET /api/protected/cart`
- **原生 SQL 語法 (INNER JOIN)**：
  ```sql
  SELECT shopping_carts.*, games.title, games.price, games.status, games.discount
  FROM shopping_carts
  JOIN games ON shopping_carts.game_id = games.game_id
  WHERE shopping_carts.user_id = 5;
  ```

### 2. 查詢願望清單與遊戲資訊
- **對應 API**：`GET /api/protected/wishlist`
- **原生 SQL 語法 (INNER JOIN)**：
  ```sql
  SELECT wish_lists.*, games.title, games.price, games.status
  FROM wish_lists
  JOIN games ON wish_lists.game_id = games.game_id
  WHERE wish_lists.user_id = 5;
  ```

### 3. 獲取遊戲的評論列表並包含評論者名稱
- **對應 API**：`GET /api/games/:id/reviews`
- **原生 SQL 語法 (INNER JOIN)**：
  ```sql
  SELECT reviews.*, users.username, users.avatar_url
  FROM reviews
  JOIN users ON reviews.user_id = users.user_id
  WHERE reviews.game_id = 42 AND reviews.status = 'VISIBLE'
  ORDER BY reviews.created_at DESC;
  ```
- **說明**：載入遊戲評論區時，除了評論內容本身，還必須 JOIN `users` 資料表來顯示這則評論是誰留的、以及他的大頭貼。

### 4. 獲取評論的獨立回覆列表
- **對應 API**：`GET /api/social/reviews/:id/replies`
- **原生 SQL 語法 (INNER JOIN)**：
  ```sql
  SELECT review_replies.*, users.username, users.avatar_url
  FROM review_replies
  JOIN users ON review_replies.user_id = users.user_id
  WHERE review_replies.review_id = 128
  ORDER BY review_replies.created_at ASC;
  ```
- **說明**：與評論本身相同，留言回覆也必須關聯 `users` 表來顯示回覆者。

### 5. 商店搜尋由特定開發者發布的遊戲
- **對應 API**：`GET /api/games?developer={username}`
- **原生 SQL 語法 (INNER JOIN)**：
  ```sql
  SELECT games.* 
  FROM games
  JOIN users filter_developers ON filter_developers.user_id = games.developer_id
  WHERE games.status = 'ACTIVE' AND filter_developers.username ILIKE '%CDProjekt%';
  ```
