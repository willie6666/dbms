# 資料庫展示 SQL 語法大全 (Adminer 實機展示用)

這份文件為您整理了 `sql_a` 到 `sql_g` 裡面提到的**每一項**資料庫操作。裡面的 `?` 和變數都已經替換為**實際的測試數值** (例如 `user_id = 1`, `game_id = 1`)，您可以直接把程式碼區塊中的 SQL 複製，並貼到 Adminer 的 SQL 命令列中執行，向老師完美展示您的資料庫架構！

> [!TIP]
> 這些數值都是基於 `02_init_data.sql` 與 `03_init_super_data.sql` 產生的資料，絕大多數都能直接跑出結果。如果查無資料，您可以自行把 ID 換成其他數字 (例如 101, 102 等 Super Data)。

---

## A. 單一資料表查詢 (sql_a.md)

### 1. 取得全域標籤列表
```sql
SELECT tag_id AS "標籤編號", tag_name AS "標籤名稱" FROM tags ORDER BY tag_id ASC;
```

### 2. 使用者登入檢查 (依 Email 查詢)
```sql
SELECT * FROM users WHERE email = 'admin@vaporauror.com' LIMIT 1;
```

### 3. 使用者登出 (狀態確認)
*(此功能純前端與 Token 處理，無實際資料庫 SQL 操作)*

### 4. 取得單一遊戲基本資訊
```sql
SELECT * FROM games WHERE game_id = 1 AND status != 'TAKEN_DOWN' LIMIT 1;
```

### 5. 查看自己的好友名單
```sql
SELECT * FROM friendships 
WHERE (sender_id = 1 OR receiver_id = 1) AND status = 'ACCEPTED';
```

### 6. 查看未處理的好友邀請
```sql
SELECT * FROM friendships WHERE receiver_id = 1 AND status = 'PENDING';
```

### 7. 查看自己的黑名單
```sql
SELECT * FROM blacklists WHERE blocker_id = 1;
```

### 8. 管理員獲取全站使用者名單
```sql
SELECT user_id, username, email, role, status FROM users ORDER BY registration_date DESC;
```

### 9. 客服獲取所有待處理退款單
```sql
SELECT * FROM refund_requests WHERE status = 'PENDING' ORDER BY created_at DESC;
```

### 10. 開發者獲取自己發布的所有遊戲
```sql
SELECT * FROM games WHERE developer_id = 1 ORDER BY game_id DESC;
```

### 11. 檢查交易明細是否存在 (申請退款前)
```sql
SELECT * FROM transaction_items WHERE item_id = 1 LIMIT 1;
```

### 12. 取得個人退款歷史紀錄
```sql
SELECT * FROM refund_requests WHERE buyer_id = 1 ORDER BY created_at DESC;
```

### 13. 驗證玩家是否擁有該遊戲遊玩權限
```sql
SELECT * FROM game_licenses 
WHERE user_id = 1 AND game_id = 1 AND status = 'ACTIVE' LIMIT 1;
```

### 14. 顯示與某使用者的對話紀錄
```sql
SELECT * FROM messages 
WHERE (sender_id = 1 AND receiver_id = 2) OR (sender_id = 2 AND receiver_id = 1) 
ORDER BY sent_at ASC;
```

---

## B. 雙資料表 JOIN (sql_b.md)

### 1. 查詢購物車內容與遊戲資訊
```sql
SELECT c.cart_id, g.title, g.price, c.added_at
FROM shopping_carts c
JOIN games g ON c.game_id = g.game_id
WHERE c.user_id = 1;
```

### 2. 查詢願望清單與遊戲資訊
```sql
SELECT w.wishlist_id, g.title, g.overall_rating
FROM wish_lists w
JOIN games g ON w.game_id = g.game_id
WHERE w.user_id = 1;
```

### 3. 獲取遊戲的評論列表並包含評論者名稱
```sql
SELECT r.content, r.attitude, u.username, r.created_at
FROM reviews r
JOIN users u ON r.user_id = u.user_id
WHERE r.game_id = 1 AND r.status = 'VISIBLE'
ORDER BY r.created_at DESC;
```

