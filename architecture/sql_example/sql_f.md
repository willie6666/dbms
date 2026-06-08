# f. 新增或刪除資料的操作 (INSERT, UPDATE, DELETE)

本文件列出系統進行資料異動 (Data Manipulation) 時使用的 SQL 操作。由於這涵蓋了系統中所有的狀態變更 API，以下依業務模組進行分類。

---

## 模組一：使用者與身分驗證 (User & Auth)

### 1. 玩家註冊 (INSERT)
- **對應 API**：`POST /api/auth/register`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO users (username, email, password_hash, role) 
  VALUES ('JohnDoe', 'john@test.com', '$2a$10$...', 'PLAYER') RETURNING user_id;
  ```

### 2. 更新個人檔案 (UPDATE)
- **對應 API**：`PUT /api/users/profile`
- **原生 SQL 語法**：
  ```sql
  UPDATE users SET bio = 'Hello World!', avatar_url = 'img.png' WHERE user_id = 5;
  ```

---

## 模組二：商城與交易 (Store, Cart, Wishlist, Transactions)

### 3. 將遊戲加入購物車 (INSERT)
- **對應 API**：`POST /api/protected/cart`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO shopping_carts (user_id, game_id) VALUES (5, 42);
  ```

### 4. 移除單一購物車商品 (DELETE)
- **對應 API**：`DELETE /api/protected/cart/:id`
- **原生 SQL 語法**：
  ```sql
  DELETE FROM shopping_carts WHERE user_id = 5 AND game_id = 42;
  ```

### 5. 清空購物車 (DELETE)
- **對應 API**：`DELETE /api/protected/cart`
- **原生 SQL 語法**：
  ```sql
  DELETE FROM shopping_carts WHERE user_id = 5;
  ```

### 6. 購物車結帳 (巨型 Transaction: INSERT + DELETE)
- **對應 API**：`POST /api/shopping/checkout`
- **原生 SQL 語法** (包裹在 `BEGIN` 與 `COMMIT` 之間)：
  ```sql
  INSERT INTO transactions (user_id, total_amount) VALUES (5, 1200.00) RETURNING transaction_id;
  
  INSERT INTO transaction_items (transaction_id, game_id, purchase_price) VALUES (99, 42, 1200.00);
  
  INSERT INTO game_licenses (user_id, game_id, transaction_item_id, status) VALUES (5, 42, 105, 'ACTIVE');
  
  DELETE FROM shopping_carts WHERE user_id = 5;
  ```

### 7. 將遊戲加入願望清單 (INSERT)
- **對應 API**：`POST /api/protected/wishlist`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO wish_lists (user_id, game_id) VALUES (5, 42);
  ```

### 8. 移除願望清單 (DELETE)
- **對應 API**：`DELETE /api/protected/wishlist/:id`
- **原生 SQL 語法**：
  ```sql
  DELETE FROM wish_lists WHERE user_id = 5 AND game_id = 42;
  ```

### 9. 玩家申請退款 (INSERT)
- **對應 API**：`POST /api/social/refunds`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO refund_requests (transaction_item_id, user_id, reason, status) 
  VALUES (105, 5, '不好玩', 'PENDING');
  ```

---

## 模組三：社群互動 (Social, Friends, Reviews)

### 10. 送出好友邀請 (INSERT)
- **對應 API**：`POST /api/social/friends/requests`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO friend_requests (sender_id, receiver_id, status) VALUES (5, 10, 'PENDING');
  ```

### 11. 接受好友邀請 (UPDATE + INSERT)
- **對應 API**：`PUT /api/social/friends/requests/:id` (附帶 status=ACCEPTED)
- **原生 SQL 語法**：
  ```sql
  UPDATE friend_requests SET status = 'ACCEPTED' WHERE request_id = 33;
  INSERT INTO friendships (user_id1, user_id2) VALUES (5, 10);
  ```

### 12. 刪除好友 (DELETE)
- **對應 API**：`DELETE /api/social/friends/:id`
- **原生 SQL 語法**：
  ```sql
  DELETE FROM friendships WHERE (user_id1 = 5 AND user_id2 = 10) OR (user_id1 = 10 AND user_id2 = 5);
  ```

### 13. 加入黑名單並解除好友 (INSERT + DELETE)
- **對應 API**：`POST /api/social/blocklist`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO blocklist (blocker_id, blocked_id) VALUES (5, 10);
  DELETE FROM friendships WHERE (user_id1 = 5 AND user_id2 = 10) OR (user_id1 = 10 AND user_id2 = 5);
  ```

