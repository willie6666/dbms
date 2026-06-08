# 8. 遊戲搜尋與探索 (Search & Discovery)

本文件描述玩家在商店首頁 (`index.html`) 如何透過各種維度（關鍵字、標籤、價格區間、開發者等）進行過濾與搜尋，找到感興趣的遊戲。

---

## 1. 關鍵字模糊搜尋 (Keyword Fuzzy Search)

- **起點**：玩家在導覽列或首頁的 Search Box 輸入關鍵字（例如："Action" 或 "StudioAurora"）。
- **流程**：
  1. 前端將關鍵字放在 `q` 參數，呼叫 `GET /api/games?q={keyword}`。
  2. 後端組建 SQL 查詢，針對以下四個維度進行 `ILIKE` 模糊比對：
     - `games.title` (遊戲標題)
     - `games.description` (遊戲簡介)
     - `q_tags.tag_name` (遊戲綁定的標籤名稱)
     - `q_developers.username` (開發者的名稱)
  3. 只要上述任一欄位包含關鍵字，該遊戲就會被篩選出來。
- **終點**：前端收到結果並重新渲染遊戲卡片列表。

---

## 2. 進階過濾器與排序 (Advanced Filters & Sorting)

- **起點**：玩家在首頁左側的過濾面板 (Filter Panel) 操作勾選框或拉桿。
- **支援的過濾條件**：
  - **指定開發者** (`?developer={name}`)：精準尋找某位開發者的所有作品。
  - **指定標籤** (`?tag={name}`)：精準篩選帶有特定分類標籤 (如：RPG, Simulation) 的遊戲。
  - **價格區間** (`?min_price={min}&max_price={max}`)：設定最低與最高價格。
  - **隱藏已擁有** (`?hide_owned=true`)：若玩家已登入 (Header 帶有 JWT)，後端會自動透過子查詢 (`NOT EXISTS`) 過濾掉該玩家在 `game_licenses` 中已經擁有 (`ACTIVE`) 的遊戲。
- **排序條件** (`?sort={type}`)：
  - `price_asc` (預設)：價格從低到高排列。
  - `price_desc`：價格從高到低排列。
- **流程**：
  1. 前端將玩家的所有過濾條件打包成 Query String (`api.js` 內的 `apiGetGames` 負責處理)。
  2. 後端利用 GORM 動態組合 `Where` 條件。
  3. **基礎防護**：無論如何過濾，後端最底層一定會加上 `games.status = 'ACTIVE'`，確保任何已被下架 (`TAKEN_DOWN`) 的遊戲絕不會出現在搜尋結果中。
- **終點**：精準呈現符合玩家需求、且可正常購買的遊戲清單。