### 4. 獲取評論的獨立回覆列表
```sql
SELECT rr.content, u.username, rr.created_at
FROM review_replies rr
JOIN users u ON rr.user_id = u.user_id
WHERE rr.review_id = 1 AND rr.status = 'VISIBLE'
ORDER BY rr.created_at ASC;
```

### 5. 商店搜尋由特定開發者發布的遊戲
```sql
SELECT g.game_id, g.title, u.username AS developer_name
FROM games g
JOIN users u ON g.developer_id = u.user_id
WHERE u.username ILIKE '%Studio%' AND g.status = 'ACTIVE';
```

---

## C. 多資料表 JOIN (sql_c.md)

### 1. 透過特定標籤搜尋遊戲 (多對多關聯)
```sql
SELECT g.title, g.price, t.tag_name
FROM games g
JOIN game_tags gt ON g.game_id = gt.game_id
JOIN tags t ON gt.tag_id = t.tag_id
WHERE t.tag_name ILIKE '%RPG%' AND g.status = 'ACTIVE';
```

### 2. 查詢玩家的遊戲庫並包含遊戲封面圖
```sql
SELECT g.title, l.acquired_date, m.file_url AS cover_image
FROM game_licenses l
JOIN games g ON l.game_id = g.game_id
LEFT JOIN game_media m ON g.game_id = m.game_id AND m.media_type = 'media'
WHERE l.user_id = 1 AND l.status = 'ACTIVE';
```

### 3. 查詢歷史訂單明細與對應的遊戲名稱
```sql
SELECT t.receipt_number, t.transaction_date, ti.purchase_price, g.title
FROM transactions t
JOIN transaction_items ti ON t.transaction_id = ti.transaction_id
JOIN games g ON ti.game_id = g.game_id
WHERE t.user_id = 1
ORDER BY t.transaction_date DESC;
```

---

## D. 進階過濾與搜尋 (sql_d.md)

### 1. 使用 NOT EXISTS 過濾已購買的遊戲
```sql
SELECT * FROM games g 
WHERE status = 'ACTIVE' 
AND NOT EXISTS (
    SELECT 1 FROM game_licenses l 
    WHERE l.game_id = g.game_id AND l.user_id = 1
);
```

### 2. 使用 ILIKE 進行多欄位模糊搜尋 (Keyword Search)
```sql
SELECT title, description FROM games 
WHERE (title ILIKE '%city%' OR description ILIKE '%city%') 
AND status = 'ACTIVE';
```

### 3. 使用 >= 與 <= 進行價格區間過濾
```sql
SELECT title, price FROM games 
WHERE price >= 100 AND price <= 1000 
ORDER BY price ASC;
```

### 4. 後台管理員使用 ILIKE 搜尋使用者
```sql
SELECT username, email FROM users 
WHERE username ILIKE '%player%' OR email ILIKE '%player%';
```

### 5. 確保審核通過的遊戲才顯示 (使用 !=)
```sql
SELECT title, status FROM games WHERE status != 'TAKEN_DOWN';
```

---

## E. 聚合與分組排序 (sql_e.md)

### 1. 遊戲評分統計 (COUNT, AVG 與 CASE WHEN)
```sql
SELECT 
    game_id, 
    COUNT(*) AS total_reviews,
    SUM(CASE WHEN attitude = 'POSITIVE' THEN 1 ELSE 0 END) AS positive_count
FROM reviews 
WHERE status = 'VISIBLE'
GROUP BY game_id;
```

### 2. 商店首頁的分群與自訂排序 (GROUP BY & ORDER BY)
```sql
SELECT developer_id, COUNT(game_id) AS game_count 
FROM games 
GROUP BY developer_id 
ORDER BY game_count DESC;
```

### 3. 結帳前的購物車商品清單校驗 (使用 IN)
```sql
SELECT * FROM games 
WHERE game_id IN (
    SELECT game_id FROM shopping_carts WHERE user_id = 1
) AND status = 'ACTIVE';
```

