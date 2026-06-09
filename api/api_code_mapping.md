# VaporAuror API 與後端程式碼對應表 (API Code Mapping)

這份文件詳細列出了本專案中所有 RESTful API 端點，以及它們在 Go 語言後端中對應的**路由註冊位置 (Router)** 與**負責處理邏輯的控制器函式 (Controller Function)**。

> **程式碼來源說明**：
> - 路由註冊統一集中於 `backend/routes/routes.go`。
> - 業務邏輯實作統一集中於 `backend/controllers/` 目錄下的各個 Controller 檔案。

---

## 1. 使用者與權限 (Users & Auth)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **POST** | `/api/auth/register` | `auth.POST("/register", ...)` | `controllers.Register` | 註冊新帳號 |
| **POST** | `/api/auth/login` | `auth.POST("/login", ...)` | `controllers.Login` | 使用者登入 |
| **POST** | `/api/auth/logout` | `auth.POST("/logout", ...)` | `controllers.Logout` | 使用者登出 (需 JWT) |
| **PUT** | `/api/users/profile` | `users.PUT("/profile", ...)` | `controllers.UpdateProfile` | 修改個人資料 (需 JWT) |
| **GET** | `/api/admin/users` | `admin.GET("/users", ...)` | `controllers.GetUsers` | 查看所有使用者清單 (需 ADMIN) |
| **PUT** | `/api/admin/users/{id}/suspend` | `admin.PUT("/users/:id/suspend", ...)` | `controllers.SuspendUser` | 切換帳號停權狀態 (需 ADMIN) |
| **DELETE** | `/api/admin/users/{id}` | `admin.DELETE("/users/:id", ...)` | `controllers.DeleteUser` | 移除帳號 (需 ADMIN) |
| **PUT** | `/api/admin/users/{id}/role` | `admin.PUT("/users/:id/role", ...)` | `controllers.ChangeUserRole` | 更改帳號權限 (需 ADMIN) |

---

## 2. 商店與遊戲 (Store & Games)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **GET** | `/api/games` | `api.GET("/games", ...)` | `controllers.GetGames` | 取得所有遊戲 (含搜尋/篩選) |
| **GET** | `/api/games/{id}` | `api.GET("/games/:id", ...)` | `controllers.GetGameByID` | 取得單一遊戲詳情 |
| **GET** | `/api/games/{id}/reviews` | `api.GET("/games/:id/reviews", ...)` | `controllers.GetReviews` | 取得遊戲評論 |
| **GET** | `/api/developer/games` | `developer.GET("/games", ...)` | `controllers.GetDeveloperGames` | 查看自己的遊戲列表 (需 DEV) |
| **POST** | `/api/developer/games` | `developer.POST("/games", ...)` | `controllers.UploadGame` | 建立新遊戲草稿 (需 DEV) |
| **PUT** | `/api/developer/games/{id}/publish` | `developer.PUT("/games/:id/publish", ...)` | `controllers.PublishGame` | 正式上架遊戲 (需 DEV) |
| **PUT** | `/api/developer/games/{id}` | `developer.PUT("/games/:id", ...)` | `controllers.UpdateGame` | 編輯遊戲資訊 (需 DEV) |
| **DELETE** | `/api/developer/games/{id}` | `developer.DELETE("/games/:id", ...)` | `controllers.DeleteGame` | 下架自己的遊戲 (需 DEV) |
| **DELETE** | `/api/admin/games/{id}` | `admin.DELETE("/games/:id", ...)` | `controllers.AdminDeleteGame` | 強制下架遊戲 (需 ADMIN) |
| **POST** | `/api/developer/games/{id}/media` | `developer.POST("/games/:id/media", ...)` | `controllers.UploadMedia` | 上傳遊戲圖片或主檔 (需 DEV) |
| **DELETE** | `/api/developer/games/{id}/media/{id}`| `developer.DELETE("/games/:id/media/:media_id", ...)`| `controllers.DeleteMedia` | 刪除遊戲素材 (需 DEV) |
| **GET** | `/api/developer/games/{id}/stats` | `developer.GET("/games/:id/stats", ...)` | `controllers.GetGameStats` | 查看遊戲銷售量與收入 (需 DEV) |
| **GET** | `/api/tags` | `api.GET("/tags", ...)` | `controllers.GetTags` | 查看標籤列表 |
| **POST** | `/api/developer/tags` | `developer.POST("/tags", ...)` | `controllers.CreateTag` | 建立標籤 (需 DEV) |
| **POST** | `/api/developer/games/{id}/tags` | `developer.POST("/games/:id/tags", ...)` | `controllers.AddTagToGame` | 貼上標籤 (需 DEV) |
| **DELETE** | `/api/developer/games/{id}/tags/{id}` | `developer.DELETE("/games/:id/tags/:tag_id", ...)`| `controllers.RemoveTagFromGame` | 移除標籤 (需 DEV) |

---

