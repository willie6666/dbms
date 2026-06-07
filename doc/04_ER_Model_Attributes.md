# 六. 資料實體、屬性與關係設計(ER model)

## (1) 實體類型與屬性
| 模組分類 | 實體 | 屬性 | 文字說明 |
| --- | --- | --- | --- |
| 使用者與權限 | USERS | user_id (PK), username, email, password_hash, registration_date, last_visit_ip, role(NULL, USERS, CSR, ADMIN, DEVELOPER), status(ONLINE,OFFLINE), permission(ACTIVE,DEACTIVE,DELETED) | 這是使用者 |
| 商店與遊戲 | GAMES | game_id (PK), developer_id (FK), title, price, overall_rating | 這是遊戲 (developer_id refer USERS.user_id) |
| 商店與遊戲 | TAG | tag_id (PK), tag_name | 這是標籤 |
| 商店與遊戲 | GAMES_TAG (weak) | (game_id, tag_id) (PK + FK) | 這是某遊戲擁有的標籤 (game_id refer GAMES.game_id, tag_id refer TAG.tag_id) |
| 商店與遊戲 | GAMES_MEDIA (weak) | media_id (PK), game_id (FK), file_url, media_type (media, game_file) | 這是遊戲媒體 (game_id refer GAMES.game_id) |
| 訂單與交易 | TRANSACTION | transaction_id(PK), user_id (FK), total_amount, transaction_date, receipt_number | 這是交易 (user_id refer USERS.user_id) |
| 訂單與交易 | TRANSACTION_ITEM (weak) | item_id(PK), transaction_id (FK), game_id (FK), purchase_price | 這是交易品項 (transaction_id refer TRANSACTION.transaction_id, game_id refer GAMES.game_id) |
| 訂單與交易 | SHOPPING_CART | cart_id(PK), user_id (FK), game_id (FK), added_at | 這是購物車 (user_id refer USERS.user_id, game_id refer GAMES.game_id) |
| 遊戲庫與授權 | GAMES_LICENSE | license_id(PK), game_id (FK), user_id (FK), transaction_item_id (FK), acquired_date, status(ACTIVE,REVOKED, FROZEN) | 這是遊戲許可證 (user_id refer USERS.user_id, game_id refer GAMES.game_id, transaction_item_id refer TRANSACTION_ITEM.item_id) |
| 遊戲庫與授權 | WISH_LIST | wishlist_id (PK), user_id (FK), game_id (FK), added_at | 這是願望清單 (user_id refer USERS.user_id, game_id refer GAMES.game_id) |
| 客服 | REFUND_REQUEST | refund_id (PK), buyer_id (FK), transaction_item_id (FK), handled_by (FK), reason, reject_reason, created_at, resolved_at, status(PENDING,APPROVED,REJECTED) | 這是退款請求 (buyer_id refer USERS.user_id, handled_by refer USERS.user_id, transaction_item_id refer TRANSACTION_ITEM.item_id) |
| 社交與社群 | REVIEW | review_id(PK), game_id (FK), user_id (FK), content, created_at, attitude (POSITIVE,NEGATIVE), status (VISIBLE,HIDDEN,DELETED) | 這是評論 (game_id refer GAMES.game_id, user_id refer USERS.user_id) |
| 社交與社群 | REVIEW_REPLY (遞迴) | review_reply_id (PK), review_id (FK), user_id (FK), parent_reply_id (FK), content, created_at, status (VISIBLE,HIDDEN,DELETED) | 這是評論的回復 (user_id refer USERS.user_id, review_id refer REVIEW.review_id, parent_reply_id refer REVIEW_REPLY.review_reply_id) |
| 社交與社群 | FRIENDSHIP | friendship_id(PK), sender_id (FK), receiver_id (FK), created_at, status (PENDING, ACCEPTED, DECLINED) | 這是好友關係 (sender_id,receiver_id refer USERS.user_id) |
| 社交與社群 | MESSAGE | message_id(PK), sender_id (FK), receiver_id (FK), content, sent_at, is_read | 這是訊息 (sender_id,receiver_id refer USERS.user_id) |
| 社交與社群 | BLACKLIST | blacklist_id (PK), blocker_id (FK), blocked_id (FK), created_at | 這是黑名單 (blocker_id, blocked_id refer USERS.user_id) |

## (2) 實體的關係
(略：詳情請參考 mermaid ER diagrams)
