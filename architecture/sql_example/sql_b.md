# b. 使用兩個資料表的查詢 (Two-Table JOIN Queries)

本文件列出所有需要跨越兩個關聯資料表 (Table) 來取得資料的 API 查詢。透過 JOIN，我們能將代碼 (如 ID) 轉換為人類可讀的名稱或帶出關聯狀態。

---

### 1. 查詢購物車內容與遊戲資訊
- **說明**：玩家打開購物車時，系統需從 `shopping_carts` 取得紀錄，並關聯 `games` 表以顯示遊戲名稱、價格與上架狀態。
- **對應 API**：`GET /api/protected/cart`
- **Go 實作 (GORM)**：
  ```go
  var cart []models.ShoppingCart
  database.DB.Preload("Game").Where("user_id = ?", userID).Find(&cart)
  ```
- **原生 SQL 語法 (INNER JOIN 等效邏輯)**：
  ```sql
  SELECT shopping_carts.*, games.title, games.price, games.status, games.discount
  FROM shopping_carts
  JOIN games ON shopping_carts.game_id = games.game_id
  WHERE shopping_carts.user_id = 5;
  ```

### 2. 查詢願望清單與遊戲資訊
- **說明**：玩家查看願望清單時，同樣需關聯 `games` 表，才能在畫面上呈現出他關注的遊戲之最新價格或狀態變動。
- **對應 API**：`GET /api/protected/wishlist`
- **Go 實作 (GORM)**：
  ```go
  var wishlist []models.WishList
  database.DB.Preload("Game").Where("user_id = ?", userID).Find(&wishlist)
  ```
- **原生 SQL 語法 (INNER JOIN 等效邏輯)**：
  ```sql
  SELECT wish_lists.*, games.title, games.price, games.status
  FROM wish_lists
  JOIN games ON wish_lists.game_id = games.game_id
  WHERE wish_lists.user_id = 5;
  ```

### 3. 獲取遊戲的評論列表並包含評論者名稱
- **說明**：載入遊戲評論區時，除了評論內容本身 (`reviews` 表)，還必須 JOIN `users` 資料表來顯示這則評論是誰留的，以及該玩家的暱稱。
- **對應 API**：`GET /api/games/:id/reviews`
- **Go 實作 (GORM)**：
  ```go
  var reviews []models.Review
  database.DB.Preload("User").Where("game_id = ? AND status = 'VISIBLE'", gameID).Order("created_at DESC").Find(&reviews)
  ```
- **原生 SQL 語法 (INNER JOIN 等效邏輯)**：
  ```sql
  SELECT reviews.*, users.username, users.avatar_url
  FROM reviews
  JOIN users ON reviews.user_id = users.user_id
  WHERE reviews.game_id = 42 AND reviews.status = 'VISIBLE'
  ORDER BY reviews.created_at DESC;
  ```

### 4. 獲取評論的獨立回覆列表
- **說明**：如同遊戲評論，展開某篇評論底下的「樓中樓回覆」時，也需要 JOIN `users` 表來取得回覆者的名字，方便使用者辨識討論對象。
- **對應 API**：`GET /api/social/reviews/:id/replies`
- **Go 實作 (GORM)**：
  ```go
  var replies []models.ReviewReply
  database.DB.Preload("User").Where("review_id = ?", reviewID).Order("created_at ASC").Find(&replies)
  ```
- **原生 SQL 語法 (INNER JOIN 等效邏輯)**：
  ```sql
  SELECT review_replies.*, users.username, users.avatar_url
  FROM review_replies
  JOIN users ON review_replies.user_id = users.user_id
  WHERE review_replies.review_id = 128
  ORDER BY review_replies.created_at ASC;
  ```

### 5. 商店搜尋由特定開發者發布的遊戲
- **說明**：玩家若想查看某間工作室或開發者所發布的所有作品，系統會將 `games` 表與 `users` 表 JOIN，根據開發者的名稱字串進行篩選。
- **對應 API**：`GET /api/games?developer={username}`
- **Go 實作 (GORM)**：
  ```go
  query = query.
      Joins("JOIN users filter_developers ON filter_developers.user_id = games.developer_id").
      Where("filter_developers.username ILIKE ?", "%"+developer+"%")
  ```
- **原生 SQL 語法 (真實 INNER JOIN)**：
  ```sql
  SELECT games.* 
  FROM games
  JOIN users filter_developers ON filter_developers.user_id = games.developer_id
  WHERE games.status = 'ACTIVE' AND filter_developers.username ILIKE '%CDProjekt%';
  ```
