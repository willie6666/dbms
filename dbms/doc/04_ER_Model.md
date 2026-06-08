# 資料實體、屬性與關係設計 (ER Model)

## 六、實體類型與屬性

### 1. 使用者與權限
| 實體 | 屬性 | 說明 |
| --- | --- | --- |
| **USER** | `user_id` (PK)<br>`username`<br>`email`<br>`password_hash`<br>`registration_date`<br>`last_visit_ip`<br>`role` (NULL, USER, CSR, ADMIN, DEVELOPER)<br>`status` (ONLINE, OFFLINE)<br>`permission` (ACTIVE, DEACTIVE, DELETED) | 這是使用者 |

### 2. 商店與遊戲
| 實體 | 屬性 | 說明 |
| --- | --- | --- |
| **GAME** | `game_id` (PK)<br>`developer_id` (FK) -> `USER.user_id`<br>`title`<br>`price`<br>`overall_rating` | 這是遊戲 |
| **TAG** | `tag_id` (PK)<br>`tag_name` | 這是標籤 |
| **GAME_TAG** (weak) | `game_id` (PK, FK) -> `GAME.game_id`<br>`tag_id` (PK, FK) -> `TAG.tag_id` | 這是某遊戲擁有的標籤 |
| **GAME_MEDIA** (weak) | `media_id` (PK)<br>`game_id` (FK) -> `GAME.game_id`<br>`file_url`<br>`media_type` (media, game_file) | 這是遊戲媒體 |

### 3. 訂單與交易
| 實體 | 屬性 | 說明 |
| --- | --- | --- |
| **TRANSACTION** | `transaction_id` (PK)<br>`user_id` (FK) -> `USER.user_id`<br>`total_amount`<br>`transaction_date`<br>`receipt_number` | 這是交易 |
| **TRANSACTION_ITEM** (weak) | `item_id` (PK)<br>`transaction_id` (FK) -> `TRANSACTION.transaction_id`<br>`game_id` (FK) -> `GAME.game_id`<br>`purchase_price` | 這是交易品項 |
| **SHOPPING_CART** | `cart_id` (PK)<br>`user_id` (FK) -> `USER.user_id`<br>`game_id` (FK) -> `GAME.game_id`<br>`added_at` | 這是購物車 |

### 4. 遊戲庫與授權
| 實體 | 屬性 | 說明 |
| --- | --- | --- |
| **GAME_LICENSE** | `license_id` (PK)<br>`game_id` (FK) -> `GAME.game_id`<br>`user_id` (FK) -> `USER.user_id`<br>`transaction_item_id` (FK) -> `TRANSACTION_ITEM.item_id`<br>`acquired_date`<br>`status` (ACTIVE, REVOKED, FROZEN) | 這是遊戲許可證 |
| **WISH_LIST** | `wishlist_id` (PK)<br>`user_id` (FK) -> `USER.user_id`<br>`game_id` (FK) -> `GAME.game_id`<br>`added_at` | 這是願望清單 |

### 5. 客服
| 實體 | 屬性 | 說明 |
| --- | --- | --- |
| **REFUND_REQUEST** | `refund_id` (PK)<br>`buyer_id` (FK) -> `USER.user_id`<br>`transaction_item_id` (FK) -> `TRANSACTION_ITEM.item_id`<br>`handled_by` (FK) -> `USER.user_id`<br>`reason`<br>`reject_reason`<br>`created_at`<br>`resolved_at`<br>`status` (PENDING, APPROVED, REJECTED) | 這是退款請求 |