## 3. 訂單、購物車與客服 (Transactions & Carts)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **GET** | `/api/protected/cart` | `protected.GET("/cart", ...)` | `controllers.GetCart` | 查看購物車內容 |
| **POST** | `/api/protected/cart` | `protected.POST("/cart", ...)` | `controllers.AddToCart` | 將遊戲加入購物車 |
| **DELETE** | `/api/protected/cart/{game_id}` | `protected.DELETE("/cart/:game_id", ...)` | `controllers.RemoveFromCart` | 移除購物車項目 |
| **POST** | `/api/protected/checkout` | `protected.POST("/checkout", ...)` | `controllers.Checkout` | 結帳購買 |
| **GET** | `/api/protected/transactions` | `protected.GET("/transactions", ...)` | `controllers.GetTransactions` | 查看購買紀錄 |
| **GET** | `/api/protected/refunds` | `protected.GET("/refunds", ...)` | `controllers.GetMyRefunds` | 取得個人退款歷史紀錄 |
| **POST** | `/api/social/refunds` | `social.POST("/refunds", ...)` | `controllers.ApplyRefund` | 申請遊戲退款 |
| **GET** | `/api/csr/refunds` | `csr.GET("/refunds", ...)` | `controllers.GetRefundRequests` | 取得所有退款申請 (需 CSR) |
| **PUT** | `/api/csr/refunds/{id}` | `csr.PUT("/refunds/:id", ...)` | `controllers.ProcessRefund` | 同意/拒絕玩家退款申請 (需 CSR)|

---

## 4. 遊戲庫與願望清單 (Library & Wishlist)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **GET** | `/api/protected/library` | `protected.GET("/library", ...)` | `controllers.GetLibrary` | 顯示個人遊戲庫 |
| **GET** | `/api/protected/library/{game_id}/play` | `protected.GET("/library/:game_id/play", ...)` | `controllers.PlayGame` | 玩遊戲 (驗證授權) |
| **GET** | `/api/protected/library/{game_id}/download` | `protected.GET("/library/:game_id/download", ...)` | `controllers.DownloadGame` | 下載遊戲 (直接串流檔案) |
| **GET** | `/api/protected/wishlist` | `protected.GET("/wishlist", ...)` | `controllers.GetWishlist` | 查看願望清單 |
| **POST** | `/api/protected/wishlist` | `protected.POST("/wishlist", ...)` | `controllers.AddToWishlist` | 加入願望清單 |
| **DELETE** | `/api/protected/wishlist/{game_id}`| `protected.DELETE("/wishlist/:game_id", ...)`| `controllers.RemoveFromWishlist`| 移除願望清單 |

---

## 5. 社交、評論與通訊 (Social & Reviews)

| HTTP 方法 | API 網址路徑 | 路由註冊 (Router) | 對應的控制器函式 (Controller) | 備註功能 |
|---|---|---|---|---|
| **POST** | `/api/social/games/{id}/reviews` | `social.POST("/games/:id/reviews", ...)` | `controllers.PostReview` | 對遊戲發表評價 |
| **POST** | `/api/social/reviews/{id}/replies` | `social.POST("/reviews/:review_id/replies",...)`| `controllers.ReplyToReview` | 樓中樓回覆評論 |
| **DELETE** | `/api/social/reviews/replies/{id}` | `social.DELETE("/reviews/replies/:reply_id",...)`| `controllers.DeleteReviewReply`| 刪除樓中樓回覆 |
| **GET** | `/api/social/friends` | `social.GET("/friends", ...)` | `controllers.GetFriends` | 查看好友列表 |
| **GET** | `/api/social/friends/requests` | `social.GET("/friends/requests", ...)` | `controllers.GetFriendRequests` | 查看待審核的好友邀請 |
| **POST** | `/api/social/friends/request` | `social.POST("/friends/request", ...)` | `controllers.SendFriendRequest` | 發送好友邀請 |
| **DELETE** | `/api/social/friends/request/{id}` | `social.DELETE("/friends/request/:id", ...)`| `controllers.RevokeFriendRequest`| 收回好友邀請 / 解除好友 |
| **PUT** | `/api/social/friends/request/{id}/accept`| `social.PUT("/friends/request/:id/accept",...)`| `controllers.AcceptFriendRequest`| 接受好友邀請 |
| **PUT** | `/api/social/friends/request/{id}/decline`| `social.PUT("/friends/request/:id/decline",...)`| `controllers.DeclineFriendRequest`|拒絕好友邀請 |
| **GET** | `/api/social/blacklist` | `social.GET("/blacklist", ...)` | `controllers.GetBlacklist` | 查看黑名單列表 |
| **POST** | `/api/social/blacklist` | `social.POST("/blacklist", ...)` | `controllers.AddBlacklist` | 將玩家加入黑名單 |
| **DELETE** | `/api/social/blacklist/{user_id}` | `social.DELETE("/blacklist/:user_id", ...)`| `controllers.RemoveBlacklist` | 將玩家移除黑名單 |
| **POST** | `/api/social/messages` | `social.POST("/messages", ...)` | `controllers.SendMessage` | 傳輸文字通訊給對方 |
| **GET** | `/api/social/messages/{user_id}` | `social.GET("/messages/:user_id", ...)` | `controllers.GetMessages` | 顯示與某使用者的對話紀錄 |
