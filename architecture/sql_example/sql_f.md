# f. 新增或刪除資料的操作 (INSERT, UPDATE, DELETE)

本文件列出系統進行資料異動 (Data Manipulation) 時使用的 SQL 操作。由於這涵蓋了系統中所有的狀態變更 API，以下依業務模組進行分類。

---

## 模組一：使用者與身分驗證 (User & Auth)

### 1. 玩家註冊 (INSERT)
- **對應 API**：`POST /api/auth/register`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&user)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO users (username, email, password_hash, role) 
  VALUES ('JohnDoe', 'john@test.com', '$2a$10$...', 'PLAYER') RETURNING user_id;
  ```

### 2. 更新個人檔案 (UPDATE)
- **對應 API**：`PUT /api/users/profile`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Save(&user)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE users SET bio = 'Hello World!', avatar_url = 'img.png', updated_at = NOW() WHERE user_id = 5;
  ```

---

## 模組二：商城與交易 (Store, Cart, Wishlist, Transactions)

### 3. 將遊戲加入購物車 (INSERT)
- **對應 API**：`POST /api/protected/cart`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&cartItem)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO shopping_carts (user_id, game_id) VALUES (5, 42);
  ```

### 4. 移除單一購物車商品 (DELETE)
- **對應 API**：`DELETE /api/protected/cart/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Where("user_id = ? AND game_id = ?", userID, gameID).Delete(&models.ShoppingCart{})
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM shopping_carts WHERE user_id = 5 AND game_id = 42;
  ```

### 5. 清空購物車 (DELETE)
- **對應 API**：`DELETE /api/protected/cart`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Where("user_id = ?", userID).Delete(&models.ShoppingCart{})
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM shopping_carts WHERE user_id = 5;
  ```

### 6. 購物車結帳 (巨型 Transaction: INSERT + DELETE)
- **對應 API**：`POST /api/shopping/checkout`
- **Go 實作 (GORM)**：
  ```go
  tx := database.DB.Begin()
  tx.Create(&transaction)
  tx.Create(&transactionItem)
  tx.Create(&license)
  tx.Where("user_id = ?", userID).Delete(&models.ShoppingCart{})
  tx.Commit()
  ```
- **原生 SQL 語法**：
  ```sql
  BEGIN;
  INSERT INTO transactions (user_id, total_amount) VALUES (5, 1200.00) RETURNING transaction_id;
  INSERT INTO transaction_items (transaction_id, game_id, purchase_price) VALUES (99, 42, 1200.00);
  INSERT INTO game_licenses (user_id, game_id, transaction_item_id, status) VALUES (5, 42, 105, 'ACTIVE');
  DELETE FROM shopping_carts WHERE user_id = 5;
  COMMIT;
  ```

### 7. 將遊戲加入願望清單 (INSERT)
- **對應 API**：`POST /api/protected/wishlist`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&wishlistItem)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO wish_lists (user_id, game_id) VALUES (5, 42);
  ```

### 8. 移除願望清單 (DELETE)
- **對應 API**：`DELETE /api/protected/wishlist/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Where("user_id = ? AND game_id = ?", userID, gameID).Delete(&models.WishList{})
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM wish_lists WHERE user_id = 5 AND game_id = 42;
  ```

### 9. 玩家申請退款 (INSERT)
- **對應 API**：`POST /api/social/refunds`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&refundRequest)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO refund_requests (transaction_item_id, user_id, reason, status) 
  VALUES (105, 5, '不好玩', 'PENDING');
  ```

---

## 模組三：社群互動 (Social, Friends, Reviews)

