# 10. 遊戲庫與遊玩下載 (Game Library & Play)

本文件描述玩家購買遊戲後，在「個人遊戲庫 (`library.html`)」中瀏覽已購遊戲、下載檔案以及啟動遊戲的完整業務邏輯。由後端 `library_controller.go` 負責處理。

---

## 1. 遊戲庫列表 (Library Listing)

- **起點**：登入的玩家點擊導覽列進入「遊戲庫 (`library.html`)」。
- **流程 (`GET /api/protected/library`)**：
  1. 系統透過 JWT 確認使用者身分 (`userID`)。
  2. 後端查詢 `game_licenses` 表，過濾條件為：
     - `user_id = ?`
     - `status = 'ACTIVE'` (確保只有處於生效狀態的授權才會顯示)。
  3. 透過 GORM Preload 帶出關聯的 `Game` 主檔與 `Game.Media` (包含封面圖等)，一併回傳給前端。
- **防護機制**：
  - 若玩家的遊戲授權被管理員撤銷 (`REVOKED`) 或因退款被凍結 (`FROZEN`)，該遊戲將**自動從遊戲庫中消失**，前端無法取得該資料。
- **終點**：玩家在頁面上看到網格排列的已擁有遊戲清單。

---

## 2. 下載遊戲檔案 (Download Game File)

- **起點**：玩家在遊戲庫中，針對特定遊戲點擊「下載」按鈕。
- **流程 (`GET /api/protected/library/:game_id/download`)**：
  1. **授權驗證**：後端嚴格檢查 `game_licenses` 中，該玩家對該 `game_id` 是否擁有 `status = 'ACTIVE'` 的紀錄。若無授權，直接回傳 HTTP 403 Forbidden。
  2. **尋找檔案實體**：在 `game_media` 表中尋找 `game_id = ?` 且 `media_type = 'game_file'` 的紀錄，取得 `file_url`。
  3. **路徑安全防護 (Path Traversal Prevention)**：
     - 系統要求 `file_url` 必須以 `/downloads/` 開頭。
     - 系統會移除前綴並透過 `filepath.Clean` 清理路徑。
     - **關鍵防護**：若清理後的路徑包含 `..` (回退目錄) 或以 `/` 開頭，後端會立即阻擋並回傳 `Invalid game file path`，防止駭客利用目錄遍歷漏洞讀取伺服器敏感檔案。
  4. **回傳實體檔案**：安全驗證通過後，將路徑映射至伺服器實體目錄 `assets/game-files/...` 並將二進位檔案以 Attachment 形式串流回傳給前端。
- **終點**：瀏覽器觸發檔案下載。

---

## 3. 啟動遊戲 (Play Game)

- **起點**：玩家在遊戲庫中，針對特定遊戲點擊「遊玩」按鈕。
- **流程 (`GET /api/protected/library/:game_id/play`)**：
  1. **授權驗證**：如同下載功能，後端嚴格檢查 `game_licenses` 中該玩家對該遊戲是否有 `ACTIVE` 狀態的授權。
  2. **核發憑證**：驗證通過後，後端會模擬回傳一組啟動憑證 (Token)，例如 `mock-play-token-12345`。
- **終點**：前端收到成功訊息與憑證後，以模擬方式彈出提示「啟動成功」，象徵遊戲用戶端已接手該憑證並開始執行遊戲。
