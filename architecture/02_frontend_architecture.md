# VaporAuror 前端架構與技術解析 (Frontend Architecture)

本文件詳細解析 VaporAuror 專案的前端底層架構，採用純原生（Vanilla JS）的多頁面架構（MPA），職責切割清晰，確保前後端溝通順暢。

---

## 1. 前端技術棧與整體架構概述
- **架構類型**: 多頁面應用程式 (MPA, Multi-Page Application)
- **核心語言**: HTML5, CSS3 (搭配 Bulma CSS 框架), Vanilla JavaScript (原生 JS，無使用 React / Vue 等框架)
- **狀態管理**: 依賴瀏覽器的 `localStorage` (儲存 JWT Token 與 User Info)
- **畫面渲染機制**: Caddy Web Server 提供靜態 HTML 骨架，前端 JS 載入後發起 AJAX (Fetch) 請求，取得 JSON 資料後，透過 DOM API 動態寫入內容。

---

## 2. 檔案目錄結構與職責劃分 (Tree Architecture)

```text
frontend/
├── index.html                  # 系統總入口 (商店首頁：顯示熱門與所有遊戲)
├── assets/                     # 靜態資源共用區 (前端的核心引擎)
│   ├── css/
│   │   └── style.css           # 全域樣式表 (載入 Bulma CSS 框架並定義額外色票、版面與 UI)
│   └── js/
│       ├── api.js              # [API 溝通層] 集中封裝所有與後端的 RESTful API 呼叫
│       └── main.js             # [共用邏輯層] 全域導覽列渲染、身分權限驗證、Toast 提示組件
└── pages/                      # 功能畫面區 (依業務邏輯高度模組化)
    ├── auth/                   # [身分驗證模組]
    │   ├── login.html          # 登入畫面
    │   └── register.html       # 註冊畫面
    ├── store/                  # [商店模組]
    │   ├── search.html         # 關鍵字與標籤搜尋結果頁面
    │   └── game_detail.html    # 遊戲詳細介紹、評論展示與加入購物車操作
    ├── user/                   # [玩家專屬模組]
    │   ├── library.html        # 個人遊戲庫 (遊玩、下載)
    │   ├── cart.html           # 購物車結帳頁面
    │   ├── profile.html        # 個人資料修改
    │   ├── history.html        # 交易歷史紀錄
    │   ├── refund_request.html # 退款申請頁面
    │   ├── wishlist.html       # 願望清單
    │   └── social.html         # 好友、訊息與黑名單管理
    └── dashboard/              # [後台管理模組 (嚴格依角色隔離)]
        ├── admin_dashboard.html # 系統管理員後台 (管理所有帳號權限與強制下架遊戲)
        ├── csr_dashboard.html   # 客服人員後台 (審核玩家退款申請)
        ├── dev_dashboard.html   # 開發者後台 (上架新遊戲、查看銷售數據)
        └── edit_game.html       # 開發者專屬：遊戲內容與素材編輯器
```

---

## 3. 核心模組深度解析 (Core Modules Details)

### 🔴 API 溝通層 (API Layer) - `assets/js/api.js`
- **職責**: 作為前端與後端溝通的「唯一閘口」。所有跨頁面的 API 請求都必須且只能透過此檔案發出。
- **核心機制**:
  - **自動授權 `authFetch()`**: 封裝了原生的 `fetch()` API。每次發送受保護的請求時，會自動從 `localStorage` 取出 JWT Token，並注入 HTTP Header (`Authorization: Bearer <token>`) 中。
  - **401 攔截器**: 負責全域錯誤攔截。如果後端回傳 HTTP Status `401 Unauthorized` (Token 過期或被竄改)，會自動強制登出，清空快取並導向 `login.html`。
  - **同源代理**: 設定 `API_BASE = ''`，所有 `/api/*` 的請求都送往與前端相同的網域 (`localhost:3000`)，再由底層的 Caddy 伺服器反向代理至 Go Server (`backend:8000`)，徹底免除了 CORS 的設定麻煩。
  - **業務函數映射**: 提供具語意化的函數如 `apiGetGames()`, `apiAddToCart()`, `apiApproveRefund()`，隱藏底層的 URL 路徑與 HTTP Method 差異，讓各頁面的程式碼保持乾淨。

### 🟡 共用邏輯與狀態層 (Shared Logic Layer) - `assets/js/main.js`
- **職責**: 處理全站共用的 DOM 操作與全域狀態快取。
- **核心機制**:
  - **動態導覽列 `renderHeader()`**: 根據當前使用者的登入狀態（遊客、已登入）與角色權限（`USERS`, `ADMIN`, `CSR`, `DEVELOPER`），動態生成右上角的選單按鈕（例如：只有開發者才會看到「開發者後台」按鈕）。
  - **全域提示組件 `showToast(msg, type)`**: 畫面的右下角彈出通知 (Success/Error)，處理非同步請求完成後的使用者回饋。
  - **狀態管理 API**: 提供 `getToken()`, `getCurrentUser()`, `getCurrentRole()` 等封裝函式，避免直接操作 localStorage 造成的拼字錯誤。