### 10. 送出好友邀請 (INSERT)
- **對應 API**：`POST /api/social/friends/requests`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&friendReq)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO friend_requests (sender_id, receiver_id, status) VALUES (5, 10, 'PENDING');
  ```

### 11. 接受好友邀請 (UPDATE + INSERT)
- **對應 API**：`PUT /api/social/friends/requests/:id` (附帶 status=ACCEPTED)
- **Go 實作 (GORM)**：
  ```go
  database.DB.Model(&request).Update("status", "ACCEPTED")
  database.DB.Create(&friendship)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE friend_requests SET status = 'ACCEPTED' WHERE request_id = 33;
  INSERT INTO friendships (user_id1, user_id2) VALUES (5, 10);
  ```

### 12. 拒絕好友邀請或收回邀請 (DELETE)
- **對應 API**：`PUT /api/social/friends/requests/:id` (附帶 status=REJECTED) 與 `DELETE /api/social/friends/request/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Delete(&request)
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM friend_requests WHERE request_id = 33;
  ```

### 13. 刪除好友 (DELETE)
- **對應 API**：`DELETE /api/social/friends/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", myID, peerID, peerID, myID).Delete(&models.Friendship{})
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM friendships WHERE (user_id1 = 5 AND user_id2 = 10) OR (user_id1 = 10 AND user_id2 = 5);
  ```

### 14. 加入黑名單並自動解除好友 (INSERT + DELETE)
- **對應 API**：`POST /api/social/blacklist`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&blockRecord)
  database.DB.Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", userID, targetID, targetID, userID).Delete(&models.Friendship{})
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO blocklist (blocker_id, blocked_id) VALUES (5, 10);
  DELETE FROM friendships WHERE (user_id1 = 5 AND user_id2 = 10) OR (user_id1 = 10 AND user_id2 = 5);
  ```

### 15. 移除黑名單 (DELETE)
- **對應 API**：`DELETE /api/social/blacklist/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Where("blocker_id = ? AND blocked_id = ?", userID, targetID).Delete(&models.Blocklist{})
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM blocklist WHERE blocker_id = 5 AND blocked_id = 10;
  ```

### 16. 發送對話訊息 (INSERT)
- **對應 API**：`POST /api/social/messages`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&msg)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO messages (sender_id, receiver_id, content) VALUES (5, 10, 'Hello!');
  ```

### 17. 標記訊息為已讀 (UPDATE)
- **對應 API**：`GET /api/social/messages/{user_id}` (撈取同時會順便把未讀標記為已讀)
- **Go 實作 (GORM)**：
  ```go
  database.DB.Model(&models.Message{}).Where("sender_id = ? AND receiver_id = ? AND is_read = false", peerID, myID).Update("is_read", true)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE messages SET is_read = true WHERE sender_id = 10 AND receiver_id = 5 AND is_read = false;
  ```

### 18. 發布遊戲評論 (INSERT)
- **對應 API**：`POST /api/games/:id/reviews`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&review)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO reviews (game_id, user_id, rating, attitude, content, status) 
  VALUES (42, 5, 0, 'POSITIVE', '超好玩', 'VISIBLE');
  ```

### 19. 修改遊戲評論 (UPDATE)
- **對應 API**：`PUT /api/social/reviews/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Save(&review)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE reviews SET content = '玩久了有點膩', attitude = 'NEGATIVE' WHERE review_id = 128 AND user_id = 5;
  ```

### 20. 刪除遊戲評論 (UPDATE 軟刪除)
- **對應 API**：`DELETE /api/social/reviews/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Model(&review).Update("status", "HIDDEN")
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE reviews SET status = 'HIDDEN' WHERE review_id = 128 AND user_id = 5;
  ```

### 21. 回覆別人的評論 (INSERT)
- **對應 API**：`POST /api/social/reviews/:id/replies`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&reply)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO review_replies (review_id, user_id, content) VALUES (128, 10, '+1 認同');
  ```

### 22. 編輯與刪除回覆 (UPDATE / DELETE)
- **對應 API**：`PUT /api/social/replies/:id` 與 `DELETE /api/social/replies/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Save(&reply)
  database.DB.Delete(&reply)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE review_replies SET content = '修改後的留言' WHERE reply_id = 56;
  DELETE FROM review_replies WHERE reply_id = 56;
  ```

---

## 模組四：開發者功能 (Developer)

### 23. 發行新遊戲 (INSERT)
- **對應 API**：`POST /api/developer/games`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&game)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO games (title, description, price, developer_id, status) VALUES ('新遊戲', '...', 500, 5, 'DRAFT');
  ```

