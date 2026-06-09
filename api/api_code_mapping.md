# VaporAuror API 與後端程式碼對應表 (API Code Mapping)

這份文件詳細列出了本專案中所有 RESTful API 端點，以及它們在 Go 語言後端中對應的**路由註冊位置 (Router)** 與**負責處理邏輯的控制器函式 (Controller Function)**。

> 💡 **超便利功能**：您可以直接點擊 `對應的控制器函式` 欄位中的連結，系統會自動幫您跳轉到該檔案的確切行數，非常方便您截圖！

> **程式碼來源說明**：
> - 路由註冊統一集中於 `backend/routes/routes.go`。
> - 業務邏輯實作統一集中於 `backend/controllers/` 目錄下的各個 Controller 檔案。

---

## 1. 使用者與權限 (Users & Auth)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **POST** | `/api/auth/register` | `auth.POST("/register", ...)` | [auth_controller.go:27](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/auth_controller.go#L27) | 註冊新帳號 |
| **POST** | `/api/auth/login` | `auth.POST("/login", ...)` | [auth_controller.go:91](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/auth_controller.go#L91) | 使用者登入 |
| **POST** | `/api/auth/logout` | `auth.POST("/logout", ...)` | [auth_controller.go:143](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/auth_controller.go#L143) | 使用者登出 (需 JWT) |
| **PUT** | `/api/users/profile` | `users.PUT("/profile", ...)` | [user_controller.go:25](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/user_controller.go#L25) | 修改個人資料 (需 JWT) |
| **GET** | `/api/admin/users` | `admin.GET("/users", ...)` | [user_controller.go:13](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/user_controller.go#L13) | 查看所有使用者清單 (需 ADMIN) |
| **PUT** | `/api/admin/users/{id}/suspend` | `admin.PUT("/users/:id/suspend", ...)` | [admin_controller.go:12](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/admin_controller.go#L12) | 切換帳號停權狀態 (需 ADMIN) |
| **DELETE** | `/api/admin/users/{id}` | `admin.DELETE("/users/:id", ...)` | [admin_controller.go:38](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/admin_controller.go#L38) | 移除帳號 (需 ADMIN) |
| **PUT** | `/api/admin/users/{id}/role` | `admin.PUT("/users/:id/role", ...)` | [admin_controller.go:69](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/admin_controller.go#L69) | 更改帳號權限 (需 ADMIN) |

---

## 2. 商店與遊戲 (Store & Games)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **GET** | `/api/games` | `api.GET("/games", ...)` | [game_controller.go:37](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/game_controller.go#L37) | 取得所有遊戲 (含搜尋/篩選) |
| **GET** | `/api/games/{id}` | `api.GET("/games/:id", ...)` | [game_controller.go:112](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/game_controller.go#L112) | 取得單一遊戲詳情 |
| **GET** | `/api/games/{id}/reviews` | `api.GET("/games/:id/reviews", ...)` | [social_controller.go:48](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L48) | 取得遊戲評論 |
| **GET** | `/api/developer/games` | `developer.GET("/games", ...)` | [developer_controller.go:29](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L29) | 查看自己的遊戲列表 (需 DEV) |
| **POST** | `/api/developer/games` | `developer.POST("/games", ...)` | [developer_controller.go:48](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L48) | 建立新遊戲草稿 (需 DEV) |
| **PUT** | `/api/developer/games/{id}/publish` | `developer.PUT("/games/:id/publish", ...)` | [developer_controller.go:79](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L79) | 正式上架遊戲 (需 DEV) |
| **PUT** | `/api/developer/games/{id}` | `developer.PUT("/games/:id", ...)` | [developer_controller.go:110](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L110) | 編輯遊戲資訊 (需 DEV) |
| **DELETE** | `/api/developer/games/{id}` | `developer.DELETE("/games/:id", ...)` | [developer_controller.go:145](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L145) | 下架自己的遊戲 (需 DEV) |
| **DELETE** | `/api/admin/games/{id}` | `admin.DELETE("/games/:id", ...)` | [admin_controller.go:88](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/admin_controller.go#L88) | 強制下架遊戲 (需 ADMIN) |
| **POST** | `/api/developer/games/{id}/media` | `developer.POST("/games/:id/media", ...)` | [developer_controller.go:181](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L181) | 上傳遊戲圖片或主檔 (需 DEV) |
| **DELETE** | `/api/developer/games/{id}/media/{id}`| `developer.DELETE("/games/:id/media/:media_id", ...)`| [developer_controller.go:278](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L278) | 刪除遊戲素材 (需 DEV) |
| **GET** | `/api/developer/games/{id}/stats` | `developer.GET("/games/:id/stats", ...)` | [developer_controller.go:330](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L330) | 查看遊戲銷售量與收入 (需 DEV) |
| **GET** | `/api/tags` | `api.GET("/tags", ...)` | [developer_controller.go:361](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L361) | 查看標籤列表 |
| **POST** | `/api/developer/tags` | `developer.POST("/tags", ...)` | [developer_controller.go:371](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L371) | 建立標籤 (需 DEV) |
| **POST** | `/api/developer/games/{id}/tags` | `developer.POST("/games/:id/tags", ...)` | [developer_controller.go:402](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L402) | 貼上標籤 (需 DEV) |
| **DELETE** | `/api/developer/games/{id}/tags/{id}` | `developer.DELETE("/games/:id/tags/:tag_id", ...)`| [developer_controller.go:436](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/developer_controller.go#L436) | 移除標籤 (需 DEV) |

---

## 3. 訂單、購物車與客服 (Transactions & Carts)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **GET** | `/api/protected/cart` | `protected.GET("/cart", ...)` | [cart_controller.go:16](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/cart_controller.go#L16) | 查看購物車內容 |
| **POST** | `/api/protected/cart` | `protected.POST("/cart", ...)` | [cart_controller.go:31](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/cart_controller.go#L31) | 將遊戲加入購物車 |
| **DELETE** | `/api/protected/cart/{game_id}` | `protected.DELETE("/cart/:game_id", ...)` | [cart_controller.go:81](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/cart_controller.go#L81) | 移除購物車項目 |
| **POST** | `/api/protected/checkout` | `protected.POST("/checkout", ...)` | [transaction_controller.go:15](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/transaction_controller.go#L15) | 結帳購買 |
| **GET** | `/api/protected/transactions` | `protected.GET("/transactions", ...)` | [transaction_controller.go:99](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/transaction_controller.go#L99) | 查看購買紀錄 |
| **GET** | `/api/protected/refunds` | `protected.GET("/refunds", ...)` | [social_controller.go:660](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L660) | 取得個人退款歷史紀錄 |
| **POST** | `/api/social/refunds` | `social.POST("/refunds", ...)` | [social_controller.go:166](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L166) | 申請遊戲退款 |
| **GET** | `/api/csr/refunds` | `csr.GET("/refunds", ...)` | [csr_controller.go:14](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/csr_controller.go#L14) | 取得所有退款申請 (需 CSR) |
| **PUT** | `/api/csr/refunds/{id}` | `csr.PUT("/refunds/:id", ...)` | [csr_controller.go:44](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/csr_controller.go#L44) | 同意/拒絕玩家退款申請 (需 CSR)|

---

## 4. 遊戲庫與願望清單 (Library & Wishlist)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **GET** | `/api/protected/library` | `protected.GET("/library", ...)` | [library_controller.go:15](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/library_controller.go#L15) | 顯示個人遊戲庫 |
| **GET** | `/api/protected/library/{game_id}/play` | `protected.GET("/library/:game_id/play", ...)` | [library_controller.go:97](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/library_controller.go#L97) | 玩遊戲 (驗證授權) |
| **GET** | `/api/protected/library/{game_id}/download` | `protected.GET("/library/:game_id/download", ...)` | [library_controller.go:112](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/library_controller.go#L112) | 下載遊戲 (直接串流檔案) |
| **GET** | `/api/protected/wishlist` | `protected.GET("/wishlist", ...)` | [library_controller.go:30](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/library_controller.go#L30) | 查看願望清單 |
| **POST** | `/api/protected/wishlist` | `protected.POST("/wishlist", ...)` | [library_controller.go:44](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/library_controller.go#L44) | 加入願望清單 |
| **DELETE** | `/api/protected/wishlist/{game_id}`| `protected.DELETE("/wishlist/:game_id", ...)`| [library_controller.go:83](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/library_controller.go#L83)| 移除願望清單 |

---

## 5. 社交、評論與通訊 (Social & Reviews)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **POST** | `/api/social/games/{id}/reviews` | `social.POST("/games/:id/reviews", ...)` | [social_controller.go:118](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L118) | 對遊戲發表評價 |
| **POST** | `/api/social/reviews/{id}/replies` | `social.POST("/reviews/:review_id/replies",...)`| [social_controller.go:422](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L422) | 樓中樓回覆評論 |
| **DELETE** | `/api/social/reviews/replies/{id}` | `social.DELETE("/reviews/replies/:reply_id",...)`| [social_controller.go:458](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L458)| 刪除樓中樓回覆 |
| **GET** | `/api/social/friends` | `social.GET("/friends", ...)` | [social_controller.go:221](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L221) | 查看好友列表 |
| **GET** | `/api/social/friends/requests` | `social.GET("/friends/requests", ...)` | [social_controller.go:600](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L600) | 查看待審核的好友邀請 |
| **POST** | `/api/social/friends/request` | `social.POST("/friends/request", ...)` | [social_controller.go:292](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L292) | 發送好友邀請 |
| **DELETE** | `/api/social/friends/request/{id}` | `social.DELETE("/friends/request/:id", ...)`| [social_controller.go:527](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L527)| 收回好友邀請 / 解除好友 |
| **PUT** | `/api/social/friends/request/{id}/accept`| `social.PUT("/friends/request/:id/accept",...)`| [social_controller.go:483](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L483)| 接受好友邀請 |
| **PUT** | `/api/social/friends/request/{id}/decline`| `social.PUT("/friends/request/:id/decline",...)`| [social_controller.go:505](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L505)|拒絕好友邀請 |
| **GET** | `/api/social/blacklist` | `social.GET("/blacklist", ...)` | [social_controller.go:634](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L634) | 查看黑名單列表 |
| **POST** | `/api/social/blacklist` | `social.POST("/blacklist", ...)` | [social_controller.go:548](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L548) | 將玩家加入黑名單 |
| **DELETE** | `/api/social/blacklist/{user_id}` | `social.DELETE("/blacklist/:user_id", ...)`| [social_controller.go:586](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L586) | 將玩家移除黑名單 |
| **POST** | `/api/social/messages` | `social.POST("/messages", ...)` | [social_controller.go:365](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L365) | 傳輸文字通訊給對方 |
| **GET** | `/api/social/messages/{user_id}` | `social.GET("/messages/:user_id", ...)` | [social_controller.go:394](file:///c:/Users/HP/Downloads/dbms-git/dbms/backend/controllers/social_controller.go#L394) | 顯示與某使用者的對話紀錄 |
