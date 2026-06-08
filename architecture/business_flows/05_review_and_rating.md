# 5. 遊戲評論與社群互動 (Review & Rating)

本文件說明玩家如何發表評論、這些評論如何影響遊戲的整體評分，以及社群與管理員如何與這些評論互動。

---

## 1. 發表遊戲評論 (Posting a Review)

- **起點**：玩家在 `game_detail.html` 底部查看評論 (前端呼叫 `GET /api/games/:id/reviews`)，並點擊「寫評論」。
- **流程**：
  1. 呼叫 `POST /api/social/games/:id/reviews`，傳入 `attitude` (`POSITIVE` 或 `NEGATIVE`) 與 `content` (文字內容)。
  2. **購買限制檢查**：後端強制檢查該玩家是否在 `game_licenses` 中擁有這款遊戲，且 `status = 'ACTIVE'`。**未購買或是已退款、已吊銷授權的玩家無法發表評論**。這確保了「壓倒性好評」等標籤的公信力。
  3. 將評論寫入 `reviews` 表，預設狀態為 `'VISIBLE'`。
  4. **觸發自動計分 (Trigger Rating Update)**：
     - 在寫入評論成功後，後端會自動呼叫 `updateGameOverallRating(game_id)`。
     - 該函式會統計此遊戲所有 `VISIBLE` 狀態的 `POSITIVE` 數量，並計算比例。
     - 根據正面評價比例，將結果四捨五入後存回 `games.overall_rating` (滿分 5.0 的小數)。
     - (註：前端會將 0.0 ~ 5.0 的數字再轉換成如「壓倒性好評」、「褒貶不一」的文字)。
- **終點**：玩家發布成功，頁面重新載入，分數與評價標籤隨之更新。

---

## 2. 評論回覆 (Replying to Reviews)

- **起點**：其他玩家在看評論時，點擊某篇評論下的「回覆」。
- **流程**：
  1. 呼叫 `POST /api/social/reviews/:review_id/replies`。
  2. 不需要擁有該遊戲的授權，任何 `ACTIVE` 的玩家都可以參與討論。
  3. 將回覆寫入 `review_replies` 表。
  4. **刪除回覆**：若玩家想刪除自己的回覆，可呼叫 `DELETE /api/social/reviews/replies/:reply_id`。
- **終點**：回覆將以巢狀或列表形式附掛在該評論下方。
