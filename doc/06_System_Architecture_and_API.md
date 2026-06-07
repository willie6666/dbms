# 十. 系統架構與 API 規劃

## 1. 專案三層架構概述
* **前端 (Frontend)**
  * 技術棧: HTML, CSS, Vanilla JS
  * 職責: 負責刻畫使用者介面、監聽使用者互動、儲存狀態 (JWT Token) 並向後端發起非同步請求。
* **連線層 (API/Auth)**
  * 技術棧: RESTful API, JSON, JWT (JSON Web Token)
  * 職責: 前後端溝通橋樑，統一使用 JSON 格式，透過 HTTP Status Codes 表達結果。
* **後端 (Backend)**
  * go.Gin: 輕量級 Web 框架。負責攔截 HTTP 請求、路由分發、以及執行中介軟體 (AuthMiddleware, RoleMiddleware)。
  * go.Controller: 業務邏輯層。將請求參數翻譯成 Go 邏輯，執行驗證、加解密、呼叫資料庫。
  * go.GORM: ORM (物件關聯對映) 層。負責將 Go 的 Object 翻譯成 PostgreSQL 的 SQL，並將回傳轉譯回 Go Object。
  * go.Driver (pgx): 底層驅動。負責建立 TCP 連線，送往資料庫並取回結果。
* **資料庫 (Database)**
  * 技術棧: PostgreSQL

## 2. 前端架構圖
* 架構類型: 多頁面應用程式 (MPA, Multi-Page Application)
* 核心語言: HTML5, CSS3, Vanilla JavaScript (無框架)
* 狀態管理: 依賴瀏覽器的 localStorage
* 渲染機制: 後端伺服器提供靜態 HTML 骨架，前端 JS 載入後發起 AJAX 請求，透過 DOM API 動態寫入內容。
* API 溝通層: `assets/js/api.js`
* 共用邏輯與狀態層: `assets/js/main.js`

## 3. 後端架構圖
* 架構類型: RESTful API Server
* 核心語言: Golang (Go 1.21+)
* 網路框架: Gin Web Framework
* 資料庫通訊: GORM
* 身分驗證: JWT 無狀態驗證機制
* 密碼加密: Bcrypt
* 網路與路由層: `routes` and `main.go`
* 中介軟體防護層: `middleware` (Auth, Role)
* 業務邏輯與控制層: `controllers`
* 資料映射層: `models`

## 4. API 流程圖與連結資料庫
### SQL 查詢流程
1. browser (Frontend 送出 HTTP 請求)
2. RESTful API (送到 go server 的 port)
3. go.Gin (攔截 轉送任務)
4. go.Controller (將 /api/* 翻譯成 task)
5. go.GORM (把 task 翻譯成 PostgreSQL 語法)
6. go.Driver (去資料庫 撈資料 raw data)
7. go.GORM (把 raw 轉譯成 go 語法的 object)
8. go.Gin (打包成 json 送回去)
9. browser (Frontend 接收 json 資料)

### 註冊與登入流程
* Bcrypt 加密 Hash 後寫入 DB。
* 登入比對 Hash 後，伺服器簽發 JWT 金鑰給 Client。
* Client 將 JWT 儲存在 LocalStorage 供後續請求使用。

### Role 驗證流程
* Auth Middleware: 解碼 JWT，檢查是否過期及合法。
* Role Middleware: 檢查階級是否符合該 API 規定 (如 DEVELOPER, ADMIN)。

## 5. RESTful API 設計與清單
(詳細 API 請參考專案內的 `api/api_spec.md` 與 `api/api_list.txt`)
包含：
1. 使用者與權限 (Users & Auth)
2. 商店與遊戲 (Store & Games)
3. 訂單、購物車與客服 (Transactions & Carts)
4. 遊戲庫與願望清單 (Library & Wishlist)
5. 社交、評論與通訊 (Social & Reviews)

# 十一. 心得
# reply_reply 多了一個create time 相關東西季的修正
(我們已成功在系統與文件中修正 review_replies 遺漏 created_at 的問題)
