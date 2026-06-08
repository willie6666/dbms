# 6. 退款申請與客服審核 (Refund & CSR)

本文件說明玩家如何對不滿意的遊戲提出退款，以及客服人員 (`CSR`) 如何審核、退回款項並撤銷遊戲授權的完整閉環機制。

---

## 1. 玩家發起退款申請 (Submit Refund Request)

- **起點**：玩家在「我的遊戲庫 (`library.html`)」中，針對某款遊戲點擊「申請退款」。
- **流程**：
  1. 點擊後，前端會先透過 `GET /api/protected/transactions` 找出這款遊戲當時的 `transaction_item_id` (交易明細 ID)。
  2. 玩家填寫退款理由，呼叫 `POST /api/social/refunds` 提交。
  3. 後端檢查防護：
     - 這筆交易明細是否屬於該玩家。
     - 是否已經有正在處理中 (`PENDING`) 或已完成的退款單？避免重複申請。
  4. 通過後，在 `refund_requests` 資料表建立一筆狀態為 `'PENDING'` 的紀錄。
- **後續追蹤**：
  - 玩家可在歷史紀錄中呼叫 `GET /api/protected/refunds` 查看已申請的退款單與狀態。
- **終點**：畫面上顯示申請已送出，玩家必須等待客服人員處理。此時遊戲仍在玩家的庫中 (授權依舊是 `ACTIVE`)，但他們知道退款已在排程。

---

## 2. 客服人員後台審核 (CSR Moderation)

- **起點**：具有 `CSR` 權限的客服人員登入後，進入客服專用後台 (`csr_dashboard.html`)。
- **流程 (列表檢視)**：
  - 呼叫 `GET /api/csr/refunds` 獲取所有待處理的案件列表，包含購買價格、申請原因、買家名稱等。
- **流程 (拒絕退款)**：
  1. 客服發現遊玩時間過長或理由不合理，點擊「拒絕」。
  2. 呼叫 `PUT /api/csr/refunds/:id`，附上 `{ "status": "REJECTED", "reject_reason": "..." }`。
  3. `refund_requests.status` 轉為 `'REJECTED'`，並記錄 `resolved_at` 與處理人 (`handled_by`)。
  4. 玩家的授權不變。
- **流程 (核准退款 - 最核心的金流與授權邏輯)**：
  1. 客服點擊「核准退款」，呼叫 `PUT /api/csr/refunds/:id`，附上 `{ "status": "APPROVED" }`。
  2. 後端開啟資料庫交易 (`DB.Begin()`) 確保以下步驟完全成功或一起失敗。
  3. **撤銷授權**：透過 `transaction_item_id` 找到關聯的 `game_licenses` 紀錄，將其 `status` 從 `'ACTIVE'` 改為 `'REVOKED'`。
  4. **結案**：`refund_requests.status` 轉為 `'APPROVED'`，寫入處理時間與人員。
  5. 交易提交 (`DB.Commit()`)。
- **終點**：退款單狀態完結。當玩家再次打開「我的遊戲庫」時，因為授權已經變成 `REVOKED`，該遊戲會直接從畫面上消失。退款閉環完成。
