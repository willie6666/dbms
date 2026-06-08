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
  5. 接收方登入後，呼叫 `GET /api/social/friends/requests/pending` 會看到該邀請。
  6. 接收方呼叫 `PUT /api/social/friends/requests/:id`，可以選擇傳入 `ACCEPTED` 或 `DECLINED`。
- **終點**：若 `ACCEPTED`，兩人結為好友 (`status` 更新為 `ACCEPTED`)，此後可以互傳私訊。

---

## 2. 私密訊息 (Direct Messages)

- **起點**：雙方為好友狀態，其中一方在聊天室發送訊息。
- **流程**：
  1. 呼叫 `POST /api/social/messages`。
  2. **關係檢查**：後端嚴格檢查這兩人目前的 `friendships.status` 是否為 `ACCEPTED`。如果不是好友，不可發送。
  3. **黑名單二次防護**：由於建立好友後仍可能被單方面封鎖，所以發送訊息前會再次檢查 `blocklist`。若被封鎖，回傳「無法發送訊息」。
  4. 訊息寫入 `messages` 資料表。
- **終點**：接收方呼叫 `GET /api/social/messages/:friend_id` 即可看見對話。

---

## 3. 黑名單與關係解除 (Blocklist & Unfriending)

- **起點**：玩家 A 不想再收到玩家 B 的任何訊息，在社群頁面對 B 點擊「封鎖」。
- **流程**：
  1. 呼叫 `POST /api/social/blocklist`，傳入 B 的 ID。
  2. 系統在 `blocklist` 寫入一筆 `blocker_id = A`, `blocked_id = B` 的紀錄。
  3. **關係清理 (Clean up)**：
     - 如果 A 和 B 目前是好友，或是有一方發出了邀請，系統會自動到 `friendships` 表中尋找並 **刪除 (DELETE)** 該筆好友紀錄，強制斬斷好友關係。
     - (注意：歷史私訊保留，但未來無法再發送新訊息)。
- **終點**：B 將無法搜尋到 A、無法加 A 好友、無法傳訊息給 A，達到完全防護的閉環效果。