### 4. 開發者總銷售額與銷量統計 (SUM, COUNT)
```sql
SELECT 
    g.developer_id, 
    COUNT(ti.item_id) AS total_sold_copies, 
    SUM(ti.purchase_price) AS total_revenue
FROM transaction_items ti
JOIN games g ON ti.game_id = g.game_id
WHERE g.developer_id = 1
GROUP BY g.developer_id;
```

### 5. 後台使用者列表排序 (多欄位 ORDER BY)
```sql
SELECT username, role, registration_date 
FROM users 
ORDER BY role ASC, registration_date DESC;
```

### 6. 退款紀錄列表排序
```sql
SELECT refund_id, status, created_at 
FROM refund_requests 
ORDER BY status ASC, created_at DESC;
```

---

## F. 資料變更操作 - 新增/修改/刪除 (sql_f.md)

> [!WARNING]
> 下列操作會改變資料庫狀態，展示時請注意執行順序。

### 1. 玩家註冊 (INSERT)
```sql
INSERT INTO users (username, email, password_hash, role) 
VALUES ('DemoUser123', 'demo123@example.com', 'hashed_pwd', 'USERS');
```

### 2. 更新個人檔案 (UPDATE)
```sql
UPDATE users SET username = 'SuperDemoUser' WHERE user_id = 4;
```

### 3. 將遊戲加入購物車 (INSERT)
```sql
INSERT INTO shopping_carts (user_id, game_id) VALUES (4, 1);
```

### 4. 移除單一購物車商品 (DELETE)
```sql
DELETE FROM shopping_carts WHERE user_id = 4 AND game_id = 1;
```

### 5. 清空購物車 (DELETE)
```sql
DELETE FROM shopping_carts WHERE user_id = 4;
```

### 6. 購物車結帳 (巨型 Transaction 模擬)
*(在 Adminer 中可逐行執行，模擬結帳流程)*
```sql
BEGIN;
-- 1. 新增訂單
INSERT INTO transactions (user_id, total_amount, receipt_number) 
VALUES (4, 1200.00, 'REC-DEMO-001');

-- 2. 新增訂單明細
INSERT INTO transaction_items (transaction_id, game_id, purchase_price) 
VALUES (currval('transactions_transaction_id_seq'), 1, 1200.00);

-- 3. 發放遊戲授權
INSERT INTO game_licenses (game_id, user_id, transaction_item_id, status) 
VALUES (1, 4, currval('transaction_items_item_id_seq'), 'ACTIVE');

-- 4. 清空購物車
DELETE FROM shopping_carts WHERE user_id = 4;
COMMIT;
```

### 7. 將遊戲加入願望清單 (INSERT)
```sql
INSERT INTO wish_lists (user_id, game_id) VALUES (4, 2) ON CONFLICT DO NOTHING;
```

### 8. 移除願望清單 (DELETE)
```sql
DELETE FROM wish_lists WHERE user_id = 4 AND game_id = 2;
```

### 9. 玩家申請退款 (INSERT)
```sql
INSERT INTO refund_requests (buyer_id, transaction_item_id, reason, status) 
VALUES (4, 1, '遊戲畫面卡頓', 'PENDING');
```

### 10. 送出好友邀請 (INSERT)
```sql
INSERT INTO friendships (sender_id, receiver_id, status) VALUES (4, 1, 'PENDING');
```

### 11. 接受好友邀請 (UPDATE)
```sql
UPDATE friendships SET status = 'ACCEPTED' WHERE sender_id = 4 AND receiver_id = 1;
```

### 12. 拒絕好友邀請 (UPDATE)
```sql
UPDATE friendships SET status = 'DECLINED' WHERE sender_id = 4 AND receiver_id = 1;
```

### 13. 收回邀請或刪除好友 (DELETE)
```sql
DELETE FROM friendships WHERE (sender_id = 4 AND receiver_id = 1) OR (sender_id = 1 AND receiver_id = 4);
```

### 14. 加入黑名單 (INSERT)
```sql
INSERT INTO blacklists (blocker_id, blocked_id) VALUES (4, 2);
```

### 15. 移除黑名單 (DELETE)
```sql
DELETE FROM blacklists WHERE blocker_id = 4 AND blocked_id = 2;
```

