# VaporAuror 三層架構與技術流程 (System Architecture & Technical Flow)

本文件詳細解析了 VaporAuror 專案的三層架構設計，並使用標準的資料流格式，深入剖析各項核心功能的實際運作流程。

---

## 1. 專案三層架構概述 (Three-Tier Architecture)

VaporAuror 採用經典的現代 Web 三層架構，職責分明：

- **前端 (Presentation Layer)**
  - **技術棧**: `HTML`, `CSS`, `Vanilla JS`
  - **職責**: 負責刻畫使用者介面、監聽使用者互動 (點擊、表單提交)、儲存狀態 (如 `localStorage` 中的 JWT Token) 並向後端發起非同步請求。
- **連線層 (API Layer)**
  - **技術棧**: `RESTful API`, `JSON`, `JWT (JSON Web Token)`
  - **職責**: 作為前後端溝通的橋樑，統一使用 JSON 格式傳遞資料，並透過 HTTP Status Codes 表達結果狀態。
- **後端 (Business & Data Access Layer)**
  - **`go.Gin`**: 輕量級 Web 框架。負責攔截 HTTP 請求、路由分發 (Routing)、以及執行中介軟體 (如 JWT 驗證 `AuthMiddleware`、權限控制 `RoleMiddleware`)。
  - **`go.Controller`**: 業務邏輯層。負責將 HTTP 請求的參數翻譯成 Go 程式邏輯，執行驗證、加解密、以及呼叫資料庫。
  - **`go.GORM`**: ORM (物件關聯對映) 層。負責將 Go 的 Object 語法翻譯成 PostgreSQL 的 SQL 語法，並將回傳的 Raw Data 轉譯回 Go Object。
  - **`go.Driver` (`pgx`)**: 底層驅動。負責建立 TCP 連線，將 SQL 語法送往資料庫並取回結果。
- **資料庫 (Data Layer)**
  - **技術棧**: `PostgreSQL`
  - **職責**: 持久化儲存所有的關聯式資料，確保資料的 ACID 特性。

---

## 2. 核心功能技術流程 (Core Technical Flows)

以下將專案中的關鍵行為，以「前端 -> API -> Router -> Controller -> ORM -> DB -> Response」的標準流程進行解構。

### 2.1 使用者登入流程 (User Login Flow)
```text
browser (Frontend 收集表單送出 HTTP POST {email, password})
    -> RESTful API (送到 /api/auth/login)
    -> go.Gin (攔截，轉送給 AuthController 的 Login 功能)
    -> go.Controller (收到請求，驗證格式後，準備透過 email 查詢使用者)
    -> go.GORM (將查詢翻譯成 SELECT * FROM users WHERE email = ?)
    -> go.Driver (去 PostgreSQL 撈取 raw data)
    -> go.GORM (把 raw data 轉譯成 Go 語法的 User Object)
    -> go.Controller (取得密碼雜湊值，呼叫 Bcrypt 套件進行【密碼比對】)
    -> go.Controller (密碼正確，呼叫 JWT 套件利用 HMAC 演算法【簽發 Token】)
    -> go.Gin (將 Token 與 User 資料打包成 json 送回去)
browser (Frontend 接收 json 資料，將 token 存入 localStorage 並跳轉首頁)
```

### 2.2 瀏覽/搜尋商店遊戲 (Browse Games Flow)
```text
browser (Frontend 進入首頁或輸入關鍵字，JS 送出 HTTP GET ?search=Cyber)
    -> RESTful API (送到 /api/games?search=Cyber)
    -> go.Gin (攔截，轉送給 GameController 的 GetGames 功能)
    -> go.Controller (解析 query 參數，發現有 search 關鍵字)
    -> go.GORM (翻譯成帶有 LIKE 的 SQL: SELECT * FROM games WHERE title ILIKE '%Cyber%')
    -> go.Driver (去 PostgreSQL 進行模糊搜尋撈出 raw data)
    -> go.GORM (把 raw data 轉譯成 []Game Object 陣列)
    -> go.Gin (打包成 json 陣列送回去)
browser (Frontend 接收 json 資料，透過 JS 動態生成 HTML DOM 顯示遊戲卡片)
```