### 14. 移除黑名單 (DELETE)
- **對應 API**：`DELETE /api/social/blocklist/:id`
- **原生 SQL 語法**：
  ```sql
  DELETE FROM blocklist WHERE blocker_id = 5 AND blocked_id = 10;
  ```

### 15. 發布遊戲評論 (INSERT)
- **對應 API**：`POST /api/games/:id/reviews`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO reviews (game_id, user_id, rating, attitude, content, status) 
  VALUES (42, 5, 0, 'POSITIVE', '超好玩', 'VISIBLE');
  ```

### 16. 修改遊戲評論 (UPDATE)
- **對應 API**：`PUT /api/social/reviews/:id`
- **原生 SQL 語法**：
  ```sql
  UPDATE reviews SET content = '玩久了有點膩', attitude = 'NEGATIVE' WHERE review_id = 128 AND user_id = 5;
  ```

### 17. 刪除遊戲評論 (UPDATE 軟刪除)
- **對應 API**：`DELETE /api/social/reviews/:id`
- **原生 SQL 語法**：
  ```sql
  UPDATE reviews SET status = 'HIDDEN' WHERE review_id = 128 AND user_id = 5;
  ```

### 18. 回覆別人的評論 (INSERT)
- **對應 API**：`POST /api/social/reviews/:id/replies`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO review_replies (review_id, user_id, content) VALUES (128, 10, '+1 認同');
  ```

### 19. 編輯/刪除回覆 (UPDATE / DELETE)
- **對應 API**：`PUT /api/social/replies/:id` 與 `DELETE /api/social/replies/:id`
- **原生 SQL 語法**：
  ```sql
  UPDATE review_replies SET content = '修改後的留言' WHERE reply_id = 56;
  DELETE FROM review_replies WHERE reply_id = 56;
  ```

---

## 模組四：開發者功能 (Developer)

### 20. 發行新遊戲 (INSERT)
- **對應 API**：`POST /api/developer/games`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO games (title, description, price, developer_id, status) VALUES ('新遊戲', '...', 500, 5, 'ACTIVE');
  ```

### 21. 編輯遊戲資訊 (UPDATE)
- **對應 API**：`PUT /api/developer/games/:id`
- **原生 SQL 語法**：
  ```sql
  UPDATE games SET title = '新標題', price = 400 WHERE game_id = 42 AND developer_id = 5;
  ```

### 22. 下架遊戲 (UPDATE)
- **對應 API**：`PUT /api/developer/games/:id/takedown`
- **原生 SQL 語法**：
  ```sql
  UPDATE games SET status = 'TAKEN_DOWN' WHERE game_id = 42 AND developer_id = 5;
  ```

### 23. 上傳遊戲圖片或檔案 (INSERT)
- **對應 API**：`POST /api/developer/games/:id/media`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO game_media (game_id, media_type, file_url) VALUES (42, 'SCREENSHOT', '/media/123.jpg');
  ```

### 24. 刪除媒體檔案 (DELETE)
- **對應 API**：`DELETE /api/developer/media/:id`
- **原生 SQL 語法**：
  ```sql
  DELETE FROM game_media WHERE media_id = 77;
  ```

### 25. 為遊戲新增/移除標籤 (INSERT / DELETE)
- **對應 API**：`POST /api/developer/games/:id/tags` 與 `DELETE /api/developer/games/:id/tags/:tag_id`
- **原生 SQL 語法**：
  ```sql
  INSERT INTO game_tags (game_id, tag_id) VALUES (42, 3);
  DELETE FROM game_tags WHERE game_id = 42 AND tag_id = 3;
  ```

---

## 模組五：管理員與客服 (Admin & CSR)

### 26. 管理員封鎖/解封玩家 (UPDATE)
- **對應 API**：`PUT /api/admin/users/:id/ban` 與 `unban`
- **原生 SQL 語法**：
  ```sql
  UPDATE users SET role = 'BANNED' WHERE user_id = 10;
  UPDATE users SET role = 'PLAYER' WHERE user_id = 10;
  ```

### 27. 客服核准/拒絕退款單 (UPDATE)
- **對應 API**：`PUT /api/csr/refunds/:id`
- **原生 SQL 語法**：
  ```sql
  UPDATE refund_requests SET status = 'APPROVED', resolved_at = NOW() WHERE request_id = 88;
  UPDATE game_licenses SET status = 'REVOKED' WHERE transaction_item_id = 105; 
  ```
