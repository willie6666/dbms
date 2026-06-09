# 4. 社交互動與黑名單機制 (Social & Blocklist)

本文件描述平台上的玩家如何互相交流，以及透過黑名單機制建立起防騷擾的防線。

---

## 1. 好友邀請與確認 (Friend Requests)

- **起點**：玩家在社群頁面 (`community.html`) 搜尋其他使用者的 Username，點擊「加入好友」。
- **流程**：
  1. 呼叫 `POST /api/social/friends/requests` 傳送 `receiver_id`。
  2. **黑名單防護檢查**：後端會去查詢 `blocklist` 資料表，檢查這兩人之中「是否有一方封鎖了另一方」。若有，直接回傳 `403 Forbidden`，並提示「無法發送好友邀請」。
  3. **帳號有效性檢查**：確認接收方是否為 `ACTIVE` 狀態 (不發給已刪除或停權的使用者)。
  4. 若無攔截，在 `friendships` 寫入一筆 `status = 'PENDING'` 的邀請紀錄。
  5. 接收方登入後，呼叫 `GET /api/social/friends/requests` 會看到該邀請。
  6. 接收方可呼叫 `PUT /api/social/friends/request/:id/accept` 或 `PUT /api/social/friends/request/:id/decline` 來接受或拒絕。
  7. **管理好友**：雙方皆可呼叫 `GET /api/social/friends` 查看已成立的好友列表；若欲解除好友或收回邀請，可呼叫 `DELETE /api/social/friends/request/:id`。
- **終點**：若接受邀請，兩人結為好友 (`status` 更新為 `ACCEPTED`)，此後可以互傳私訊。

---

## 2. 私密訊息 (Direct Messages)

- **起點**：雙方為好友狀態，其中一方在聊天室發送訊息。
- **流程**：
  1. 呼叫 `POST /api/social/messages`。
  2. **關係檢查**：後端嚴格檢查這兩人目前的 `friendships.status` 是否為 `ACCEPTED`。如果不是好友，不可發送。
  3. **黑名單二次防護**：由於建立好友後仍可能被單方面封鎖，所以發送訊息前會再次檢查 `blocklist`。若被封鎖，回傳「無法發送訊息」。
  4. 訊息寫入 `messages` 資料表。
- **終點**：接收方呼叫 `GET /api/social/messages/:friend_id` 即可看見對話。
- **前端即時互動與未讀提示**：
  - 系統採用**一秒全域輪詢 (Polling)** 與 JSON Diffing 技術，無需重整即可即時刷新對話。
  - 當有尚未讀取的訊息 (`is_read = false`)，好友列表的大頭貼會亮起**綠色底色**。點擊開啟聊天室時，會自動將訊息標為已讀，綠底隨即消失。

---

## 3. 黑名單與關係解除 (Blocklist & Unfriending)

- **起點**：玩家 A 不想再收到玩家 B 的任何訊息，在社群頁面對 B 點擊「封鎖」。
- **流程**：
  1. 呼叫 `POST /api/social/blocklist`，傳入 B 的 ID。
  2. 系統在 `blocklist` 寫入一筆 `blocker_id = A`, `blocked_id = B` 的紀錄。
  3. **關係保留與介面隱藏 (Soft Block & UI Hide)**：
     - 加入黑名單是一種「軟封鎖」機制，它並不會自動到 `friendships` 刪除好友紀錄。
     - 但為了介面整潔，**已封鎖的使用者會從好友列表 (`GetFriends`) 中被完全隱藏**，僅會出現在黑名單頁籤中。同時未來的私訊與邀請都會被強制攔截。
  4. **黑名單管理**：玩家可透過 `GET /api/social/blacklist` 查看黑名單列表，若欲解除封鎖，可呼叫 `DELETE /api/social/blacklist/:user_id`。
- **終點**：B 將無法搜尋到 A、無法加 A 好友、無法傳訊息給 A，達到防止騷擾的防護效果。