### 16. 發送對話訊息 (INSERT)
```sql
INSERT INTO messages (sender_id, receiver_id, content) VALUES (4, 1, 'Hello Demo!');
```

### 17. 標記訊息為已讀 (UPDATE)
```sql
UPDATE messages SET is_read = TRUE WHERE receiver_id = 4 AND sender_id = 1;
```

### 18. 發布遊戲評論 (INSERT)
```sql
INSERT INTO reviews (game_id, user_id, content, attitude) 
VALUES (1, 4, '非常好玩的展示遊戲！', 'POSITIVE');
```

### 19. 修改遊戲評論 (UPDATE)
```sql
UPDATE reviews SET content = '玩久了覺得還好' WHERE user_id = 4 AND game_id = 1;
```

### 20. 刪除遊戲評論 (UPDATE 軟刪除)
```sql
UPDATE reviews SET status = 'DELETED' WHERE user_id = 4 AND game_id = 1;
```

### 21. 回覆別人的評論 (INSERT)
```sql
INSERT INTO review_replies (review_id, user_id, content) 
VALUES (1, 4, '我同意你的看法！');
```

### 22. 刪除回覆 (DELETE / UPDATE 軟刪除)
```sql
UPDATE review_replies SET status = 'DELETED' WHERE review_reply_id = 1;
```

### 23. 發行新遊戲草稿 (INSERT)
```sql
INSERT INTO games (developer_id, title, price, status) 
VALUES (2, 'Demo Game Project', 500.00, 'DRAFT');
```

### 24. 正式上架或編輯遊戲資訊 (UPDATE)
```sql
UPDATE games SET status = 'ACTIVE', description = 'It is ready!' WHERE game_id = 5;
```

### 25. 下架自己的遊戲 (UPDATE)
```sql
UPDATE games SET status = 'TAKEN_DOWN' WHERE game_id = 5;
```

### 26. 上傳遊戲圖片或檔案 (INSERT)
```sql
INSERT INTO game_media (game_id, file_url, media_type) 
VALUES (5, '/media/images/5/demo.png', 'media');
```

### 27. 刪除媒體檔案 (DELETE)
```sql
DELETE FROM game_media WHERE media_id = 1;
```

### 28. 為遊戲新增標籤 (INSERT)
```sql
INSERT INTO game_tags (game_id, tag_id) VALUES (1, 100);
```

### 29. 建立全域新標籤與移除遊戲標籤 (INSERT / DELETE)
```sql
INSERT INTO tags (tag_name) VALUES ('Demo_Tag');
DELETE FROM game_tags WHERE game_id = 1 AND tag_id = 100;
```

### 30. 管理員停權/切換玩家身分 (UPDATE)
```sql
UPDATE users SET permission = 'DEACTIVE' WHERE user_id = 4;
```

### 31. 管理員強制刪除帳號 (UPDATE 軟刪除)
```sql
UPDATE users SET permission = 'DELETED' WHERE user_id = 4;
```

### 32. 管理員強制下架遊戲與終極撤銷 (UPDATE)
```sql
UPDATE games SET status = 'TAKEN_DOWN' WHERE game_id = 1;
```

### 33. 客服同意或拒絕退款單 (UPDATE)
```sql
UPDATE refund_requests SET status = 'APPROVED', handled_by = 3, resolved_at = CURRENT_TIMESTAMP WHERE refund_id = 1;
```

---

## G. 資料結構操作 (sql_g.md)

> [!CAUTION]
> DDL 語法會變更資料表結構，在 Adminer 展示完最好不要輕易執行，或者執行後重新 `docker compose down -v` 重置。

### 1. 建立新的資料表 (CREATE TABLE)
```sql
CREATE TABLE demo_table (
    demo_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);
```

### 2. 新增外部鍵約束 (ALTER TABLE)
```sql
ALTER TABLE demo_table ADD COLUMN user_id INT REFERENCES users(user_id);
```

### 3. 刪除整張資料表 (DROP TABLE)
```sql
DROP TABLE demo_table;
```