### 🔵 視圖與渲染層 (View Layer) - `pages/**/*.html`
- **職責**: 負責骨架繪製與畫面獨有的業務邏輯 (Page-specific logic)。
- **核心機制**:
  - **高內聚低耦合**: 每個 `.html` 檔案各自獨立，保有自己的專屬 JavaScript 邏輯，例如 `cart.html` 內的 JS 只負責解析購物車陣列並迴圈印出 `div.cart-item`，絕不干涉其他頁面。
  - **事件綁定**: 透過 `onclick`, `onsubmit` 攔截使用者的表單送出與點擊，收集資料後交由 `api.js` 處理。

---

## 4. 頁面生命週期與渲染流程 (Page Lifecycle)

當使用者在瀏覽器中打開任何一個頁面（例如「遊戲詳情頁」）時，前端的完整運作流程如下：

```text
1. [載入階段] 瀏覽器向 Server 請求 HTML 檔案。
2. [解析階段] 瀏覽器解析 DOM Tree，並載入 style.css 畫出基礎外觀與顏色。
3. [依賴載入] 執行位於 HTML 文件底部的 <script src="/assets/js/api.js"> 與 main.js。
4. [初始化] 觸發原生事件 document.addEventListener("DOMContentLoaded")。
5. [共用渲染] main.js 自動執行 renderHeader()，讀取 Token 判斷身分並畫出導覽列。
6. [資料請求] 該頁面專屬的腳本被觸發，向 api.js 呼叫函式 (如 apiGetGameDetails(id))。
7. [等待回覆] api.js 將請求包裝後送往同源 `/api/...`，由 Caddy 代理到後端。
8. [動態更新] 收到後端 JSON 回應後，透過 document.getElementById() 將數據 (如價格、評論) 動態寫入畫面。
9. [等待互動] 頁面渲染完畢，靜待使用者的點擊或表單輸入操作。
```

---

## 5. 前端進階技術與模式 (Advanced Frontend Patterns)

隨著系統演進，前端導入了以下幾項關鍵技術與架構模式，以處理更複雜的互動與安全性：

### 1. 狀態防護與邏輯解耦層 (UI Protection Layer)
- 雖然在後端架構上（例如：願望清單）採行了**「絕對解耦」**，允許任意狀態共存。但前端在渲染 UI 時，扮演了重要的**「狀態防禦者」**角色。
- **實作範例**：在 `wishlist.html` 中，前端不僅會抓取願望清單，還會同步呼叫 `apiGetLibrary()` 與 `apiGetCart()`。透過前端的陣列比對 (Array.some)，動態將按鈕轉換為「已在遊戲庫」或「已在購物車」的不可點擊狀態，有效避免玩家送出無效的重複購買請求，提升 UX 並減少後端無效的 Request。

### 2. 即時狀態輪詢 (Real-time Polling)
- 針對社群功能 (如 `social.html`) 中的即時訊息與未讀通知，前端採用了輕量級的**短輪詢 (Short Polling)** 機制。
- **實作機制**：利用 `setInterval` (預設為 1 秒)，定期向後端發送請求 (`GET /api/social/messages`)。若偵測到狀態改變（如新訊息），則局部更新 DOM 元素（如將大頭貼底色轉為綠色或渲染新對話氣泡），實現低延遲的即時通訊體驗，而無須引入複雜的 WebSocket。

### 3. 安全的 Markdown 渲染 (Secure Markdown Rendering)
- 遊戲介紹 (`description`) 支援豐富的 Markdown 語法。為了防範 XSS (跨站腳本攻擊)，前端在顯示這些內容時（如 `game_detail.html`, `edit_game.html`）採用了雙層防護機制：
- **解析器**: 使用 `marked.js` 將 Markdown 文本轉換為 HTML。
- **消毒器**: 緊接著將產生的 HTML 餵給 `DOMPurify` 進行嚴格的過濾，拔除所有具備潛在威脅的 `<script>` 或 `onerror` 等屬性後，才透過 `innerHTML` 寫入畫面中。

### 4. 權限與身分解耦防護 (Role & ID Resolution)
- 系統在前端負責大量的身分驗證顯示邏輯（例如：是否解鎖開發者的官方評論介面）。由於後端不同 API 回傳的 Payload 鍵值存在差異（登入時回傳 `user.id`，但在某些模組可能使用 `user_id`），前端在進行身分核對時（如判斷當前登入者是否為遊戲原作者），採取了高包容性的屬性探測機制 `(user.id || user.user_id) == developerId`。這確保了開發者能順暢地解鎖 `AUTHOR` 特權，在未購買自己遊戲的情況下，依然能發布官方評論或進行樓中樓回覆。

### 5. 特權徽章動態渲染 (Privilege Badge Rendering)
- 配合後端「無痕修改 Schema 的隱藏字首機制」，前端負責在呈現層將「隱藏資訊」轉化為精美的「視覺組件」。
- **實作機制**：當透過 API 取得評論列表時，前端能直接讀取到乾淨的內文與獨立的 `posted_as_role` 屬性（這是後端解析 `[ROLE:XXX]` 字首後動態剝離產生的）。前端會根據此屬性，動態在評論者的名稱旁渲染專屬色系的特權徽章 (例如：`ADMIN` 為紅色，`AUTHOR` 為橘黃色，`CSR` 為藍色)，以凸顯官方留言的權威性，同時完美隱藏了底層字首的實作細節。
