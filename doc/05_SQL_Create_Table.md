# 八. 實作資料庫 Create Table 語句
(詳細 SQL 請參考 `db/01_init_table.sql`，以下為重點綱要)

* `users`: PK `user_id`, UNIQUE `username`, `email`
* `games`: PK `game_id`, FK `developer_id` (ON DELETE CASCADE)
* `tags`: PK `tag_id`, UNIQUE `tag_name`
* `game_tags`: PK `(game_id, tag_id)`, FK `game_id`, `tag_id` (ON DELETE CASCADE)
* `game_media`: PK `media_id`, FK `game_id` (ON DELETE CASCADE), `media_type` in ('media', 'game_file')
* `transactions`: PK `transaction_id`, FK `user_id` (ON DELETE CASCADE), UNIQUE `receipt_number`
* `transaction_items`: PK `item_id`, FK `transaction_id`, `game_id` (ON DELETE CASCADE)
* `shopping_carts`: PK `cart_id`, FK `user_id`, `game_id` (ON DELETE CASCADE)
* `game_licenses`: PK `license_id`, FK `game_id`, `user_id`, `transaction_item_id` (ON DELETE CASCADE)
* `wish_lists`: PK `wishlist_id`, FK `user_id`, `game_id` (ON DELETE CASCADE), UNIQUE `(user_id, game_id)`
* `refund_requests`: PK `refund_id`, FK `buyer_id`, `transaction_item_id` (ON DELETE CASCADE), FK `handled_by` (ON DELETE SET NULL)
* `reviews`: PK `review_id`, FK `game_id`, `user_id` (ON DELETE CASCADE)
* `review_replies`: PK `review_reply_id`, FK `review_id`, `user_id`, `parent_reply_id` (ON DELETE CASCADE)
* `friendships`: PK `friendship_id`, FK `sender_id`, `receiver_id` (ON DELETE CASCADE), UNIQUE `(sender_id, receiver_id)`
* `messages`: PK `message_id`, FK `sender_id`, `receiver_id` (ON DELETE CASCADE)
* `blacklists`: PK `blacklist_id`, FK `blocker_id`, `blocked_id` (ON DELETE CASCADE), UNIQUE `(blocker_id, blocked_id)`

# 九. CRUD 對應表(後端 to 資料庫 SQL)
a. 僅使用單一資料表的查詢。
b. 使用兩個資料表的查詢。
c. 使用三個資料表的查詢。
d. 進階查詢 I（使用 EXISTS、NOT EXISTS、NULL、UNION、>=、LIKE 等）。
e. 進階查詢 II（使用 ORDER BY、IN、MAX/MIN/AVG/SUM/COUNT、GROUP BY、HAVING等）。
f. 新增或刪除資料的操作。
g. 修改資料表結構的操作。
