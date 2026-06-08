# 2. 開發者遊戲上架與管理 (Developer Game Management)

這份文件描述了具有 `DEVELOPER` 權限的帳號，如何在 VaporAuror 平台上發布、管理與下架自己的遊戲。

---

## 1. 遊戲草稿建立 (Create Game Draft)

- **起點**：開發者在 `dev_dashboard.html` 查看自己的遊戲列表 (呼叫 `GET /api/developer/games`)，並點擊「新增遊戲」。
- **流程**：
  1. 輸入基本的 `Title` (標題)、`Price` (定價)、`Description` (支援 Markdown 格式的介紹)。
  2. 呼叫 `POST /api/developer/games`。
  3. 後端驗證使用者身分是否為 `DEVELOPER`，並將開發者 ID 綁定到該款遊戲。
  4. 遊戲建立完成，初始狀態預設為 `'DRAFT'` (草稿)，並在 `games` 資料表生成唯一 `game_id`。此狀態下遊戲不會出現在商店首頁。
- **終點**：畫面上顯示「建立成功」，自動跳轉進入 `edit_game.html?id={game_id}` 的進階編輯介面，開發者可繼續上傳素材與標籤。

---

## 2. 遊戲素材與檔案上傳 (Upload Media & Game Files)

- **起點**：在 `edit_game.html` 內的「上傳多媒體素材」與「上傳遊戲主檔」區塊。
- **流程**：
  1. 開發者選擇圖片 (封面圖、宣傳圖)、影片 (預告片) 或是 ZIP/EXE 遊戲主程式檔。
  2. 呼叫 `POST /api/developer/games/:id/media` (使用 `multipart/form-data`)。
  3. 後端根據檔案類型 (`media_type`) 進行不同的實體儲存處理：
     - **圖片或影片 (`media` / `thumbnail`)**：會將檔案內容進行 SHA-256 雜湊 (Hash) 產生新的檔名，並存入 `assets/images/{game_id}/{hash}.{ext}`，避免檔名衝突。
     - **遊戲主檔 (`game_file`)**：**不會進行 Hash**，而是保留開發者上傳的**原始檔名**，存入專屬資料夾 `assets/game-files/{game_id}/{original_filename}`。
  4. 實體路徑會轉換成對應的虛擬路由：圖片對應至 `/media/images/{game_id}/{hash}.{ext}`，遊戲檔案則對應至受保護的下載路由 `/downloads/{game_id}/{original_filename}`，並寫入 `game_media` 表。
  5. **刪除素材**：若開發者欲刪除，可呼叫 `DELETE /api/developer/games/:id/media/:media_id` 移除不要的素材。
- **終點**：前端收到新素材的 URL，並立刻渲染預覽圖或顯示檔案名稱。

---

## 3. 標籤管理 (Tags Management)

- **起點**：開發者在編輯頁面想要為遊戲加上標籤 (例如：RPG, Action)。
- **流程**：
  1. 系統可以列出目前資料庫中所有的全域標籤 (`GET /api/tags`)。
  2. 若無適合標籤，開發者可呼叫 `POST /api/developer/tags` 創建新標籤。
  3. 呼叫 `POST /api/developer/games/:id/tags` 將 `tag_id` 與 `game_id` 綁定。這會在 `game_tags` (多對多關聯表) 中新增一筆紀錄。
  4. 若需解除標籤，呼叫 `DELETE /api/developer/games/:id/tags/:tag_id` 將其移除。
- **終點**：遊戲分類變得精準，玩家在首頁可以透過標籤篩選出這款遊戲。

---

## 4. 正式發佈上架 (Publish Game)

- **起點**：開發者在 `dev_dashboard.html` 看到草稿狀態的遊戲，點擊「正式發佈」。
- **流程**：
  1. 呼叫 `PUT /api/developer/games/:id/publish`。
  2. **【重要約束點】**：後端嚴格驗證該遊戲**是否至少擁有 1 個標籤 (Tag)**。
  3. 若條件不符，退回 HTTP 400，並提示開發者前往編輯頁面補齊資料。
  4. 驗證通過後，將遊戲狀態由 `'DRAFT'` 更新為 `'ACTIVE'`。
- **終點**：遊戲正式上架，出現在商店首頁，玩家可開始瀏覽與購買。

---

## 5. 遊戲資訊更新與自主下架 (Update & Take Down)

- **起點**：開發者想要修改價格、介紹，或是決定停止販售這款遊戲。
- **流程 (更新資訊)**：
  - 呼叫 `PUT /api/developer/games/:id`，傳入新的 `Price` 與 `Description`。
- **流程 (自主下架 Delete)**：
  1. 呼叫 `DELETE /api/developer/games/:id`。
  2. **重要安全檢查**：後端會驗證這個 `game_id` 的 `DeveloperID` 是否與當前登入者相符，防止跨權限刪除。
  3. 通過驗證後，將該遊戲的狀態 `games.status` 設為 `'TAKEN_DOWN'` (軟刪除)。
  4. **與管理員下架的差異**：開發者自主下架時，**不會**連帶撤銷玩家的 `game_licenses`。這是保障消費者的基本權益，也就是「買過的玩家依然可以在遊戲庫中找到它並下載遊玩，只是商店不再開放新玩家購買」。
- **終點**：遊戲在商店首頁隱藏，進入詳細頁面會顯示「此遊戲已下架」，且無法加入購物車。

---

## 6. 遊戲銷售數據統計 (Sales Stats)

- **起點**：開發者想要知道自己的遊戲賣得如何。
- **流程**：
  - 呼叫 `GET /api/developer/games/:id/stats`。
  - 後端會從 `transaction_items` 加總銷售數量與總收入。
- **終點**：在儀表板呈現銷售業績數據。
