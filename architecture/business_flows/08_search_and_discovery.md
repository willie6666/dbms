# 8. 遊戲搜尋與探索 (Search & Discovery)

本文件描述玩家在商店首頁 (`index.html`) 如何透過前端的三大控制面板（頂部搜尋列、左側過濾面板、右側進階控制）進行過濾與搜尋，並詳述其對應的底層邏輯。

---

## 1. 頂部搜尋列 (Top Search Bar)

- **UI 位置**：頁面正上方、最醒目的橫向搜尋框。
- **功能**：關鍵字模糊搜尋 (Keyword Fuzzy Search)。玩家輸入文字後點擊「搜尋」按鈕。點擊進入個別遊戲時，前端會呼叫 `GET /api/games/{id}` 顯示詳情。
- **底層邏輯**：
  1. 前端將關鍵字放入 `q` 參數，呼叫 `GET /api/games?q={keyword}`。
  2. 後端組建 SQL 查詢，針對以下四個維度進行 `ILIKE` 模糊比對：
     - `games.title` (遊戲標題)
     - `games.description` (遊戲簡介)
     - `q_tags.tag_name` (遊戲綁定的標籤名稱)
     - `q_developers.username` (開發者的名稱)
  3. 只要上述任一欄位包含關鍵字，該遊戲就會被篩選出來。

---

## 2. 左側過濾面板 (Left Filter Panel)

- **UI 位置**：頁面左側的垂直區塊，包含「Tag 瀏覽」與「價格篩選」兩大區塊。
- **功能與邏輯**：
  - **Tag 瀏覽 (標籤過濾)**：
    - 前端會列出如 `RPG`, `Action`, `Racing`, `Simulation` 等常見標籤（或從 `/api/tags` 動態獲取）。
    - 當玩家點擊特定標籤，前端發送 `?tag={name}`。
    - 後端透過 `JOIN game_tags` 與 `tags` 資料表，精準過濾出綁定該標籤的遊戲。選擇「所有遊戲」則會清除 `tag` 參數。
  - **開發者篩選 (Developer Filter)**：
    - 當玩家從特定開發者頁面進入，或點擊開發者名稱時，發送 `?developer={username}`。
    - 後端會透過 `ILIKE` 執行「精確的字串比對」(不加上 `%` 萬用字元)，確保輸入 `dev_1` 不會錯誤搜出 `dev_10`。
  - **價格篩選 (Price Filter)**：
    - 提供預設的價格區間選項，例如：
      - 「免費遊玩」：發送 `?max_price=0`。
      - 「NT$ 300 以下」：發送 `?max_price=300`。
      - 「NT$ 301 - 900」：發送 `?min_price=301&max_price=900`。
    - 後端收到 `min_price` 與 `max_price` 後，會直接在 SQL 中加上 `games.price >= ?` 與 `games.price <= ?` 的條件。

---

## 3. 右側進階控制面板 (Right Advanced Controls)

- **UI 位置**：遊戲列表右上角，包含一個「隱藏已購買」的開關 (Toggle) 與一個「排序方式」的下拉選單。
- **功能與邏輯**：
  - **隱藏已購買 (Hide Owned)**：
    - 這是一個 ON/OFF 切換開關。
    - 當開啟時，前端發送 `?hide_owned=true`，且必須在 Request Header 帶上使用者的 JWT Token。
    - 後端偵測到該參數與 Token 後，會自動利用子查詢 (`NOT EXISTS`) 到 `game_licenses` 檢查，**並額外比對 `games.developer_id != 當前使用者`**。只要該玩家擁有的授權是 `ACTIVE`，或是該遊戲正是由登入者自己開發的，該遊戲就會從搜尋結果中被剔除，讓玩家專注於探索尚未擁有的新遊戲。
  - **排序方式 (Sorting Dropdown)**：
    - 提供如「價格由低到高」或「價格由高到低」的選項。
    - 選擇後前端發送 `?sort=price_asc` 或 `?sort=price_desc`。
    - 後端將其轉換為 SQL 的 `ORDER BY games.price ASC/DESC`。

---

## 4. 綜合搜尋與最底層防護

- **綜合查詢**：上述三個面板的操作可以**疊加使用**。例如，玩家可以同時搜尋 "Action" (頂部)，勾選 "NT$ 300 以下" (左側)，並開啟 "隱藏已購買" (右側)。前端會將這些參數全部串接起來 (`?q=Action&max_price=300&hide_owned=true&sort=price_asc`) 一次發送給後端。
- **最底層防護**：無論玩家如何組合過濾條件，後端的查詢產生器最底層一定會加上 `WHERE games.status = 'ACTIVE'`。這確保了任何被管理員「強制下架」(`TAKEN_DOWN`) 或開發者自主下架的遊戲，絕對不會出現在任何搜尋結果中。
