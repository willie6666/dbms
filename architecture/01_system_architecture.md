# VaporAuror 三層架構與技術流程 (System Architecture & Technical Flow)

本文件詳細解析了 VaporAuror 專案的三層架構設計，並使用標準的資料流格式，深入剖析各項核心功能的實際運作流程。

---

## 1. 專案三層架構概述 (Three-Tier Architecture)

VaporAuror 採用經典的現代 Web 三層架構，職責分明：

- **基礎設施 (Infrastructure)**
  - **技術棧**: `Docker`, `Docker Compose`
  - **職責**: 負責將前端、後端與資料庫等所有服務打包成獨立的容器 (Container)，提供統一且隔離的執行環境，並透過內部網路 (Docker Network) 讓各層級服務能安全互相溝通。
- **表示層 / 前端 (Presentation Layer)**
  - **技術棧**: `Caddy`, `HTML`, `CSS (bulma)`, `Vanilla JS`
  - **職責**: 
    - **Caddy**: 作為面向瀏覽器的單一入口，負責託管與分發靜態網頁檔案，並將 API 與媒體檔案的請求反向代理 (Reverse Proxy) 至後端。
    - **Vanilla JS**: 負責刻畫使用者介面、監聽使用者互動 (點擊、表單提交)、儲存狀態 (JWT Token) 並向後端發起非同步請求。
- **連線層 (Connection Layer)**
  - **技術棧**: `RESTful API`, `JSON`, `JWT (JSON Web Token)`
  - **職責**: 作為前後端溝通的橋樑，統一使用 JSON 格式傳遞資料，並透過 HTTP Status Codes 表達結果狀態。
- **應用層 / 後端 (Application Layer)**
  - **`go.Gin`**: 輕量級 Web 框架。負責攔截 HTTP 請求、路由分發、以及執行中介軟體 (JWT 驗證 `AuthMiddleware`、權限控制 `RoleMiddleware`)。
  - **`go.Controller`**: 業務邏輯層。負責將 HTTP 請求的參數翻譯成 Go 程式邏輯，執行驗證、加解密、以及呼叫資料庫。
  - **`go.GORM`**: ORM (物件關聯對映) 層。負責將 Go 的 Object 語法翻譯成 PostgreSQL 的 SQL 語法，並將回傳的 Raw Data 轉譯回 Go Object。
  - **`go.Driver` (`pgx`)**: 底層驅動。負責建立 TCP 連線，將 SQL 語法送往資料庫並取回結果。
- **資料庫層 (Database Layer)**
  - **技術棧**: `PostgreSQL`
  - **職責**: 負責資料的持久化儲存、關聯查詢，確保資料的一致性與完整性。

---

## 2. 三層架構文字圖表 (表示層 / 應用層 / 資料庫層)

```text
               [ 使用者 / 瀏覽器 Browser ]
                           │
                           ▼ HTTP / HTTPS 請求
=====================================================================
【 表示層 (Presentation Layer) 】

       ┌──────────────── Caddy (Web Server) ────────────────┐
       │                                                    │
       ▼ (託管靜態網頁)                              ▼ (反向代理)
[ HTML / CSS(Bulma) / JS ]                 ( /api/* , /media/* )
(UI 渲染、使用者互動、JWT儲存)                       │
===========================│========================│================
【 連線層 (Bridge) 】      ▼                        ▼
                   RESTful API (JSON 格式資料傳遞 + JWT 驗證)
===========================│========================│================
【 應用層 (Application Layer) 】                    │
                                                    ▼
       ┌────────────────── Go Backend ──────────────────┐
       │                                                │
       │  1. go.Gin (路由分發、Auth/Role 中介軟體)        │
       │                         │                      │
       │  2. go.Controller (核心業務邏輯、加解密)         │
       │                         │                      │
       │  3. go.GORM (ORM 物件關聯對映)                  │
       │                         │                      │
       │  4. go.Driver (pgx) (底層 TCP 連線驅動)         │
       └─────────────────────────┬──────────────────────┘
                                 │
                                 ▼ SQL 語法 (查詢 / 寫入)
=====================================================================
【 資料庫層 (Database Layer) 】
                                 │
                                 ▼
                     ┌───────────────────────┐
                     │ PostgreSQL (資料持久化) │
                     └───────────────────────┘
=====================================================================
 * 基礎設施保護傘：上述三層之 Caddy、Go Backend 與 PostgreSQL 
   皆受到 [ Docker Compose ] 之容器化隔離環境包覆與統一管理。
```

---

## 3. 核心功能技術流程 (Core Technical Flows)

以下將專案中的關鍵行為，以「前端 -> API -> Router -> Controller -> ORM -> DB -> Response」的標準流程進行解構。

### 3.1 使用者登入流程 (User Login Flow)
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

### 3.2 瀏覽/搜尋商店遊戲 (Browse Games Flow)
```text
browser (Frontend 進入首頁或輸入搜尋條件，JS 送出 HTTP GET ?q=Cyber&developer=DevA&hide_owned=true)
    -> RESTful API (送到 /api/games?q=Cyber&developer=DevA&hide_owned=true)
    -> go.Gin (攔截，轉送給 GameController 的 GetGames 功能)
    -> go.Controller (解析 query 參數，包含關鍵字模糊搜尋、developer 精確配對，並解析隱藏已購買功能)
    -> go.GORM (翻譯成帶有 ILIKE 的 SQL，針對 developer 使用精確字串比對)
    -> go.GORM (因 hide_owned=true，額外加入條件過濾「已在 game_licenses 獲得授權的遊戲」與「當前登入者自己開發的遊戲」)
    -> go.Driver (去 PostgreSQL 進行複合條件搜尋撈出 raw data)
    -> go.GORM (把 raw data 轉譯成 []Game Object 陣列)
    -> go.Gin (打包成 json 陣列送回去)
browser (Frontend 接收 json 資料，透過 JS 動態生成 HTML DOM 顯示遊戲卡片)
```

### 3.3 加入購物車流程 (Add to Cart Flow)
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

### 3.4 購物車結帳流程 (Checkout Flow - 牽涉 Transaction)
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

### 3.5 客服退款審核流程 (CSR Refund Approval Flow)
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

### 3.6 發表特權評論與隱藏字首機制 (Privileged Review & Hidden Prefix Flow)
```text
browser (Frontend 偵測到使用者具備開發者權限，解鎖身分選單。送出帶有 post_as_role="AUTHOR" 的 HTTP POST)
    -> RESTful API (送到 /api/social/games/{id}/reviews)
    -> go.Gin (攔截並通過 AuthMiddleware 解析身分)
    -> go.Controller (發現附帶了特權發布請求，跳過購買 game_licenses 的擁有權檢查)
    -> go.Controller (在存入資料庫前，將字串加工：content = "[ROLE:AUTHOR]" + content)
    -> go.GORM (將加工後的字串寫入 reviews 資料表，完美避開修改 Database Schema)
    -> go.Driver (去 PostgreSQL 寫入)
    -> go.Gin (回傳成功)
browser (後續透過 GET 取得評論，後端自動剝離字首並產生 posted_as_role 欄位回傳)
    -> browser (前端利用 posted_as_role，在該則評論旁動態渲染出橘黃色的「AUTHOR」官方權威徽章)
```