### 6. 社交與社群
| 實體 | 屬性 | 說明 |
| --- | --- | --- |
| **REVIEW** | `review_id` (PK)<br>`game_id` (FK) -> `GAME.game_id`<br>`user_id` (FK) -> `USER.user_id`<br>`content`<br>`created_at`<br>`attitude` (POSITIVE, NEGATIVE)<br>`status` (VISIBLE, HIDDEN, DELETED) | 這是評論 |
| **REVIEW_REPLY** (遞迴) | `review_reply_id` (PK)<br>`review_id` (FK) -> `REVIEW.review_id`<br>`user_id` (FK) -> `USER.user_id`<br>`parent_reply_id` (FK) -> `REVIEW_REPLY.review_reply_id`<br>`content`<br>`status` (VISIBLE, HIDDEN, DELETED) | 這是評論的回覆<br>*(parent_reply_id 指向上一層)* |
| **FRIENDSHIP** | `friendship_id` (PK)<br>`sender_id` (FK) -> `USER.user_id`<br>`receiver_id` (FK) -> `USER.user_id`<br>`created_at`<br>`status` (PENDING, ACCEPTED, DECLINED) | 這是好友關係 |
| **MESSAGE** | `message_id` (PK)<br>`sender_id` (FK) -> `USER.user_id`<br>`receiver_id` (FK) -> `USER.user_id`<br>`content`<br>`sent_at`<br>`is_read` | 這是訊息 |
| **BLACKLIST** | `blacklist_id` (PK)<br>`blocker_id` (FK) -> `USER.user_id`<br>`blocked_id` (FK) -> `USER.user_id`<br>`created_at` | 這是黑名單 |

---

## 七、實體的關係

### 1. 使用者與遊戲 / 權限相關
- **USER to GAME**: 1:N（USER 部分；GAME 完全）
  - *說明：一個 Developer 可上架多個遊戲；每個 GAME 必須屬於一個 developer_id。*
- **USER to USER 管理關係**: 1:N（ADMIN 部分；USER 部分）
  - *說明：ADMIN 可停權 / 改權限多個使用者。*

### 2. 遊戲與標籤 / 媒體
- **GAME to GAME_MEDIA**: 1:N（GAME 完全；GAME_MEDIA 完全）
  - *說明：一個遊戲可有多個媒體或遊戲檔案；每筆 GAME_MEDIA 一定屬於一個 GAME。*
- **GAME to GAME_TAG**: 1:N（GAME 完全；GAME_TAG 完全）
- **TAG to GAME_TAG**: 1:N（TAG 部分；GAME_TAG 完全）
  - *說明：(M:N 拆分) 一個遊戲至少應有一個標籤；一個標籤可以暫時沒有任何遊戲使用。*

### 3. 使用者、交易與購買品項
- **USER to TRANSACTION**: 1:N（USER 部分；TRANSACTION 完全）
  - *說明：一個使用者可有多筆交易；每筆交易必須屬於一個使用者。*
- **TRANSACTION to TRANSACTION_ITEM**: 1:N（TRANSACTION 完全；TRANSACTION_ITEM 完全）
  - *說明：一筆交易至少有一個交易品項；每個交易品項必須屬於一筆交易。*
- **GAME to TRANSACTION_ITEM**: 1:N（GAME 部分；TRANSACTION_ITEM 完全）
  - *說明：一款遊戲可出現在多筆交易品項中；每個交易品項必須對應一款遊戲。*

### 4. 購物車
- **USER to SHOPPING_CART**: 1:N（USER 部分；SHOPPING_CART 完全）
  - *說明：一個使用者可有多個購物車項目；每個購物車項目必須屬於一個使用者。*
- **GAME to SHOPPING_CART**: 1:N（GAME 部分；SHOPPING_CART 完全）
  - *說明：一個遊戲可被多個使用者加入購物車；每個購物車項目必須對應一款遊戲。*
- **USER to GAME**: M:N 
  - *說明：透過 SHOPPING_CART 實作，表示「使用者將遊戲加入購物車」。*

### 5. 遊戲庫與授權
- **USER to GAME_LICENSE**: 1:N（USER 部分；GAME_LICENSE 完全）
  - *說明：一個使用者可擁有多個遊戲授權；每個授權必須屬於一個使用者。*
- **GAME to GAME_LICENSE**: 1:N（GAME 部分；GAME_LICENSE 完全）
  - *說明：一款遊戲可被多個使用者擁有；每個授權必須對應一款遊戲。*
