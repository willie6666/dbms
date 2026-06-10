# 5. 遊戲評論與社群互動 (Review & Rating)

本文件說明玩家如何發表評論、這些評論如何影響遊戲的整體評分，以及社群與管理員如何與這些評論互動。

---

## 1. 發表遊戲評論 (Posting a Review)

- **起點**：玩家在 `game_detail.html` 底部查看評論 (前端呼叫 `GET /api/games/:id/reviews`)，並點擊「寫評論」。
- **流程**：
  1. 呼叫 `POST /api/social/games/:id/reviews`，傳入 `attitude` (`POSITIVE` 或 `NEGATIVE`) 與 `content` (文字內容)。可選傳入 `post_as_role` 指定發布身分。
  2. **購買限制檢查與身分特權 (Bypass Ownership)**：
     - 若為一般玩家 (預設)，後端強制檢查該玩家是否在 `game_licenses` 中擁有這款遊戲，且 `status = 'ACTIVE'`。**未購買或是已退款、已吊銷授權的玩家無法發表一般評論**。
     - 若選擇以特權身分發布 (例如 `ADMIN`, `CSR` 或 `AUTHOR`) 且具備對應權限，後端會**略過**購買檢查。
       *(註：前端在解鎖開發者的 `AUTHOR` 選項時，運用了 `(user.id || user.user_id) == developerId` 高包容性探測機制，確保身份不漏接)*
  3. **隱藏字首機制與寫入**：將評論寫入 `reviews` 表。如果為特權發布，後端會在文字最前方自動加入隱藏前綴 (例如 `[ROLE:ADMIN]`)，全程**不更動 Database Schema**。狀態預設為 `'VISIBLE'`。
  4. **讀取評論**：當前端呼叫 `GET /api/games/:id/reviews`，後端會自動剔除前綴字首，並將身分轉成 `posted_as_role` 交給前端渲染專屬的身分標籤 (Badge)。
  5. **觸發自動計分 (Trigger Rating Update)**：
     - 在寫入評論成功後，後端會自動呼叫 `updateGameOverallRating(game_id)`。
     - 該函式會統計此遊戲所有 `VISIBLE` 狀態的 `POSITIVE` 數量，並計算比例。
     - 根據正面評價比例，將結果四捨五入後存回 `games.overall_rating` (滿分 5.0 的小數)。
     - (註：前端會將 0.0 ~ 5.0 的數字再轉換成如「壓倒性好評」、「褒貶不一」的文字)。
- **終點**：玩家發布成功，頁面重新載入，分數與評價標籤隨之更新。

---

## 2. 評論回覆 (Replying to Reviews)

- **起點**：其他玩家在看評論時，點擊某篇評論下的「回覆」。
- **流程**：
  1. 呼叫 `POST /api/social/reviews/:review_id/replies`，傳入 `content`，並可選傳入 `post_as_role` 指定發布身分。
  2. 不需要擁有該遊戲的授權，任何 `ACTIVE` 的玩家都可以參與討論。如果選擇以 `ADMIN`, `CSR` 或 `AUTHOR` 身分發布且具備對應權限，同樣會啟動「隱藏字首機制」。
  3. 將回覆寫入 `review_replies` 表 (附帶隱藏字首，若為特權)。
  4. **刪除回覆**：若玩家想刪除自己的回覆，可呼叫 `DELETE /api/social/reviews/replies/:reply_id`。
- **終點**：回覆將以巢狀或列表形式附掛在該評論下方。