### 2.3 加入購物車流程 (Add to Cart Flow)
```text
browser (Frontend 點擊「加入購物車」，JS 從 localStorage 提取 Token，送出 HTTP POST {game_id})
    -> RESTful API (送到 /api/protected/cart)
    -> go.Gin (進入 AuthMiddleware 驗證 JWT Token 的合法性與時效)
    -> go.Gin (Token 合法，解析出 user_id 掛載到 Context，轉交 CartController)
    -> go.Controller (收到 game_id 與 user_id，檢查是否已經在購物車或已擁有)
    -> go.GORM (翻譯成 SELECT ... WHERE user_id=? AND game_id=? 檢查重複)
    -> go.Driver (去 PostgreSQL 查詢確認未重複)
    -> go.Controller (準備建立新的 ShoppingCart Record)
    -> go.GORM (把 task 翻譯成 INSERT INTO shopping_carts 語法)
    -> go.Driver (去 PostgreSQL 寫入新資料)
    -> go.GORM (確認寫入成功，取得自動生成的 cart_item_id)
    -> go.Gin (打包成功訊息 json 送回去)
browser (Frontend 接收 json 資料，顯示「已加入購物車」並將按鈕反灰)
```

### 2.4 購物車結帳流程 (Checkout Flow - 牽涉 Transaction)
```text
browser (Frontend 點擊「確認結帳」，JS 帶上 Token 送出 HTTP POST)
    -> RESTful API (送到 /api/protected/checkout)
    -> go.Gin (攔截並通過 AuthMiddleware 驗證)
    -> go.Controller (開始處理 Checkout 邏輯)
    -> go.GORM (發送 BEGIN 語法，開啟【資料庫交易 Transaction】，確保 ACID)
    -> go.GORM (撈取該 user_id 的所有 shopping_carts 項目，並計算總額)
    -> go.GORM (翻譯成 INSERT INTO transactions 寫入主訂單)
    -> go.GORM (翻譯成 INSERT INTO transaction_items 寫入每筆訂單明細)
    -> go.GORM (翻譯成 INSERT INTO game_licenses 寫入遊戲庫授權，狀態為 ACTIVE)
    -> go.GORM (翻譯成 DELETE FROM shopping_carts WHERE user_id = ? 清空購物車)
    -> go.Driver (在 PostgreSQL 中一次性執行上述所有變更)
    -> go.GORM (確認全數成功，發送 COMMIT 語法提交交易；若有錯則 ROLLBACK)
    -> go.Gin (打包成功訊息 json 送回去)
browser (Frontend 接收 json 資料，清空畫面並提示前往遊戲庫查看)
```

### 2.5 客服退款審核流程 (CSR Refund Approval Flow)
```text
browser (Frontend CSR 管理員點擊「核准退款」，JS 帶上 CSR Token 送出 HTTP PUT {status: "APPROVED"})
    -> RESTful API (送到 /api/csr/refunds/{id})
    -> go.Gin (攔截並通過 AuthMiddleware 解析出 user_id 與 role)
    -> go.Gin (進入 RoleMiddleware 驗證，確認 Role == 'CSR' 或 'ADMIN'，放行)
    -> go.Controller (準備更新退款單狀態與收回遊戲授權)
    -> go.GORM (開啟【資料庫交易 Transaction】)
    -> go.GORM (翻譯成 UPDATE refund_requests SET status='APPROVED')
    -> go.GORM (透過 transaction_item_id 找到對應的 game_licenses)
    -> go.GORM (翻譯成 UPDATE game_licenses SET status='REVOKED' 收回遊戲遊玩權限)
    -> go.Driver (去 PostgreSQL 執行更新)
    -> go.GORM (發送 COMMIT 提交)
    -> go.Gin (打包成功訊息 json 送回去)
browser (Frontend 接收 json 資料，移除畫面上的待處理卡片)
```