### 24. 上架或編輯遊戲資訊 (UPDATE)
- **對應 API**：`PUT /api/developer/games/:id/publish` 與 `PUT /api/developer/games/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Save(&game)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE games SET status = 'ACTIVE' WHERE game_id = 42 AND developer_id = 5;
  UPDATE games SET title = '新標題', price = 400 WHERE game_id = 42 AND developer_id = 5;
  ```

### 25. 下架自己的遊戲 (UPDATE)
- **對應 API**：`DELETE /api/developer/games/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Model(&game).Update("status", "TAKEN_DOWN")
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE games SET status = 'TAKEN_DOWN' WHERE game_id = 42 AND developer_id = 5;
  ```

### 26. 上傳遊戲圖片或檔案 (INSERT)
- **對應 API**：`POST /api/developer/games/:id/media`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&media)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO game_media (game_id, media_type, file_url) VALUES (42, 'SCREENSHOT', '/media/123.jpg');
  ```

### 27. 刪除媒體檔案 (DELETE)
- **對應 API**：`DELETE /api/developer/games/:id/media/:media_id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Delete(&media)
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM game_media WHERE media_id = 77;
  ```

### 28. 為遊戲新增標籤 (INSERT)
- **對應 API**：`POST /api/developer/games/:id/tags`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&gameTag)
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO game_tags (game_id, tag_id) VALUES (42, 3);
  ```

### 29. 建立全域新標籤與移除遊戲標籤 (INSERT / DELETE)
- **對應 API**：`POST /api/developer/tags` 與 `DELETE /api/developer/games/:id/tags/:tag_id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Create(&newTag)
  database.DB.Where("game_id = ? AND tag_id = ?", gameID, tagID).Delete(&models.GameTag{})
  ```
- **原生 SQL 語法**：
  ```sql
  INSERT INTO tags (tag_name) VALUES ('MOBA');
  DELETE FROM game_tags WHERE game_id = 42 AND tag_id = 3;
  ```

---

## 模組五：管理員與客服 (Admin & CSR)

### 30. 管理員停權/解封/切換玩家身分 (UPDATE)
- **對應 API**：`PUT /api/admin/users/:id/suspend` 與 `PUT /api/admin/users/:id/role`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Save(&user)
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE users SET status = 'DEACTIVE' WHERE user_id = 10;
  UPDATE users SET role = 'BANNED' WHERE user_id = 10;
  ```

### 31. 管理員強制刪除帳號 (DELETE)
- **對應 API**：`DELETE /api/admin/users/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Delete(&user)
  ```
- **原生 SQL 語法**：
  ```sql
  DELETE FROM users WHERE user_id = 10;
  ```

### 32. 管理員強制下架遊戲 (UPDATE)
- **對應 API**：`DELETE /api/admin/games/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Model(&game).Update("status", "TAKEN_DOWN")
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE games SET status = 'TAKEN_DOWN' WHERE game_id = 42;
  ```

### 33. 客服同意或拒絕退款單 (UPDATE)
- **對應 API**：`PUT /api/csr/refunds/:id`
- **Go 實作 (GORM)**：
  ```go
  database.DB.Save(&request)
  database.DB.Model(&license).Update("status", "REVOKED")
  ```
- **原生 SQL 語法**：
  ```sql
  UPDATE refund_requests SET status = 'APPROVED', resolved_at = NOW() WHERE request_id = 88;
  UPDATE game_licenses SET status = 'REVOKED' WHERE transaction_item_id = 105; 
  ```