- **TRANSACTION_ITEM to GAME_LICENSE**: 1:1（TRANSACTION_ITEM 完全；GAME_LICENSE 完全）
  - *說明：一筆交易品項產生一個遊戲授權、一個授權應來自一筆交易品項。*

### 6. 願望清單
- **USER to WISH_LIST**: 1:N（USER 部分；WISH_LIST 完全）
- **GAME to WISH_LIST**: 1:N（GAME 部分；WISH_LIST 完全）
- **USER to GAME**: M:N
  - *說明：透過 WISH_LIST 實作，表示「使用者收藏遊戲」。*

### 7. 退款請求
- **USER to REFUND_REQUEST (Buyer)**: 1:N（USER 部分；REFUND_REQUEST 完全）
  - *說明：buyer_id 表示申請退款的使用者；一個使用者可提出多筆退款請求。*
- **TRANSACTION_ITEM to REFUND_REQUEST**: 1:0..1（TRANSACTION_ITEM 部分；REFUND_REQUEST 完全）
  - *說明：一個交易品項可以沒有退款請求，也可以有一筆退款請求；每筆退款請求必須對應一個交易品項。*
- **USER to REFUND_REQUEST (Handler)**: 1:N（USER 部分；REFUND_REQUEST 部分）
  - *說明：handled_by 表示處理退款的 CSR / ADMIN。*

### 8. 評論與回覆
- **USER to REVIEW**: 1:N（USER 部分；REVIEW 完全）
  - *說明：一個使用者可發表多篇評論；每篇評論必須由一個使用者發表。*
- **GAME to REVIEW**: 1:N（GAME 部分；REVIEW 完全）
  - *說明：一款遊戲可有多篇評論；每篇評論必須屬於一款遊戲。*
- **REVIEW to REVIEW_REPLY**: 1:N（REVIEW 部分；REVIEW_REPLY 完全）
  - *說明：一篇評論可有多個回覆；每個回覆必須屬於一篇一樓評論。*
- **USER to REVIEW_REPLY**: 1:N（USER 部分；REVIEW_REPLY 完全）
  - *說明：一個使用者可發表多個評論回覆；每個回覆必須由一個使用者發表。*
- **REVIEW_REPLY to REVIEW_REPLY**: 1:N 遞迴關係（Parent Reply 部分；Child Reply 部分）
  - *說明：一個回覆可被多個回覆接續。*

### 9. 好友關係
- **USER to FRIENDSHIP (Sender)**: 1:N（sender USER 部分；FRIENDSHIP 完全）
  - *說明：一個使用者可發送多個好友邀請。*
- **USER to FRIENDSHIP (Receiver)**: 1:N（receiver USER 部分；FRIENDSHIP 完全）
  - *說明：一個使用者可接收多個好友邀請。*
- **USER to USER**: M:N
  - *說明：透過 FRIENDSHIP 實作，表示使用者之間的好友邀請 / 好友關係。*

### 10. 訊息
- **USER to MESSAGE (Sender)**: 1:N（sender USER 部分；MESSAGE 完全）
  - *說明：一個使用者可發送多則訊息。*
- **USER to MESSAGE (Receiver)**: 1:N（receiver USER 部分；MESSAGE 完全）
  - *說明：一個使用者可接收多則訊息。*
- **USER to USER**: M:N
  - *說明：透過 MESSAGE 實作，表示使用者之間傳送訊息。*

### 11. 黑名單
- **USER to BLACKLIST (Blocker)**: 1:N（blocker USER 部分；BLACKLIST 完全）
  - *說明：一個使用者可封鎖多個人。*
- **USER to BLACKLIST (Blocked)**: 1:N（blocked USER 部分；BLACKLIST 完全）
  - *說明：一個使用者可被多個人封鎖。*
- **USER to USER**: M:N
  - *說明：透過 BLACKLIST 實作，表示使用者封鎖其他使用者。*
